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
	"fmt"
	"os"
	"path"
	"runtime"
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

func TestFormatFile(t *testing.T) {
	for _, file := range []string{
		"base",
		"dev",
		"docker",
		"prod",
	} {
		formatted, err := FormatFile(fileToBeFormatted(file))
		if err != nil {
			t.Fatal(err)
		}
		expectedContent := itesting.ReadFile(t, expectedFile(file))
		itesting.AssertEqual(t, expectedContent, formatted)
	}
}

func TestFormatFileInvalid(t *testing.T) {
	data := Data{
		// Parsing error
		{
			file:     "invalid",
			expected: "cannot unmarshal",
		},
		// Not exists
		{
			file:     "not-exists",
			expected: "no such file or directory",
		},
	}
	for _, d := range data {
		formatted, err := FormatFile(fileToBeFormatted(d.file))
		itesting.AssertError(t, d.expected, err)
		if formatted != "" {
			t.Fatalf("Should not have formatted")
		}
	}
}

func TestFormatFileInPlace(t *testing.T) {
	for _, file := range []string{
		"base",
		"dev",
		"docker",
		"prod",
	} {
		tmpFile := itesting.TempFilePath(t, file+"-*.yml")
		err := io.Copy(fileToBeFormatted(file), tmpFile)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile)
		err = FormatFileInPlace(tmpFile)
		if err != nil {
			t.Fatal(err)
		}
		// Formatted file
		formattedContent := itesting.ReadFile(t, tmpFile)
		expectedContent := itesting.ReadFile(t, expectedFile(file))
		itesting.AssertEqual(t, expectedContent, formattedContent)
		// Backup file not exists
		itesting.AssertFalse(t, io.Exists(tmpFile+".bak"))
	}
}

func TestFormatFileInPlaceInvalid(t *testing.T) {
	data := Data{
		// Parsing error
		{
			file:     "invalid",
			expected: "cannot unmarshal",
		},
		// Not exists
		{
			// Do not existing file, use one directly in tmp (that doesn't exists)
			tmpPath:  "tmp/not-exists",
			expected: "no such file or directory",
		},
	}
	for _, d := range data {
		copiedToTmp := false
		if d.tmpPath == "" {
			// Use real testdata
			d.tmpPath = itesting.TempFilePath(t, d.file+"-*.yml")
			err := io.Copy(fileToBeFormatted(d.file), d.tmpPath)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(d.tmpPath)
			copiedToTmp = true
		}
		err := FormatFileInPlace(d.tmpPath)
		itesting.AssertError(t, d.expected, err)
		if copiedToTmp {
			// Formatted file still has original content
			formattedContent := itesting.ReadFile(t, d.tmpPath)
			originalContent := itesting.ReadFile(t, fileToBeFormatted(d.file))
			itesting.AssertEqual(t, originalContent, formattedContent)
		}
		// Backup file not exists
		itesting.AssertFalse(t, io.Exists(d.tmpPath+".bak"))
	}
}

func TestFormatFileInPlaceB(t *testing.T) {
	for _, file := range []string{
		"base",
		"dev",
		"docker",
		"prod",
	} {
		tmpFile := itesting.TempFilePath(t, file+"-*.yml")
		err := io.Copy(fileToBeFormatted(file), tmpFile)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile)
		defer os.Remove(tmpFile + ".bak")
		err = FormatFileInPlaceB(tmpFile, ".bak")
		if err != nil {
			t.Fatal(err)
		}
		// Formatted file
		formattedContent := itesting.ReadFile(t, tmpFile)
		expectedContent := itesting.ReadFile(t, expectedFile(file))
		itesting.AssertEqual(t, expectedContent, formattedContent)
		// Backup file
		itesting.AssertTrue(t, io.Exists(tmpFile+".bak"))
		originalContent := itesting.ReadFile(t, fileToBeFormatted(file))
		backupContent := itesting.ReadFile(t, tmpFile+".bak")
		itesting.AssertEqual(t, originalContent, backupContent)
	}
}

func TestFormatFileInPlaceBInvalid(t *testing.T) {
	data := Data{
		// Parsing error
		{
			file:              "invalid",
			expected:          "cannot unmarshal",
			shouldBackupExist: true,
		},
		// Not exists
		{
			// Do not use existing file, use one directly in tmp (that doesn't exists)
			tmpPath:           "tmp/not-exists",
			expected:          "no such file or directory",
			shouldBackupExist: false,
		},
	}
	for _, d := range data {
		copiedToTmp := false
		if d.tmpPath == "" {
			// Use real testdata
			d.tmpPath = itesting.TempFilePath(t, d.file+"-*.yml")
			err := io.Copy(fileToBeFormatted(d.file), d.tmpPath)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(d.tmpPath)
			copiedToTmp = true
		}
		err := FormatFileInPlaceB(d.tmpPath, ".bak")
		itesting.AssertError(t, d.expected, err)
		if copiedToTmp {
			// Copied file still has original content
			formattedContent := itesting.ReadFile(t, d.tmpPath)
			originalContent := itesting.ReadFile(t, fileToBeFormatted(d.file))
			itesting.AssertEqual(t, originalContent, formattedContent)
		}
		// Backup file
		if d.shouldBackupExist {
			itesting.AssertTrue(t, io.Exists(d.tmpPath+".bak"))
			originalContent := itesting.ReadFile(t, fileToBeFormatted(d.file))
			backupContent := itesting.ReadFile(t, d.tmpPath+".bak")
			itesting.AssertEqual(t, originalContent, backupContent)
		} else {
			itesting.AssertFalse(t, io.Exists(d.tmpPath+".bak"))
		}
	}
}

func TestFormatFilesInPlace(t *testing.T) {
	files := []string{
		"base",
		"dev",
		"docker",
		"prod",
	}
	var tmpFiles []string
	tmpDir, err := os.MkdirTemp("tmp", "format-*")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		tmpFile := fmt.Sprintf("%s/%s.yml", tmpDir, file)
		err := io.Copy(fileToBeFormatted(file), tmpFile)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile)
		tmpFiles = append(tmpFiles, tmpFile)
	}
	err = FormatFilesInPlace(tmpFiles)
	if err != nil {
		t.Fatal(err)
	}
	for i, file := range files {
		// Formatted file
		formattedContent := itesting.ReadFile(t, tmpFiles[i])
		expectedContent := itesting.ReadFile(t, expectedFile(file))
		itesting.AssertEqual(t, expectedContent, formattedContent)
		// Backup file
		itesting.AssertFalse(t, io.Exists(tmpFiles[i]+".bak"))
	}
}

func TestFormatFilesInPlaceInvalid(t *testing.T) {
	data := Data{
		// OK
		{
			file: "base",
		},
		{
			file: "dev",
		},
		// Parsing error
		{
			file:     "invalid",
			invalid:  true,
			expected: "cannot unmarshal",
		},
		// Not exists
		{
			// Do not existing file, use one directly in tmp (that doesn't exists)
			tmpPath: "tmp/not-exists",
			invalid: true,
		},
		// Will fail because there is a previous error
		{
			file:    "prod",
			invalid: true,
		},
	}
	// Prepare tmp folder
	for i := 0; i < len(data); i++ {
		d := &data[i]
		d.copiedToTmp = false
		if d.tmpPath == "" {
			// Use real testdata
			d.tmpPath = itesting.TempFilePath(t, d.file+"-*.yml")
			err := io.Copy(fileToBeFormatted(d.file), d.tmpPath)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(d.tmpPath)
			d.copiedToTmp = true
		}
	}
	err := FormatFilesInPlace(data.tmpPaths())
	for _, d := range data {
		if d.invalid {
			itesting.AssertError(t, d.expected, err)
		}
		if d.copiedToTmp {
			if d.invalid {
				// Formatted file still has original content
				formattedContent := itesting.ReadFile(t, d.tmpPath)
				originalContent := itesting.ReadFile(t, fileToBeFormatted(d.file))
				itesting.AssertEqual(t, originalContent, formattedContent)
			} else {
				// Formatted file
				formattedContent := itesting.ReadFile(t, d.tmpPath)
				expectedContent := itesting.ReadFile(t, expectedFile(d.file))
				itesting.AssertEqual(t, expectedContent, formattedContent)
			}
		}
		// Backup file does not exists
		itesting.AssertFalse(t, io.Exists(d.tmpPath+".bak"))
	}
}

func TestFormatFilesInPlaceB(t *testing.T) {
	files := []string{
		"base",
		"dev",
		"docker",
		"prod",
	}
	var tmpFiles []string
	tmpDir, err := os.MkdirTemp("tmp", "format-*")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		tmpFile := fmt.Sprintf("%s/%s.yml", tmpDir, file)
		err := io.Copy(fileToBeFormatted(file), tmpFile)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile)
		tmpFiles = append(tmpFiles, tmpFile)
	}
	err = FormatFilesInPlaceB(tmpFiles, ".bak")
	if err != nil {
		t.Fatal(err)
	}
	for i, file := range files {
		// Formatted file
		formattedContent := itesting.ReadFile(t, tmpFiles[i])
		expectedContent := itesting.ReadFile(t, expectedFile(file))
		itesting.AssertEqual(t, expectedContent, formattedContent)
		// Backup file
		itesting.AssertTrue(t, io.Exists(tmpFiles[i]+".bak"))
		originalContent := itesting.ReadFile(t, fileToBeFormatted(file))
		backupContent := itesting.ReadFile(t, tmpFiles[i]+".bak")
		itesting.AssertEqual(t, originalContent, backupContent)
	}
}

func TestFormatFilesInPlaceBInvalid(t *testing.T) {
	data := Data{
		// OK
		{
			file:              "base",
			shouldBackupExist: true,
		},
		{
			file:              "dev",
			shouldBackupExist: true,
		},
		// Parsing error
		{
			file:              "invalid",
			invalid:           true,
			expected:          "cannot unmarshal",
			shouldBackupExist: true,
		},
		// Not exists
		{
			// Do not existing file, use one directly in tmp (that doesn't exists)
			tmpPath: "tmp/not-exists",
			invalid: true,
		},
		// Will fail because there is a previous error
		{
			file:    "prod",
			invalid: true,
		},
	}
	// Prepare tmp folder
	for i := 0; i < len(data); i++ {
		d := &data[i]
		d.copiedToTmp = false
		if d.tmpPath == "" {
			// Use real testdata
			d.tmpPath = itesting.TempFilePath(t, d.file+"-*.yml")
			err := io.Copy(fileToBeFormatted(d.file), d.tmpPath)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(d.tmpPath)
			d.copiedToTmp = true
		}
	}
	err := FormatFilesInPlaceB(data.tmpPaths(), ".bak")
	for _, d := range data {
		if d.invalid {
			itesting.AssertError(t, d.expected, err)
		}
		if d.copiedToTmp {
			if d.invalid {
				// Formatted file still has original content
				formattedContent := itesting.ReadFile(t, d.tmpPath)
				originalContent := itesting.ReadFile(t, fileToBeFormatted(d.file))
				itesting.AssertEqual(t, originalContent, formattedContent)
			} else {
				// Formatted file
				formattedContent := itesting.ReadFile(t, d.tmpPath)
				expectedContent := itesting.ReadFile(t, expectedFile(d.file))
				itesting.AssertEqual(t, expectedContent, formattedContent)
			}
		}
		if d.shouldBackupExist {
			// Backup file
			itesting.AssertTrue(t, io.Exists(d.tmpPath+".bak"))
			originalContent := itesting.ReadFile(t, fileToBeFormatted(d.file))
			backupContent := itesting.ReadFile(t, d.tmpPath+".bak")
			itesting.AssertEqual(t, originalContent, backupContent)
		} else {
			// Backup file does not exists
			itesting.AssertFalse(t, io.Exists(d.tmpPath+".bak"))
		}
	}
}

func expectedFile(file string) string {
	return fmt.Sprintf("testdata/formatted/%s.yml", file)
}

func fileToBeFormatted(file string) string {
	return fmt.Sprintf("testdata/%s.yml", file)
}

type Data []struct {
	file              string // Use real testdata file as base for test (will be copied to tmp)
	tmpPath           string // Custom data file for test (if left empty will be a tmp of file)
	invalid           bool   // Whether the file will fail
	expected          string // Expected error message
	shouldBackupExist bool   // Whether backup should exist
	copiedToTmp       bool   // Whether it has been copied to tmp folder
}

func (d Data) tmpPaths() []string {
	var tmpPaths []string
	for _, data := range d {
		tmpPaths = append(tmpPaths, data.tmpPath)
	}
	return tmpPaths
}
