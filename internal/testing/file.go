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

package testing

import (
	"os"
	"testing"

	"github.com/amplia-iiot/yutil/internal/io"
)

// TempFilePath return a temp file path for the given filename pattern for a file
// that does not exist. Fails on error.
func TempFilePath(t *testing.T, pattern string) string {
	tmp, err := os.CreateTemp("tmp", pattern)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	return tmp.Name()
}

// ReadFile returns the content of a file. Fails on error.
func ReadFile(t *testing.T, file string) string {
	content, err := io.ReadAsString(file)
	if err != nil {
		t.Fatal(err)
	}
	return content
}
