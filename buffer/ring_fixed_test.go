/*
Copyright 2025 The Kubernetes Authors.

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

package buffer

import (
	"io"
	"reflect"
	"testing"
)

func TestTypedRingFixedNew(t *testing.T) {
	t.Parallel()

	buf, err := NewTypedRingFixed[int](10)
	if err != nil {
		t.Errorf("size 10: unexpected error: %v", err)
	}
	if buf.Size() != 10 {
		t.Errorf("Size() = %d, want 10", buf.Size())
	}
	if _, err := NewTypedRingFixed[int](1); err != nil {
		t.Errorf("size 1: unexpected error: %v", err)
	}
	if _, err := NewTypedRingFixed[int](0); err != ErrInvalidSize {
		t.Errorf("size 0: expected ErrInvalidSize, got %v", err)
	}
	if _, err := NewTypedRingFixed[int](-1); err != ErrInvalidSize {
		t.Errorf("size -1: expected ErrInvalidSize, got %v", err)
	}
}

func TestTypedRingFixedWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		size        int
		writes      [][]int
		wantSlice   []int
		wantLen     int
		wantWritten int64
	}{
		{"short write", 10, [][]int{{1, 2, 3}}, []int{1, 2, 3}, 3, 3},
		{"full write", 3, [][]int{{1, 2, 3}}, []int{1, 2, 3}, 3, 3},
		{"long write", 3, [][]int{{1, 2, 3, 4, 5}}, []int{3, 4, 5}, 3, 5},
		{"empty write", 10, [][]int{{1, 2}, {}}, []int{1, 2}, 2, 2},
		{"overwrite", 5, [][]int{{1, 2, 3, 4, 5}, {6, 7}}, []int{3, 4, 5, 6, 7}, 5, 7},
		{"multiple small", 5, [][]int{{1}, {2}, {3}, {4}, {5}, {6}}, []int{2, 3, 4, 5, 6}, 5, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, _ := NewTypedRingFixed[int](tt.size)
			for _, w := range tt.writes {
				_, _ = buf.Write(w)
			}
			if got := buf.Slice(); !reflect.DeepEqual(got, tt.wantSlice) {
				t.Errorf("Slice() = %v, want %v", got, tt.wantSlice)
			}
			if got := buf.Len(); got != tt.wantLen {
				t.Errorf("Len() = %d, want %d", got, tt.wantLen)
			}
			if got := buf.TotalWritten(); got != tt.wantWritten {
				t.Errorf("TotalWritten() = %d, want %d", got, tt.wantWritten)
			}
		})
	}
}

func TestTypedRingFixedReset(t *testing.T) {
	t.Parallel()

	buf, _ := NewTypedRingFixed[int](4)
	_, _ = buf.Write([]int{1, 2, 3, 4, 5, 6})
	buf.Reset()

	if buf.Len() != 0 || buf.TotalWritten() != 0 || buf.Slice() != nil {
		t.Errorf("after Reset: Len=%d, TotalWritten=%d, Slice=%v; want 0, 0, nil",
			buf.Len(), buf.TotalWritten(), buf.Slice())
	}

	_, _ = buf.Write([]int{7, 8, 9, 10, 11})
	if got := buf.Slice(); !reflect.DeepEqual(got, []int{8, 9, 10, 11}) {
		t.Errorf("after Reset+Write: Slice() = %v, want %v", got, []int{8, 9, 10, 11})
	}
}

func BenchmarkTypedRingFixed_Write(b *testing.B) {
	b.ReportAllocs()
	buf, _ := NewTypedRingFixed[int](1024)
	data := make([]int, 100)

	for i := 0; i < b.N; i++ {
		_, _ = buf.Write(data)
	}
}

func BenchmarkTypedRingFixed_Slice(b *testing.B) {
	b.ReportAllocs()
	buf, _ := NewTypedRingFixed[int](1024)
	_, _ = buf.Write(make([]int, 500))  // Partial fill, cursor at 500
	_, _ = buf.Write(make([]int, 1024)) // Wrapped, cursor at 500

	for i := 0; i < b.N; i++ {
		_ = buf.Slice()
	}
}

func TestTypedRingFixedByteWriterInterface(t *testing.T) {
	t.Parallel()
	var _ io.Writer = &TypedRingFixed[byte]{}
}

func TestTypedRingFixedByteWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		size        int
		writes      []string
		wantString  string
		wantLen     int
		wantWritten int64
	}{
		{"short write", 1024, []string{"hello world"}, "hello world", 11, 11},
		{"full write", 11, []string{"hello world"}, "hello world", 11, 11},
		{"long write", 6, []string{"hello world"}, " world", 6, 11},
		{"huge write", 3, []string{"hello world"}, "rld", 3, 11},
		{"empty write", 10, []string{"hello", ""}, "hello", 5, 5},
		{"size one", 1, []string{"a", "b", "xyz"}, "z", 1, 5},
		{"overwrite", 10, []string{"0123456789", "abc"}, "3456789abc", 10, 13},
		{"multiple small", 10, []string{"aa", "bb", "cc", "dd", "ee", "ff"}, "bbccddeeff", 10, 12},
		{"many single bytes", 3, []string{"h", "e", "l", "l", "o", " ", "w", "o", "r", "l", "d"}, "rld", 3, 11},
		{
			"multi part",
			16,
			[]string{"hello world\n", "this is a test\n", "my cool input\n"},
			"t\nmy cool input\n",
			16,
			41,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, _ := NewTypedRingFixed[byte](tt.size)
			for _, w := range tt.writes {
				_, _ = buf.Write([]byte(w))
			}
			if got := string(buf.Slice()); got != tt.wantString {
				t.Errorf("Slice() = %q, want %q", got, tt.wantString)
			}
			if got := buf.Len(); got != tt.wantLen {
				t.Errorf("Len() = %d, want %d", got, tt.wantLen)
			}
			if got := buf.TotalWritten(); got != tt.wantWritten {
				t.Errorf("TotalWritten() = %d, want %d", got, tt.wantWritten)
			}
		})
	}
}

func TestTypedRingFixedByteWriteReturnValue(t *testing.T) {
	t.Parallel()

	buf, _ := NewTypedRingFixed[byte](5)

	// Write returns original length even when truncated
	if n, err := buf.Write([]byte("hello world")); n != 11 || err != nil {
		t.Errorf("Write() = (%d, %v), want (11, nil)", n, err)
	}
}

func TestTypedRingFixedByteReset(t *testing.T) {
	t.Parallel()

	buf, _ := NewTypedRingFixed[byte](4)
	_, _ = buf.Write([]byte("hello world\n"))
	_, _ = buf.Write([]byte("this is a test\n"))
	buf.Reset()

	if buf.Len() != 0 || buf.TotalWritten() != 0 || buf.Slice() != nil {
		t.Errorf("after Reset: Len=%d, TotalWritten=%d, Slice=%v; want 0, 0, nil",
			buf.Len(), buf.TotalWritten(), buf.Slice())
	}

	// Write after reset
	_, _ = buf.Write([]byte("hello"))
	if got := string(buf.Slice()); got != "ello" {
		t.Errorf("after Reset+Write: Slice() = %q, want %q", got, "ello")
	}
}

func TestTypedRingFixedByteSlice(t *testing.T) {
	t.Parallel()

	buf, _ := NewTypedRingFixed[byte](10)

	// Empty
	if buf.Slice() != nil {
		t.Errorf("empty buffer: Slice() = %v, want nil", buf.Slice())
	}

	// Partial fill - returns slice of internal buffer
	_, _ = buf.Write([]byte("hello"))
	if got := string(buf.Slice()); got != "hello" {
		t.Errorf("partial fill: Slice() = %q, want %q", got, "hello")
	}

	// Exact fill at cursor 0 - returns internal buffer directly
	buf.Reset()
	_, _ = buf.Write([]byte("0123456789"))
	if got := string(buf.Slice()); got != "0123456789" {
		t.Errorf("exact fill: Slice() = %q, want %q", got, "0123456789")
	}

	// Wrapped - returns new slice with reordered data
	_, _ = buf.Write([]byte("ab"))
	if got := string(buf.Slice()); got != "23456789ab" {
		t.Errorf("wrapped: Slice() = %q, want %q", got, "23456789ab")
	}
}

func BenchmarkTypedRingFixedByte_Write(b *testing.B) {
	b.ReportAllocs()
	buf, _ := NewTypedRingFixed[byte](1024)
	data := make([]byte, 100)

	for i := 0; i < b.N; i++ {
		_, _ = buf.Write(data)
	}
}

func BenchmarkTypedRingFixedByte_Slice(b *testing.B) {
	b.ReportAllocs()
	buf, _ := NewTypedRingFixed[byte](1024)
	_, _ = buf.Write(make([]byte, 500))  // Partial fill, cursor at 500
	_, _ = buf.Write(make([]byte, 1024)) // Wrapped, cursor at 500

	for i := 0; i < b.N; i++ {
		_ = buf.Slice()
	}
}
