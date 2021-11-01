/*
Copyright 2021 The Kubernetes Authors.

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

package dump

import (
	"github.com/davecgh/go-spew/spew"
)

var prettyPrintConfig = &spew.ConfigState{
	Indent:                  "  ",
	DisableMethods:          true,
	DisablePointerAddresses: true,
	DisableCapacities:       true,
}

// PrettyPrintObject wrap the spew.Sdump with Indent, and disabled methods like error() and String()
// The output may change over time, so for guaranteed output please take more direct control
func PrettyPrintObject(a interface{}) string {
	return prettyPrintConfig.Sdump(a)
}
