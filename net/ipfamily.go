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

// IPFamilyOf returns the IP family of netIP.
func IPFamilyOf(netIP net.IP) IPFamily {
	switch {
	case netIP == nil:
		return IPFamilyUnknown
	case netIP.To4() != nil:
		return IPv4
	default:
		return IPv6
	}
}

// IPFamilyOfString returns the IP family of ip.
func IPFamilyOfString(ip string) IPFamily {
	return IPFamilyOf(ParseIPSloppy(ip))
}

// IPFamilyOfCIDR returns the IP family of cidr.
func IPFamilyOfCIDR(cidr *net.IPNet) IPFamily {
	if cidr == nil {
		return IPFamilyUnknown
	}
	return IPFamilyOf(cidr.IP)
}

// IPFamilyOfCIDRString returns the IP family of cidr.
func IPFamilyOfCIDRString(cidr string) IPFamily {
	ip, _, _ := ParseCIDRSloppy(cidr)
	return IPFamilyOf(ip)
}

// IsIPv6 returns whether netIP is IPv6.
func IsIPv6(netIP net.IP) bool {
	return IPFamilyOf(netIP) == IPv6
}

// IsIPv6String returns whether ip is IPv6.
func IsIPv6String(ip string) bool {
	return IPFamilyOfString(ip) == IPv6
}

// IsIPv6CIDR returns whether cidr is ipv6
func IsIPv6CIDR(cidr *net.IPNet) bool {
	return IPFamilyOfCIDR(cidr) == IPv6
}

// IsIPv6CIDRString returns whether cidr is IPv6.
func IsIPv6CIDRString(cidr string) bool {
	return IPFamilyOfCIDRString(cidr) == IPv6
}

// IsIPv4 returns whether netIP is IPv4.
func IsIPv4(netIP net.IP) bool {
	return IPFamilyOf(netIP) == IPv4
}

// IsIPv4String returns whether ip is IPv4.
func IsIPv4String(ip string) bool {
	return IPFamilyOfString(ip) == IPv4
}

// IsIPv4CIDR returns whether cidr is ipv4
func IsIPv4CIDR(cidr *net.IPNet) bool {
	return IPFamilyOfCIDR(cidr) == IPv4
}

// IsIPv4CIDRString returns whether cidr is IPv4.
func IsIPv4CIDRString(cidr string) bool {
	return IPFamilyOfCIDRString(cidr) == IPv4
}
