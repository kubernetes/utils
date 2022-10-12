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

package netip

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/netip"
	"strconv"
)

// ParseCIDRs parses a list of cidrs and return error if any is invalid.
// order is maintained
func ParseCIDRs(cidrsString []string) ([]netip.Prefix, error) {
	cidrs := make([]netip.Prefix, 0, len(cidrsString))
	for _, cidrString := range cidrsString {
		cidr, err := netip.ParsePrefix(cidrString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cidr value:%q with error:%v", cidrString, err)
		}
		cidrs = append(cidrs, cidr)
	}
	return cidrs, nil
}

// IsDualStackIPs returns if a slice of ips is:
// - all are valid ips
// - at least one ip from each family (v4 or v6)
func IsDualStackIPs(ips []netip.Addr) (bool, error) {
	v4Found := false
	v6Found := false
	for _, ip := range ips {
		if !ip.IsValid() {
			return false, fmt.Errorf("ip %v is invalid", ip)
		}

		if v4Found && v6Found {
			continue
		}

		if ip.Is6() {
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
	parsedIPs := make([]netip.Addr, 0, len(ips))
	for _, ip := range ips {
		parsedIP, err := netip.ParseAddr(ip)
		if err != nil {
			return false, err
		}

		parsedIPs = append(parsedIPs, parsedIP)
	}
	return IsDualStackIPs(parsedIPs)
}

// IsDualStackCIDRs returns if
// - all are valid cidrs
// - at least one cidr from each family (v4 or v6)
func IsDualStackCIDRs(cidrs []netip.Prefix) (bool, error) {
	v4Found := false
	v6Found := false
	for _, cidr := range cidrs {
		if !cidr.IsValid() {
			return false, fmt.Errorf("cidr %v is invalid", cidr)
		}

		if v4Found && v6Found {
			continue
		}

		if cidr.Addr().Is6() {
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

// IsIPv6String returns if ip is IPv6.
func IsIPv6String(ip string) bool {
	ipa, _ := netip.ParseAddr(ip)
	return ipa.Is6()
}

// IsIPv6CIDRString returns if cidr is IPv6.
// This assumes cidr is a valid CIDR.
func IsIPv6CIDRString(cidr string) bool {
	ip, _ := netip.ParsePrefix(cidr)
	return ip.Addr().Is6()
}

// IsIPv4String returns if ip is IPv4.
func IsIPv4String(ip string) bool {
	ipa, _ := netip.ParseAddr(ip)
	return ipa.Is4()
}

// IsIPv4CIDRString returns if cidr is IPv4.
// This assumes cidr is a valid CIDR.
func IsIPv4CIDRString(cidr string) bool {
	ipp, _ := netip.ParsePrefix(cidr)
	return ipp.Addr().Is4()
}

// ParsePort parses a string representing an IP port.  If the string is not a
// valid port number, this returns an error.
func ParsePort(port string, allowZero bool) (int, error) {
	portInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, err
	}
	if portInt == 0 && !allowZero {
		return 0, errors.New("0 is not a valid port number")
	}
	return int(portInt), nil
}

// BigForIP creates a big.Int based on the provided netip.Addr
func BigForIP(ip netip.Addr) *big.Int {
	ipv6 := ip.As16()
	return big.NewInt(0).SetBytes(ipv6[:])
}

// AddIPOffset adds the provided integer offset to a base big.Int representing a netip.Addr
// NOTE: If you started with a v4 address and overflow it, you get a v6 result.
func AddIPOffset(base *big.Int, offset int) netip.Addr {
	r := big.NewInt(0).Add(base, big.NewInt(int64(offset))).Bytes()
	r = append(make([]byte, 16), r...)
	ipa, _ := netip.AddrFromSlice(r[len(r)-16:])
	if ipa.Is4In6() {
		return netip.AddrFrom4(ipa.As4())
	}
	return ipa
}

// RangeSize returns the size of a range in valid addresses.
// returns the size of the subnet (or math.MaxInt64 if the range size would overflow int64)
func RangeSize(subnet netip.Prefix) int64 {
	ones, bits := subnet.Bits(), subnet.Addr().BitLen()
	if bits == 32 && (bits-ones) >= 31 || bits == 128 && (bits-ones) >= 127 {
		return 0
	}
	// this checks that we are not overflowing an int64
	if bits-ones >= 63 {
		return math.MaxInt64
	}
	return int64(1) << uint(bits-ones)
}

// GetIndexedIP returns a net.IP that is subnet.IP + index in the contiguous IP space.
func GetIndexedIP(subnet netip.Prefix, index int) (netip.Addr, error) {
	ip := AddIPOffset(BigForIP(subnet.Addr()), index)
	if !subnet.Contains(ip) {
		return netip.Addr{}, fmt.Errorf("can't generate IP with index %d from subnet. subnet too small. subnet: %q", index, subnet)
	}
	return ip, nil
}
