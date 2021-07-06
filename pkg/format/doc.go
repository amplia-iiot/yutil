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

// Package format provides primitives for formatting yaml files and content.
//
// Formatting a yaml includes sorting its nodes alphabetically and cleaning
// the format of the values.
//
// Strings that do not need quotes to remain a primitive string lose the quotes.
// When quotes are needed, single quotes are preferred for strings with special
// characters. For strings containing a number, boolean or null values, double
// quotes are used. Unicode escape sequences in a string are replaced with the
// character.
//
// The proper formatting for null is null, not Null. The same happens to boolean
// values, lowercase is used when formatting.
//
// Arrays maintain the order of elements, and each element appears on a new line.
//
// Comments are removed.
package format
