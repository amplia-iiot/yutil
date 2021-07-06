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
	"errors"
	"fmt"

	"github.com/amplia-iiot/yutil/internal/io"
	"github.com/amplia-iiot/yutil/pkg/format"
	"github.com/spf13/cobra"
)

type formatOptions struct {
	outputFile string
	inPlace    bool
	suffix     string
}

var fOptions formatOptions

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format [FILE...]",
	Short: "Format a yaml file",
	Long: `Format a yaml file ordering its keys alphabetically and
cleaning it.

For example:

yutil format file.yml
yutil format file.yml -o file.formatted.yml
cat file.yml | yutil format > file.formatted.yml
echo "this is not a yaml" | yutil --no-input format file.yml > file.formatted.yml
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if inPlaceEnabled(cmd) {
			if canAccessStdin() {
				return errors.New("stdin not compatible with in place format")
			}
			if fOptions.outputFile != "" {
				return errors.New("output option not compatible with in place format")
			}
		} else {
			if canAccessStdin() && len(args) != 0 {
				return errors.New("only one yaml can be formatted to output, stdin is active")
			} else if !canAccessStdin() && len(args) == 0 {
				if stdinBlocked() {
					return errors.New("requires one file to be formatted, stdin is blocked")
				} else {
					return errors.New("requires one file to be formatted")
				}
			} else if !canAccessStdin() && len(args) != 1 {
				return errors.New("only one file can be formatted to output")
			}
		}
		for _, file := range args {
			if !io.Exists(file) {
				return fmt.Errorf("file %s does not exist", file)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if inPlaceEnabled(cmd) {
			if fOptions.suffix == "" {
				err = format.FormatFilesInPlace(args)
			} else {
				err = format.FormatFilesInPlaceB(args, fOptions.suffix)
			}
		} else {
			var formatted string
			if canAccessStdin() {
				formatted, err = format.FormatStdin()
			} else {
				formatted, err = format.FormatFile(args[0])
			}
			if err != nil {
				panic(err)
			}
			if len(fOptions.outputFile) > 0 {
				err = io.WriteToFile(fOptions.outputFile, formatted)
			} else {
				err = io.WriteToStdout(formatted)
			}
		}
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)

	formatCmd.Flags().StringVarP(&fOptions.outputFile, "output", "o", "", "format yaml to output file instead of stdout (not compatible in place format)")
	formatCmd.Flags().BoolVarP(&fOptions.inPlace, "in-place", "i", false, "format yaml files in place (makes backup if suffix is supplied)")
	formatCmd.Flags().StringVarP(&fOptions.suffix, "suffix", "s", "", "format yaml files in place making a backup with the given suffix (-i is not necessary if suffix is passed)")
}

// Whether in place format is enabled
func inPlaceEnabled(cmd *cobra.Command) bool {
	return fOptions.inPlace || cmd.Flags().Changed("suffix")
}
