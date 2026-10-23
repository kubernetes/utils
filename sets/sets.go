/*
Copyright The Kubernetes Authors.

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

package sets

import (
	"cmp"

	"k8s.io/utils/set"
)

// Set is an alias for set.Set so that code using apimachinery sets.Set can be
// migrated without having to change all the corresponding code.
type Set[E comparable] = set.Set[E]

// Empty is public since it is used by some internal API objects for conversions between external
// string arrays and internal sets, and conversion logic requires public types today.
type Empty = set.Empty

// New creates a new set.
// Deprecated: use set.New instead.
//
//go:fix inline
func New[E comparable](items ...E) Set[E] {
	return set.New[E](items...) // The parameter type is required for inlining when items is empty
}

// KeySet creates a Set[E] from a keys of a map[E](? extends any).
// Deprecated: use set.KeySet instead.
//
//go:fix inline
func KeySet[E comparable, A any](theMap map[E]A) Set[E] {
	return set.KeySet(theMap)
}

// List returns the contents as a sorted E slice.
// Deprecated: use set.List instead.
//
//go:fix inline
func List[E cmp.Ordered](s Set[E]) []E {
	return set.List(s)
}
