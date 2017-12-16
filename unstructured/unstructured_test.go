package unstructured

import "testing"

func TestField(t *testing.T) {
	u, err := Unstructured{Data: map[interface{}]interface{}{"field": 5}}.Field("field")
	if err != nil {
		t.Fatal(err)
	}
	if u.Data.(int) != 5 {
		t.Errorf(`Expected "u" to be 5, got %d`, u.Data.(int))
	}
}

func TestFieldMissing(t *testing.T) {
	_, err := Unstructured{Data: map[interface{}]interface{}{"field": 5}}.Field("not-there")
	if err != missingField {
		t.Fatalf(`Expected error to be %q, got: %q`, missingField, err)
	}
}

func TestFieldInvalidType(t *testing.T) {
	_, err := Unstructured{Data: []interface{}{"field"}}.Field("field")
	if err != invalidType {
		t.Fatalf(`Expected error to be %q, got: %q`, missingField, err)
	}
}

func TestSetField(t *testing.T) {
	u := Unstructured{Data: map[interface{}]interface{}{"field": 5}}
	err := u.SetField("field", 1)
	if err != nil {
		t.Fatal(err)
	}
	d, err := u.Field("field")
	if err != nil {
		t.Fatal(err)
	}
	if d.Data.(int) != 1 {
		t.Errorf(`Expected "u" to be 5, got %d`, d.Data.(int))
	}
}

func TestSetFieldInvalid(t *testing.T) {
	u := Unstructured{Data: []interface{}{"field"}}
	err := u.SetField("field", 1)
	if err != invalidType {
		t.Fatalf(`Expected error to be %q, got: %q`, missingField, err)
	}
}

func TestAt(t *testing.T) {
	u, err := Unstructured{Data: []interface{}{"field", 5}}.At(1)
	if err != nil {
		t.Fatal(err)
	}
	if u.Data.(int) != 5 {
		t.Errorf(`Expected "u" to be 5, got %d`, u.Data.(int))
	}
}

func TestAtOutOfBOund(t *testing.T) {
	_, err := Unstructured{Data: []interface{}{"field", 5}}.At(-1)
	if err != invalidIndex {
		t.Fatalf(`Expected error to be %q, got: %q`, invalidIndex, err)
	}

	_, err = Unstructured{Data: []interface{}{"field", 5}}.At(5)
	if err != invalidIndex {
		t.Fatalf(`Expected error to be %q, got: %q`, invalidIndex, err)
	}
}

func TestAtInvalidType(t *testing.T) {
	_, err := Unstructured{Data: map[interface{}]interface{}{"field": 5}}.At(0)
	if err != invalidType {
		t.Fatalf(`Expected error to be %q, got: %q`, invalidType, err)
	}
}

func TestSetAt(t *testing.T) {
	u := Unstructured{Data: []interface{}{"field", 5}}
	err := u.SetAt(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	d, err := u.At(1)
	if err != nil {
		t.Fatal(err)
	}
	if d.Data.(int) != 1 {
		t.Errorf(`Expected "u" to be 5, got %d`, d.Data.(int))
	}
}

func TestSetAtInvalid(t *testing.T) {
	u := Unstructured{Data: map[interface{}]interface{}{"field": 5}}
	err := u.SetAt(0, 1)
	if err != invalidType {
		t.Fatalf(`Expected error to be %q, got: %q`, invalidIndex, err)
	}
}
