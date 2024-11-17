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

type ipOrString interface {
	net.IP | string
}

type cidrOrString interface {
	*net.IPNet | string
}

// IPFamilyOf returns the IP family of val (or IPFamilyUnknown if val is nil or invalid).
// IPv6-encoded IPv4 addresses (e.g., "::ffff:1.2.3.4") are considered IPv4. val can be a
// net.IP or a string containing a single IP address.
//
// Note that "k8s.io/utils/net/v2".IPFamily intentionally has identical values to
// "k8s.io/api/core/v1".IPFamily and "k8s.io/discovery/v1".AddressType, so you can cast
// the return value of this function to those types.
func IPFamilyOf[T ipOrString](val T) IPFamily {
	switch typedVal := interface{}(val).(type) {
	case net.IP:
		switch {
		case typedVal.To4() != nil:
			return IPv4
		case typedVal.To16() != nil:
			return IPv6
		}
	case string:
		return IPFamilyOf(ParseIPSloppy(typedVal))
	}

	return IPFamilyUnknown
}

// IsIPv4 returns true if IPFamilyOf(val) is IPv4 (and false if it is IPv6 or invalid).
func IsIPv4[T ipOrString](val T) bool {
	return IPFamilyOf(val) == IPv4
}

// IsIPv6 returns true if IPFamilyOf(val) is IPv6 (and false if it is IPv4 or invalid).
func IsIPv6[T ipOrString](val T) bool {
	return IPFamilyOf(val) == IPv6
}

// IsDualStack returns true if vals contains at least one IPv4 address and at least one
// IPv6 address (and no invalid values).
func IsDualStack[T ipOrString](vals []T) bool {
	v4Found := false
	v6Found := false
	for _, val := range vals {
		switch IPFamilyOf(val) {
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

// IsDualStackPair returns true if vals contains exactly 1 IPv4 address and 1 IPv6 address
// (in either order).
func IsDualStackPair[T ipOrString](vals []T) bool {
	return len(vals) == 2 && IsDualStack(vals)
}

// IPFamilyOfCIDR returns the IP family of val (or IPFamilyUnknown if val is nil or
// invalid). IPv6-encoded IPv4 addresses (e.g., "::ffff:1.2.3.0/120") are considered IPv4.
// val can be a *net.IPNet or a string containing a single CIDR value.
//
// Note that "k8s.io/utils/net/v2".IPFamily intentionally has identical values to
// "k8s.io/api/core/v1".IPFamily and "k8s.io/discovery/v1".AddressType, so you can cast
// the return value of this function to those types.
func IPFamilyOfCIDR[T cidrOrString](val T) IPFamily {
	switch typedVal := interface{}(val).(type) {
	case *net.IPNet:
		if typedVal != nil {
			family := IPFamilyOf(typedVal.IP)
			// An IPv6 CIDR must have a 128-bit mask. An IPv4 CIDR must have a
			// 32- or 128-bit mask. (Any other mask length is invalid.)
			_, masklen := typedVal.Mask.Size()
			if masklen == 128 || (family == IPv4 && masklen == 32) {
				return family
			}
		}
	case string:
		parsedIP, _, _ := ParseCIDRSloppy(typedVal)
		return IPFamilyOf(parsedIP)
	}

	return IPFamilyUnknown
}

// IsIPv4CIDR returns true if IPFamilyOfCIDR(val) is IPv4 (and false if it is IPv6 or invalid).
func IsIPv4CIDR[T cidrOrString](val T) bool {
	return IPFamilyOfCIDR(val) == IPv4
}

// IsIPv6CIDR returns true if IPFamilyOfCIDR(val) is IPv6 (and false if it is IPv4 or invalid).
func IsIPv6CIDR[T cidrOrString](val T) bool {
	return IPFamilyOfCIDR(val) == IPv6
}

// IsDualStackCIDRs returns true if vals contains at least one IPv4 CIDR value and at
// least one IPv6 CIDR value (and no invalid values).
func IsDualStackCIDRs[T cidrOrString](vals []T) bool {
	v4Found := false
	v6Found := false
	for _, val := range vals {
		switch IPFamilyOfCIDR(val) {
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

// IsDualStackCIDRPair returns true if vals contains exactly 1 IPv4 CIDR value and 1 IPv6
// CIDR value (in either order).
func IsDualStackCIDRPair[T cidrOrString](vals []T) bool {
	return len(vals) == 2 && IsDualStackCIDRs(vals)
}

// OtherIPFamily returns the other IP family from ipFamily.
//
// Note that "k8s.io/utils/net/v2".IPFamily intentionally has identical values to
// "k8s.io/api/core/v1".IPFamily and "k8s.io/discovery/v1".AddressType, so you can cast
// the input/output values of this function between these types.
func OtherIPFamily(ipFamily IPFamily) IPFamily {
	switch ipFamily {
	case IPv4:
		return IPv6
	case IPv6:
		return IPv4
	default:
		return IPFamilyUnknown
	}
}
