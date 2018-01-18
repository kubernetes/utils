package unstructured

import (
	"errors"
	"testing"
)

func TestBroken(t *testing.T) {
	_, err := broken{Error: errors.New("broken...")}.
		Field("SomeField").
		SetField("SomeOtherField", 1).
		At(1).
		SetAt(0, 12).Data()
	if err.Error() != "broken..." {
		t.Fatalf(`Expected error to be "broken...", got: %v`, err)
	}
}
