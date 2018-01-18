package unstructured

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

func ExampleNewChainable() {
	data := map[interface{}]interface{}{
		"int":    1,
		"string": "my string",
		"list": []interface{}{
			"number1",
			map[interface{}]interface{}{
				"deep": "deep",
				"very": "deep",
			},
			3,
		},
	}
	NewChainable(data).
		SetField("int", 5).
		SetField("string", "new string").
		Field("list").
		SetAt(0, "number 0").
		At(1).SetField("deeper", "string")

	y, _ := yaml.Marshal(data)
	fmt.Println(string(y))
	// Output: int: 5
	// list:
	// - number 0
	// - deep: deep
	//   deeper: string
	//   very: deep
	// - 3
	// string: new string
}

func ExampleNewChainable_broken() {
	data := map[interface{}]interface{}{
		"list": []interface{}{
			map[interface{}]interface{}{
				"int": 5,
			},
		},
	}
	_, err := NewChainable(data).
		Field("list").
		At(0).
		Field("int").
		SetField("foo", "bar"). // This is an int, not a map. value is not set.
		SetAt(0, "some-value"). // We had an erorr already, value is not set.
		Data()
	fmt.Println(err)
	// Output: error setting key "foo" in .list[0].int: invalid type
}
