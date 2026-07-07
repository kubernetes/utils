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

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

type testEvent struct {
	Key         string
	ModRevision int64
}

type iterMethod string

const (
	ascendGTE  iterMethod = "AscendGreaterOrEqual"
	ascendLT   iterMethod = "AscendLessThan"
	descendGT  iterMethod = "DescendGreaterThan"
	descendLTE iterMethod = "DescendLessOrEqual"
)

func TestPeekLatestAndOldest(t *testing.T) {
	tests := []struct {
		name          string
		capacity      int
		revs          []int64
		wantLatestRev int64
		wantOldestRev int64
	}{
		{
			name:          "empty_buffer",
			capacity:      4,
			revs:          nil,
			wantLatestRev: 0,
			wantOldestRev: 0,
		},
		{
			name:          "single_element",
			capacity:      8,
			revs:          []int64{1},
			wantLatestRev: 1,
			wantOldestRev: 1,
		},
		{
			name:          "ascending_fill",
			capacity:      4,
			revs:          []int64{1, 2, 3, 4},
			wantLatestRev: 4,
			wantOldestRev: 1,
		},
		{
			name:          "overwrite_when_full",
			capacity:      3,
			revs:          []int64{5, 6, 7, 8},
			wantLatestRev: 8,
			wantOldestRev: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := mustNewTestRingBuffer(t, tt.capacity)
			for _, r := range tt.revs {
				batch, err := makeEventBatch(r, "k", 1)
				if err != nil {
					t.Fatalf("makeEventBatch(%d, k, 1) failed: %v", r, err)
				}
				rb.Append(batch)
			}

			if got := rb.PeekLatest(); got != tt.wantLatestRev {
				t.Fatalf("PeekLatest()=%d, want=%d", got, tt.wantLatestRev)
			}
			if got := rb.PeekOldest(); got != tt.wantOldestRev {
				t.Fatalf("PeekOldest()=%d, want=%d", got, tt.wantOldestRev)
			}
		})
	}
}

func TestIterationMethods(t *testing.T) {
	type iterTestCase struct {
		method            iterMethod
		pivot             int64
		wantIterRevisions []int64
	}
	tests := []struct {
		name           string
		capacity       int
		setupRevisions []int64
		cases          []iterTestCase
	}{
		{
			name:           "empty_buffer",
			capacity:       4,
			setupRevisions: nil,
			cases: []iterTestCase{
				{ascendGTE, 0, []int64{}},
				{ascendLT, 10, []int64{}},
				{descendGT, 0, []int64{}},
				{descendLTE, 10, []int64{}},
			},
		},
		{
			name:           "basic_filtering",
			capacity:       5,
			setupRevisions: []int64{1, 2, 3},
			cases: []iterTestCase{
				{ascendGTE, 0, []int64{1, 2, 3}},
				{ascendGTE, 2, []int64{2, 3}},
				{ascendGTE, 100, []int64{}},
				{ascendLT, 3, []int64{1, 2}},
				{ascendLT, 1, []int64{}},
				{ascendLT, 100, []int64{1, 2, 3}},
				{descendGT, 1, []int64{3, 2}},
				{descendGT, 3, []int64{}},
				{descendGT, 0, []int64{3, 2, 1}},
				{descendLTE, 2, []int64{2, 1}},
				{descendLTE, 3, []int64{3, 2, 1}},
				{descendLTE, 0, []int64{}},
			},
		},
		{
			name:           "overflowed stores only entries within capacity",
			capacity:       3,
			setupRevisions: []int64{20, 21, 22, 23, 24},
			cases: []iterTestCase{
				{ascendGTE, 23, []int64{23, 24}},
				{ascendGTE, 0, []int64{22, 23, 24}},
				{ascendLT, 23, []int64{22}},
				{ascendLT, 25, []int64{22, 23, 24}},
				{descendGT, 22, []int64{24, 23}},
				{descendGT, 25, []int64{}},
				{descendLTE, 23, []int64{23, 22}},
				{descendLTE, 24, []int64{24, 23, 22}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := setupRingBuffer(t, tt.capacity, tt.setupRevisions)

			for _, tc := range tt.cases {
				t.Run(fmt.Sprintf("%s_pivot_%d", tc.method, tc.pivot), func(t *testing.T) {
					got := collectRevisions(rb, tc.method, tc.pivot)
					if !slices.Equal(tc.wantIterRevisions, got) {
						t.Fatalf("%s(%d)=%v, want %v", tc.method, tc.pivot, got, tc.wantIterRevisions)
					}
				})
			}
		})
	}
}

func TestIterationWithBatching(t *testing.T) {
	rb := mustNewTestRingBuffer(t, 6)
	batchA := []testEvent{
		{Key: "key-a", ModRevision: 5},
	}
	batchB := []testEvent{
		{Key: "key-b-1", ModRevision: 10},
		{Key: "key-b-2", ModRevision: 10},
		{Key: "key-b-3", ModRevision: 10},
	}
	batchC := []testEvent{
		{Key: "key-c", ModRevision: 12},
	}
	rb.Append(batchA)
	rb.Append(batchB)
	rb.Append(batchC)

	tests := []struct {
		name   string
		method iterMethod
		pivot  int64
		want   [][]testEvent
	}{
		{
			name:   "ascending_gte_includes_batched_revision",
			method: ascendGTE,
			pivot:  10,
			want: [][]testEvent{
				{
					{Key: "key-b-1", ModRevision: 10},
					{Key: "key-b-2", ModRevision: 10},
					{Key: "key-b-3", ModRevision: 10},
				},
				{
					{Key: "key-c", ModRevision: 12},
				},
			},
		},
		{
			name:   "ascending_lt_stops_before_batched_revision",
			method: ascendLT,
			pivot:  10,
			want: [][]testEvent{
				{
					{Key: "key-a", ModRevision: 5},
				},
			},
		},
		{
			name:   "all_revisions_with_proper_batch_sizes",
			method: ascendGTE,
			pivot:  0,
			want: [][]testEvent{
				{
					{Key: "key-a", ModRevision: 5},
				},
				{
					{Key: "key-b-1", ModRevision: 10},
					{Key: "key-b-2", ModRevision: 10},
					{Key: "key-b-3", ModRevision: 10},
				},
				{
					{Key: "key-c", ModRevision: 12},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got [][]testEvent

			rb.iterate(tt.method, tt.pivot, func(rev int64, events []testEvent) bool {
				got = append(got, events)
				return true
			})

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("events=%v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterationEarlyStop(t *testing.T) {
	rb := setupRingBuffer(t, 5, []int64{5, 10, 15, 20})
	tests := []struct {
		name      string
		method    iterMethod
		pivot     int64
		stopAfter int
		want      []int64
	}{
		{
			name:      "find_first_match_ascending",
			method:    ascendGTE,
			pivot:     10,
			stopAfter: 1,
			want:      []int64{10},
		},
		{
			name:      "find_first_two_ascending_lt",
			method:    ascendLT,
			pivot:     20,
			stopAfter: 2,
			want:      []int64{5, 10},
		},
		{
			name:      "find_first_two_descending_gt",
			method:    descendGT,
			pivot:     5,
			stopAfter: 2,
			want:      []int64{20, 15},
		},
		{
			name:      "find_first_match_descending_lte",
			method:    descendLTE,
			pivot:     15,
			stopAfter: 1,
			want:      []int64{15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var collected []int64
			callCount := 0

			rb.iterate(tt.method, tt.pivot, func(rev int64, events []testEvent) bool {
				collected = append(collected, rev)
				callCount++
				return callCount < tt.stopAfter
			})

			if !slices.Equal(tt.want, collected) {
				t.Fatalf("collected=%v, want %v", collected, tt.want)
			}
			if callCount != tt.stopAfter {
				t.Fatalf("callback calls=%d, want %d", callCount, tt.stopAfter)
			}
		})
	}
}

func TestAtomicOrdered(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		inputs   []struct {
			rev  int64
			key  string
			size int
		}
		wantRev  []int64
		wantSize []int
	}{
		{
			name:     "unfiltered",
			capacity: 5,
			inputs: []struct {
				rev  int64
				key  string
				size int
			}{
				{5, "a", 1},
				{10, "b", 3},
				{15, "c", 7},
				{20, "d", 11},
			},
			wantRev:  []int64{5, 10, 15, 20},
			wantSize: []int{1, 3, 7, 11},
		},
		{
			name:     "across_wrap",
			capacity: 3,
			inputs: []struct {
				rev  int64
				key  string
				size int
			}{
				{1, "a", 2},
				{2, "b", 1},
				{3, "c", 3},
				{4, "d", 7},
			},
			wantRev:  []int64{2, 3, 4},
			wantSize: []int{1, 3, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rb := mustNewTestRingBuffer(t, tt.capacity)
			for _, in := range tt.inputs {
				batch, err := makeEventBatch(in.rev, in.key, in.size)
				if err != nil {
					t.Fatalf("makeEventBatch(%d, %s, %d) failed: %v", in.rev, in.key, in.size, err)
				}
				rb.Append(batch)
			}

			gotRevs := []int64{}
			gotSizes := []int{}
			rb.AscendGreaterOrEqual(0, func(rev int64, events []testEvent) bool {
				gotRevs = append(gotRevs, rev)
				gotSizes = append(gotSizes, len(events))
				return true
			})

			if !slices.Equal(gotRevs, tt.wantRev) {
				t.Fatalf("revisions=%v, want %v", gotRevs, tt.wantRev)
			}
			if !slices.Equal(gotSizes, tt.wantSize) {
				t.Fatalf("sizes=%v, want %v", gotSizes, tt.wantSize)
			}
		})
	}
}

func TestRebaseHistory(t *testing.T) {
	tests := []struct {
		name string
		revs []int64
	}{
		{
			name: "rebase_empty_buffer",
			revs: nil,
		},
		{
			name: "rebase_after_data",
			revs: []int64{7, 8, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rb := mustNewTestRingBuffer(t, 4)
			for _, r := range tt.revs {
				batch, err := makeEventBatch(r, "k", 1)
				if err != nil {
					t.Fatalf("makeEventBatch(%d, k, 1) failed: %v", r, err)
				}
				rb.Append(batch)
			}

			rb.RebaseHistory()

			if got := rb.PeekOldest(); got != 0 {
				t.Fatalf("PeekOldest()=%d, want=0", got)
			}
			if got := rb.PeekLatest(); got != 0 {
				t.Fatalf("PeekLatest()=%d, want=0", got)
			}

			gotRevs := []int64{}
			rb.AscendGreaterOrEqual(0, func(rev int64, events []testEvent) bool {
				gotRevs = append(gotRevs, rev)
				return true
			})
			if len(gotRevs) != 0 {
				t.Fatalf("len(revisions)=%d, want 0", len(gotRevs))
			}
		})
	}
}

func TestFull(t *testing.T) {
	tests := []struct {
		name         string
		capacity     int
		numAppends   int
		expectedFull bool
	}{
		{
			name:         "empty_buffer",
			capacity:     3,
			numAppends:   0,
			expectedFull: false,
		},
		{
			name:         "partially_filled",
			capacity:     5,
			numAppends:   3,
			expectedFull: false,
		},
		{
			name:         "exactly_at_capacity",
			capacity:     3,
			numAppends:   3,
			expectedFull: true,
		},
		{
			name:         "beyond_capacity_wrapping",
			capacity:     3,
			numAppends:   5,
			expectedFull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := mustNewTestRingBuffer(t, tt.capacity)
			for i := 1; i <= tt.numAppends; i++ {
				batch, err := makeEventBatch(int64(i), "k", 1)
				if err != nil {
					t.Fatalf("makeEventBatch(%d, k, 1) failed: %v", i, err)
				}
				rb.Append(batch)
			}

			if got := rb.full(); got != tt.expectedFull {
				t.Fatalf("full()=%t, want=%t", got, tt.expectedFull)
			}
		})
	}
}

func TestLen(t *testing.T) {
	rb := setupRingBuffer(t, 3, []int64{1, 2, 3, 4})

	if got := rb.Len(); got != 3 {
		t.Fatalf("Len()=%d, want=3", got)
	}

	rb.RemoveLess(4)
	if got := rb.Len(); got != 1 {
		t.Fatalf("Len() after RemoveLess=%d, want=1", got)
	}

	rb.RebaseHistory()
	if got := rb.Len(); got != 0 {
		t.Fatalf("Len() after RebaseHistory=%d, want=0", got)
	}
}

func TestReplaceLatest(t *testing.T) {
	t.Run("replace_latest_entry", func(t *testing.T) {
		rb := mustNewTestRingBuffer(t, 4)
		for _, rev := range []int64{5, 10} {
			batch, err := makeEventBatch(rev, "key", 1)
			if err != nil {
				t.Fatalf("makeEventBatch(%d, key, 1) failed: %v", rev, err)
			}
			rb.Append(batch)
		}

		replacement, err := makeEventBatch(12, "replacement", 2)
		if err != nil {
			t.Fatalf("makeEventBatch(12, replacement, 2) failed: %v", err)
		}
		rb.ReplaceLatest(replacement)

		if got := rb.Len(); got != 2 {
			t.Fatalf("Len()=%d, want=2", got)
		}
		if got := rb.PeekLatest(); got != 12 {
			t.Fatalf("PeekLatest()=%d, want=12", got)
		}

		gotRevs := collectRevisions(rb, ascendGTE, 0)
		if !slices.Equal([]int64{5, 12}, gotRevs) {
			t.Fatalf("revisions=%v, want [5 12]", gotRevs)
		}

		var gotBatch [][]testEvent
		rb.AscendGreaterOrEqual(0, func(rev int64, events []testEvent) bool {
			gotBatch = append(gotBatch, events)
			return true
		})
		if len(gotBatch[1]) != 2 {
			t.Fatalf("replacement batch size=%d, want 2", len(gotBatch[1]))
		}
	})

	t.Run("panic_on_empty_buffer", func(t *testing.T) {
		rb := mustNewTestRingBuffer(t, 2)
		replacement, err := makeEventBatch(1, "replacement", 1)
		if err != nil {
			t.Fatalf("makeEventBatch(1, replacement, 1) failed: %v", err)
		}

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("ReplaceLatest() did not panic on empty buffer")
			}
		}()
		rb.ReplaceLatest(replacement)
	})
}

func TestRemoveLess(t *testing.T) {
	tests := []struct {
		name         string
		capacity     int
		setupRevs    []int64
		removeBefore int64
		wantRevs     []int64
		wantLen      int
		wantOldest   int64
		wantLatest   int64
	}{
		{
			name:         "empty_buffer",
			capacity:     4,
			removeBefore: 10,
			wantRevs:     []int64{},
			wantLen:      0,
			wantOldest:   0,
			wantLatest:   0,
		},
		{
			name:         "no_entries_removed",
			capacity:     4,
			setupRevs:    []int64{5, 10, 15},
			removeBefore: 5,
			wantRevs:     []int64{5, 10, 15},
			wantLen:      3,
			wantOldest:   5,
			wantLatest:   15,
		},
		{
			name:         "remove_prefix",
			capacity:     5,
			setupRevs:    []int64{5, 10, 15, 20},
			removeBefore: 15,
			wantRevs:     []int64{15, 20},
			wantLen:      2,
			wantOldest:   15,
			wantLatest:   20,
		},
		{
			name:         "remove_all_entries",
			capacity:     3,
			setupRevs:    []int64{10, 11, 12},
			removeBefore: 20,
			wantRevs:     []int64{},
			wantLen:      0,
			wantOldest:   0,
			wantLatest:   0,
		},
		{
			name:         "wrapped_buffer",
			capacity:     3,
			setupRevs:    []int64{20, 21, 22, 23, 24},
			removeBefore: 24,
			wantRevs:     []int64{24},
			wantLen:      1,
			wantOldest:   24,
			wantLatest:   24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := setupRingBuffer(t, tt.capacity, tt.setupRevs)

			rb.RemoveLess(tt.removeBefore)

			if got := rb.Len(); got != tt.wantLen {
				t.Fatalf("Len()=%d, want=%d", got, tt.wantLen)
			}
			if got := rb.PeekOldest(); got != tt.wantOldest {
				t.Fatalf("PeekOldest()=%d, want=%d", got, tt.wantOldest)
			}
			if got := rb.PeekLatest(); got != tt.wantLatest {
				t.Fatalf("PeekLatest()=%d, want=%d", got, tt.wantLatest)
			}

			gotRevs := collectRevisions(rb, ascendGTE, 0)
			if !slices.Equal(tt.wantRevs, gotRevs) {
				t.Fatalf("remaining revisions=%v, want %v", gotRevs, tt.wantRevs)
			}
		})
	}
}

func TestEnsureCapacity(t *testing.T) {
	t.Run("small_growth_rounds_up_to_default_ring_size", func(t *testing.T) {
		rb := setupRingBuffer(t, 1, []int64{1})

		rb.EnsureCapacity(2)

		if got := len(rb.buffer); got != defaultRingSize {
			t.Fatalf("len(buffer)=%d, want=%d", got, defaultRingSize)
		}
		if rb.head != rb.size {
			t.Fatalf("head=%d, want size=%d", rb.head, rb.size)
		}
		if rb.tail != 0 {
			t.Fatalf("tail=%d, want=0", rb.tail)
		}
		gotRevs := collectRevisions(rb, ascendGTE, 0)
		if !slices.Equal([]int64{1}, gotRevs) {
			t.Fatalf("revisions after EnsureCapacity=%v, want [1]", gotRevs)
		}
	})

	t.Run("grow_preserves_logical_order", func(t *testing.T) {
		rb := setupRingBuffer(t, 3, []int64{1, 2, 3, 4})

		rb.EnsureCapacity(5)

		if got := len(rb.buffer); got != defaultRingSize {
			t.Fatalf("len(buffer)=%d, want=%d", got, defaultRingSize)
		}
		if rb.tail != 0 {
			t.Fatalf("tail=%d, want=0", rb.tail)
		}
		if rb.head != rb.size {
			t.Fatalf("head=%d, want size=%d", rb.head, rb.size)
		}

		gotRevs := collectRevisions(rb, ascendGTE, 0)
		if !slices.Equal([]int64{2, 3, 4}, gotRevs) {
			t.Fatalf("revisions after EnsureCapacity=%v, want [2 3 4]", gotRevs)
		}
	})

	t.Run("noop_when_capacity_is_sufficient", func(t *testing.T) {
		rb := setupRingBuffer(t, 4, []int64{1, 2})
		beforeCap := len(rb.buffer)
		beforeHead, beforeTail, beforeSize := rb.head, rb.tail, rb.size

		rb.EnsureCapacity(4)

		if got := len(rb.buffer); got != beforeCap {
			t.Fatalf("len(buffer)=%d, want=%d", got, beforeCap)
		}
		if rb.head != beforeHead || rb.tail != beforeTail || rb.size != beforeSize {
			t.Fatalf("metadata changed unexpectedly: got head=%d tail=%d size=%d, want head=%d tail=%d size=%d",
				rb.head, rb.tail, rb.size, beforeHead, beforeTail, beforeSize)
		}
	})
}

func mustNewTestRingBuffer(t *testing.T, capacity int) *RevisionRingBuffer[[]testEvent] {
	t.Helper()

	rb, err := NewRevisionRingBuffer(capacity, func(batch []testEvent) int64 {
		return batch[0].ModRevision
	})
	if err != nil {
		t.Fatalf("NewRevisionRingBuffer(%d) failed: %v", capacity, err)
	}
	return rb
}

func (r *RevisionRingBuffer[T]) iterate(method iterMethod, pivot int64, fn IterFunc[T]) {
	switch method {
	case ascendGTE:
		r.AscendGreaterOrEqual(pivot, fn)
	case ascendLT:
		r.AscendLessThan(pivot, fn)
	case descendGT:
		r.DescendGreaterThan(pivot, fn)
	case descendLTE:
		r.DescendLessOrEqual(pivot, fn)
	default:
		panic(fmt.Sprintf("unknown iteration method: %s", method))
	}
}

func setupRingBuffer(t *testing.T, capacity int, revs []int64) *RevisionRingBuffer[[]testEvent] {
	t.Helper()

	rb := mustNewTestRingBuffer(t, capacity)
	for _, r := range revs {
		batch, err := makeEventBatch(r, "key", 1)
		if err != nil {
			t.Fatalf("makeEventBatch(%d, key, 1) failed: %v", r, err)
		}
		rb.Append(batch)
	}
	return rb
}

func collectRevisions(rb *RevisionRingBuffer[[]testEvent], method iterMethod, pivot int64) []int64 {
	revs := []int64{}
	rb.iterate(method, pivot, func(rev int64, events []testEvent) bool {
		revs = append(revs, rev)
		return true
	})
	return revs
}

func makeEventBatch(rev int64, key string, batchSize int) ([]testEvent, error) {
	if batchSize < 0 {
		return nil, fmt.Errorf("invalid batchSize %d", batchSize)
	}
	events := make([]testEvent, batchSize)
	for i := range events {
		events[i] = testEvent{
			Key:         fmt.Sprintf("%s-%d", key, i),
			ModRevision: rev,
		}
	}
	return events, nil
}
