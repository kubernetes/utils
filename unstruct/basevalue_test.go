package unstruct_test

import (
	"reflect"
	"testing"

	"k8s.io/utils/unstruct"
)

func TestSet(t *testing.T) {
	u := unstruct.New(nil)
	u.Set(5)

	if !reflect.DeepEqual(u.Data(), 5) {
		t.Fatalf("Expected 5, got %v", u.Data())
	}
}

func TestParent(t *testing.T) {
	data := map[string]interface{}{
		"field": map[string]interface{}{
			"slice": []interface{}{5},
		},
	}
	u := unstruct.New(data)

	if u.Parent() != nil {
		t.Fatal("Root shouldn't have a parent.")
	}
	if u.Map().Parent() != nil {
		t.Fatal("Root shouldn't have a parent.")
	}

	field := u.Map().Field("field")
	if !reflect.DeepEqual(u.Data(), field.Parent().Data()) {
		t.Fatalf("Expected %v, got %v", u.Data(), field.Parent().Data())
	}

	m := field.Map()
	if !reflect.DeepEqual(u.Data(), m.Parent().Data()) {
		t.Fatalf("Expected %v, got %v", u.Data(), m.Parent().Data())
	}

	slice := m.Field("slice")
	if !reflect.DeepEqual(m.Data(), slice.Parent().Data()) {
		t.Fatalf("Expected %v, got %v", m.Data(), field.Parent().Data())
	}

	s := slice.Slice()
	if !reflect.DeepEqual(m.Data(), slice.Parent().Data()) {
		t.Fatalf("Expected %v, got %v", m.Data(), field.Parent().Data())
	}

	five := s.At(0)
	if !reflect.DeepEqual(slice.Data(), five.Parent().Data()) {
		t.Fatalf("Expected %v, got %v", slice.Data(), five.Parent().Data())
	}
}
