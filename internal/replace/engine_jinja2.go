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
	j2 "github.com/kluctl/go-jinja2"
)

var jinja2 = jinja2Engine{}

type jinja2Engine struct {
}

func (e jinja2Engine) Replace(content string, replacements map[string]interface{}) (replaced string, err error) {
	tmpl, err := j2.NewJinja2("", 1, j2.WithGlobals(replacements))
	if err != nil {
		return
	}
	defer tmpl.Close()

	replaced, err = tmpl.RenderString(content, j2.WithGlobals(replacements))
	return
}