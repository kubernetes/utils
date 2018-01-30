package unstruct

type sliceValue struct {
	parent Value
	index  int
}

var _ Value = &sliceValue{}

func (s *sliceValue) Data() interface{} {
	return s.parent.Slice().Data()[s.index]
}

func (s *sliceValue) Parent() Value {
	return s.parent
}

func (s *sliceValue) Set(value interface{}) Value {
	s.parent.Slice().Data()[s.index] = value
	return s
}

func (s *sliceValue) Map() Map {
	return newMap(s)
}

func (s *sliceValue) Slice() Slice {
	return newSlice(s)
}
