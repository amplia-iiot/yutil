/*
Copyright (c) 2023 Adrian Haasler García <dev@ahaasler.com>

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

	"github.com/amplia-iiot/yutil/internal/io"
	"github.com/amplia-iiot/yutil/pkg/replace"
	"github.com/spf13/cobra"
)

type replaceOptions struct {
	golang           bool
	jinja2           bool
	directory        string
	node             string
	include          []string
	exclude          []string
	replacementFiles []string
	includeEnv       bool
	extensions       []string
}

func (o replaceOptions) engine() replace.Engine {
	if o.jinja2 {
		return replace.Jinja2
	}
	return replace.Golang
}

var rOptions replaceOptions

// replaceCmd represents the replace command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace files using a template engine",
	Long: `Replace files with a template engine, using one or
more replacement files that will be merged. Files should be ordered
in ascending level of importance in the hierarchy. A yaml
node in the last file replaces values in any previous file. Stdin is
treated as first replacement file.

By default all files in the current directory and subdirectories are passed
through the template engine. Use include and exclude to filter files.
If no engine is picked the default golang template with slim-sprig functions
will be used.

The extension/s is/are used for including those files by default (unless
include flag is used) and renaming replaced files accordingly.

For example:

yutil replace -r base.yml -r changes.yml -d directory
cat base.yml | yutil replace
yutil replace -r config.yml -n root_node --jinja2
yutil replace -r config.yml -e .go -e .gotempl --env
yutil replace -r config.yml --include 'directory/*.conf'
yutil replace -r config.yml --jinja2 -d directory --exclude '*/secret/*'
echo "this is not a yaml" | yutil --no-input replace -r base.yml -r changes.yml
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		engine := rOptions.engine()
		extensions := rOptions.extensions
		include := rOptions.include
		switch engine {
		case replace.Golang:
			if len(extensions) == 0 {
				extensions = append(extensions, ".tmpl")
			}
		case replace.Jinja2:
			if len(extensions) == 0 {
				extensions = append(extensions, ".j2")
			}
		}
		if len(include) == 0 {
			for _, ext := range extensions {
				include = append(include, fmt.Sprintf("*%s.*", ext), fmt.Sprintf("*%s", ext))
			}
		}
		if !io.Exists(rOptions.directory) {
			return fmt.Errorf("directory %s does not exist", rOptions.directory)
		}
		for _, f := range rOptions.replacementFiles {
			if !io.Exists(f) {
				return fmt.Errorf("replacement file %s does not exist", f)
			}
		}
		err := replace.Replace(
			engine,
			replace.WithDirectory(rOptions.directory),
			replace.WithReplacementFiles(rOptions.replacementFiles...),
			replace.WithRootNode(rOptions.node),
			replace.WithExtension(extensions...),
			replace.WithInclude(include...),
			replace.WithExclude(rOptions.exclude...),
			replace.WithIncludeEnvironmentInReplacements(rOptions.includeEnv),
			replace.WithIncludeStdinInReplacements(canAccessStdin()),
		)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)

	replaceCmd.Flags().BoolVar(&rOptions.golang, "golang", false, "use golang template engine (default), automatically sets include and extension config for .tmpl files unless overriden")
	replaceCmd.Flags().BoolVar(&rOptions.jinja2, "jinja2", false, "use jinja2 template engine, automatically sets include and extension config for .j2 files unless overriden")
	replaceCmd.MarkFlagsMutuallyExclusive("golang", "jinja2")
	replaceCmd.Flags().StringVarP(&rOptions.directory, "directory", "d", ".", "pick root directory to search for files to be replaced (defaults to current directory)")
	replaceCmd.Flags().StringSliceVarP(&rOptions.replacementFiles, "replacements", "r", []string{}, "replacement files (multiple files will be merged)")
	replaceCmd.Flags().StringVarP(&rOptions.node, "node", "n", "", "only include replacements from inside this node")
	replaceCmd.Flags().BoolVar(&rOptions.includeEnv, "env", false, "include environment variables as input for the template engine (available inside the 'env' node)")
	replaceCmd.Flags().StringSliceVarP(&rOptions.extensions, "extension", "e", []string{}, "define the extension/s of the files to replace and then remove in the file name when saving (normally you should include the dot), automatically sets default include config unless overriden (*<ext> and *<ext>.*)")
	replaceCmd.Flags().StringSliceVar(&rOptions.include, "include", []string{}, "include files that match the filter/s")
	replaceCmd.Flags().StringSliceVar(&rOptions.exclude, "exclude", []string{}, "exclude files that match the filter/s (takes precedence over include)")
	onViperInitialize(func() {
		bindViperC(replaceCmd, "golang", "replace.golang")
		bindViperC(replaceCmd, "jinja2", "replace.jinja2")
		bindViperC(replaceCmd, "directory", "replace.directory")
		bindViperC(replaceCmd, "replacements", "replace.replacements")
		bindViperC(replaceCmd, "node", "replace.node")
		bindViperC(replaceCmd, "env", "replace.env")
		bindViperC(replaceCmd, "extension", "replace.extension")
		bindViperC(replaceCmd, "include", "replace.include")
		bindViperC(replaceCmd, "exclude", "replace.exclude")
	})
}
