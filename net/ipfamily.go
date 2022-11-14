/*
Copyright 2018 The Kubernetes Authors.

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
)

// IPFamily refers to a specific family if not empty, i.e. "4" or "6".
type IPFamily string

// Constants for valid IPFamilys:
const (
	IPFamilyUnknown IPFamily = ""

	IPv4 IPFamily = "4"
	IPv6 IPFamily = "6"
)

// IsDualStackIPs returns if a slice of ips is:
// - all are valid ips
// - at least one ip from each family (v4 or v6)
func IsDualStackIPs(ips []net.IP) (bool, error) {
	v4Found := false
	v6Found := false
	for _, ip := range ips {
		if ip == nil {
			return false, fmt.Errorf("ip %v is invalid", ip)
		}

		if v4Found && v6Found {
			continue
		}

		if IsIPv6(ip) {
			v6Found = true
			continue
		}

		v4Found = true
	}

	return (v4Found && v6Found), nil
}

// IsDualStackIPStrings returns if
// - all are valid ips
// - at least one ip from each family (v4 or v6)
func IsDualStackIPStrings(ips []string) (bool, error) {
	parsedIPs := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		parsedIP := ParseIPSloppy(ip)
		parsedIPs = append(parsedIPs, parsedIP)
	}
	return IsDualStackIPs(parsedIPs)
}

// IsDualStackCIDRs returns if
// - all are valid cidrs
// - at least one cidr from each family (v4 or v6)
func IsDualStackCIDRs(cidrs []*net.IPNet) (bool, error) {
	v4Found := false
	v6Found := false
	for _, cidr := range cidrs {
		if cidr == nil {
			return false, fmt.Errorf("cidr %v is invalid", cidr)
		}

		if v4Found && v6Found {
			continue
		}

		if IsIPv6(cidr.IP) {
			v6Found = true
			continue
		}
		v4Found = true
	}

	return v4Found && v6Found, nil
}

// IsDualStackCIDRStrings returns if
// - all are valid cidrs
// - at least one cidr from each family (v4 or v6)
func IsDualStackCIDRStrings(cidrs []string) (bool, error) {
	parsedCIDRs, err := ParseCIDRs(cidrs)
	if err != nil {
		return false, err
	}
	return IsDualStackCIDRs(parsedCIDRs)
}

// IsIPv6 returns true if ip is IPv6, or false if it is IPv4, nil, or invalid.
func IsIPv6(ip net.IP) bool {
	return ip.To16() != nil && ip.To4() == nil
}

// IsIPv6String returns true if ip contains a single IPv6 address and nothing else. It
// returns false if ip is an empty string, an IPv4 address, or anything else that is not a
// single IPv6 address.
func IsIPv6String(ip string) bool {
	return IsIPv6(ParseIPSloppy(ip))
}

// IsIPv6CIDRString returns true if cidr contains a single IPv6 CIDR and nothing else. It
// returns false if cidr is an empty string, an IPv4 CIDR, or anything else that is not a
// single valid IPv6 CIDR.
func IsIPv6CIDRString(cidr string) bool {
	ip, _, _ := ParseCIDRSloppy(cidr)
	return IsIPv6(ip)
}

// IsIPv6CIDR returns true if a cidr is a valid IPv6 CIDR. It returns false if cidr is
// nil or an IPv4 CIDR. Its behavior is not defined if cidr is invalid.
func IsIPv6CIDR(cidr *net.IPNet) bool {
	if cidr == nil {
		return false
	}
	ip := cidr.IP
	return IsIPv6(ip)
}

// IsIPv4 returns true if ip is IPv4, or false if it is IPv6, nil, or invalid.
func IsIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

// IsIPv4String returns true if ip contains a single IPv4 address and nothing else. It
// returns false if ip is an empty string, an IPv6 address, or anything else that is not a
// single IPv4 address.
func IsIPv4String(ip string) bool {
	return IsIPv4(ParseIPSloppy(ip))
}

// IsIPv4CIDR returns true if cidr is a valid IPv4 CIDR. It returns false if cidr is nil
// or an IPv6 CIDR. Its behavior is not defined if cidr is invalid.
func IsIPv4CIDR(cidr *net.IPNet) bool {
	if cidr == nil {
		return false
	}
	ip := cidr.IP
	return IsIPv4(ip)
}

// IsIPv4CIDRString returns true if cidr contains a single IPv4 CIDR and nothing else. It
// returns false if cidr is an empty string, an IPv6 CIDR, or anything else that is not a
// single valid IPv4 CIDR.
func IsIPv4CIDRString(cidr string) bool {
	ip, _, _ := ParseCIDRSloppy(cidr)
	return IsIPv4(ip)
}
