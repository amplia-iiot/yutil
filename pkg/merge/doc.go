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

// Package merge provides primitives for deep merging yaml files and content.
//
// Map nodes will be merged containing all child nodes, replacing leaf values for
// duplicated less important nodes. The most important node is usually the last
// given yaml (file or content) and takes precedence over any previous yaml,
// check the specific method to be sure.
//
// Primitive value nodes like strings, booleans, integers, etc. are considered
// leaf nodes and will be replaced on less important leaf nodes.
//
// Array nodes are also considered leaf nodes and will be replaced, they will not
// be merged. Merging "data: [1]" with "data: [2]" will result in "data: [2]",
// not "data: [1, 2]".
//
// Leaf nodes can replace map and array nodes entirely if they are more important
// in the merge hierarchy, replacing a complex node with a primitive value if
// they are on the same key path.
//
// The process of merging also formats the merged yaml (see format).
package merge
