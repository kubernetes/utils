package unstructured

import (
	"fmt"
)

func ExampleUnstructured_Field() {
	data := map[interface{}]interface{}{
		"int":    1,
		"string": "my string",
	}

	fmt.Println(Unstructured{Data: data}.Field("string"))
	// Output: {my string} <nil>
}

func ExampleUnstructured_Field_missingkey() {
	data := map[interface{}]interface{}{
		"int":    1,
		"string": "my string",
	}

	fmt.Println(Unstructured{Data: data}.Field("randomkey"))
	// Output: {<nil>} field not found
}

func ExampleUnstructured_Field_invalidtype() {
	data := []interface{}{
		"string",
		5,
	}

	fmt.Println(Unstructured{Data: data}.Field("randomkey"))
	// Output: {<nil>} invalid type
}

func ExampleUnstructured_SetField() {
	data := map[interface{}]interface{}{
		"int":    1,
		"string": "my string",
	}

	Unstructured{Data: data}.SetField("string", "my new string")
	fmt.Println(Unstructured{Data: data}.Field("string"))
	// Output: {my new string} <nil>
}

func ExampleUnstructured_At() {
	data := []interface{}{
		"string",
		5,
	}

	fmt.Println(Unstructured{Data: data}.At(1))
	// Output: {5} <nil>
}

func ExampleUnstructured_SetAt() {
	data := []interface{}{
		"string",
		5,
	}

	Unstructured{Data: data}.SetAt(1, 4.2)
	fmt.Println(Unstructured{Data: data}.At(1))
	// Output: {4.2} <nil>
}
