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
	"testing"

	itesting "github.com/amplia-iiot/yutil/internal/testing"
)

func TestFormatContent(t *testing.T) {
	for _, i := range []struct {
		content  string
		expected string
	}{
		// Formatted to multiple lines
		{
			content: `data: {one: 1, two: 2}`,
			expected: `data:
  one: 1
  two: 2
`,
		},
		// Keys are alphabetically ordered
		{
			content: `data: {b: b, c: c, a: a}`,
			expected: `data:
  a: a
  b: b
  c: c
`,
		},
		// Null should be formatted `null`
		{
			content: `data:`,
			expected: `data: null
`,
		},
		{
			content: `data: !!null null`,
			expected: `data: null
`,
		},
		{
			content: `data: null`,
			expected: `data: null
`,
		},
		{
			content: `data: Null`,
			expected: `data: null
`,
		},
		// Comments are lost
		{
			content: `with-comment: 42 # The meaning of life`,
			expected: `with-comment: 42
`,
		},
		// Numbers
		{
			content: `data: !!float -1`,
			expected: `data: -1
`,
		},
		{
			content: `data: !!float 0`,
			expected: `data: 0
`,
		},
		{
			content: `data: !!float 2.3e4`,
			expected: `data: 23000
`,
		},
		{
			content: `data: !!float .inf`,
			expected: `data: .inf
`,
		},
		{
			content: `data: !!float .nan`,
			expected: `data: .nan
`,
		},
		// Booleans
		{
			content: `data: false`,
			expected: `data: false
`,
		},
		{
			content: `data: true`,
			expected: `data: true
`,
		},
		{
			content: `data: False`,
			expected: `data: false
`,
		},
		{
			content: `data: !!bool "true"`,
			expected: `data: true
`,
		},
		// String do not have quotes
		{
			content: `data: "one"`,
			expected: `data: one
`,
		},
		{
			content: `data: 'this is a string'`,
			expected: `data: this is a string
`,
		},
		// Multi-line
		{
			content: `plain:
  This unquoted scalar
  spans many lines.`,
			expected: `plain: This unquoted scalar spans many lines.
`,
		},
		{
			content: `quoted: "So does this
  quoted scalar.\n"`,
			expected: `quoted: |
  So does this quoted scalar.
`,
		},
		{
			content: `quoted: "So does this quoted scalar.\n"`,
			expected: `quoted: |
  So does this quoted scalar.
`,
		},
		// Strings with escaped quote chars don't need to be quoted
		{
			content: `data: "Hello \"World\""`,
			expected: `data: Hello "World"
`,
		},
		{
			content: `data: 'Hello "World"'`,
			expected: `data: Hello "World"
`,
		},
		{
			content: `data: "Hello \'World\'"`,
			expected: `data: Hello 'World'
`,
		},
		// Unless starting with quote or containing special chars (with simple quotes)
		{
			content: `data: "\'Hello\' World"`,
			expected: `data: '''Hello'' World'
`,
		},
		{
			content: `data: "\"Hello\" World"`,
			expected: `data: '"Hello" World'
`,
		},
		{
			content: `data: "{"`,
			expected: `data: '{'
`,
		},
		{
			content: `tie-fighter: '|\-*-/|'`,
			expected: `tie-fighter: '|\-*-/|'
`,
		},
		{
			content: `tie-fighter: "|\\-*-/|"`,
			expected: `tie-fighter: '|\-*-/|'
`,
		},
		{
			content: `not-a-comment: '# Not a ''comment''.'`,
			expected: `not-a-comment: '# Not a ''comment''.'
`,
		},
		// Unicode is formatted
		{
			content: `unicode: "Sosa did fine.\u263A"`,
			expected: `unicode: Sosa did fine.☺
`,
		},
		// Strings containing data that otherwise will be another type do need to be quoted with double quotes
		{
			content: `data: "null"`,
			expected: `data: "null"
`,
		},
		{
			content: `data: 'null'`,
			expected: `data: "null"
`,
		},
		{
			content: `data: !!str null`,
			expected: `data: "null"
`,
		},
		{
			content: `data: "123"`,
			expected: `data: "123"
`,
		},
		{
			content: `data: '123'`,
			expected: `data: "123"
`,
		},
		{
			content: `data: !!str 123`,
			expected: `data: "123"
`,
		},
		{
			content: `data: '.nan'`,
			expected: `data: ".nan"
`,
		},
		{
			content: `data: !!str .nan`,
			expected: `data: ".nan"
`,
		},
		{
			content: `data: '.inf'`,
			expected: `data: ".inf"
`,
		},
		{
			content: `data: !!str .inf`,
			expected: `data: ".inf"
`,
		},
		{
			content: `data: "false"`,
			expected: `data: "false"
`,
		},
		{
			content: `data: 'false'`,
			expected: `data: "false"
`,
		},
		{
			content: `data: 'False'`,
			expected: `data: "False"
`,
		},
		{
			content: `data: !!str false`,
			expected: `data: "false"
`,
		},
		// Arrays are not reordered
		{
			content: `data: [b, c, a]`,
			expected: `data:
- b
- c
- a
`,
		},
	} {
		formatted, err := FormatContent(i.content)
		if err != nil {
			t.Fatal(err)
		}
		itesting.AssertEqual(t, i.expected, formatted)
	}
}
