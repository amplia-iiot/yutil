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
	"errors"
	"fmt"
	"path/filepath"

	iio "github.com/amplia-iiot/yutil/internal/io"
)

type EngineType int

const (
	Golang EngineType = iota
	Jinja2
)

func (e EngineType) String() string {
	return [...]string{"Golang", "Jinja2"}[e]
}

type FileNameRenamer func(string) string

type Options struct {
	Engine          EngineType
	Directory       string
	Include         []string
	Exclude         []string
	Replacements    map[string]interface{}
	FileNameRenamer FileNameRenamer
}

func (o *Options) sanitize() {
	if o.Directory == "" {
		o.Directory = "."
	}
	if o.Include == nil {
		o.Include = []string{}
	}
	if o.Exclude == nil {
		o.Exclude = []string{}
	}
	if o.Replacements == nil {
		o.Replacements = map[string]interface{}{}
	} else {
		o.Replacements = sanitizeReplacements(&o.Replacements)
	}
	if o.FileNameRenamer == nil {
		o.FileNameRenamer = func(s string) string {
			return s
		}
	}
}

func sanitizeReplacements(node *map[string]any) map[string]interface{} {
	reps := map[string]interface{}{}
	for name, v := range *node {
		if n, ok := v.(map[any]any); ok {
			reps[name] = sanitizeNode(n)
		} else {
			reps[name] = v
		}
	}
	return reps
}

func sanitizeNode(node any) any {
	if n, ok := (node).(map[any]any); ok {
		m := map[string]interface{}{}
		for k, v := range n {
			if name, ok := k.(string); ok {
				m[name] = sanitizeNode(v)
			}
		}
		return m
	} else if n, ok := (node).([]any); ok {
		s := []interface{}{}
		for _, v := range n {
			s = append(s, sanitizeNode(v))
		}
		return s
	}
	return node
}

func (o *Options) engine() (Engine, error) {
	switch o.Engine {
	case Golang:
		return golang, nil
	case Jinja2:
		return jinja2, nil
	}
	return nil, fmt.Errorf("unsupported engine: %s", o.Engine)
}

func Replace(opts Options) (err error) {
	opts.sanitize()
	engine, err := opts.engine()
	if err != nil {
		return
	}
	files, err := iio.ListFiles(opts.Directory, opts.Include, opts.Exclude)
	if err != nil {
		return
	}
	errs := []error{}
	for _, file := range files {
		content, err := iio.ReadAsString(file)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		replaced, err := engine.Replace(content, opts.Replacements)
		if err != nil {
			errs = append(errs, fmt.Errorf("error on file %s: %w", file, err))
			continue
		}
		path := filepath.Dir(file)
		name := filepath.Base(file)
		renamed := opts.FileNameRenamer(name)
		err = iio.WriteToFile(filepath.Join(path, renamed), replaced)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return errors.Join(errs...)
}
