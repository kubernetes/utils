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
	"testing"
)

func TestIsDualStack(t *testing.T) {
	testCases := []struct {
		desc           string
		ips            []string
		expectedResult bool
	}{
		{
			desc:           "should fail because length is not at least 2",
			ips:            []string{"1.1.1.1"},
			expectedResult: false,
		},
		{
			desc:           "should fail because length is not at least 2",
			ips:            []string{},
			expectedResult: false,
		},
		{
			desc:           "should fail because all are v4",
			ips:            []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			expectedResult: false,
		},
		{
			desc:           "should fail because all are v6",
			ips:            []string{"fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff", "fd92:20ba:ca:34f7:ffff:ffff:ffff:fff0", "fd92:20ba:ca:34f7:ffff:ffff:ffff:fff1"},
			expectedResult: false,
		},
		{
			desc:           "should fail because 2nd ip is invalid",
			ips:            []string{"1.1.1.1", "not-a-valid-ip"},
			expectedResult: false,
		},
		{
			desc:           "should fail because 1st ip is invalid",
			ips:            []string{"not-a-valid-ip", "fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff"},
			expectedResult: false,
		},
		{
			desc:           "should fail despite dual-stack because 3rd ip is invalid",
			ips:            []string{"1.1.1.1", "fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff", "not-a-valid-ip"},
			expectedResult: false,
		},
		{
			desc:           "dual-stack ipv4-primary",
			ips:            []string{"1.1.1.1", "fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff"},
			expectedResult: true,
		},
		{
			desc:           "dual-stack, multiple ipv6",
			ips:            []string{"fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff", "1.1.1.1", "fd92:20ba:ca:34f7:ffff:ffff:ffff:fff0"},
			expectedResult: true,
		},
		{
			desc:           "dual-stack, multiple ipv4",
			ips:            []string{"1.1.1.1", "fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff", "10.0.0.0"},
			expectedResult: true,
		},
		{
			desc:           "dual-stack, ipv6-primary",
			ips:            []string{"fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff", "1.1.1.1"},
			expectedResult: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			netips := make([]net.IP, len(tc.ips))
			for i := range tc.ips {
				netips[i] = ParseIPSloppy(tc.ips[i])
			}

			dualStack := IsDualStack(tc.ips)
			if dualStack != tc.expectedResult {
				t.Errorf("expected %v, []string got %v", tc.expectedResult, dualStack)
			}
			if IsDualStackPair(tc.ips) != (dualStack && len(tc.ips) == 2) {
				t.Errorf("IsDualStackIPPair gave wrong result for []string")
			}

			dualStack = IsDualStack(netips)
			if dualStack != tc.expectedResult {
				t.Errorf("expected %v []net.IP got %v", tc.expectedResult, dualStack)
			}
			if IsDualStackPair(netips) != (dualStack && len(tc.ips) == 2) {
				t.Errorf("IsDualStackIPPair gave wrong result for []net.IP")
			}
		})
	}
}

func TestIsDualStackCIDRs(t *testing.T) {
	testCases := []struct {
		desc           string
		cidrs          []string
		expectedResult bool
	}{
		{
			desc:           "should fail because length is not at least 2",
			cidrs:          []string{"10.10.10.10/8"},
			expectedResult: false,
		},
		{
			desc:           "should fail because length is not at least 2",
			cidrs:          []string{},
			expectedResult: false,
		},
		{
			desc:           "should fail because all cidrs are v4",
			cidrs:          []string{"10.10.10.10/8", "20.20.20.20/8", "30.30.30.30/8"},
			expectedResult: false,
		},
		{
			desc:           "should fail because all cidrs are v6",
			cidrs:          []string{"2000::/10", "3000::/10"},
			expectedResult: false,
		},
		{
			desc:           "should fail because 2nd cidr is invalid",
			cidrs:          []string{"10.10.10.10/8", "not-a-valid-cidr"},
			expectedResult: false,
		},
		{
			desc:           "should fail because 1st cidr is invalid",
			cidrs:          []string{"not-a-valid-ip", "2000::/10"},
			expectedResult: false,
		},
		{
			desc:           "dual-stack, ipv4-primary",
			cidrs:          []string{"10.10.10.10/8", "2000::/10"},
			expectedResult: true,
		},
		{
			desc:           "dual-stack, ipv6-primary",
			cidrs:          []string{"2000::/10", "10.10.10.10/8"},
			expectedResult: true,
		},
		{
			desc:           "dual-stack, multiple IPv6",
			cidrs:          []string{"2000::/10", "10.10.10.10/8", "3000::/10"},
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ipnets := make([]*net.IPNet, len(tc.cidrs))
			for i := range tc.cidrs {
				_, ipnets[i], _ = ParseCIDRSloppy(tc.cidrs[i])
			}

			dualStack := IsDualStackCIDRs(tc.cidrs)
			if dualStack != tc.expectedResult {
				t.Errorf("expected %v []string got %v", tc.expectedResult, dualStack)
			}
			if IsDualStackCIDRPair(tc.cidrs) != (dualStack && len(tc.cidrs) == 2) {
				t.Errorf("IsDualStackCIDRPair gave wrong result for []string")
			}

			dualStack = IsDualStackCIDRs(ipnets)
			if dualStack != tc.expectedResult {
				t.Errorf("expected %v []*net.IPNet got %v", tc.expectedResult, dualStack)
			}
			if IsDualStackCIDRPair(ipnets) != (dualStack && len(tc.cidrs) == 2) {
				t.Errorf("IsDualStackCIDRPair gave wrong result for []*net.IPNet")
			}
		})
	}
}

func checkOneIPFamily(t *testing.T, ip string, expectedFamily, family IPFamily, isIPv4, isIPv6 bool) {
	t.Helper()
	if family != expectedFamily {
		t.Errorf("Expect %q family %q, got %q", ip, expectedFamily, family)
	}
	if isIPv4 != (expectedFamily == IPv4) {
		t.Errorf("Expect %q ipv4 %v, got %v", ip, expectedFamily == IPv4, isIPv6)
	}
	if isIPv6 != (expectedFamily == IPv6) {
		t.Errorf("Expect %q ipv6 %v, got %v", ip, expectedFamily == IPv6, isIPv6)
	}
}

func TestIPFamilyOf(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestIPs {
		if tc.skipFamily {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for _, str := range tc.strings {
				family := IPFamilyOf(str)
				isIPv4 := IsIPv4(str)
				isIPv6 := IsIPv6(str)
				checkOneIPFamily(t, str, tc.family, family, isIPv4, isIPv6)
			}
			for _, ip := range tc.ips {
				family := IPFamilyOf(ip)
				isIPv4 := IsIPv4(ip)
				isIPv6 := IsIPv6(ip)
				checkOneIPFamily(t, ip.String(), tc.family, family, isIPv4, isIPv6)
			}
		})
	}

	// See test cases in ips_test.go
	for _, tc := range badTestIPs {
		if tc.skipFamily {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for _, ip := range tc.ips {
				family := IPFamilyOf(ip)
				isIPv4 := IsIPv4(ip)
				isIPv6 := IsIPv6(ip)
				checkOneIPFamily(t, fmt.Sprintf("%#v", ip), IPFamilyUnknown, family, isIPv4, isIPv6)
			}
			for _, str := range tc.strings {
				family := IPFamilyOf(str)
				isIPv4 := IsIPv4(str)
				isIPv6 := IsIPv6(str)
				checkOneIPFamily(t, str, IPFamilyUnknown, family, isIPv4, isIPv6)
			}
		})
	}
}

func TestIPFamilyOfCIDR(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestCIDRs {
		if tc.skipFamily {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for _, str := range tc.strings {
				family := IPFamilyOfCIDR(str)
				isIPv4 := IsIPv4CIDR(str)
				isIPv6 := IsIPv6CIDR(str)
				checkOneIPFamily(t, str, tc.family, family, isIPv4, isIPv6)
			}
			for _, ipnet := range tc.ipnets {
				family := IPFamilyOfCIDR(ipnet)
				isIPv4 := IsIPv4CIDR(ipnet)
				isIPv6 := IsIPv6CIDR(ipnet)
				checkOneIPFamily(t, ipnet.String(), tc.family, family, isIPv4, isIPv6)
			}
		})
	}

	// See test cases in ips_test.go
	for _, tc := range badTestCIDRs {
		if tc.skipFamily {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for _, ipnet := range tc.ipnets {
				family := IPFamilyOfCIDR(ipnet)
				isIPv4 := IsIPv4CIDR(ipnet)
				isIPv6 := IsIPv6CIDR(ipnet)
				str := "<nil>"
				if ipnet != nil {
					str = fmt.Sprintf("%#v", *ipnet)
				}
				checkOneIPFamily(t, str, IPFamilyUnknown, family, isIPv4, isIPv6)
			}
			for _, str := range tc.strings {
				family := IPFamilyOfCIDR(str)
				isIPv4 := IsIPv4CIDR(str)
				isIPv6 := IsIPv6CIDR(str)
				checkOneIPFamily(t, str, IPFamilyUnknown, family, isIPv4, isIPv6)
			}
		})
	}
}
