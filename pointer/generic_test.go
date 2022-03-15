package pointer

import "testing"

func TestPointer(t *testing.T) {
	ptr := Pointer(5)
	if *ptr != 5 {
		t.Errorf("expected pointer value to be 5, was %v", *ptr)
	}
}
