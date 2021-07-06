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

package merge

import (
	"errors"

	"github.com/amplia-iiot/yutil/internal/io"
)

// MergeFiles returns the result of merging two yaml files. A yaml leaf node in
// the 'changes' file takes precedence over and replaces the value in the 'base'
// file.
func MergeFiles(base string, changes string) (string, error) {
	baseContent, err := io.ReadAsString(base)
	if err != nil {
		return "", err
	}
	changesContent, err := io.ReadAsString(changes)
	if err != nil {
		return "", err
	}
	return MergeContents(baseContent, changesContent)
}

// MergeAllFiles returns the result of merging all yaml files, which should be
// ordered in ascending level of importance in the hierarchy. A yaml leaf node in
// the last file takes precedence over and replaces the value in any previous
// file.
func MergeAllFiles(files []string) (string, error) {
	if len(files) < 2 {
		return "", errors.New("slice must contain at least two files")
	}
	contents := make([]string, len(files))

	for i, f := range files {
		var err error
		contents[i], err = io.ReadAsString(f)
		if err != nil {
			return "", err
		}
	}
	return MergeAllContents(contents)
}

// MergeStdinWithFiles returns the result of merging stdin as yaml content with
// all yaml files, which should be ordered in ascending level of importance in
// the hierarchy. Stdin is the least important yaml. A yaml leaf node in the last
// file takes precedence over and replaces the value in any previous file,
// including values in stdin.
func MergeStdinWithFiles(files []string) (string, error) {
	if len(files) < 1 {
		return "", errors.New("slice must contain at least one files")
	}
	var err error

	contents := make([]string, len(files)+1)

	contents[0], err = io.ReadStdin()
	if err != nil {
		return "", err
	}

	for i, f := range files {
		contents[i+1], err = io.ReadAsString(f)
		if err != nil {
			return "", err
		}
	}
	return MergeAllContents(contents)
}

// MergeAllFilesToFile writes to an output file the result of merging all yaml
// files, which should be ordered in ascending level of importance in the
// hierarchy. A yaml leaf node in the last file takes precedence over and
// replaces the value in any previous file.
func MergeAllFilesToFile(files []string, output string) error {
	merged, err := MergeAllFiles(files)
	if err != nil {
		return err
	}
	return io.WriteToFile(output, merged)
}
