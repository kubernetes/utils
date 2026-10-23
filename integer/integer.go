/*
Copyright 2016 The Kubernetes Authors.

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

package integer

import "math"

// IntMax returns the maximum of the params.
// Deprecated: for new code, use the max() builtin instead.
//
//go:fix inline
func IntMax(a, b int) int {
	return max(a, b)
}

// IntMin returns the minimum of the params.
// Deprecated: for new code, use the min() builtin instead.
//
//go:fix inline
func IntMin(a, b int) int {
	return min(a, b)
}

// Int32Max returns the maximum of the params.
// Deprecated: for new code, use the max() builtin instead.
//
//go:fix inline
func Int32Max(a, b int32) int32 {
	return max(a, b)
}

// Int32Min returns the minimum of the params.
// Deprecated: for new code, use the min() builtin instead.
//
//go:fix inline
func Int32Min(a, b int32) int32 {
	return min(a, b)
}

// Int64Max returns the maximum of the params.
// Deprecated: for new code, use the max() builtin instead.
//
//go:fix inline
func Int64Max(a, b int64) int64 {
	return max(a, b)
}

// Int64Min returns the minimum of the params.
// Deprecated: for new code, use the min() builtin instead.
//
//go:fix inline
func Int64Min(a, b int64) int64 {
	return min(a, b)
}

// RoundToInt32 rounds floats into integer numbers.
// Deprecated: use math.Round() and a cast directly.
//
//go:fix inline
func RoundToInt32(a float64) int32 {
	return int32(math.Round(a))
}
