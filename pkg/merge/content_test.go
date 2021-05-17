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
	"reflect"
	"testing"

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
		assert(t, i.expected, merged)
	}
}

func format(content string) (string, error) {
	data, err := yaml.Parse(content)
	if err != nil {
		return "", err
	}
	return yaml.Compose(data)
}

func assert(t *testing.T, expected string, got string) {
	expected, err := format(expected)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, expected, got)
}

func assertEqual(t *testing.T, expected interface{}, got interface{}) {
	// fmt.Printf("comparing %s with %s", expected, got)
	if expected == got {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", got, reflect.TypeOf(got), expected, reflect.TypeOf(expected))
}
