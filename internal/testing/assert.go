/*
Copyright (c) 2021-2023 amplia-iiot

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
	"reflect"
	"strings"
	"testing"
)

// AssertEqual fails if expected and got are not equal.
func AssertEqual(t *testing.T, expected interface{}, got interface{}) {
	if expected == got {
		return
	}
	t.Fatalf("Received %v (type %v), expected %v (type %v)", got, reflect.TypeOf(got), expected, reflect.TypeOf(expected))
}

// AssertError fails if an error is expected and does not contain the expected
// string.
func AssertError(t *testing.T, expected string, got error) {
	if got == nil && expected != "" {
		t.Fatalf("Error expected and not triggered")
	}
	if expected != "" && !strings.Contains(got.Error(), expected) {
		t.Fatalf("Error '%s' does not contain '%s'", got, expected)
	}
	if expected == "" && got != nil {
		t.Fatalf("Error '%s' not expected", got)
	}
}

// AssertTrue fails if a false is passed.
func AssertTrue(t *testing.T, got bool) {
	if !got {
		t.Fatal("True expected")
	}
}

// AssertFalse fails if a true is passed.
func AssertFalse(t *testing.T, got bool) {
	if got {
		t.Fatal("False expected")
	}
}
