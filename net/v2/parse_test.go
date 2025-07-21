/*
Copyright 2025 The Kubernetes Authors.

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
	"testing"
)

func TestParseIP(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestIPs {
		if tc.skipParse {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, str := range tc.strings {
				ip, err := ParseIP(str)
				if err != nil {
					t.Errorf("expected %q to parse, but got error %v", str, err)
				}
				if !ip.Equal(tc.ips[0]) {
					t.Errorf("expected string %d %q to parse equal to IP %#v %q but got %#v (%q)", i+1, str, tc.ips[0], tc.ips[0].String(), ip, ip.String())
				}
			}
		})
	}

	// See test cases in ips_test.go
	for _, tc := range badTestIPs {
		if tc.skipParse {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, ip := range tc.ips {
				errStr := ip.String()
				parsedIP, _ := ParseIP(errStr)
				if parsedIP != nil {
					t.Errorf("expected IP %d %#v (%q) to not re-parse but got %#v (%q)", i+1, ip, errStr, parsedIP, parsedIP.String())
				}
			}

			for i, addr := range tc.addrs {
				errStr := addr.String()
				parsedIP, _ := ParseIP(errStr)
				if parsedIP != nil {
					t.Errorf("expected Addr %d %#v (%q) to not re-parse but got %#v (%q)", i+1, addr, errStr, parsedIP, parsedIP.String())
				}
			}

			for i, str := range tc.strings {
				ip, _ := ParseIP(str)
				if ip != nil {
					t.Errorf("expected string %d %q to not parse but got %#v (%q)", i+1, str, ip, ip.String())
				}
			}
		})
	}
}

func TestParseAddr(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestIPs {
		if tc.skipParse {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, str := range tc.strings {
				addr, err := ParseAddr(str)
				if err != nil {
					t.Errorf("expected %q to parse, but got error %v", str, err)
				}
				if addr != tc.addrs[0] {
					t.Errorf("expected string %d %q to parse equal to Addr %#v %q but got %#v (%q)", i+1, str, tc.addrs[0], tc.addrs[0].String(), addr, addr.String())
				}
			}
		})
	}

	for _, tc := range badTestIPs {
		if tc.skipParse {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, ip := range tc.ips {
				errStr := ip.String()
				parsedAddr, err := ParseAddr(errStr)
				if err == nil {
					t.Errorf("expected IP %d %#v (%q) to not re-parse but got %#v (%q)", i+1, ip, errStr, parsedAddr, parsedAddr.String())
				}
			}

			for i, addr := range tc.addrs {
				errStr := addr.String()
				parsedAddr, err := ParseAddr(errStr)
				if err == nil {
					t.Errorf("expected Addr %d %#v (%q) to not re-parse but got %#v (%q)", i+1, addr, errStr, parsedAddr, parsedAddr.String())
				}
			}

			for i, str := range tc.strings {
				addr, err := ParseAddr(str)
				if err == nil {
					t.Errorf("expected string %d %q to not parse but got %#v (%q)", i+1, str, addr, addr.String())
				}
			}
		})
	}
}

func TestParseIPNet(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestCIDRs {
		if tc.skipParse {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, str := range tc.strings {
				ipnet, err := ParseIPNet(str)
				if err != nil {
					t.Errorf("expected %q to parse, but got error %v", str, err)
				}
				ifaddr, err := ParseIPAsIPNet(str)
				if err != nil {
					t.Errorf("expected %q to parse via ParseIPAsIPNet, but got error %v", str, err)
				}

				if tc.ifaddr {
					// The test case expects ParseIPNet and
					// ParseIPAsIPNet to parse to different values.
					if ipnet.String() == ifaddr.String() {
						t.Errorf("expected %q to parse differently with ParseIPNet and ParseIPAsIPNet but got %q for both", str, ipnet.String())
					}
					// In this case, it's the ParseIPAsIPNet value
					// that should re-stringify correctly. (ParseIPNet
					// will have discarded the trailing bits.)
					ipnet = ifaddr
				} else {
					// Some strings might parse to the same value and
					// others might parse to different values.
					// However, in all cases, the ParseIPAsIPNet value
					// should be the same as the ParseIPNet value
					// after masking it.
					if !ipnet.IP.Equal(ifaddr.IP.Mask(ifaddr.Mask)) {
						t.Errorf("expected %q to parse similarly with ParseIPNet and ParseIPAsIPNet but got IPs %q and %q->%q", str, ipnet.IP, ifaddr, ifaddr.IP.Mask(ifaddr.Mask))
					}
				}

				if ipnet.String() != tc.ipnets[0].String() {
					t.Errorf("expected string %d %q to parse and re-stringify to %q but got %q", i+1, str, tc.ipnets[0].String(), ipnet.String())
				}
			}
		})
	}

	// See test cases in ips_test.go
	for _, tc := range badTestCIDRs {
		if tc.skipParse {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, ipnet := range tc.ipnets {
				errStr := ipnet.String()
				parsedIPNet, err := ParseIPNet(errStr)
				if err == nil {
					t.Errorf("expected IPNet %d %q to not parse but got %#v (%q)", i+1, errStr, *parsedIPNet, parsedIPNet.String())
				}
			}

			for i, prefix := range tc.prefixes {
				errStr := prefix.String()
				parsedIPNet, err := ParseIPNet(errStr)
				if err == nil {
					t.Errorf("expected Prefix %d %#v %q to not parse but got %#v (%q)", i+1, prefix, errStr, *parsedIPNet, parsedIPNet.String())
				}
			}

			for i, str := range tc.strings {
				ipnet, err := ParseIPNet(str)
				if err == nil {
					t.Errorf("expected string %d %q to not parse but got %#v (%q)", i+1, str, *ipnet, ipnet.String())
				}
			}
		})
	}
}

func TestParsePrefix(t *testing.T) {
	// See test cases in ips_test.go
	for _, tc := range goodTestCIDRs {
		if tc.skipParse {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, str := range tc.strings {
				prefix, err := ParsePrefix(str)
				if err != nil {
					t.Errorf("expected %q to parse, but got error %v", str, err)
				}
				ifaddr, err := ParseAddrAsPrefix(str)
				if err != nil {
					t.Errorf("expected %q to parse via ParseAddrAsPrefix, but got error %v", str, err)
				}

				if tc.ifaddr {
					// The test case expects ParsePrefix and
					// ParseAddrAsPrefix to parse to different values.
					if prefix == ifaddr {
						t.Errorf("expected %q to parse differently with ParsePrefix and ParseAddrAsPrefix but got %#v %q for both", str, prefix, prefix)
					}
					// In this case, it's the ParseAddrAsPrefix value
					// that should re-stringify correctly. (ParsePrefix
					// will have discarded the trailing bits.)
					prefix = ifaddr
				} else {
					// Some strings might parse to the same value and
					// others might parse to different values.
					// However, in all cases, the ParseAddrAsPrefix
					// value should be the same as the ParsePrefix
					// value after masking it.
					if prefix != ifaddr.Masked() {
						t.Errorf("expected %q to parse similarly with ParsePrefix and ParseAddrAsPrefix but got %q and %q->%q", str, prefix, ifaddr, ifaddr.Masked())
					}
				}

				if prefix != tc.prefixes[0] {
					t.Errorf("expected string %d %q to parse equal to Prefix %#v %q but got %#v (%q)", i+1, str, tc.prefixes[0], tc.prefixes[0].String(), prefix, prefix.String())
				}
			}
		})
	}

	// See test cases in ips_test.go
	for _, tc := range badTestCIDRs {
		if tc.skipParse {
			continue
		}
		t.Run(tc.desc, func(t *testing.T) {
			for i, ipnet := range tc.ipnets {
				errStr := ipnet.String()
				parsedPrefix, err := ParsePrefix(errStr)
				if err == nil {
					t.Errorf("expected IPNet %d %q to not parse but got %#v (%q)", i+1, errStr, parsedPrefix, parsedPrefix.String())
				}
			}

			for i, prefix := range tc.prefixes {
				errStr := prefix.String()
				parsedPrefix, err := ParsePrefix(errStr)
				if err == nil {
					t.Errorf("expected Prefix %d %q to not parse but got %#v (%q)", i+1, errStr, parsedPrefix, parsedPrefix.String())
				}
			}

			for i, str := range tc.strings {
				prefix, err := ParsePrefix(str)
				if err == nil {
					t.Errorf("expected string %d %q to not parse but got %#v (%q)", i+1, str, prefix, prefix.String())
				}
			}
		})
	}
}
