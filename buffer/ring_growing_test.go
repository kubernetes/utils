/*
Copyright 2017 The Kubernetes Authors.

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

package buffer

import (
	"reflect"
	"testing"
)

func TestGrowth(t *testing.T) {
	t.Parallel()
	x := 10
	g := NewRingGrowing(1)
	for i := 0; i < x; i++ {
		if e, a := i, g.readable; !reflect.DeepEqual(e, a) {
			t.Fatalf("expected equal, got %#v, %#v", e, a)
		}
		g.WriteOne(i)
	}
	read := 0
	for g.readable > 0 {
		v, ok := g.ReadOne()
		if !ok {
			t.Fatal("expected true")
		}
		if read != v {
			t.Fatalf("expected %#v==%#v", read, v)
		}
		read++
	}
	if x != read {
		t.Fatalf("expected to have read %d items: %d", x, read)
	}
	if g.readable != 0 {
		t.Fatalf("expected readable to be zero: %d", g.readable)
	}
	if 16 != g.n {
		t.Fatalf("expected N to be 16: %d", g.n)
	}
}

func TestEmpty(t *testing.T) {
	t.Parallel()
	g := NewRingGrowing(1)
	_, ok := g.ReadOne()
	if ok != false {
		t.Fatal("expected false")
	}
}
