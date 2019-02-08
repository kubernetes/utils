/*
Copyright 2019 The Kubernetes Authors.

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

package semantic

import (
	"k8s.io/utils/semantic/internal/third_party/forked/golang/reflect"
)

// Equality provides an extensible semantic deep-equal equality. Semantic here means that
// empty and nil slices and maps are identified.
type Equality interface {
	// AddFuncs is a shortcut for multiple calls to AddFunc.
	AddFuncs(funcs ...interface{}) error
	// AddFunc uses func as an equality function: it must take
	// two parameters of the same type, and return a boolean.
	AddFunc(eqFunc interface{}) error
	// DeepEqual is like reflect.DeepEqual, but focused on semantic equality
	// instead of memory equality.
	//
	// It will uses the registered equality functions if it finds types that match.
	//
	// An empty slice *is* equal to a nil slice for our purposes; same for maps.
	//
	// Unexported field members cannot be compared and will cause an imformative panic.
	DeepEqual(a1, a2 interface{}) bool
	// DeepDerivative is similar to DeepEqual except that unset fields in a1 are
	// ignored (not compared). This allows us to focus on the fields that matter to
	// the semantic comparison.
	//
	// The unset fields include a nil pointer and an empty string.
	DeepDerivative(a1, a2 interface{}) bool
}

// NewEquality creates a semantic equality object with the given equality funcs.
func NewEquality(funcs ...interface{}) (Equality, error) {
	eq := reflect.Equalities{}
	if err := eq.AddFuncs(funcs...); err != nil {
		return nil, err
	}
	return eq, nil
}

// NewEqualityOrDie creates a semantic equality object with the given equality funcs and panics on errors.
func NewEqualityOrDie(funcs ...interface{}) Equality {
	eq, err := NewEquality(funcs...)
	if err != nil {
		panic(err)
	}
	return eq
}
