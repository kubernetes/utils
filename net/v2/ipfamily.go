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
	"net"
)

// IPFamily refers to the IP family of an address or CIDR value. Its values are
// intentionally identical to those of "k8s.io/api/core/v1".IPFamily and
// "k8s.io/discovery/v1".AddressType, so you can cast values between these types.
type IPFamily string

const (
	// IPv4 indicates an IPv4 IP or CIDR.
	IPv4 IPFamily = "IPv4"
	// IPv6 indicates an IPv4 IP or CIDR.
	IPv6 IPFamily = "IPv6"

	// IPFamilyUnknown indicates an unspecified or invalid IP family.
	IPFamilyUnknown IPFamily = ""
)

// IsDualStackIPs returns true if:
// - all elements of ips are valid
// - at least one IP from each family (v4 and v6) is present
func IsDualStackIPs(ips []net.IP) bool {
	v4Found := false
	v6Found := false
	for _, ip := range ips {
		switch IPFamilyOf(ip) {
		case IPv4:
			v4Found = true
		case IPv6:
			v6Found = true
		default:
			return false
		}
	}

	return (v4Found && v6Found)
}

// IsDualStackIPStrings returns true if:
// - all elements of ips can be parsed as IPs
// - at least one IP from each family (v4 and v6) is present
func IsDualStackIPStrings(ips []string) bool {
	parsedIPs := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		parsedIP := ParseIPSloppy(ip)
		if parsedIP == nil {
			return false
		}
		parsedIPs = append(parsedIPs, parsedIP)
	}
	return IsDualStackIPs(parsedIPs)
}

// IsDualStackCIDRs returns true if:
// - all elements of cidrs are non-nil
// - at least one CIDR from each family (v4 and v6) is present
func IsDualStackCIDRs(cidrs []*net.IPNet) bool {
	v4Found := false
	v6Found := false
	for _, cidr := range cidrs {
		switch IPFamilyOfCIDR(cidr) {
		case IPv4:
			v4Found = true
		case IPv6:
			v6Found = true
		default:
			return false
		}
	}

	return (v4Found && v6Found)
}

// IsDualStackCIDRStrings returns if
// - all elements of cidrs can be parsed as CIDRs
// - at least one CIDR from each family (v4 and v6) is present
func IsDualStackCIDRStrings(cidrs []string) bool {
	parsedCIDRs, err := ParseCIDRs(cidrs)
	if err != nil {
		return false
	}
	return IsDualStackCIDRs(parsedCIDRs)
}

// IPFamilyOf returns the IP family of ip, or IPFamilyUnknown if it is invalid.
func IPFamilyOf(ip net.IP) IPFamily {
	switch {
	case ip.To4() != nil:
		return IPv4
	case ip.To16() != nil:
		return IPv6
	default:
		return IPFamilyUnknown
	}
}

// IPFamilyOfString returns the IP family of ip, or IPFamilyUnknown if ip cannot
// be parsed as an IP.
func IPFamilyOfString(ip string) IPFamily {
	return IPFamilyOf(ParseIPSloppy(ip))
}

// IPFamilyOfCIDR returns the IP family of cidr.
func IPFamilyOfCIDR(cidr *net.IPNet) IPFamily {
	if cidr != nil {
		family := IPFamilyOf(cidr.IP)
		// An IPv6 CIDR must have a 128-bit mask. An IPv4 CIDR must have a
		// 32- or 128-bit mask. (Any other mask length is invalid.)
		_, masklen := cidr.Mask.Size()
		if masklen == 128 || (family == IPv4 && masklen == 32) {
			return family
		}
	}
	return IPFamilyUnknown
}

// IPFamilyOfCIDRString returns the IP family of cidr.
func IPFamilyOfCIDRString(cidr string) IPFamily {
	ip, _, _ := ParseCIDRSloppy(cidr)
	return IPFamilyOf(ip)
}

// IsIPv6 returns true if netIP is IPv6 (and false if it is IPv4, nil, or invalid).
func IsIPv6(netIP net.IP) bool {
	return IPFamilyOf(netIP) == IPv6
}

// IsIPv6String returns true if ip contains a single IPv6 address and nothing else. It
// returns false if ip is an empty string, an IPv4 address, or anything else that is not a
// single IPv6 address.
func IsIPv6String(ip string) bool {
	return IPFamilyOfString(ip) == IPv6
}

// IsIPv6CIDR returns true if a cidr is a valid IPv6 CIDR. It returns false if cidr is
// nil or an IPv4 CIDR. Its behavior is not defined if cidr is invalid.
func IsIPv6CIDR(cidr *net.IPNet) bool {
	return IPFamilyOfCIDR(cidr) == IPv6
}

// IsIPv6CIDRString returns true if cidr contains a single IPv6 CIDR and nothing else. It
// returns false if cidr is an empty string, an IPv4 CIDR, or anything else that is not a
// single valid IPv6 CIDR.
func IsIPv6CIDRString(cidr string) bool {
	return IPFamilyOfCIDRString(cidr) == IPv6
}

// IsIPv4 returns true if netIP is IPv4 (and false if it is IPv6, nil, or invalid).
func IsIPv4(netIP net.IP) bool {
	return IPFamilyOf(netIP) == IPv4
}

// IsIPv4String returns true if ip contains a single IPv4 address and nothing else. It
// returns false if ip is an empty string, an IPv6 address, or anything else that is not a
// single IPv4 address.
func IsIPv4String(ip string) bool {
	return IPFamilyOfString(ip) == IPv4
}

// IsIPv4CIDR returns true if cidr is a valid IPv4 CIDR. It returns false if cidr is nil
// or an IPv6 CIDR. Its behavior is not defined if cidr is invalid.
func IsIPv4CIDR(cidr *net.IPNet) bool {
	return IPFamilyOfCIDR(cidr) == IPv4
}

// IsIPv4CIDRString returns true if cidr contains a single IPv4 CIDR and nothing else. It
// returns false if cidr is an empty string, an IPv6 CIDR, or anything else that is not a
// single valid IPv4 CIDR.
func IsIPv4CIDRString(cidr string) bool {
	return IPFamilyOfCIDRString(cidr) == IPv4
}
