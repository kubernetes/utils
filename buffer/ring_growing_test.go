/*
Copyright 2017 The Kubernetes Authors.

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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrowth(t *testing.T) {
	t.Parallel()
	x := 10
	g := NewRingGrowing(1)
	for i := 0; i < x; i++ {
		assert.Equal(t, i, g.readable)
		g.WriteOne(i)
	}
	read := 0
	for g.readable > 0 {
		v, ok := g.ReadOne()
		assert.True(t, ok)
		assert.Equal(t, read, v)
		read++
	}
	assert.Equalf(t, x, read, "expected to have read %d items: %d", x, read)
	assert.Zerof(t, g.readable, "expected readable to be zero: %d", g.readable)
	assert.Equalf(t, 16, g.n, "expected N to be 16: %d", g.n)
}

func TestEmpty(t *testing.T) {
	t.Parallel()
	g := NewRingGrowing(1)
	_, ok := g.ReadOne()
	assert.False(t, ok)
}

// Write to this variable to prevent the compiler from optimizing the benchmark code away.
var codeRemovalStopper bool

// Benchmark only power-of-2 buffer sizes in the hypothesis that they will be
// the most commonly-used ones.
var bufferSizes = []int{512, 1024, 4096, 32768, 65536, 1024 * 1024}

// BenchmarkFillBuffer benchmarks the time it takes to fill an empty buffer,
// without ever growing the buffer.
func BenchmarkFillBuffer(b *testing.B) {
	for _, size := range bufferSizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			r := NewRingGrowing(size)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				// Reset the buffer so that it can be refilled.
				r.reset()
				for j := 0; j < size; j++ {
					r.WriteOne(j)
				}
				_, codeRemovalStopper = r.ReadOne()
			}

		})
	}
}

// BenchmarkEmptyBuffer benchmarks the time it takes to empty a full buffer.
func BenchmarkEmptyBuffer(b *testing.B) {
	for _, size := range bufferSizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			r := NewRingGrowing(size)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				// Fill the buffer so that it can be emptied.
				b.StopTimer()
				r.fillWithInts()
				b.StartTimer()

				// Now, empty the buffer.
				for j := 0; j < size; j++ {
					_, codeRemovalStopper = r.ReadOne()
				}
			}
		})
	}
}

// BenchmarkReadAndWrite benchmarks the time it takes to pass over a full buffer
// reading and writing one item at a time.
func BenchmarkReadAndWrite(b *testing.B) {
	for _, size := range bufferSizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			r := NewRingGrowing(size)
			r.fillWithInts()

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					_, codeRemovalStopper = r.ReadOne()
					r.WriteOne(j)
				}
			}
		})
	}
}

// BenchmarkGrowth benchmarks the time it takes to grow a full buffer when a new
// item is received. The growth is benchmarked when the read offset of the
// buffer is (1) at the beginning of the ring and (2) not at the beginning of
// the ring, because some implementations might handle these two cases
// differently. Notice that an implementation might perform better than another
// only because it grows less.
func BenchmarkGrowth(b *testing.B) {
	benchmarkGrowth := func(initialSize, readOffset int) func(b *testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r := NewRingGrowing(initialSize)
				r.fillWithInts()
				r.setReadOffset(readOffset)
				b.StartTimer()

				// The buffer is full, so the next write triggers a growth.
				r.WriteOne(i)
				// This read time will be negligible wrt the write, which needs
				// to allocate new memory.
				_, codeRemovalStopper = r.ReadOne()
			}
		}
	}

	for _, initialSize := range bufferSizes {
		b.Run(fmt.Sprintf("start_of_ring/initial_size=%d", initialSize), benchmarkGrowth(initialSize, 0))
		// Place the read offset in the middle of the ring, but anywhere other
		// than the beginning would work.
		b.Run(fmt.Sprintf("middle_of_ring/initial_size=%d", initialSize), benchmarkGrowth(initialSize, initialSize/2))
	}
}

// BenchmarkWritesWithGrowth benchmarks writing to a buffer way more items than
// its initial size can accomodate, triggering many growths.
func BenchmarkWritesWithGrowth(b *testing.B) {
	// Not sure whether these are sensible numbers.
	initialSize := 1024
	numberOfWrites := 65536

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := NewRingGrowing(initialSize)
		b.StartTimer()

		for j := 0; j < numberOfWrites; j++ {
			r.WriteOne(j)
		}
		_, codeRemovalStopper = r.ReadOne()
	}
}
