package unstruct_test

import (
	"fmt"

	"k8s.io/utils/unstruct"

	"gopkg.in/yaml.v2"
)

func Example() {
	u := unstruct.New(map[string]interface{}{
		"int": 1,
		"list": []interface{}{
			"A string.",
			2,
		},
	})

	u.Map().Field("int").Set("Number one")
	u.Map().Field("string").Set("my string")
	s := u.Map().Field("list").Slice()
	s.Append("Another string.")
	s.Append([]interface{}{})
	s.At(3).Slice().Append(4)

	y, _ := yaml.Marshal(u.Data())
	fmt.Println(string(y))
	// Output: int: Number one
	// list:
	// - A string.
	// - 2
	// - Another string.
	// - - 4
	// string: my string
}

// This example demonstrates how the interface lets you chain calls
// to easily set and change values in a unstruct type.
func Example_chaining() {
	u := unstruct.New(nil)

	u.Map().
		Field("int").Set("Number one").Parent().Map().
		Field("string").Set("my string").Parent().Map().
		Field("list").Slice().
		Append("A string.").Parent().Slice().
		Append(2).Parent().Slice().
		Append("Another string.").Parent().Slice().
		Append(nil).Slice().Append(4)

	y, _ := yaml.Marshal(u.Data())
	fmt.Println(string(y))
	// Output: int: Number one
	// list:
	// - A string.
	// - 2
	// - Another string.
	// - - 4
	// string: my string
}
