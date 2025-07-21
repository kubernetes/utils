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

package net

import (
	"fmt"
	"net"

	forkednet "k8s.io/utils/internal/third_party/forked/golang/net"
)

// ParseIP parses an IPv4 or IPv6 address to a net.IP. This accepts both fully-valid IP
// addresses and irregular/ambiguous forms, making it usable for both validated and
// non-validated input strings. It should be used instead of net.ParseIP (which rejects
// some strings we need to accept for backward compatibility) and the old
// netutilsv1.ParseIPSloppy.
func ParseIP(ipStr string) (net.IP, error) {
	// Note: if we want to get rid of forkednet, we should be able to use some
	// invocation of regexp.ReplaceAllString to get rid of leading 0s in ipStr.
	ip := forkednet.ParseIP(ipStr)
	if ip != nil {
		return ip, nil
	}

	if ipStr == "" {
		return nil, fmt.Errorf("expected an IP address")
	}
	// NB: we use forkednet.ParseCIDR directly, not ParseIPNet, to avoid recursing
	// between ParseIP and ParseIPNet.
	if _, _, err := forkednet.ParseCIDR(ipStr); err == nil {
		return nil, fmt.Errorf("expected an IP address, got a CIDR value")
	}
	return nil, fmt.Errorf("not a valid IP address")
}

// ParseIPNet parses an IPv4 or IPv6 CIDR string representing a subnet or mask, to a
// *net.IPNet. This accepts both fully-valid CIDR values and irregular/ambiguous forms,
// making it usable for both validated and non-validated input strings. It should be used
// instead of net.ParseCIDR (which rejects some strings that we need to accept for
// backward-compatibility) and the old netutilsv1.ParseCIDRSloppy.
//
// The return value is equivalent to the second return value from net.ParseCIDR. Note that
// this means that if the CIDR string has bits set beyond the prefix length (e.g., the "5"
// in "192.168.1.5/24"), those bits are simply discarded.
func ParseIPNet(cidrStr string) (*net.IPNet, error) {
	// Note: if we want to get rid of forkednet, we should be able to use some
	// invocation of regexp.ReplaceAllString to get rid of leading 0s in cidrStr.
	if _, ipnet, err := forkednet.ParseCIDR(cidrStr); err == nil {
		return ipnet, nil
	}

	if cidrStr == "" {
		return nil, fmt.Errorf("expected a CIDR value")
	}
	// NB: we use forkednet.ParseIP directly, not our own ParseIP, to avoid recursing
	// between ParseIPNet and ParseIP.
	if forkednet.ParseIP(cidrStr) != nil {
		return nil, fmt.Errorf("expected a CIDR value, but got IP address")
	}
	return nil, fmt.Errorf("not a valid CIDR value")
}
