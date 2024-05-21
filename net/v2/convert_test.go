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
	"testing"
)

func TestAddrFromIP(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestIPs {
		if tc.skipConvert {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, ip := range tc.ips {
				addr := AddrFromIP(ip)
				if tc.addrs[0] != addr {
					t.Errorf("IP %d %#v %s converted to addr %q, but expected %q", i+1, ip, ip, addr, tc.addrs[0])
				}

				// No net.IP should convert to an IPv4-mapped IPv6 netip.Addr
				if addr.Is4In6() {
					t.Errorf("AddrFromIP() converted IP %d %#v %s to IPv4-mapped IPv6 Addr %#v %s", i+1, ip, ip, addr, addr)
				}
				// And thus every value should round-trip.
				rtIP := IPFromAddr(addr)
				if !ip.Equal(rtIP) {
					t.Errorf("IP %d %#v %s round-tripped to %#v %s", i+1, ip, ip, rtIP, rtIP)
				}
			}
		})
	}

	// See test cases in ips_test.go
	for _, tc := range badTestIPs {
		if tc.skipConvert {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, ip := range tc.ips {
				addr := AddrFromIP(ip)
				if addr.IsValid() {
					t.Errorf("Expected IP %d %#v to convert to invalid netip.Addr but got %#v %s", i+1, ip, addr, addr)
				}
			}
		})
	}
}

func TestIPFromAddr(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestIPs {
		if tc.skipConvert {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, addr := range tc.addrs {
				ip := IPFromAddr(addr)
				if !ip.Equal(tc.ips[0]) {
					t.Errorf("addr %d %#v %s converted to ip %q, but expected %q", i, addr, addr, ip, tc.ips[0])
				}

				// As long as addr is not IPv4-mapped IPv6, it should round-trip.
				if !addr.Is4In6() {
					rtAddr := AddrFromIP(ip)
					if addr != rtAddr {
						t.Errorf("Addr %d %#v %s round-tripped to %#v %s", i+1, addr, addr, rtAddr, rtAddr)
					}
				}
			}
		})
	}

	// Conversion of IPv4-mapped IPv6 is asymmetric because netip.Addr distinguishes
	// plain IPv4 from IPv4-mapped IPv6, while net.IP does not. The "IPv4-mapped IPv6"
	// test case in goodTestIPs covers most of the cases, but goodTestIPs has no way
	// to describe the asymmetric part.
	t.Run("IPv4-mapped IPv6 conversion from netip.Addr", func(t *testing.T) {
		addr := netip.AddrFrom16([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, 1, 2, 3, 4})
		if !addr.Is4In6() {
			panic("failed to create IPv4-mapped IPv6 netip.Addr?")
		}

		ip := IPFromAddr(addr)
		expectedIP := net.IP{1, 2, 3, 4}
		if !ip.Equal(expectedIP) {
			t.Errorf("netip.Addr %q converted to %q, expected %q", addr, ip, expectedIP)
		}
		rtAddr := AddrFromIP(ip)
		if rtAddr == addr {
			t.Errorf("IPv4-mapped IPv6 netip.Addr unexpectedly round-tripped through net.IP!")
		}
	})

	// See test cases in ips_test.go
	for _, tc := range badTestIPs {
		if tc.skipConvert {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, addr := range tc.addrs {
				ip := IPFromAddr(addr)
				if ip != nil {
					t.Errorf("Expected Addr %d %#v to convert to invalid net.IP but got %#v %s", i+1, addr, ip, ip)
				}
			}
		})
	}
}

func TestPrefixFromIPNet(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestCIDRs {
		if tc.skipConvert {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, ipnet := range tc.ipnets {
				prefix := PrefixFromIPNet(ipnet)
				if tc.prefixes[0] != prefix {
					t.Errorf("IPNet %d %#v %s converted to prefix %q, but expected %q", i+1, *ipnet, ipnet, prefix, tc.prefixes[0])
				}

				// No net.IPNet should convert to an IPv4-mapped IPv6 netip.Prefix
				if prefix.Addr().Is4In6() {
					t.Errorf("PrefixFromIPNet() converted IPNet %d %#v %s to IPv4-mapped IPv6 prefix %#v %s", i+1, *ipnet, ipnet, prefix, prefix)
				}
				// And thus every value should round-trip.
				rtIPNet := IPNetFromPrefix(prefix)
				if rtIPNet.String() != ipnet.String() {
					t.Errorf("IPNet %d %#v %s round-tripped to %#v %s", i+1, *ipnet, ipnet, *rtIPNet, rtIPNet)
				}
			}
		})
	}

	// See test cases in ips_test.go
	for _, tc := range badTestCIDRs {
		if tc.skipConvert {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, ipnet := range tc.ipnets {
				prefix := PrefixFromIPNet(ipnet)
				if prefix.IsValid() {
					str := "<nil>"
					if ipnet != nil {
						str = fmt.Sprintf("%#v", *ipnet)
					}
					t.Errorf("Expected IPNet %d %s to convert to invalid netip.Prefix but got %#v %s", i+1, str, prefix, prefix)
				}
			}
		})
	}
}

func TestIPNetFromPrefix(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestCIDRs {
		if tc.skipConvert {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, prefix := range tc.prefixes {
				ipnet := IPNetFromPrefix(prefix)
				if ipnet.String() != tc.ipnets[0].String() {
					t.Errorf("prefix %d %#v %s converted to ipnet %q, but expected %q", i, prefix, prefix, ipnet, tc.ipnets[0])
				}

				// As long as addr is not IPv4-mapped IPv6, it should round-trip.
				if !prefix.Addr().Is4In6() {
					rtPrefix := PrefixFromIPNet(ipnet)
					if prefix != rtPrefix {
						t.Errorf("prefix %d %#v %s round-tripped to %#v %s", i+1, prefix, prefix, rtPrefix, rtPrefix)
					}
				}
			}
		})
	}

	// Conversion of IPv4-mapped IPv6 is asymmetric because netip.Addr distinguishes
	// plain IPv4 from IPv4-mapped IPv6, while net.IP does not. The "IPv4-mapped IPv6"
	// test case in goodTestCIDRs covers most of the cases, but goodTestCIDRs has no way
	// to describe the asymmetric part.
	t.Run("IPv4-mapped IPv6 conversion from netip.Prefix", func(t *testing.T) {
		prefix := netip.MustParsePrefix("::ffff:1.2.3.0/120")
		if !prefix.Addr().Is4In6() {
			panic("failed to create IPv4-mapped IPv6 netip.Addr?")
		}

		ipnet := IPNetFromPrefix(prefix)
		expected := "1.2.3.0/24"
		if ipnet.String() != expected {
			t.Errorf("netip.Prefix %q converted to %q, expected %q", prefix, ipnet.String(), expected)
		}
		rtPrefix := PrefixFromIPNet(ipnet)
		if rtPrefix == prefix {
			t.Errorf("IPv4-mapped IPv6 netip.Prefix unexpectedly round-tripped through net.IPNet!")
		}
	})

	// See test cases in ips_test.go
	for _, tc := range badTestCIDRs {
		if tc.skipConvert {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, prefix := range tc.prefixes {
				ipnet := IPNetFromPrefix(prefix)
				if ipnet != nil {
					t.Errorf("Expected Prefix %d %#v to convert to invalid net.IPNet but got %#v %s", i+1, prefix, *ipnet, ipnet)
				}
			}
		})
	}
}
