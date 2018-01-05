package unstruct_test

import (
	"reflect"
	"testing"

	"k8s.io/utils/unstruct"
)

func TestConvertRootMap(t *testing.T) {
	u := unstruct.New(map[string]interface{}{"one": 1, "two": 2, "three": 3})

	if u.Map() == nil {
		t.Fatal("Couldn't convert slice to map")
	}
	if !reflect.DeepEqual(u.Map().Data(), u.Data()) {
		t.Fatalf("Map data %v different from Value data %v",
			u.Map().Data(),
			u.Data())
	}
}

func TestSetNewField(t *testing.T) {
	u := unstruct.New(map[string]interface{}{"one": 1, "two": 2, "three": 3})

	u.Map().Field("four").Set(4)
	if !reflect.DeepEqual(u.Data(), map[string]interface{}{"one": 1, "two": 2, "three": 3, "four": 4}) {
		t.Fatalf("Failure to Append(4), expected %v, got %v",
			map[string]interface{}{"one": 1, "two": 2, "three": 3, "four": 4}, u.Data())
	}
}

func TestGetAndSetElementRootMap(t *testing.T) {
	u := unstruct.New(map[string]interface{}{"one": 1, "two": 2, "three": 3})
	v := u.Map().Field("two")
	if !reflect.DeepEqual(v.Data(), 2) {
		t.Fatalf("Expected data to be 2, got %v", v.Data())
	}
	v.Set("2")
	if !reflect.DeepEqual(u.Data(), map[string]interface{}{"one": 1, "two": "2", "three": 3}) {
		t.Fatalf("Expected data to be %v, got %v",
			map[string]interface{}{"one": 1, "two": "2", "three": 3},
			u.Data())
	}
}

func TestMapInSlice(t *testing.T) {
	u := unstruct.New([]interface{}{map[string]interface{}{"one": 1}})

	u.Slice().Append(2)
	expected := []interface{}{map[string]interface{}{"one": 1}, 2}
	if !reflect.DeepEqual(expected, u.Data()) {
		t.Fatalf("Expected %v, got %v", expected, u.Data())
	}

	u.Slice().At(0).Map().Field("three").Set(3)
	expected = []interface{}{map[string]interface{}{"one": 1, "three": 3}, 2}
	if !reflect.DeepEqual(expected, u.Data()) {
		t.Fatalf("Expected %v, got %v", expected, u.Data())
	}
}

func TestMapAsSlice(t *testing.T) {
	u := unstruct.New(map[string]interface{}{"one": 1, "two": 2, "three": 3})
	if u.Slice() != nil {
		t.Fatal("Map shouldn't be converted to Slice.")
	}
}

func TestNilMapConversion(t *testing.T) {
	u := unstruct.New(nil)

	u.Map().Field("int").Set(1)
	expected := map[string]interface{}{"int": 1}
	if !reflect.DeepEqual(expected, u.Data()) {
		t.Fatalf("Expected %v, got %v", expected, u.Data())
	}
}
