// +build !linux

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

package procfs

import (
	"fmt"
	"syscall"
)

type ProcFS struct{}

func NewProcFS() ProcInterface {
	return &ProcFS{}
}

// GetFullContainerName gets the container name given the root process id of the container.
func (pfs *ProcFS) GetFullContainerName(pid int) (string, error) {
	return "", fmt.Errorf("GetFullContainerName is unsupported in this build")
}

// PKill finds process(es) using a regular expression
// and sends a specified signal to each process
func PKill(name string, sig syscall.Signal) error {
	return fmt.Errorf("PKill is unsupported in this build")
}

// PidOf finds process(es) with a specified name (regexp match)
// and returns their pid(s)
func PidOf(name string) ([]int, error) {
	return []int{}, fmt.Errorf("PidOf is unsupported in this build")
}
