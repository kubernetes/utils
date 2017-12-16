package unstructured

import "fmt"

type chainer struct {
	data interface{}

	// path is used to build the path as we walk through the
	// chain. This is used to build the error when we have one.
	path string
}

// NewChainable returns a Chainable for the given object.
func NewChainable(value interface{}) Chainable {
	return chainer{data: value}
}

func (c chainer) Data() (interface{}, error) {
	return c.data, nil
}

func (c chainer) Field(key string) Chainable {
	f, err := Unstructured{c.data}.Field(key)
	if err != nil {
		return broken{
			Error: fmt.Errorf("error getting key %q in %s: %v",
				key, c.path, err),
		}
	}
	return chainer{data: f.Data, path: c.path + fmt.Sprintf(".%s", key)}
}

func (c chainer) SetField(key string, value interface{}) Chainable {
	err := Unstructured{c.data}.SetField(key, value)
	if err != nil {
		return broken{
			Error: fmt.Errorf("error setting key %q in %s: %v",
				key, c.path, err),
		}
	}
	return c
}

func (c chainer) At(index int) Chainable {
	f, err := Unstructured{c.data}.At(index)
	if err != nil {
		return broken{
			Error: fmt.Errorf(`error getting item "%d" in %q: %v`,
				index, c.path, err),
		}
	}
	return chainer{data: f.Data, path: c.path + fmt.Sprintf("[%d]", index)}
}

func (c chainer) SetAt(index int, value interface{}) Chainable {
	err := Unstructured{c.data}.SetAt(index, value)
	if err != nil {
		return broken{
			Error: fmt.Errorf(`error setting item "%d" in %q: %v`,
				index, c.path, err),
		}
	}
	return c
}
