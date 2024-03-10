/*
Copyright 2024 The Kubernetes Authors.

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

package bidirectionalmap

import (
	"k8s.io/utils/genericinterfaces"
	"k8s.io/utils/set"
)

// BidirectionalMap is a bidirectional map.
type BidirectionalMap[X genericinterfaces.Ordered, Y genericinterfaces.Ordered] struct {
	right map[X]set.Set[Y]
	left  map[Y]set.Set[X]
}

// NewBidirectionalMap creates a new BidirectionalMap.
func NewBidirectionalMap[X genericinterfaces.Ordered, Y genericinterfaces.Ordered]() *BidirectionalMap[X, Y] {
	return &BidirectionalMap[X, Y]{
		right: make(map[X]set.Set[Y]),
		left:  make(map[Y]set.Set[X]),
	}
}

// InsertRight inserts a new item into the right map, return true if the key-value was not already
// present in the map, false otherwise
func (bdm *BidirectionalMap[X, Y]) InsertRight(x X, y Y) bool {
	if bdm.right[x] == nil {
		bdm.right[x] = set.New[Y]()
	}
	if bdm.right[x].Has(y) {
		return false
	}
	if bdm.left[y] == nil {
		bdm.left[y] = set.New[X]()
	}
	bdm.right[x].Insert(y)
	bdm.left[y].Insert(x)
	return true
}

// InsertLeft inserts a new item into the left map, return true if the key-value was not already
// present in the map, false otherwise
func (bdm *BidirectionalMap[X, Y]) InsertLeft(y Y, x X) bool {
	return bdm.InsertRight(x, y)
}

// GetRight returns a value from the right map.
func (bdm *BidirectionalMap[X, Y]) GetRight(x X) set.Set[Y] {
	return bdm.right[x]
}

// GetLeft returns a value from left map.
func (bdm *BidirectionalMap[X, Y]) GetLeft(y Y) set.Set[X] {
	return bdm.left[y]
}

// DeleteRightKey deletes the key from the right map and removes
// the inverse mapping from the left map.
func (bdm *BidirectionalMap[X, Y]) DeleteRightKey(x X) {
	if leftValues, ok := bdm.right[x]; ok {
		delete(bdm.right, x)
		for y := range leftValues {
			bdm.left[y].Delete(x)
			if bdm.left[y].Len() == 0 {
				delete(bdm.left, y)
			}
		}
	}
}

// DeleteLeftKey deletes the key from the left map and removes
// the inverse mapping from the right map.
func (bdm *BidirectionalMap[X, Y]) DeleteLeftKey(y Y) {
	if rightValues, ok := bdm.left[y]; ok {
		delete(bdm.left, y)
		for x := range rightValues {
			bdm.right[x].Delete(y)
			if bdm.right[x].Len() == 0 {
				delete(bdm.right, x)
			}
		}
	}
}

// GetRightKeys returns the keys from the right map.
func (bdm *BidirectionalMap[X, Y]) GetRightKeys() set.Set[X] {
	return set.KeySet[X](bdm.right)
}

// GetLeftKeys returns the keys from the left map.
func (bdm *BidirectionalMap[X, Y]) GetLeftKeys() set.Set[Y] {
	return set.KeySet[Y](bdm.left)
}
