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
	"sort"
)

type Slice []string

func NewSlice(strs ...string) Slice {
	var slice Slice
	for _, str := range strs {
		slice = append(slice, str)
	}
	return slice
}

// Copy copies the content of a Slice to a new Slice
func (s Slice) Copy() Slice {
	if s == nil {
		return nil
	}
	c := make(Slice, len(s))
	copy(c, s)
	return c
}

// Equiv verifies if a slice is equivalent to other.
// It does not sort the slice.
func (rslice Slice) Equiv(lslice Slice) bool {
	if len(rslice) != len(lslice) {
		return false
	}
	for k, v := range lslice {
		if v != rslice[k] {
			return false
		}
	}
	return true
}

// Sort sorts an slice into a new slice
func (s Slice) Sort() Slice {
	c := s.Copy()
	sort.Strings(c)
	return c
}
