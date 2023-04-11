/*
Copyright 2014 The Kubernetes Authors.

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
	"reflect"
	"testing"
)

func TestStringSet(t *testing.T) {
	s := NewSet[string]()
	s2 := NewSet[string]()
	if len(s) != 0 {
		t.Errorf("Expected len=0: %d", len(s))
	}
	s.Insert("a", "b")
	if len(s) != 2 {
		t.Errorf("Expected len=2: %d", len(s))
	}
	s.Insert("c")
	if s.Has("d") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Has("a") {
		t.Errorf("Missing contents: %#v", s)
	}
	s.Delete("a")
	if s.Has("a") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s.Insert("a")
	if s.HasAll("a", "b", "d") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.HasAll("a", "b") {
		t.Errorf("Missing contents: %#v", s)
	}
	s2.Insert("a", "b", "d")
	if s.IsSuperset(s2) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s2.Delete("d")
	if !s.IsSuperset(s2) {
		t.Errorf("Missing contents: %#v", s)
	}
}

func TestStringSetDeleteMultiples(t *testing.T) {
	s := NewSet[string]()
	s.Insert("a", "b", "c")
	if len(s) != 3 {
		t.Errorf("Expected len=3: %d", len(s))
	}

	s.Delete("a", "c")
	if len(s) != 1 {
		t.Errorf("Expected len=1: %d", len(s))
	}
	if s.Has("a") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if s.Has("c") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Has("b") {
		t.Errorf("Missing contents: %#v", s)
	}
}

func TestNewStringSet(t *testing.T) {
	s := NewSet[string]("a", "b", "c")
	if len(s) != 3 {
		t.Errorf("Expected len=3: %d", len(s))
	}
	if !s.Has("a") || !s.Has("b") || !s.Has("c") {
		t.Errorf("Unexpected contents: %#v", s)
	}
}

func TestStringSetList(t *testing.T) {
	s := NewSet[string]("z", "y", "x", "a")
	if !reflect.DeepEqual(s.List(), []string{"a", "x", "y", "z"}) {
		t.Errorf("List gave unexpected result: %#v", s.List())
	}
}

func TestStringSetDifference(t *testing.T) {
	a := NewSet[string]("1", "2", "3")
	b := NewSet[string]("1", "2", "4", "5")
	c := a.Difference(b)
	d := b.Difference(a)
	if len(c) != 1 {
		t.Errorf("Expected len=1: %d", len(c))
	}
	if !c.Has("3") {
		t.Errorf("Unexpected contents: %#v", c.List())
	}
	if len(d) != 2 {
		t.Errorf("Expected len=2: %d", len(d))
	}
	if !d.Has("4") || !d.Has("5") {
		t.Errorf("Unexpected contents: %#v", d.List())
	}
}

func TestStringSetHasAny(t *testing.T) {
	a := NewSet[string]("1", "2", "3")

	if !a.HasAny("1", "4") {
		t.Errorf("expected true, got false")
	}

	if a.HasAny("0", "4") {
		t.Errorf("expected false, got true")
	}
}

func TestStringSetEquals(t *testing.T) {
	// Simple case (order doesn't matter)
	a := NewSet[string]("1", "2")
	b := NewSet[string]("2", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// It is a set; duplicates are ignored
	b = NewSet[string]("2", "2", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// Edge cases around empty sets / empty strings
	a = NewSet[string]()
	b = NewSet[string]()
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	b = NewSet[string]("1", "2", "3")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	b = NewSet[string]("1", "2", "")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	// Check for equality after mutation
	a = NewSet[string]()
	a.Insert("1")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert("2")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert("")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	a.Delete("")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}
}

func TestStringUnion(t *testing.T) {
	tests := []struct {
		s1       Set[string]
		s2       Set[string]
		expected Set[string]
	}{
		{
			NewSet[string]("1", "2", "3", "4"),
			NewSet[string]("3", "4", "5", "6"),
			NewSet[string]("1", "2", "3", "4", "5", "6"),
		},
		{
			NewSet[string]("1", "2", "3", "4"),
			NewSet[string](),
			NewSet[string]("1", "2", "3", "4"),
		},
		{
			NewSet[string](),
			NewSet[string]("1", "2", "3", "4"),
			NewSet[string]("1", "2", "3", "4"),
		},
		{
			NewSet[string](),
			NewSet[string](),
			NewSet[string](),
		},
	}

	for _, test := range tests {
		union := test.s1.Union(test.s2)
		if union.Len() != test.expected.Len() {
			t.Errorf("Expected union.Len()=%d but got %d", test.expected.Len(), union.Len())
		}

		if !union.Equal(test.expected) {
			t.Errorf("Expected union.Equal(expected) but not true.  union:%v expected:%v", union.List(), test.expected.List())
		}
	}
}

func TestStringIntersection(t *testing.T) {
	tests := []struct {
		s1       Set[string]
		s2       Set[string]
		expected Set[string]
	}{
		{
			NewSet[string]("1", "2", "3", "4"),
			NewSet[string]("3", "4", "5", "6"),
			NewSet[string]("3", "4"),
		},
		{
			NewSet[string]("1", "2", "3", "4"),
			NewSet[string]("1", "2", "3", "4"),
			NewSet[string]("1", "2", "3", "4"),
		},
		{
			NewSet[string]("1", "2", "3", "4"),
			NewSet[string](),
			NewSet[string](),
		},
		{
			NewSet[string](),
			NewSet[string]("1", "2", "3", "4"),
			NewSet[string](),
		},
		{
			NewSet[string](),
			NewSet[string](),
			NewSet[string](),
		},
	}

	for _, test := range tests {
		intersection := test.s1.Intersection(test.s2)
		if intersection.Len() != test.expected.Len() {
			t.Errorf("Expected intersection.Len()=%d but got %d", test.expected.Len(), intersection.Len())
		}

		if !intersection.Equal(test.expected) {
			t.Errorf("Expected intersection.Equal(expected) but not true.  intersection:%v expected:%v", intersection.List(), test.expected.List())
		}
	}
}

func TestNewSetFromMapKeys(t *testing.T) {
	m := map[string]string{
		"hallo":   "world",
		"goodbye": "and goodnight",
	}
	expected := []string{"goodbye", "hallo"}
	gotList := NewSetFromMapKeys(m).List() // List() returns a sorted list
	if len(gotList) != len(m) {
		t.Fatalf("got %v elements, wanted %v", len(gotList), len(m))
	}
	for i, entry := range NewSetFromMapKeys(m).List() {
		if entry != expected[i] {
			t.Errorf("got %v, expected %v", entry, expected[i])
		}
	}
}
