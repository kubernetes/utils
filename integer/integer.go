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

// IntMax returns the maximum of the params
func IntMax(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// IntMin returns the minimum of the params
func IntMin(a, b int) int {
	if b < a {
		return b
	}
	return a
}

// IntBounded bound the value in range [lower, upper]
func IntBounded(value, lower, upper int) int {
	return IntMin(IntMax(value, lower), upper)
}

// Uint8Max returns the maximum of the params
func Uint8Max(a, b uint8) uint8 {
	if b > a {
		return b
	}
	return a
}

// Uint8Min returns the minimum of the params
func Uint8Min(a, b uint8) uint8 {
	if b < a {
		return b
	}
	return a
}

// Uint8Bounded bound the value in range [lower, upper]
func Uint8Bounded(value, lower, upper uint8) uint8 {
	return Uint8Min(Uint8Max(value, lower), upper)
}

// Int32Max returns the maximum of the params
func Int32Max(a, b int32) int32 {
	if b > a {
		return b
	}
	return a
}

// Int32Min returns the minimum of the params
func Int32Min(a, b int32) int32 {
	if b < a {
		return b
	}
	return a
}

// Int32Bounded bound the value in range [lower, upper]
func Int32Bounded(value, lower, upper int32) int32 {
	return Int32Min(Int32Max(value, lower), upper)
}

// Int64Max returns the maximum of the params
func Int64Max(a, b int64) int64 {
	if b > a {
		return b
	}
	return a
}

// Int64Min returns the minimum of the params
func Int64Min(a, b int64) int64 {
	if b < a {
		return b
	}
	return a
}

// Int64Bounded bound the value in range [lower, upper]
func Int64Bounded(value, lower, upper int64) int64 {
	return Int64Min(Int64Max(value, lower), upper)
}

// RoundToInt32 rounds floats into integer numbers.
func RoundToInt32(a float64) int32 {
	if a < 0 {
		return int32(a - 0.5)
	}
	return int32(a + 0.5)
}
