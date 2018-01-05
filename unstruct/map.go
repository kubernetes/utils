package unstruct

import "sort"

type object struct {
	value Value
}

var _ Map = &object{}

func newMap(value Value) Map {
	if value.Data() == nil {
		value.Set(map[string]interface{}{})
	}
	_, ok := value.Data().(map[string]interface{})
	if !ok {
		return nil
	}
	return &object{value: value}
}

func (o *object) Data() map[string]interface{} {
	d, ok := o.value.Data().(map[string]interface{})
	if !ok {
		return nil
	}
	return d
}

func (o *object) Parent() Value {
	return o.value.Parent()
}

func (o *object) Field(key string) Value {
	if !o.HasField(key) {
		o.Data()[key] = nil
	}
	return &mapValue{parent: o.value, key: key}
}

func (o *object) HasField(key string) bool {
	_, ok := o.Data()[key]
	return ok
}

func (o *object) Fields() []string {
	fields := []string{}
	for key := range o.Data() {
		fields = append(fields, key)
	}
	sort.Strings(fields)
	return fields
}
