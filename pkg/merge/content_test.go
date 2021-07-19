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
	"testing"

	itesting "github.com/amplia-iiot/yutil/internal/testing"
	"github.com/amplia-iiot/yutil/internal/yaml"
)

func TestMergeContents(t *testing.T) {
	for _, i := range []struct {
		base     string
		changes  string
		expected string
	}{
		{
			base:     `data: {one: 1, two: 2  }`,
			changes:  `data: {        two: two}`,
			expected: `data: {one: 1, two: two}`,
		},
		{
			base:     `data: {one: 1, two: 2          }`,
			changes:  `data: {                three: 3}`,
			expected: `data: {one: 1, two: 2, three: 3}`,
		},
		// Strings are formatted without quotes
		{
			base:     `data: {one: 1, two: 2          }`,
			changes:  `data: {                three: 'three'}`,
			expected: `data: {one: 1, two: 2, three: three}`,
		},
		// Unless containing special chars
		{
			base:     `data: {one: 1, two: 2          }`,
			changes:  `data: {                three: "{"}`,
			expected: `data: {one: 1, two: 2, three: "{"}`,
		},
		// Arrays are replaced, not joined
		{
			base:     `data: {simple: one, extra: extra, array: [{name: one}, {name: two}]}`,
			changes:  `data: {simple: two,               array: [{name: three}]}`,
			expected: `data: {simple: two, extra: extra, array: [{name: three}]}`,
		},
		// Primitive values can replace complex values
		{
			base:     `data: {one: {two: 2}}`,
			changes:  `data: {one: 1}`,
			expected: `data: {one: 1}`,
		},
		{
			base:     `data: {one: [1]}`,
			changes:  `data: {one: 1}`,
			expected: `data: {one: 1}`,
		},
		// Multiple root estructures
		{
			base:     "data: {one: one}                             \ndata3: {value: 1}",
			changes:  "data: {          two: two}\ndata2: {two: two}\ndata3: {value: 2}",
			expected: "data: {one: one, two: two}\ndata2: {two: two}\ndata3: {value: 2}",
		},
	} {
		merged, err := MergeContents(i.base, i.changes)
		if err != nil {
			t.Fatal(err)
		}
		itesting.AssertEqual(t, format(t, i.expected), merged)
	}
}

func TestMergeContentsInvalid(t *testing.T) {
	for _, i := range []struct {
		base     string
		changes  string
		expected string
	}{
		{
			base:     `;`,
			changes:  `data: 2`,
			expected: `cannot unmarshal`,
		},
		{
			base:     `data: 1`,
			changes:  `;`,
			expected: `cannot unmarshal`,
		},
	} {
		merged, err := MergeContents(i.base, i.changes)
		itesting.AssertError(t, i.expected, err)
		if merged != "" {
			t.Fatalf("Should not have merged")
		}
	}
}

func TestMergeContentsMergeError(t *testing.T) {
	// Mock merge internal function
	originalMerge := yaml.Merge
	defer func() { yaml.Merge = originalMerge }()
	yaml.Merge = func(base, changes map[string]interface{}) (map[string]interface{}, error) {
		return nil, errors.New("merging error")
	}
	merged, err := MergeContents("", "")
	itesting.AssertError(t, "merging error", err)
	if merged != "" {
		t.Fatalf("Should not have merged")
	}
}

func TestMergeAllContentsInvalid(t *testing.T) {
	for _, i := range []struct {
		contents []string
		expected string
	}{
		// At least two
		{
			contents: make([]string, 0),
			expected: `slice must contain at least two contents`,
		},
		{
			contents: make([]string, 1),
			expected: `slice must contain at least two contents`,
		},
		// Parsing error
		{
			contents: []string{
				"data: 1",
				"data: 2",
				";",
			},
			expected: `cannot unmarshal`,
		},
		{
			contents: []string{
				";",
				"data: 2",
				"data: 3",
			},
			expected: `cannot unmarshal`,
		},
	} {
		merged, err := MergeAllContents(i.contents)
		itesting.AssertError(t, i.expected, err)
		if merged != "" {
			t.Fatalf("Should not have merged")
		}
	}
}

func format(t *testing.T, content string) string {
	data, err := yaml.Parse(content)
	if err != nil {
		t.Fatalf("Error formatting '%s'. %s", content, err)
	}
	var formatted string
	formatted, err = yaml.Compose(data)
	if err != nil {
		t.Fatalf("Error formatting '%s'. %s", content, err)
	}
	return formatted
}
