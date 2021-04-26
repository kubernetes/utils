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

import "testing"

func TestIntMax(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"a should be bigger", args{2, 1}, 2},
		{"b should be bigger", args{1, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntMax(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("IntMax() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntMin(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"a should be smaller", args{1, 2}, 1},
		{"b should be smaller", args{2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntMin(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("IntMin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntBounded(t *testing.T) {
	type args struct {
		value int
		lower int
		upper int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"unchanged", args{2, 1, 3}, 2},
		{"changed to lower", args{0, 1, 3}, 1},
		{"changed to upper", args{4, 1, 3}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntBounded(tt.args.value, tt.args.lower, tt.args.upper); got != tt.want {
				t.Errorf("IntBounded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8Max(t *testing.T) {
	type args struct {
		a uint8
		b uint8
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"a should be bigger", args{2, 1}, 2},
		{"b should be bigger", args{1, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint8Max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Uint8Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8Min(t *testing.T) {
	type args struct {
		a uint8
		b uint8
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"a should be smaller", args{1, 2}, 1},
		{"b should be smaller", args{2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint8Min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Uint8Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8Bounded(t *testing.T) {
	type args struct {
		value uint8
		lower uint8
		upper uint8
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"unchanged", args{2, 1, 3}, 2},
		{"changed to lower", args{0, 1, 3}, 1},
		{"changed to upper", args{4, 1, 3}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint8Bounded(tt.args.value, tt.args.lower, tt.args.upper); got != tt.want {
				t.Errorf("Uint8Bounded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32Max(t *testing.T) {
	type args struct {
		a int32
		b int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{"a should be bigger", args{2, 1}, 2},
		{"b should be bigger", args{1, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int32Max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Int32Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32Min(t *testing.T) {
	type args struct {
		a int32
		b int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{"a should be smaller", args{1, 2}, 1},
		{"b should be smaller", args{2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int32Min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Int32Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32Bounded(t *testing.T) {
	type args struct {
		value int32
		lower int32
		upper int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{"unchanged", args{2, 1, 3}, 2},
		{"changed to lower", args{0, 1, 3}, 1},
		{"changed to upper", args{4, 1, 3}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int32Bounded(tt.args.value, tt.args.lower, tt.args.upper); got != tt.want {
				t.Errorf("Int32Bounded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Max(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"a should be bigger", args{2, 1}, 2},
		{"b should be bigger", args{1, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64Max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Int64Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Min(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"a should be smaller", args{1, 2}, 1},
		{"b should be smaller", args{2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64Min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Int64Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Bounded(t *testing.T) {
	type args struct {
		value int64
		lower int64
		upper int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"unchanged", args{2, 1, 3}, 2},
		{"changed to lower", args{0, 1, 3}, 1},
		{"changed to upper", args{4, 1, 3}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64Bounded(tt.args.value, tt.args.lower, tt.args.upper); got != tt.want {
				t.Errorf("Int64Bounded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundToInt32(t *testing.T) {
	type args struct {
		a float64
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{"round 5.5 to 6", args{5.5}, 6},
		{"round -3.7 to -4", args{-3.7}, -4},
		{"round 3.49 to 3", args{3.49}, 3},
		{"round 5.5 to 6", args{5.5}, 6},
		{"round -7.9 to -8", args{-7.9}, -8},
		{"round -4.499999 to -4", args{-4.499999}, -4},
		{"round 0 to 0", args{0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundToInt32(tt.args.a); got != tt.want {
				t.Errorf("RoundToInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}
