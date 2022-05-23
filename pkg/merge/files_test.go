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
	itesting "github.com/amplia-iiot/yutil/internal/testing"
)

func init() {
	// Go to root folder to access testdata/
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..", "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	// Create tmp folder
	if !io.Exists("tmp") {
		err = os.Mkdir("tmp", 0700)
		if err != nil {
			panic(err)
		}
	}
}

func TestMergeFiles(t *testing.T) {
	for _, i := range []struct {
		base    string
		changes string
	}{
		{"base", "dev"},
		{"base", "prod"},
	} {
		merged, err := MergeFiles(fileToBeMerged(i.base), fileToBeMerged(i.changes))
		if err != nil {
			t.Fatal(err)
		}
		expectedContent := itesting.ReadFile(t, expectedFile([]string{i.base, i.changes}))
		itesting.AssertEqual(t, format(t, expectedContent), merged)
	}
}

func TestMergeAllFiles(t *testing.T) {
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
		expectedContent := itesting.ReadFile(t, expectedFile(files))
		itesting.AssertEqual(t, format(t, expectedContent), merged)
	}
}

func TestMergeFilesInvalid(t *testing.T) {
	for _, i := range []struct {
		base     string
		changes  string
		expected string
	}{
		// Parsing error
		{
			base:     "invalid",
			changes:  "dev",
			expected: "cannot unmarshal",
		},
		{
			base:     "base",
			changes:  "invalid",
			expected: "cannot unmarshal",
		},
		// Not exists
		{
			base:     "not-exists",
			changes:  "dev",
			expected: "no such file or directory",
		},
		{
			base:     "base",
			changes:  "not-exists",
			expected: "no such file or directory",
		},
	} {
		merged, err := MergeFiles(fileToBeMerged(i.base), fileToBeMerged(i.changes))
		itesting.AssertError(t, i.expected, err)
		if merged != "" {
			t.Fatalf("Should not have merged")
		}
	}
}

func TestMergeAllFilesInvalid(t *testing.T) {
	for _, i := range []struct {
		files    []string
		expected string
	}{
		// At least two
		{
			files:    []string{},
			expected: "",
		},
		{
			files:    []string{"base"},
			expected: "",
		},
		// Parsing error
		{
			files:    []string{"base", "invalid"},
			expected: "cannot unmarshal",
		},
		{
			files:    []string{"invalid", "prod"},
			expected: "cannot unmarshal",
		},
		{
			files:    []string{"invalid", "prod", "docker"},
			expected: "cannot unmarshal",
		},
		{
			files:    []string{"base", "prod", "invalid"},
			expected: "cannot unmarshal",
		},
		// Not exists
		{
			files:    []string{"base", "not-exists"},
			expected: "no such file or directory",
		},
		{
			files:    []string{"base", "prod", "not-exists"},
			expected: "no such file or directory",
		},
	} {
		merged, err := MergeAllFiles(filesToBeMerged(i.files))
		itesting.AssertError(t, i.expected, err)
		if merged != "" {
			t.Fatalf("Should not have merged")
		}
	}
}

func TestMergeStdinWithFiles(t *testing.T) {
	for _, i := range []struct {
		stdin    string
		files    []string
		expected string
	}{
		{
			stdin: "app: {env: {test: true}}",
			files: []string{"base"},
			expected: `app:
  api:
    url: http://example.com
    version: v1
  cluster:
    hosts:
    - http://one.example.com
    - http://two.example.com
  description: YAML utils
  env:
    test: true
  long-description: Common functionality for working with YAML files
  name: yutil
  version: 1.0.0
`,
		},
		{
			stdin: "app: {env: {test: true}}",
			files: []string{"base", "dev"},
			expected: `app:
  api:
    url: http://localhost:8080
    version: v1-dev
  cluster:
    hosts:
    - http://localhost:8081
    - http://localhost:8082
  description: YAML utils
  env:
    dev: true
    test: true
  long-description: Common functionality for working with YAML files
  name: yutil
  version: 1.0.0-alpha
`,
		},
		{
			stdin: "{app: {env: {test: true}}, extra: extra}",
			files: []string{"base", "prod", "docker"},
			expected: `app:
  api:
    url: http://service
    version: v1
  cluster:
    hosts:
    - http://service-1
    - http://service-2
  description: YAML utils
  env:
    docker: true
    prod: true
    test: true
  long-description: Common functionality for working with YAML files
  name: yutil
  version: 1.0.0
extra: extra
`,
		},
	} {
		itesting.SimulateStdinContent(t, i.stdin, func() {
			merged, err := MergeStdinWithFiles(filesToBeMerged(i.files))
			if err != nil {
				t.Fatal(err)
			}
			itesting.AssertEqual(t, i.expected, merged)
		})
	}
}

func TestMergeStdinWithFilesInvalid(t *testing.T) {
	for _, i := range []struct {
		stdin     string
		stdinFile os.File
		files     []string
		expected  string
	}{
		// At least one file
		{
			stdin:    "app: {env: {test: true}}",
			files:    []string{},
			expected: "slice must contain at least one file",
		},
		// Parsing error
		{
			stdin:    ";",
			files:    []string{"base"},
			expected: "cannot unmarshal",
		},
		// Not exists
		{
			stdin:    "app: {env: {test: true}}",
			files:    []string{"not-exists"},
			expected: "no such file or directory",
		},
		// Stdin error
		{
			stdinFile: *os.Stderr,
			files:     []string{"not-exists"},
			expected:  "bad file descriptor",
		},
	} {
		test := func() {
			merged, err := MergeStdinWithFiles(filesToBeMerged(i.files))
			itesting.AssertError(t, i.expected, err)
			if merged != "" {
				t.Fatalf("Should not have merged")
			}
		}
		if i.stdin != "" {
			itesting.SimulateStdinContent(t, i.stdin, test)
		} else {
			itesting.SimulateStdinFile(i.stdinFile, test)
		}
	}
}

func TestMergeAllFilesToFile(t *testing.T) {
	for _, files := range [][]string{
		{"base", "dev"},
		{"base", "prod"},
		{"base", "dev", "docker"},
		{"base", "prod", "docker"},
		// Empty values (https://github.com/amplia-iiot/yutil/issues/3)
		{"issue-3-base", "issue-3-changes"},
	} {
		tmpPath := itesting.TempFilePath(t, "merged-*.yml")
		defer os.Remove(tmpPath)
		err := MergeAllFilesToFile(filesToBeMerged(files), tmpPath)
		if err != nil {
			t.Fatal(err)
		}
		expectedContent := itesting.ReadFile(t, expectedFile(files))
		mergedContent := itesting.ReadFile(t, tmpPath)
		itesting.AssertEqual(t, format(t, expectedContent), mergedContent)
	}
}

func TestMergeAllFilesToFileInvalid(t *testing.T) {
	for _, i := range []struct {
		files    []string
		expected string
	}{
		// At least two file
		{
			files:    []string{},
			expected: "slice must contain at least two files",
		},
		{
			files:    []string{"base"},
			expected: "slice must contain at least two files",
		},
		// Parsing error
		{
			files:    []string{"base", "invalid"},
			expected: "cannot unmarshal",
		},
		// Not exists
		{
			files:    []string{"base", "not-exists"},
			expected: "no such file or directory",
		},
	} {
		tmpPath := itesting.TempFilePath(t, "merged-*.yml")
		defer os.Remove(tmpPath)
		err := MergeAllFilesToFile(filesToBeMerged(i.files), tmpPath)
		itesting.AssertError(t, i.expected, err)
		if io.Exists(tmpPath) {
			t.Fatalf("Should not have merged")
		}
	}
}

func expectedFile(files []string) string {
	// fmt.Printf("expected: %v", files)
	return fmt.Sprintf("testdata/merged/%s.yml", strings.Join(files, "-"))
}

func fileToBeMerged(file string) string {
	return fmt.Sprintf("testdata/%s.yml", file)
}

func filesToBeMerged(files []string) []string {
	completeFiles := make([]string, len(files))
	for i, file := range files {
		completeFiles[i] = fileToBeMerged(file)
	}
	return completeFiles
}
