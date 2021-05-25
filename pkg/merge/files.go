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

func MergeAllFilesToFile(files []string, output string) error {
	merged, err := MergeAllFiles(files)
	if err != nil {
		return err
	}
	return io.WriteToFile(output, merged)
}
