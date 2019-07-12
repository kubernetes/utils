/*
Copyright 2019 The Kubernetes Authors.

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

package ipvs

import (
	"fmt"
	"strings"

	"k8s.io/utils/exec"
	"k8s.io/utils/version"
)

// IPVS required kernel modules.
const (
	// ModIPVS is the kernel module "ip_vs"
	ModIPVS string = "ip_vs"
	// ModIPVSRR is the kernel module "ip_vs_rr"
	ModIPVSRR string = "ip_vs_rr"
	// ModIPVSWRR is the kernel module "ip_vs_wrr"
	ModIPVSWRR string = "ip_vs_wrr"
	// ModIPVSSH is the kernel module "ip_vs_sh"
	ModIPVSSH string = "ip_vs_sh"
	// ModNfConntrackIPV4 is the module "nf_conntrack_ipv4"
	ModNfConntrackIPV4 string = "nf_conntrack_ipv4"
	// ModNfConntrack is the kernel module "nf_conntrack"
	ModNfConntrack string = "nf_conntrack"
)

// GetKernelVersionAndIPVSMods returns the linux kernel version and the required ipvs modules
func GetKernelVersionAndIPVSMods(Executor exec.Interface) (kernelVersion string, ipvsModules []string, err error) {
	kernelVersionFile := "/proc/sys/kernel/osrelease"
	out, err := Executor.Command("cut", "-f1", "-d", " ", kernelVersionFile).CombinedOutput()
	if err != nil {
		return "", nil, fmt.Errorf("error getting os release kernel version: %v(%s)", err, out)
	}
	kernelVersion = strings.TrimSpace(string(out))
	// parse kernel version
	ver1, err := version.ParseGeneric(kernelVersion)
	if err != nil {
		return kernelVersion, nil, fmt.Errorf("error parsing kernel version: %v(%s)", err, kernelVersion)
	}
	// "nf_conntrack_ipv4" has been removed since v4.19
	// see https://github.com/torvalds/linux/commit/a0ae2562c6c4b2721d9fddba63b7286c13517d9f
	ver2, _ := version.ParseGeneric("4.19")
	// get required ipvs modules
	if ver1.LessThan(ver2) {
		ipvsModules = append(ipvsModules, ModIPVS, ModIPVSRR, ModIPVSWRR, ModIPVSSH, ModNfConntrackIPV4)
	} else {
		ipvsModules = append(ipvsModules, ModIPVS, ModIPVSRR, ModIPVSWRR, ModIPVSSH, ModNfConntrack)
	}

	return kernelVersion, ipvsModules, nil
}
