package unstruct

type slice struct {
	value Value
}

var _ Slice = &slice{}

func newSlice(value Value) Slice {
	if value.Data() == nil {
		value.Set([]interface{}{})
	}
	_, ok := value.Data().([]interface{})
	if !ok {
		return nil
	}
	return &slice{value: value}
}

func (s *slice) Data() []interface{} {
	d, ok := s.value.Data().([]interface{})
	if !ok {
		return nil
	}
	return d
}

func (s *slice) Parent() Value {
	return s.value.Parent()
}

func (s *slice) Length() int {
	return len(s.Data())
}

func (s *slice) At(index int) Value {
	if index < 0 || index >= s.Length() {
		return nil
	}
	return &sliceValue{parent: s.value, index: index}
}

func (s *slice) Append(value interface{}) Value {
	s.value.Set(append(s.Data(), value))
	return s.At(s.Length() - 1)
}

func (s *slice) InsertAt(index int, value interface{}) Slice {
	if index < 0 || index >= s.Length() {
		return nil
	}
	a := s.Data()
	a = append(a, nil)
	copy(a[index+1:], a[index:])
	a[index] = value
	s.value.Set(a)
	return s

}
