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
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/amplia-iiot/yutil/internal/io"
)

func init() {
	// Go to root folder to access testdata/
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..", "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestMergeFiles(t *testing.T) {
	for _, files := range [][]string{
		{"base", "dev"},
		{"base", "prod"},
		{"base", "dev", "docker"},
		{"base", "prod", "docker"},
	} {
		merged, err := MergeAllFiles(filesToBeMerged(files))
		if err != nil {
			t.Fatal(err)
		}
		expectedContent, err := io.ReadAsString(expectedFile(files))
		if err != nil {
			t.Fatal(err)
		}
		assert(t, expectedContent, merged)
	}
}

func expectedFile(files []string) string {
	// fmt.Printf("expected: %v", files)
	return fmt.Sprintf("testdata/merged/%s.yml", strings.Join(files, "-"))
}

func filesToBeMerged(files []string) []string {
	completeFiles := make([]string, len(files))
	for i, file := range files {
		completeFiles[i] = fmt.Sprintf("testdata/%s.yml", file)
	}
	return completeFiles
}
