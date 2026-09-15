/*
Copyright 2026 The Kubernetes Authors.

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

package lru

import "testing"

func TestCacheLen(t *testing.T) {
	c := New(10)
	if got := c.Len(); got != 0 {
		t.Fatalf("Len() on empty cache = %d, want 0", got)
	}
	c.Add("a", 1)
	c.Add("b", 2)
	if got := c.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	c.Add("a", 3)
	if got := c.Len(); got != 2 {
		t.Fatalf("Len() after overwrite = %d, want 2", got)
	}
}

func TestCacheRemoveOldest(t *testing.T) {
	c := New(10)
	c.Add("a", 1)
	c.Add("b", 2)
	c.Add("c", 3)

	c.RemoveOldest()

	if _, ok := c.Get("a"); ok {
		t.Fatalf("expected oldest key %q to be evicted", "a")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatalf("expected key %q to remain", "b")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatalf("expected key %q to remain", "c")
	}
	if got := c.Len(); got != 2 {
		t.Fatalf("Len() after RemoveOldest = %d, want 2", got)
	}
}

func TestCacheRemoveOldestEmpty(t *testing.T) {
	c := New(10)
	c.RemoveOldest()
	if got := c.Len(); got != 0 {
		t.Fatalf("Len() after RemoveOldest on empty cache = %d, want 0", got)
	}
}

func TestCacheClear(t *testing.T) {
	evicted := map[Key]bool{}
	c := NewWithEvictionFunc(10, func(key Key, value any) {
		evicted[key] = true
	})
	c.Add("a", 1)
	c.Add("b", 2)

	c.Clear()

	if got := c.Len(); got != 0 {
		t.Fatalf("Len() after Clear = %d, want 0", got)
	}
	if _, ok := c.Get("a"); ok {
		t.Fatalf("expected key %q to be gone after Clear", "a")
	}
	if !evicted["a"] || !evicted["b"] {
		t.Fatalf("expected Clear to invoke the eviction func for both keys, got %v", evicted)
	}
}

func TestCacheClearEmpty(t *testing.T) {
	c := New(10)
	c.Clear()
	if got := c.Len(); got != 0 {
		t.Fatalf("Len() after Clear on empty cache = %d, want 0", got)
	}
}
