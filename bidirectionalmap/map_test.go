/*
Copyright 2020 The Kubernetes Authors.

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

import "testing"

func TestMultipleInserts(t *testing.T) {
	bidimap := NewBidirectionalMap[string, string]()
	bidimap.InsertRight("r1", "l1")
	bidimap.InsertRight("r1", "l2")
	if bidimap.GetRight("r1").Len() != 2 {
		t.Errorf("GetRight('r1').Len() == %d, expected 2", bidimap.GetRight("r1").Len())
	}
	if bidimap.GetLeft("l2").Len() != 1 {
		t.Errorf("GetLeft('l2').Len() == %d, expected 1", bidimap.GetLeft("l2").Len())
	}
	bidimap.InsertLeft("l2", "r2")
	if bidimap.GetLeft("l2").Len() != 2 {
		t.Errorf("GetLeft('l2').Len() == %d, expected 2", bidimap.GetLeft("l2").Len())
	}
	r2Len := bidimap.GetRight("r2").Len()
	if r2Len != 1 {
		t.Errorf("GetRight('r2').Len() == %d, expected 1", r2Len)
	}
	bidimap.DeleteRightKey("r2")
	if bidimap.GetRight("r2") != nil {
		t.Errorf("GetRight('r2') should be nil")
	}
	if bidimap.GetLeft("l2").Len() != 1 {
		t.Errorf("GetLeft('l2').Len() == %d, expected 1", bidimap.GetLeft("l2").Len())
	}
}
