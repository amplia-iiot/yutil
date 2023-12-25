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
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/amplia-iiot/yutil/internal/io"
	"github.com/amplia-iiot/yutil/internal/replace"
	"github.com/amplia-iiot/yutil/internal/yaml"
	"github.com/amplia-iiot/yutil/pkg/merge"
)

type option func(o *options)

type Engine int

const (
	Golang Engine = iota
	Jinja2
)

type options struct {
	replace.Options
	rootNode                         string
	replacementFiles                 []string
	includeStdinInReplacements       bool
	includeEnvironmentInReplacements bool
	extensions                       []string
}

func (o *options) changeReplacementsRootNode() error {
	if node, ok := o.Options.Replacements[o.rootNode]; ok {
		if node, ok := node.(map[any]any); ok {
			reps := map[string]interface{}{}
			for k, v := range node {
				if name, ok := k.(string); ok {
					reps[name] = v
				}
			}
			o.Options.Replacements = reps
		} else {
			return fmt.Errorf("node %s does not contain more elements", o.rootNode)
		}
	} else {
		return fmt.Errorf("no %s node", o.rootNode)
	}
	return nil
}

// WithDirectory configures the root directory to search for files to be replaced (defaults to current directory).
func WithDirectory(directory string) option {
	return func(o *options) {
		o.Directory = directory
	}
}

// WithInclude configures the glob pattern for files to be included.
func WithInclude(pattern ...string) option {
	return func(o *options) {
		o.Include = append(o.Include, pattern...)
	}
}

// WithExclude configures the glob pattern for files to be excluded.
func WithExclude(pattern ...string) option {
	return func(o *options) {
		o.Exclude = append(o.Exclude, pattern...)
	}
}

// WithRootNode configures the root node to include only replacements from inside that node.
func WithRootNode(node string) option {
	return func(o *options) {
		o.rootNode = node
	}
}

// WithReplacementFile adds a file to be used as replacement file.
func WithReplacementFile(file string) option {
	return func(o *options) {
		o.replacementFiles = append(o.replacementFiles, file)
	}
}

// WithReplacementFiles adds multiple files to be used as replacement files (they will be merged).
func WithReplacementFiles(files ...string) option {
	return func(o *options) {
		o.replacementFiles = append(o.replacementFiles, files...)
	}
}

// IncludeStdinInReplacements includes stdin to be used as replacement file.
func IncludeStdinInReplacements() option {
	return WithIncludeStdinInReplacements(true)
}

// WithIncludeStdinInReplacements configures whether to use stdin as replacement file.
func WithIncludeStdinInReplacements(include bool) option {
	return func(o *options) {
		o.includeStdinInReplacements = include
	}
}

// IncludeEnvironmentInReplacements includes all environment variables to be used in the template engine inside the env node.
func IncludeEnvironmentInReplacements(include bool) option {
	return WithIncludeEnvironmentInReplacements(true)
}

// WithIncludeEnvironmentInReplacements configures whether to use include all environment variables to be used in the template engine inside the env node.
func WithIncludeEnvironmentInReplacements(include bool) option {
	return func(o *options) {
		o.includeEnvironmentInReplacements = include
	}
}

// WithExtension configures the extensions to be removed from a file name when it's saved after being passed through the template engine.
func WithExtension(extension ...string) option {
	return func(o *options) {
		o.extensions = append(o.extensions, extension...)
	}
}

// Replace uses the template engine to replace files following the optional configuration.
func Replace(engine Engine, opts ...option) (err error) {
	o := &options{
		Options: replace.Options{
			Directory: ".",
		},
	}
	for _, opt := range opts {
		opt(o)
	}
	var replacements string
	switch engine {
	case Golang:
		o.Engine = replace.Golang
	case Jinja2:
		o.Engine = replace.Jinja2
	}
	if o.includeStdinInReplacements {
		if len(o.replacementFiles) > 0 {
			replacements, err = merge.MergeStdinWithFiles(o.replacementFiles)
		} else {
			replacements, err = io.ReadStdin()
		}
	} else {
		switch len(o.replacementFiles) {
		case 0:
			err = fmt.Errorf("no replacement files defined")
		case 1:
			replacements, err = io.ReadAsString(o.replacementFiles[0])
		default:
			replacements, err = merge.MergeAllFiles(o.replacementFiles)
		}
	}
	if err != nil {
		return
	}
	o.Options.Replacements, err = yaml.Parse(replacements)
	if err != nil {
		return
	}
	if o.rootNode != "" {
		err = o.changeReplacementsRootNode()
		if err != nil {
			return
		}
	}
	if o.includeEnvironmentInReplacements {
		env := map[string]interface{}{}
		for _, envVar := range os.Environ() {
			if name, value, ok := strings.Cut(envVar, "="); ok {
				env[name] = value
			}
		}
		o.Options.Replacements, err = yaml.Merge(o.Options.Replacements, map[string]interface{}{"env": env})
		if err != nil {
			return
		}
	}
	if len(o.extensions) > 0 {
		sort.SliceStable(o.extensions, func(i, j int) bool {
			return len(o.extensions[i]) > len(o.extensions[j])
		})
		o.FileNameRenamer = func(s string) (res string) {
			res = s
			for _, extension := range o.extensions {
				res = strings.ReplaceAll(res, extension, "")
			}
			return
		}
	}
	return replace.Replace(o.Options)
}
