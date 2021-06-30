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

	"github.com/amplia-iiot/yutil/internal/yaml"
)

func MergeContents(base string, changes string) (string, error) {
	baseData, err := yaml.Parse(base)
	if err != nil {
		return "", err
	}
	changesData, err := yaml.Parse(changes)
	if err != nil {
		return "", err
	}
	mergedData, err := yaml.Merge(baseData, changesData)
	if err != nil {
		return "", err
	}
	return yaml.Compose(mergedData)
}

func MergeAllContents(contents []string) (string, error) {
	if len(contents) < 2 {
		return "", errors.New("slice must contain at least two contents")
	}
	var result string
	var err error
	result = contents[0]
	for i := 1; i < len(contents); i++ {
		result, err = MergeContents(result, contents[i])
		if err != nil {
			return "", err
		}
	}
	return result, nil
}
