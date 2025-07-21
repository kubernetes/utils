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
	"net/netip"

	forkednet "k8s.io/utils/internal/third_party/forked/golang/net"
)

// ParseIP parses an IPv4 or IPv6 address to a net.IP. This accepts both fully-valid IP
// addresses and irregular/ambiguous forms, making it usable for both validated and
// non-validated input strings. It should be used instead of net.ParseIP (which rejects
// some strings we need to accept for backward compatibility) and the old
// netutilsv1.ParseIPSloppy.
//
// Compare ParseAddr, which returns a netip.Addr but is otherwise identical.
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

// ParseAddr parses an IPv4 or IPv6 address to a netip.Addr. This accepts both fully-valid
// IP addresses and irregular/ambiguous forms, making it usable for both validated and
// non-validated input strings. As compared to netip.ParseAddr:
//
//   - It accepts IPv4 IPs with extra leading "0"s, for backward compatibility.
//   - It rejects IPs with attached zone identifiers (e.g. `"fe80::1234%eth0"`).
//   - It converts "IPv4-mapped IPv6" addresses (e.g. `"::ffff:1.2.3.4"`) to the
//     corresponding "plain" IPv4 values (e.g. `"1.2.3.4"`). That is, it never returns an
//     address for which `Is4In6()` would return `true`.
//
// Compare ParseIP, which returns a net.IP but is otherwise identical.
func ParseAddr(ipStr string) (netip.Addr, error) {
	// To ensure identical parsing, we use ParseIP and then convert. (If ParseIP
	// returns a nil ip, AddrFromIP will convert that to the zero/invalid netip.Addr,
	// which is what we want.)
	ip, err := ParseIP(ipStr)
	return AddrFromIP(ip), err
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
//
// Compare ParsePrefix, which returns a netip.Prefix but is otherwise identical.
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

// ParsePrefix parses an IPv4 or IPv6 CIDR string representing a subnet or mask, to a
// netip.Prefix. This accepts both fully-valid CIDR values and irregular/ambiguous forms,
// making it usable for both validated and non-validated input strings. As compared to
// netip.ParsePrefix:
//
//   - It accepts IPv4 IPs with extra leading "0"s, for backward compatibility.
//   - It converts "IPv4-mapped IPv6" addresses (e.g. `"::ffff:1.2.3.0/120"`) to the
//     corresponding "plain" IPv4 values (e.g. `"1.2.3.0/24"`). That is, it never returns
//     a prefix for which `.Addr().Is4In6()` would return `true`.
//   - When given a CIDR string with bits set beyond the prefix length, like
//     `"192.168.1.5/24"`, it discards those extra bits (the equivalent of calling
//     .Masked() on the return value of netip.ParsePrefix).
//
// Compare ParseIPNet, which returns a *net.IPNet but is otherwise identical.
func ParsePrefix(cidrStr string) (netip.Prefix, error) {
	// To ensure identical parsing, we use ParseIPNet and then convert. (If ParseIPNet
	// returns nil, PrefixFromIPNet will convert that to the zero/invalid
	// netip.Prefix, which is what we want.)
	ipnet, err := ParseIPNet(cidrStr)
	return PrefixFromIPNet(ipnet), err
}
