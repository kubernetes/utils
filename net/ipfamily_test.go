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
	"testing"
)

func TestDualStackIPs(t *testing.T) {
	testCases := []struct {
		ips            []string
		errMessage     string
		expectedResult bool
		expectError    bool
	}{
		{
			ips:            []string{"1.1.1.1"},
			errMessage:     "should fail because length is not at least 2",
			expectedResult: false,
			expectError:    false,
		},
		{
			ips:            []string{},
			errMessage:     "should fail because length is not at least 2",
			expectedResult: false,
			expectError:    false,
		},
		{
			ips:            []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			errMessage:     "should fail because all are v4",
			expectedResult: false,
			expectError:    false,
		},
		{
			ips:            []string{"fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff", "fd92:20ba:ca:34f7:ffff:ffff:ffff:fff0", "fd92:20ba:ca:34f7:ffff:ffff:ffff:fff1"},
			errMessage:     "should fail because all are v6",
			expectedResult: false,
			expectError:    false,
		},
		{
			ips:            []string{"1.1.1.1", "not-a-valid-ip"},
			errMessage:     "should fail because 2nd ip is invalid",
			expectedResult: false,
			expectError:    true,
		},
		{
			ips:            []string{"not-a-valid-ip", "fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff"},
			errMessage:     "should fail because 1st ip is invalid",
			expectedResult: false,
			expectError:    true,
		},
		{
			ips:            []string{"1.1.1.1", "fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff"},
			errMessage:     "expected success, but found failure",
			expectedResult: true,
			expectError:    false,
		},
		{
			ips:            []string{"fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff", "1.1.1.1", "fd92:20ba:ca:34f7:ffff:ffff:ffff:fff0"},
			errMessage:     "expected success, but found failure",
			expectedResult: true,
			expectError:    false,
		},
		{
			ips:            []string{"1.1.1.1", "fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff", "10.0.0.0"},
			errMessage:     "expected success, but found failure",
			expectedResult: true,
			expectError:    false,
		},
		{
			ips:            []string{"fd92:20ba:ca:34f7:ffff:ffff:ffff:ffff", "1.1.1.1"},
			errMessage:     "expected success, but found failure",
			expectedResult: true,
			expectError:    false,
		},
	}
	// for each test case, test the regular func and the string func
	for _, tc := range testCases {
		dualStack, err := IsDualStackIPStrings(tc.ips)
		if err == nil && tc.expectError {
			t.Errorf("%s", tc.errMessage)
			continue
		}
		if err != nil && !tc.expectError {
			t.Errorf("failed to run test case for %v, error: %v", tc.ips, err)
			continue
		}
		if dualStack != tc.expectedResult {
			t.Errorf("%v for %v", tc.errMessage, tc.ips)
		}
	}

	for _, tc := range testCases {
		ips := make([]net.IP, 0, len(tc.ips))
		for _, ip := range tc.ips {
			parsedIP := ParseIPSloppy(ip)
			ips = append(ips, parsedIP)
		}
		dualStack, err := IsDualStackIPs(ips)
		if err == nil && tc.expectError {
			t.Errorf("%s", tc.errMessage)
			continue
		}
		if err != nil && !tc.expectError {
			t.Errorf("failed to run test case for %v, error: %v", tc.ips, err)
			continue
		}
		if dualStack != tc.expectedResult {
			t.Errorf("%v for %v", tc.errMessage, tc.ips)
		}
	}
}

func TestDualStackCIDRs(t *testing.T) {
	testCases := []struct {
		cidrs          []string
		errMessage     string
		expectedResult bool
		expectError    bool
	}{
		{
			cidrs:          []string{"10.10.10.10/8"},
			errMessage:     "should fail because length is not at least 2",
			expectedResult: false,
			expectError:    false,
		},
		{
			cidrs:          []string{},
			errMessage:     "should fail because length is not at least 2",
			expectedResult: false,
			expectError:    false,
		},
		{
			cidrs:          []string{"10.10.10.10/8", "20.20.20.20/8", "30.30.30.30/8"},
			errMessage:     "should fail because all cidrs are v4",
			expectedResult: false,
			expectError:    false,
		},
		{
			cidrs:          []string{"2000::/10", "3000::/10"},
			errMessage:     "should fail because all cidrs are v6",
			expectedResult: false,
			expectError:    false,
		},
		{
			cidrs:          []string{"10.10.10.10/8", "not-a-valid-cidr"},
			errMessage:     "should fail because 2nd cidr is invalid",
			expectedResult: false,
			expectError:    true,
		},
		{
			cidrs:          []string{"not-a-valid-ip", "2000::/10"},
			errMessage:     "should fail because 1st cidr is invalid",
			expectedResult: false,
			expectError:    true,
		},
		{
			cidrs:          []string{"10.10.10.10/8", "2000::/10"},
			errMessage:     "expected success, but found failure",
			expectedResult: true,
			expectError:    false,
		},
		{
			cidrs:          []string{"2000::/10", "10.10.10.10/8"},
			errMessage:     "expected success, but found failure",
			expectedResult: true,
			expectError:    false,
		},
		{
			cidrs:          []string{"2000::/10", "10.10.10.10/8", "3000::/10"},
			errMessage:     "expected success, but found failure",
			expectedResult: true,
			expectError:    false,
		},
	}

	// for each test case, test the regular func and the string func
	for _, tc := range testCases {
		dualStack, err := IsDualStackCIDRStrings(tc.cidrs)
		if err == nil && tc.expectError {
			t.Errorf("%s", tc.errMessage)
			continue
		}
		if err != nil && !tc.expectError {
			t.Errorf("failed to run test case for %v, error: %v", tc.cidrs, err)
			continue
		}
		if dualStack != tc.expectedResult {
			t.Errorf("%v for %v", tc.errMessage, tc.cidrs)
		}
	}

	for _, tc := range testCases {
		cidrs := make([]*net.IPNet, 0, len(tc.cidrs))
		for _, cidr := range tc.cidrs {
			_, parsedCIDR, _ := ParseCIDRSloppy(cidr)
			cidrs = append(cidrs, parsedCIDR)
		}

		dualStack, err := IsDualStackCIDRs(cidrs)
		if err == nil && tc.expectError {
			t.Errorf("%s", tc.errMessage)
			continue
		}
		if err != nil && !tc.expectError {
			t.Errorf("failed to run test case for %v, error: %v", tc.cidrs, err)
			continue
		}
		if dualStack != tc.expectedResult {
			t.Errorf("%v for %v", tc.errMessage, tc.cidrs)
		}
	}
}

func TestIPFamilyString(t *testing.T) {
	testCases := []struct {
		ip     string
		family IPFamily
	}{
		{
			ip:     "0.0.0.0",
			family: IPv4,
		},
		{
			ip:     "255.255.255.255",
			family: IPv4,
		},
		{
			ip:     "127.0.0.1",
			family: IPv4,
		},
		{
			ip:     "192.168.0.0",
			family: IPv4,
		},
		{
			ip:     "1.2.3.4",
			family: IPv4,
		},
		{
			ip:     "bad ip",
			family: IPFamilyUnknown,
		},
		{
			// CIDR rather than IP
			ip:     "192.168.0.0/16",
			family: IPFamilyUnknown,
		},
		{
			ip:     "::",
			family: IPv6,
		},
		{
			ip:     "::1",
			family: IPv6,
		},
		{
			ip:     "fd00::600d:f00d",
			family: IPv6,
		},
		{
			ip:     "2001:db8::5",
			family: IPv6,
		},
	}
	for i := range testCases {
		family := IPFamilyOfString(testCases[i].ip)
		isIPv4 := IsIPv4String(testCases[i].ip)
		isIPv6 := IsIPv6String(testCases[i].ip)
		switch testCases[i].family {
		case IPFamilyUnknown:
			if family != IPFamilyUnknown || isIPv4 || isIPv6 {
				t.Errorf("[%d] Expected family %q, isIPv4 %v, isIPv6 %v. Got family %q, isIPv4 %v, isIPv6 %v", i+1, IPFamilyUnknown, false, false, family, isIPv4, isIPv6)
			}
		case IPv4:
			if family != IPv4 || !isIPv4 || isIPv6 {
				t.Errorf("[%d] Expected family %q, isIPv4 %v, isIPv6 %v. Got family %q, isIPv4 %v, isIPv6 %v", i+1, IPv4, true, false, family, isIPv4, isIPv6)
			}
		case IPv6:
			if family != IPv6 || isIPv4 || !isIPv6 {
				t.Errorf("[%d] Expected family %q, isIPv4 %v, isIPv6 %v. Got family %q, isIPv4 %v, isIPv6 %v", i+1, IPv6, false, true, family, isIPv4, isIPv6)
			}
		}
	}
}

func TestIPFamilyOf(t *testing.T) {
	testCases := []struct {
		ip     net.IP
		family IPFamily
	}{
		{
			ip:     net.IPv4zero,
			family: IPv4,
		},
		{
			ip:     net.IPv4bcast,
			family: IPv4,
		},
		{
			ip:     ParseIPSloppy("127.0.0.1"),
			family: IPv4,
		},
		{
			ip:     ParseIPSloppy("10.20.40.40"),
			family: IPv4,
		},
		{
			ip:     ParseIPSloppy("172.17.3.0"),
			family: IPv4,
		},
		{
			ip:     nil,
			family: IPFamilyUnknown,
		},
		{
			ip:     net.IPv6loopback,
			family: IPv6,
		},
		{
			ip:     net.IPv6zero,
			family: IPv6,
		},
		{
			ip:     ParseIPSloppy("fd00::600d:f00d"),
			family: IPv6,
		},
		{
			ip:     ParseIPSloppy("2001:db8::5"),
			family: IPv6,
		},
	}
	for i := range testCases {
		family := IPFamilyOf(testCases[i].ip)
		isIPv4 := IsIPv4(testCases[i].ip)
		isIPv6 := IsIPv6(testCases[i].ip)
		switch testCases[i].family {
		case IPFamilyUnknown:
			if family != IPFamilyUnknown || isIPv4 || isIPv6 {
				t.Errorf("[%d] Expected family %q, isIPv4 %v, isIPv6 %v. Got family %q, isIPv4 %v, isIPv6 %v", i+1, IPFamilyUnknown, false, false, family, isIPv4, isIPv6)
			}
		case IPv4:
			if family != IPv4 || !isIPv4 || isIPv6 {
				t.Errorf("[%d] Expected family %q, isIPv4 %v, isIPv6 %v. Got family %q, isIPv4 %v, isIPv6 %v", i+1, IPv4, true, false, family, isIPv4, isIPv6)
			}
		case IPv6:
			if family != IPv6 || isIPv4 || !isIPv6 {
				t.Errorf("[%d] Expected family %q, isIPv4 %v, isIPv6 %v. Got family %q, isIPv4 %v, isIPv6 %v", i+1, IPv6, false, true, family, isIPv4, isIPv6)
			}
		}
	}
}

func TestIPFamilyOfCIDR(t *testing.T) {
	testCases := []struct {
		desc   string
		cidr   string
		family IPFamily
	}{
		{
			desc:   "ipv4 CIDR 1",
			cidr:   "10.0.0.0/8",
			family: IPv4,
		},
		{
			desc:   "ipv4 CIDR 2",
			cidr:   "192.168.0.0/16",
			family: IPv4,
		},
		{
			desc:   "ipv6 CIDR 1",
			cidr:   "::/1",
			family: IPv6,
		},
		{
			desc:   "ipv6 CIDR 2",
			cidr:   "2000::/10",
			family: IPv6,
		},
		{
			desc:   "ipv6 CIDR 3",
			cidr:   "2001:db8::/32",
			family: IPv6,
		},
		{
			desc:   "Unknown IP family",
			cidr:   "invalid.cidr/mask",
			family: IPFamilyUnknown,
		},
		{
			desc:   "IP rather than CIDR",
			cidr:   "192.168.0.1",
			family: IPFamilyUnknown,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			family := IPFamilyOfCIDRString(tc.cidr)
			isIPv4 := IsIPv4CIDRString(tc.cidr)
			isIPv6 := IsIPv6CIDRString(tc.cidr)
			switch tc.family {
			case IPFamilyUnknown:
				if family != IPFamilyUnknown || isIPv4 || isIPv6 {
					t.Errorf("Expected family %q, isIPv4 %v, isIPv6 %v. Got family %q, isIPv4 %v, isIPv6 %v", IPFamilyUnknown, false, false, family, isIPv4, isIPv6)
				}
			case IPv4:
				if family != IPv4 || !isIPv4 || isIPv6 {
					t.Errorf("Expected family %q, isIPv4 %v, isIPv6 %v. Got family %q, isIPv4 %v, isIPv6 %v", IPv4, true, false, family, isIPv4, isIPv6)
				}
			case IPv6:
				if family != IPv6 || isIPv4 || !isIPv6 {
					t.Errorf("Expected family %q, isIPv4 %v, isIPv6 %v. Got family %q, isIPv4 %v, isIPv6 %v", IPv6, false, true, family, isIPv4, isIPv6)
				}
			}
		})
	}
}
