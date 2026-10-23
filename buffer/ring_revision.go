/*
Copyright 2026 The Kubernetes Authors.

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

// RevisionRingBuffer stores revision-ordered entries in a circular buffer.
// It is suitable for append-mostly histories that are queried or trimmed by revision.
// When the buffer is full, appending a new entry overwrites the oldest one.
// Not thread safe.
type RevisionRingBuffer[T any] struct {
	buffer []entry[T]
	// head is the index immediately after the last non-empty entry in the
	// buffer, i.e. the next write position.
	head, tail, size int
	revisionOf       RevisionOf[T]
}

type entry[T any] struct {
	revision int64
	item     T
}

type (
	// RevisionOf extracts the revision used to order entries in the ring.
	RevisionOf[T any] func(T) int64
	// IterFunc is called for each visited entry. Returning false stops iteration.
	IterFunc[T any] func(rev int64, item T) bool
)

// NewRevisionRingBuffer creates a revision-ordered ring buffer with fixed
// initial capacity.
func NewRevisionRingBuffer[T any](capacity int, revisionOf RevisionOf[T]) (*RevisionRingBuffer[T], error) {
	if capacity <= 0 {
		return nil, ErrInvalidSize
	}
	return &RevisionRingBuffer[T]{
		buffer:     make([]entry[T], capacity),
		revisionOf: revisionOf,
	}, nil
}

// Append adds item to the ring. When the ring is full, the oldest entry is
// overwritten.
func (r *RevisionRingBuffer[T]) Append(item T) {
	entry := entry[T]{revision: r.revisionOf(item), item: item}
	if r.full() {
		r.tail = (r.tail + 1) % len(r.buffer)
	} else {
		r.size++
	}
	r.buffer[r.head] = entry
	r.head = (r.head + 1) % len(r.buffer)
}

func (r *RevisionRingBuffer[T]) full() bool {
	return r.size == len(r.buffer)
}

// AscendGreaterOrEqual iterates through entries in ascending order starting
// from the first entry with revision >= pivot.
func (r *RevisionRingBuffer[T]) AscendGreaterOrEqual(pivot int64, iter IterFunc[T]) {
	for i := r.findFirstIndexGreaterOrEqual(pivot); i < r.size; i++ {
		entry := r.at(i)
		if !iter(entry.revision, entry.item) {
			return
		}
	}
}

// AscendLessThan iterates in ascending order over entries with revision <
// pivot.
func (r *RevisionRingBuffer[T]) AscendLessThan(pivot int64, iter IterFunc[T]) {
	for i := 0; i < r.findFirstIndexGreaterOrEqual(pivot); i++ {
		entry := r.at(i)
		if !iter(entry.revision, entry.item) {
			return
		}
	}
}

// DescendGreaterThan iterates in descending order over entries with revision >
// pivot.
func (r *RevisionRingBuffer[T]) DescendGreaterThan(pivot int64, iter IterFunc[T]) {
	for i := r.size - 1; i > r.findLastIndexLessOrEqual(pivot); i-- {
		entry := r.at(i)
		if !iter(entry.revision, entry.item) {
			return
		}
	}
}

// DescendLessOrEqual iterates in descending order over entries with revision <=
// pivot.
func (r *RevisionRingBuffer[T]) DescendLessOrEqual(pivot int64, iter IterFunc[T]) {
	for i := r.findLastIndexLessOrEqual(pivot); i >= 0; i-- {
		entry := r.at(i)
		if !iter(entry.revision, entry.item) {
			return
		}
	}
}

// PeekLatest returns the most recently-appended revision, or 0 if empty.
func (r *RevisionRingBuffer[T]) PeekLatest() int64 {
	if r.size == 0 {
		return 0
	}
	idx := (r.head - 1 + len(r.buffer)) % len(r.buffer)
	return r.buffer[idx].revision
}

// PeekOldest returns the oldest revision currently stored, or 0 if empty.
func (r *RevisionRingBuffer[T]) PeekOldest() int64 {
	if r.size == 0 {
		return 0
	}
	return r.buffer[r.tail].revision
}

// RebaseHistory clears the ring and zeroes out stored entries.
func (r *RevisionRingBuffer[T]) RebaseHistory() {
	r.head, r.tail, r.size = 0, 0, 0
	for i := range r.buffer {
		r.buffer[i] = entry[T]{}
	}
}

func (r *RevisionRingBuffer[T]) moduloIndex(index int) int {
	return (index + len(r.buffer)) % len(r.buffer)
}

func (r *RevisionRingBuffer[T]) at(logicalIndex int) entry[T] {
	return r.buffer[r.moduloIndex(r.tail+logicalIndex)]
}

func (r *RevisionRingBuffer[T]) findFirstIndexGreaterOrEqual(pivot int64) int {
	left, right := 0, r.size-1
	for left <= right {
		// Prevent overflow; see https://github.com/golang/go/blob/master/src/sort/search.go#L105.
		mid := int(uint(left+right) >> 1)
		if r.at(mid).revision >= pivot {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return left
}

func (r *RevisionRingBuffer[T]) findLastIndexLessOrEqual(pivot int64) int {
	left, right := 0, r.size-1
	for left <= right {
		// Prevent overflow; see https://github.com/golang/go/blob/master/src/sort/search.go#L105.
		mid := int(uint(left+right) >> 1)
		if r.at(mid).revision <= pivot {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return right
}

// Len returns the number of entries currently stored in the buffer.
func (r *RevisionRingBuffer[T]) Len() int {
	return r.size
}

// ReplaceLatest overwrites the most recently appended entry with the given
// item. Panics if the buffer is empty.
func (r *RevisionRingBuffer[T]) ReplaceLatest(item T) {
	if r.size == 0 {
		panic("RevisionRingBuffer.ReplaceLatest called on empty buffer")
	}
	r.buffer[r.moduloIndex(r.head-1)] = entry[T]{revision: r.revisionOf(item), item: item}
}

// RemoveLess drops all entries whose revision is strictly less than rv.
func (r *RevisionRingBuffer[T]) RemoveLess(rv int64) {
	if r.size == 0 {
		return
	}
	idx := r.findFirstIndexGreaterOrEqual(rv)
	if idx <= 0 {
		return
	}
	if idx == r.size {
		r.RebaseHistory()
		return
	}
	for i := range idx {
		r.buffer[r.moduloIndex(r.tail+i)] = entry[T]{}
	}
	r.tail = r.moduloIndex(r.tail + idx)
	r.size -= idx
}

// EnsureCapacity grows the underlying slice by doubling until it can hold at
// least minCapacity entries. After expansion the logical entries are laid out
// contiguously starting at index 0.
func (r *RevisionRingBuffer[T]) EnsureCapacity(minCapacity int) {
	if len(r.buffer) >= minCapacity {
		return
	}
	newCap := len(r.buffer)
	if newCap == 0 {
		newCap = minCapacity
	}
	for newCap < minCapacity {
		newCap *= 2
	}
	if newCap < defaultRingSize {
		newCap = defaultRingSize
	}
	newBuffer := make([]entry[T], newCap)
	to := r.tail + r.size
	if to <= len(r.buffer) {
		copy(newBuffer, r.buffer[r.tail:to])
	} else {
		copied := copy(newBuffer, r.buffer[r.tail:])
		copy(newBuffer[copied:], r.buffer[:(to%len(r.buffer))])
	}
	r.buffer = newBuffer
	r.tail = 0
	r.head = r.size
}
