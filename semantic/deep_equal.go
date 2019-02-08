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
// empty and nil slices and maps are identified (by default, can be disabled) and custom
// equality funcs can be added to follow custom equality semantics.
type Equality interface {
	// AddFuncs is a shortcut for multiple calls to AddFunc.
	AddFuncs(funcs ...interface{}) error
	// AddFunc uses func as an equality function: it must take two parameters of the same
	// type, and return a boolean.
	AddFunc(eqFunc interface{}) error
	// DeepEqual is like reflect.DeepEqual, but by applying the registered equality funcs
	// and identifies nil and empty slices and maps (if not disabled by options).
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

// EqualityOptions configures a semantic equality object.
type EqualityOptions struct {
	// DisableSemanticMaps set to true disables identification of nil and empty maps.
	DisableSemanticMaps bool
	// DisableSemanticSlices set to true disables identification of nil and empty slices.
	DisableSemanticSlices bool
}

// NewEquality creates a semantic equality object with the given equality funcs.
func NewEquality(opt EqualityOptions, funcs ...interface{}) (Equality, error) {
	eq := reflect.Equalities{
		DisableSemanticMaps:   opt.DisableSemanticMaps,
		DisableSemanticSlices: opt.DisableSemanticSlices,
	}
	if err := eq.AddFuncs(funcs...); err != nil {
		return nil, err
	}
	return &eq, nil
}

// NewEqualityOrDie creates a semantic equality object with the given equality funcs and panics on errors.
func NewEqualityOrDie(opt EqualityOptions, funcs ...interface{}) Equality {
	eq, err := NewEquality(opt, funcs...)
	if err != nil {
		panic(err)
	}
	return eq
}
