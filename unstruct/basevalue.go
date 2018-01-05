package unstruct

type baseValue struct {
	parent Value
	data   interface{}
}

var _ Value = &baseValue{}

func (b *baseValue) Data() interface{} {
	return b.data
}

func (b *baseValue) Parent() Value {
	return b.parent
}

func (b *baseValue) Set(value interface{}) Value {
	b.data = value
	return b
}

func (b *baseValue) Map() Map {
	return newMap(b)
}

func (b *baseValue) Slice() Slice {
	return newSlice(b)
}
