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

	"github.com/amplia-iiot/yutil/pkg/merge"
	"github.com/spf13/cobra"
)

var outputFile string

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge FILE [FILE...]",
	Short: "Merge yaml files",
	Long: `Merge as many yaml files as desired. Files should be ordered
in ascending level of importance in the hierarchy. A yaml
node in the last file replaces values in any previous file.

For example:

yutil merge base.yml changes.yml
yutil merge base.yml changes.yml important.yml
yutil merge base.yml changes.yml -o merged.yml
cat base.yml | yutil merge changes.yml > merged.yml
echo "this is not a yaml" | yutil --no-input merge base.yml changes.yml
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if canAccessStdin() && len(args) < 1 {
			return errors.New("requires at least one file to be merged with stdin")
		} else if !canAccessStdin() && len(args) < 2 {
			return errors.New("requires at least two files to be merged")
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
		var merged string
		if canAccessStdin() {
			merged, err = merge.MergeStdinWithFiles(args)
		} else {
			merged, err = merge.MergeAllFiles(args)
		}
		if err != nil {
			panic(err)
		}
		if len(outputFile) > 0 {
			err = io.WriteToFile(outputFile, merged)
		} else {
			err = io.WriteToStdout(merged)
		}
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "write merged yaml to output file instead of stdout")
	OnViperInitialize(func() {
		bindViperC(mergeCmd, "output", "merge.output")
	})
}
