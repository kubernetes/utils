/*
Copyright 2017 The Kubernetes Authors.

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

import "testing"

func TestLocalPortString(t *testing.T) {
	testCases := []struct {
		description string
		ip          string
		family      IPFamily
		port        int
		protocol    string
		expectedStr string
	}{
		{"IPv4 UDP", "1.2.3.4", "", 9999, "udp", "\"IPv4 UDP\" (1.2.3.4:9999/udp)"},
		{"IPv4 TCP", "5.6.7.8", "", 1053, "tcp", "\"IPv4 TCP\" (5.6.7.8:1053/tcp)"},
		{"IPv4 TCP", "", IPv4, 1053, "tcp", "\"IPv4 TCP\" (:1053/tcp4)"},
		{"IPv6 TCP", "2001:db8::1", "", 80, "tcp", "\"IPv6 TCP\" ([2001:db8::1]:80/tcp)"},
		{"IPv4 SCTP", "9.10.11.12", "", 7777, "sctp", "\"IPv4 SCTP\" (9.10.11.12:7777/sctp)"},
		{"IPv6 SCTP", "2001:db8::2", "", 80, "sctp", "\"IPv6 SCTP\" ([2001:db8::2]:80/sctp)"},
	}

	for _, tc := range testCases {
		lp := &LocalPort{
			Description: tc.description,
			IP:          tc.ip,
			IPFamily:    tc.family,
			Port:        tc.port,
			Protocol:    tc.protocol,
		}
		str := lp.String()
		if str != tc.expectedStr {
			t.Errorf("Unexpected output for %s, expected: %s, got: %s", tc.description, tc.expectedStr, str)
		}
	}
}
