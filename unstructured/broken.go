package unstructured

// broken is a Chainable that is used when the chain is broken, but
// since we can't return the error, we just go wherever needed while
// carrying that err.
type broken struct {
	Error error
}

var _ Chainable = broken{}

func (c broken) Data() (interface{}, error) {
	return nil, c.Error
}

func (c broken) Field(key string) Chainable {
	return c
}

func (c broken) SetField(key string, value interface{}) Chainable {
	return c
}

func (c broken) At(index int) Chainable {
	return c
}

func (c broken) SetAt(index int, value interface{}) Chainable {
	return c
}
