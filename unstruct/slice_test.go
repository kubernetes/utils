package unstruct_test

import (
	"reflect"
	"testing"

	"k8s.io/utils/unstruct"
)

func TestConvertRootSlice(t *testing.T) {
	u := unstruct.New([]interface{}{1, 2, 3})

	if u.Slice() == nil {
		t.Fatal("Couldn't convert slice to slice")
	}
	if !reflect.DeepEqual(u.Slice().Data(), u.Data()) {
		t.Fatalf("Slice data %v different from Value data %v",
			u.Slice().Data(),
			u.Data())
	}
}

func TestAppendRootSlice(t *testing.T) {
	u := unstruct.New([]interface{}{1, 2, 3})

	u.Slice().Append(4)
	if !reflect.DeepEqual(u.Data(), []interface{}{1, 2, 3, 4}) {
		t.Fatalf("Failure to Append(4), expected %v, got %v",
			[]interface{}{1, 2, 3, 4}, u.Data())
	}
}

func TestGetAndSetElementRootSlice(t *testing.T) {
	u := unstruct.New([]interface{}{1, 2, 3})
	v := u.Slice().At(2)
	if !reflect.DeepEqual(v.Data(), 3) {
		t.Fatalf("Expected data to be 3, got %v", v.Data())
	}
	v.Set("3")
	if !reflect.DeepEqual(u.Data(), []interface{}{1, 2, "3"}) {
		t.Fatalf("Expected data to be %v, got %v",
			[]interface{}{1, 2, "3"},
			u.Data())
	}
}

func TestSliceInSlice(t *testing.T) {
	u := unstruct.New([]interface{}{1, 2, []interface{}{3, 4}})

	u.Slice().Append(5)
	expected := []interface{}{1, 2, []interface{}{3, 4}, 5}
	if !reflect.DeepEqual(expected, u.Data()) {
		t.Fatalf("Expected %v, got %v", expected, u.Data())
	}

	u.Slice().At(2).Slice().Append(4.5)
	expected = []interface{}{1, 2, []interface{}{3, 4, 4.5}, 5}
	if !reflect.DeepEqual(expected, u.Data()) {
		t.Fatalf("Expected %v, got %v", expected, u.Data())
	}

	u.Slice().At(2).Slice().At(1).Set(3.9)
	expected = []interface{}{1, 2, []interface{}{3, 3.9, 4.5}, 5}
	if !reflect.DeepEqual(expected, u.Data()) {
		t.Fatalf("Expected %v, got %v", expected, u.Data())
	}

	u.Slice().At(0).Set("1")
	expected = []interface{}{"1", 2, []interface{}{3, 3.9, 4.5}, 5}
	if !reflect.DeepEqual(expected, u.Data()) {
		t.Fatalf("Expected %v, got %v", expected, u.Data())
	}
}

func TestSliceAsMap(t *testing.T) {
	u := unstruct.New([]interface{}{1, 2, 3})
	if u.Map() != nil {
		t.Fatal("Slice shouldn't be converted to Slice.")
	}
}

func TestNilSliceConversion(t *testing.T) {
	u := unstruct.New(nil)

	u.Slice().Append(1)
	expected := []interface{}{1}
	if !reflect.DeepEqual(expected, u.Data()) {
		t.Fatalf("Expected %v, got %v", expected, u.Data())
	}
}
