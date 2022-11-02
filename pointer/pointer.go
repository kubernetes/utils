/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pointer

import (
	"fmt"
	"reflect"
	"time"
)

// AllPtrFieldsNil tests whether all pointer fields in a struct are nil.  This is useful when,
// for example, an API struct is handled by plugins which need to distinguish
// "no plugin accepted this spec" from "this spec is empty".
//
// This function is only valid for structs and pointers to structs.  Any other
// type will cause a panic.  Passing a typed nil pointer will return true.
func AllPtrFieldsNil(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	if !v.IsValid() {
		panic(fmt.Sprintf("reflect.ValueOf() produced a non-valid Value for %#v", obj))
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && !v.Field(i).IsNil() {
			return false
		}
	}
	return true
}

// Ref returns a pointer to the value.
func Ref[T any](v T) *T {
	return &v
}

// Deref dereferences the ptr and returns it if not nil, or else returns def.
func Deref[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}

func Equal[T comparable](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return *a == *b
}

var (

	// Int returns a pointer to an int
	Int = Ref[int]

	// IntPtr is a function variable referring to Int.
	//
	// Deprecated: Use Int instead.
	IntPtr = Int // for back-compat

	// IntDeref dereferences the int ptr and returns it if not nil, or else
	// returns def.
	IntDeref = Deref[int]

	// IntPtrDerefOr is a function variable referring to IntDeref.
	//
	// Deprecated: Use IntDeref instead.
	IntPtrDerefOr = IntDeref // for back-compat

	// Int32 returns a pointer to an int32.
	Int32 = Ref[int32]

	// Int32Ptr is a function variable referring to Int32.
	//
	// Deprecated: Use Int32 instead.
	Int32Ptr = Int32 // for back-compat

	// Int32Deref dereferences the int32 ptr and returns it if not nil, or else
	// returns def.
	Int32Deref = Deref[int32]

	// Int32PtrDerefOr is a function variable referring to Int32Deref.
	//
	// Deprecated: Use Int32Deref instead.
	Int32PtrDerefOr = Int32Deref // for back-compat

	// Int32Equal returns true if both arguments are nil or both arguments
	// dereference to the same value.
	Int32Equal = Equal[int32]

	// Uint returns a pointer to an uint
	Uint = Ref[uint]

	// UintPtr is a function variable referring to Uint.
	//
	// Deprecated: Use Uint instead.
	UintPtr = Uint // for back-compat

	// UintDeref dereferences the uint ptr and returns it if not nil, or else
	// returns def.
	UintDeref = Deref[uint]

	// UintPtrDerefOr is a function variable referring to UintDeref.
	//
	// Deprecated: Use UintDeref instead.
	UintPtrDerefOr = UintDeref // for back-compat

	// Uint32 returns a pointer to an uint32.
	Uint32 = Ref[uint32]

	// Uint32Ptr is a function variable referring to Uint32.
	//
	// Deprecated: Use Uint32 instead.
	Uint32Ptr = Uint32 // for back-compat

	// Uint32Deref dereferences the uint32 ptr and returns it if not nil, or else
	// returns def.
	Uint32Deref = Deref[uint32]

	// Uint32PtrDerefOr is a function variable referring to Uint32Deref.
	//
	// Deprecated: Use Uint32Deref instead.
	Uint32PtrDerefOr = Uint32Deref // for back-compat

	// Uint32Equal returns true if both arguments are nil or both arguments
	// dereference to the same value.
	Uint32Equal = Equal[uint32]

	// Int64 returns a pointer to an int64.
	Int64 = Ref[int64]

	// Int64Ptr is a function variable referring to Int64.
	//
	// Deprecated: Use Int64 instead.
	Int64Ptr = Int64 // for back-compat

	// Int64Deref dereferences the int64 ptr and returns it if not nil, or else
	// returns def.
	Int64Deref = Deref[int64]

	// Int64PtrDerefOr is a function variable referring to Int64Deref.
	//
	// Deprecated: Use Int64Deref instead.
	Int64PtrDerefOr = Int64Deref // for back-compat

	// Int64Equal returns true if both arguments are nil or both arguments
	// dereference to the same value.
	Int64Equal = Equal[int64]

	// Uint64 returns a pointer to an uint64.
	Uint64 = Ref[uint64]

	// Uint64Ptr is a function variable referring to Uint64.
	//
	// Deprecated: Use Uint64 instead.
	Uint64Ptr = Uint64 // for back-compat

	// Uint64Deref dereferences the uint64 ptr and returns it if not nil, or else
	// returns def.
	Uint64Deref = Deref[uint64]

	// Uint64PtrDerefOr is a function variable referring to Uint64Deref.
	//
	// Deprecated: Use Uint64Deref instead.
	Uint64PtrDerefOr = Uint64Deref // for back-compat

	// Uint64Equal returns true if both arguments are nil or both arguments
	// dereference to the same value.
	Uint64Equal = Equal[uint64]

	// Bool returns a pointer to a bool.
	Bool = Ref[bool]

	// BoolPtr is a function variable referring to Bool.
	//
	// Deprecated: Use Bool instead.
	BoolPtr = Bool // for back-compat

	// BoolDeref dereferences the bool ptr and returns it if not nil, or else
	// returns def.
	BoolDeref = Deref[bool]

	// BoolPtrDerefOr is a function variable referring to BoolDeref.
	//
	// Deprecated: Use BoolDeref instead.
	BoolPtrDerefOr = BoolDeref // for back-compat

	// BoolEqual returns true if both arguments are nil or both arguments
	// dereference to the same value.
	BoolEqual = Equal[bool]

	// String returns a pointer to a string.
	String = Ref[string]

	// StringPtr is a function variable referring to String.
	//
	// Deprecated: Use String instead.
	StringPtr = String // for back-compat

	// StringDeref dereferences the string ptr and returns it if not nil, or else
	// returns def.
	StringDeref = Deref[string]

	// StringPtrDerefOr is a function variable referring to StringDeref.
	//
	// Deprecated: Use StringDeref instead.
	StringPtrDerefOr = StringDeref // for back-compat

	// StringEqual returns true if both arguments are nil or both arguments
	// dereference to the same value.
	StringEqual = Equal[string]

	// Float32 returns a pointer to a float32.
	Float32 = Ref[float32]

	// Float32Ptr is a function variable referring to Float32.
	//
	// Deprecated: Use Float32 instead.
	Float32Ptr = Float32

	// Float32Deref dereferences the float32 ptr and returns it if not nil, or else
	// returns def.
	Float32Deref = Deref[float32]

	// Float32PtrDerefOr is a function variable referring to Float32Deref.
	//
	// Deprecated: Use Float32Deref instead.
	Float32PtrDerefOr = Float32Deref // for back-compat

	// Float32Equal returns true if both arguments are nil or both arguments
	// dereference to the same value.
	Float32Equal = Equal[float32]

	// Float64 returns a pointer to a float64.
	Float64 = Ref[float64]

	// Float64Ptr is a function variable referring to Float64.
	//
	// Deprecated: Use Float64 instead.
	Float64Ptr = Float64

	// Float64Deref dereferences the float64 ptr and returns it if not nil, or else
	// returns def.
	Float64Deref = Deref[float64]

	// Float64PtrDerefOr is a function variable referring to Float64Deref.
	//
	// Deprecated: Use Float64Deref instead.
	Float64PtrDerefOr = Float64Deref // for back-compat

	// Float64Equal returns true if both arguments are nil or both arguments
	// dereference to the same value.
	Float64Equal = Equal[float64]

	// Duration returns a pointer to a time.Duration.
	Duration = Ref[time.Duration]

	// DurationDeref dereferences the time.Duration ptr and returns it if not nil, or else
	// returns def.
	DurationDeref = Deref[time.Duration]

	// DurationEqual returns true if both arguments are nil or both arguments
	// dereference to the same value.
	DurationEqual = Equal[time.Duration]
)
