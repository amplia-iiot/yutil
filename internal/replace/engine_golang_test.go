/*
Copyright (c) 2023 Adrian Haasler García <dev@ahaasler.com>

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

package replace

import (
	"testing"

	itesting "github.com/amplia-iiot/yutil/internal/testing"
)

func TestGolangReplace(t *testing.T) {
	for name, i := range map[string]struct {
		content      string
		replacements map[string]interface{}
		expected     string
		expectedErr  string
	}{
		"single replace": {
			content:      "Hello, {{ .name }}!",
			replacements: map[string]interface{}{"name": "World"},
			expected:     "Hello, World!",
		},
		"replace with underscore": {
			content:      "Hello, {{ .first_name }}!",
			replacements: map[string]interface{}{"first_name": "Adrian"},
			expected:     "Hello, Adrian!",
		},
		"map": {
			content: "{{ .root_node.child_node.leaf_node }}",
			replacements: map[string]interface{}{"root_node": map[string]interface{}{
				"child_node": map[string]interface{}{
					"leaf_node": "value",
				},
			}},
			expected: "value",
		},
		"variable": {
			content:      `{{ $data := "value" }}{{ $data }}`,
			replacements: map[string]interface{}{},
			expected:     "value",
		},
		"variable array": {
			content:      `{{ $data := list "value1" "value2" }}[{{ $data | join ", " }}]`,
			replacements: map[string]interface{}{},
			expected:     "[value1, value2]",
		},
		"array union": {
			content: `[{{ concat .one_and_two .three_and_four .five_and_six | join ", " }}]`,
			replacements: map[string]interface{}{
				"one_and_two":    []int{1, 2},
				"three_and_four": []int{3, 4},
				"five_and_six":   []int{5, 6},
			},
			expected: "[1, 2, 3, 4, 5, 6]",
		},
		"array iteration": {
			content: "{{ range $key, $value := .root_node }}{{ $key }}={{ $value }}{{ end }}",
			replacements: map[string]interface{}{
				"root_node": map[string]interface{}{
					"one": 1,
				},
			},
			expected: "one=1",
		},
	} {
		t.Run(name, func(t *testing.T) {
			result, err := golang.Replace(i.content, i.replacements)
			itesting.AssertError(t, i.expectedErr, err)
			itesting.AssertEqual(t, i.expected, result)
		})
	}
}
