/*
Copyright (c) 2021 amplia-iiot

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/amplia-iiot/yutil/internal/io"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	defaultConfigFile = ".yutil"
	envPrefix         = "YUTIL"
)

var cfgFile string
var viperInitializers []func()

var rootCmd = &cobra.Command{
	Use:   "yutil",
	Short: "YAML utils",
	Long:  `Common functionality for working with YAML files`,
}

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

func (b *BuildInfo) format() string {
	return fmt.Sprintf("yutil %s, built at %s by %s from %s", b.Version, b.Date, b.BuiltBy, b.Commit)
}

var buildInfo BuildInfo

func Execute(info BuildInfo) {
	buildInfo = info
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.yutil or $HOME/.yutil)")
	rootCmd.PersistentFlags().Bool("no-input", false, "ignore stdin input (by default stdin is read as yaml content)")
	OnViperInitialize(func() {
		bindViper(rootCmd, "no-input")
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// viper := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".yutil" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(defaultConfigFile)
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if cfgFile != "" {
		fmt.Fprintln(os.Stderr, "Error using config file,", err)
	}

	// Call viper initializers
	for _, initializer := range viperInitializers {
		initializer()
	}
}

// OnViperInitialize sets the passed functions to be run when viper is ready to be configured
func OnViperInitialize(y ...func()) {
	viperInitializers = append(viperInitializers, y...)
}

// Binds a cobra flag with the same viper config key
func bindViper(cmd *cobra.Command, name string) {
	bindViperC(cmd, name, name)
}

// Binds a cobra flag with a custom viper config key
func bindViperC(cmd *cobra.Command, cobraName string, viperName string) {
	err := viper.BindEnv(viperName, fmt.Sprintf("%s_%s", envPrefix, strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(viperName, ".", "_"), "-", "_"))))
	if err != nil {
		panic(err)
	}
	if !cmd.Flag(cobraName).Changed && viper.IsSet(viperName) {
		err := cmd.Flags().Set(cobraName, fmt.Sprintf("%v", viper.Get(viperName)))
		if err != nil {
			panic(err)
		}
	}
}

// Whether stdin is accessible (received and not blocked with --no-input)
func canAccessStdin() bool {
	noInput, err := rootCmd.Flags().GetBool("no-input")
	return err == nil && io.ReceivedStdin() && !noInput
}
