package unstructured

import (
	"errors"
)

var missingField = errors.New("field not found")
var invalidIndex = errors.New("index is out of bound")
var invalidType = errors.New("invalid type")

// Unstructured is the object that let's you manipulate Data.
type Unstructured struct {
	Data interface{}
}

// Field returns the Unstructured object in the "key" field, if this
// object is a map.
//
// If the object is not a map, or if the field doesn't exist, an error
// is returned.
func (u Unstructured) Field(key string) (Unstructured, error) {
	m, ok := u.Data.(map[interface{}]interface{})
	if !ok {
		return Unstructured{}, invalidType
	}
	d, ok := m[key]
	if !ok {
		return Unstructured{}, missingField
	}
	return Unstructured{Data: d}, nil
}

// SetField sets the value in the "key" field of the Unstructured
// object, if this object is a map.
//
// If the object is not a map, an error is returned.
func (u Unstructured) SetField(key string, value interface{}) error {
	m, ok := u.Data.(map[interface{}]interface{})
	if !ok {
		return invalidType
	}
	m[key] = value
	return nil
}

// At returns the "index"-th item in the list, if this object is an
// array.
//
// If the object is not an array, or if the index is out-of-bound, an
// error is returned.
func (u Unstructured) At(index int) (Unstructured, error) {
	a, ok := u.Data.([]interface{})
	if !ok {
		return Unstructured{}, invalidType
	}
	if index < 0 || index >= len(a) {
		return Unstructured{}, invalidIndex
	}
	return Unstructured{Data: a[index]}, nil
}

// SetAt sets the value for the index-th field of the Unstructured
// object, if this object is an array.
//
// If the object is not a map, an error is returned.
func (u Unstructured) SetAt(index int, value interface{}) error {
	a, ok := u.Data.([]interface{})
	if !ok {
		return invalidType
	}
	if index < 0 || index >= len(a) {
		return invalidIndex
	}
	a[index] = value
	return nil
}

// InsertAt inserts the value at the index-th position in the
// Unstructured object, if this object is an array. The new array is
// returned, and needs to be re-inserted where needed.
//
// If the object is not an array, or if the index is out-of-bound, an
// error is returned.
func (u *Unstructured) InsertAt(index int, value interface{}) ([]interface{}, error) {
	a, ok := u.Data.([]interface{})
	if !ok {
		return nil, invalidType
	}

	a = append(a, 0)
	copy(a[index+1:], a[index:])
	a[index] = value
	return a, nil
}
