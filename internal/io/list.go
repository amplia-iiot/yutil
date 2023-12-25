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

package io

import (
	"io/fs"
	"path/filepath"

	"github.com/gobwas/glob"
)

var ListFiles = func(dir string, include []string, exclude []string) (files []string, err error) {
	includeGlobs := toGlobs(include)
	excludeGlobs := toGlobs(exclude)
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && matchFile(path, includeGlobs, excludeGlobs) {
			files = append(files, path)
		}
		return nil
	})
	return
}

func toGlobs(strings []string) (globs []glob.Glob) {
	for _, s := range strings {
		globs = append(globs, glob.MustCompile(s))
	}
	return
}

func matchFile(file string, include []glob.Glob, exclude []glob.Glob) bool {
	fileName := filepath.Base(file)
	for _, g := range exclude {
		if g.Match(fileName) || g.Match(file) {
			return false
		}
	}
	for _, g := range include {
		if g.Match(fileName) || g.Match(file) {
			return true
		}
	}
	return len(include) == 0
}
