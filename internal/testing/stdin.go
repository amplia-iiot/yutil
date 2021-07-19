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
)

// SimulateStdinContent executes a function that reads stdin simulating the
// content. Fails on error.
func SimulateStdinContent(t *testing.T, stdin string, function func()) {
	// Create temporal file
	tmp, err := os.CreateTemp("tmp", "stdin-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	// Clean up on exit
	defer os.Remove(tmp.Name())

	// Write custom content
	if _, err := tmp.Write([]byte(stdin)); err != nil {
		t.Fatal(err)
	}

	// Reset offset for next read
	if _, err := tmp.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	SimulateStdinFile(*tmp, function)
}

// SimulateStdinFile executes a function that reads stdin using another file as
// simulated stdin. Fails on error.
func SimulateStdinFile(stdin os.File, function func()) {
	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()

	os.Stdin = &stdin

	// Execute function that uses stdin
	function()
}
