package unstruct

type mapValue struct {
	parent Value
	key    string
}

var _ Value = &mapValue{}

func (b *mapValue) Data() interface{} {
	return b.parent.Map().Data()[b.key]
}

func (b *mapValue) Parent() Value {
	return b.parent
}

func (b *mapValue) Set(value interface{}) Value {
	b.parent.Map().Data()[b.key] = value
	return b
}

func (b *mapValue) Map() Map {
	return newMap(b)
}

func (b *mapValue) Slice() Slice {
	return newSlice(b)
}
