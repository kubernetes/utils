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

package diff

import (
	"strings"
	"testing"
)

func TestObjectDiff(t *testing.T) {
	out := ObjectDiff("hello", "world")
	if !strings.Contains(out, `A: hello"`) {
		t.Fatalf("expected a's remainder in output, got %q", out)
	}
	if !strings.Contains(out, `B: world"`) {
		t.Fatalf("expected b's remainder in output, got %q", out)
	}
}

func TestObjectDiffEqual(t *testing.T) {
	out := ObjectDiff("same", "same")
	if !strings.Contains(out, `"same"`) {
		t.Fatalf("expected the shared JSON value in output, got %q", out)
	}
	if !strings.Contains(out, "A: \n\nB: ") {
		t.Fatalf("expected empty A/B remainders for identical objects, got %q", out)
	}
}

func TestObjectDiffPanicsOnUnmarshalableType(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected ObjectDiff to panic when json.Marshal fails")
		}
	}()
	ObjectDiff(func() {}, "b")
}

func TestObjectGoPrintDiff(t *testing.T) {
	out := ObjectGoPrintDiff("hello", "world")
	if !strings.Contains(out, "(string)") {
		t.Fatalf("expected go-spew header in output, got %q", out)
	}
	if !strings.Contains(out, "A: hello") {
		t.Fatalf("expected a's remainder in output, got %q", out)
	}
	if !strings.Contains(out, "B: world") {
		t.Fatalf("expected b's remainder in output, got %q", out)
	}
}

func TestObjectGoPrintSideBySide(t *testing.T) {
	type obj struct{ A int }
	out := ObjectGoPrintSideBySide(obj{A: 1}, obj{A: 2})

	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected multiple lines of side by side output, got %q", out)
	}
	if !strings.Contains(out, "A: (int) 1") {
		t.Fatalf("expected a's dump in output, got %q", out)
	}
	if !strings.Contains(out, "A: (int) 2") {
		t.Fatalf("expected b's dump in output, got %q", out)
	}
}

func TestObjectGoPrintSideBySideDifferentLineCounts(t *testing.T) {
	out := ObjectGoPrintSideBySide(1, "longer value here")
	if !strings.Contains(out, "(int) 1") {
		t.Fatalf("expected a's short dump in output, got %q", out)
	}
	if !strings.Contains(out, "longer value here") {
		t.Fatalf("expected b's longer dump in output, got %q", out)
	}
}
