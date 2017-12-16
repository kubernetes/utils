package unstructured

// Chainable is an interface to a more-easily chain Unstructured types.
type Chainable interface {
	// Data returns the data pointed by the current element, or the
	// current saved error. The error might have been carried from a
	// previous failure to reach or change an element.
	//
	// Note that the first error will preempt over all other
	// operations. In other words, if you chain multiple set
	// operations, only the ones before the first failure will happen,
	// others will be ignored.
	Data() (interface{}, error)

	// Field returns a Chainable to the value of the "key" field.
	// If the current object is not a map, this will return a
	// Chainable containing an error.
	Field(key string) Chainable
	// SetField changes the value of the "key" field in the map.
	// This operation doesn't change the Chainable, so that you can
	// set multiple fields consecutively. If the current object is
	// not a map, this will return a Chainable containing an error.
	SetField(key string, value interface{}) Chainable
	// At returns a Chainable to the value at the index-th item. If
	// the current object is not an array, this will return a
	// Chainable containing an error.
	At(index int) Chainable
	// SetAt changes the value of the index-th item in the array.
	// This operation doesn't change the Chainable, so that you can
	// set multiple field consecutively. If the current object is
	// not an array, this will return a Chainable container an
	// error.
	SetAt(index int, value interface{}) Chainable
}
