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

import "sync/atomic"

// RingGrowing is a growing ring buffer.
// Not thread safe.
type RingGrowing struct {
	data     []interface{}
	n        int32 // Size of Data
	beg      int32 // First available element
	readable int32 // Number of data items available
}

// NewRingGrowing constructs a new RingGrowing instance with provided parameters.
func NewRingGrowing(initialSize int32) *RingGrowing {
	return &RingGrowing{
		data: make([]interface{}, initialSize),
		n:    initialSize,
	}
}

// ReadOne reads (consumes) first item from the buffer if it is available, otherwise returns false.
func (r *RingGrowing) ReadOne() (data interface{}, ok bool) {
	oldReadable := atomic.LoadInt32(&r.readable)
	oldN := atomic.LoadInt32(&r.n)
	if oldReadable == 0 {
		return nil, false
	}
	if !atomic.CompareAndSwapInt32(&r.readable, oldReadable, oldReadable-1) {
		return nil, false
	}
	oldBeg := atomic.LoadInt32(&r.beg)

	if oldBeg == oldN-1 {
		if !atomic.CompareAndSwapInt32(&r.beg, oldBeg, 0) {
			return nil, false
		}
	} else {
		if !atomic.CompareAndSwapInt32(&r.beg, oldBeg, oldBeg+1) {
			return nil, false
		}
	}
	element := r.data[oldBeg]
	r.data[oldBeg] = nil // Remove reference to the object to help GC
	return element, true
}

// WriteOne adds an item to the end of the buffer, growing it if it is full.
func (r *RingGrowing) WriteOne(data interface{}) {
	oldReadable := atomic.LoadInt32(&r.readable)
	oldN := atomic.LoadInt32(&r.n)
	oldBeg := atomic.LoadInt32(&r.beg)

	if oldN == oldReadable {
		// Time to grow
		newN := oldN * 2
		newData := make([]interface{}, newN)
		to := oldBeg + oldReadable
		if to <= r.n {
			copy(newData, r.data[r.beg:to])
		} else {
			copied := copy(newData, r.data[r.beg:])
			copy(newData[copied:], r.data[:(to%r.n)])
		}
		atomic.StoreInt32(&r.n, newN)
		r.data = newData
		atomic.StoreInt32(&r.beg, 0)
	}
	r.data[(r.readable+r.beg)%r.n] = data
	atomic.StoreInt32(&r.readable, oldReadable+1)
}
