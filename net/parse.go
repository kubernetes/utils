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

package net

import (
	"net"
	"regexp"
)

var matchZeros = regexp.MustCompile(`(^|:)0*(0|[1-9][0-9]*)\.0*(0|[1-9][0-9]*)\.0*(0|[1-9][0-9]*)\.0*(0|[1-9][0-9]*)((/)0*(0|[1-9][0-9]*))?$`)

const stripZeros = `$1$2.$3.$4.$5$7$8`

// ParseIPSloppy is identical to Go's standard net.ParseIP, except that it allows
// leading '0' characters on numbers.  Go used to allow this and then changed
// the behavior in 1.17.  We're choosing to keep it for compat with potential
// stored values.
func ParseIPSloppy(str string) net.IP {
	// Try calling ParseIP directly, which will probably work.
	if ip := net.ParseIP(str); ip != nil {
		return ip
	}
	// Try stripping 0s and parsing the result.
	str = matchZeros.ReplaceAllString(str, stripZeros)
	return net.ParseIP(str)
}

// ParseCIDRSloppy is identical to Go's standard net.ParseCIDR, except that it allows
// leading '0' characters on numbers.  Go used to allow this and then changed
// the behavior in 1.17.  We're choosing to keep it for compat with potential
// stored values.
func ParseCIDRSloppy(str string) (net.IP, *net.IPNet, error) {
	// Try calling ParseCIDR directly, which will probably work.
	if ip, ipnet, err := net.ParseCIDR(str); err == nil {
		return ip, ipnet, err
	}
	// Try stripping 0s and parsing the result. Note that if str didn't have any
	// leading 0s, then this will just fail again in exactly the same way (with
	// exactly the same error) as the first call.
	str = matchZeros.ReplaceAllString(str, stripZeros)
	return net.ParseCIDR(str)
}
