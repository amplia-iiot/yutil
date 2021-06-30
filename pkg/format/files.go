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

package format

import (
	"github.com/amplia-iiot/yutil/internal/io"
)

func FormatFile(file string) (string, error) {
	content, err := io.ReadAsString(file)
	if err != nil {
		return "", err
	}
	return FormatContent(content)
}

func FormatFileInPlace(file string) error {
	formatted, err := FormatFile(file)
	if err != nil {
		return err
	}
	return io.WriteToFile(file, formatted)
}

func FormatFileInPlaceB(file, backupSuffix string) error {
	if backupSuffix != "" {
		err := io.Copy(file, file+backupSuffix)
		if err != nil {
			return err
		}
	}
	return FormatFileInPlace(file)
}

func FormatFilesInPlace(files []string) error {
	for _, file := range files {
		err := FormatFileInPlace(file)
		if err != nil {
			return err
		}
	}
	return nil
}
func FormatFilesInPlaceB(files []string, backupSuffix string) error {
	for _, file := range files {
		err := FormatFileInPlaceB(file, backupSuffix)
		if err != nil {
			return err
		}
	}
	return nil
}
