/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package strings

import (
	"testing"
)

func TestNewSlice(t *testing.T) {
	testCases := []struct {
		input []string
	}{
		{input: []string{"abcd", "abcef", "lalala", "blabla"}},
		{input: []string{"xxxx", "yyy", "", "zzzz"}},
	}

	for i, tc := range testCases {
		newSlice := NewSlice(tc.input...)
		for k, v := range newSlice {
			if v != tc.input[k] {
				t.Errorf("case[%d]: expected (%q), got (%q)", i, tc.input[k], v)
			}
		}
	}
}

func TestCopySlice(t *testing.T) {
	testCases := []struct {
		input Slice
	}{
		{input: Slice{"abcd", "abcef", "lalala", "blabla"}},
		{input: Slice{"xpto", "", "blargh", "blabla"}},
	}
	for i, tc := range testCases {
		newSlice := tc.input.Copy()
		for k, v := range newSlice {
			if v != tc.input[k] {
				t.Errorf("case[%d]: expected (%q), got (%q)", i, tc.input[k], v)
			}
		}
	}
}

func TestSortSlice(t *testing.T) {
	testCases := []struct {
		input  Slice
		output Slice
	}{
		{
			input:  Slice{"abcd", "abcef", "lalala", "blabla"},
			output: Slice{"abcd", "abcef", "blabla", "lalala"},
		},
		{
			input:  Slice{"xpto", "", "blargh", "blabla"},
			output: Slice{"", "blabla", "blargh", "xpto"},
		},
		{
			input:  Slice{"aaaa", "bbbb", "aaaa", "bbbb"},
			output: Slice{"aaaa", "aaaa", "bbbb", "bbbb"},
		},
	}
	for i, tc := range testCases {
		newSlice := tc.input.Sort()
		for k, v := range newSlice {
			if v != tc.output[k] {
				t.Errorf("case[%d]: expected (%q), got (%q)", i, tc.output[k], v)
			}
		}
	}
}

func TestEquivSlice(t *testing.T) {
	testCases := []struct {
		left     Slice
		right    Slice
		expected bool
	}{
		{
			left:     Slice{"abcd", "abcef", "blabla", "lalala"},
			right:    Slice{"abcd", "abcef", "blabla", "lalala"},
			expected: true,
		},
		{
			left:     Slice{"", "", "blabla", "lalala"},
			right:    Slice{"", "", "blabla", "lalala"},
			expected: true,
		},
		{
			left:     Slice{"a", "b", "c", "d"},
			right:    Slice{"a", "b", "c"},
			expected: false,
		},
		{
			left:     Slice{"xpto", "", "blargh", "blabla"},
			right:    Slice{"", "blabla", "blargh", "xpto"},
			expected: false,
		},
		{
			left:     Slice{"aaaa", "bbbb", "aaaa", "bbbb"},
			right:    Slice{"aaaa", "aaaa", "bbbb", "bbbb"},
			expected: false,
		},
	}
	for i, tc := range testCases {
		if compare := tc.left.Equiv(tc.right); compare != tc.expected {
			t.Errorf("case[%d]: expected (%t), got (%t)", i, tc.expected, compare)
		}
	}
}
