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
	"net"
	"net/netip"
	"reflect"
	"testing"
)

func TestBroadcastAddress(t *testing.T) {
	tests := []struct {
		name    string
		subnet  netip.Prefix
		want    netip.Addr
		wantErr bool
	}{
		{
			name:    "emty subnet",
			wantErr: true,
		},
		{
			name:   "IPv4 even mask",
			subnet: netip.MustParsePrefix("192.168.0.0/24"),
			want:   netip.MustParseAddr("192.168.0.255"),
		},
		{
			name:   "IPv4 odd mask",
			subnet: netip.MustParsePrefix("192.168.0.0/23"),
			want:   netip.MustParseAddr("192.168.1.255"),
		},
		{
			name:   "IPv6 even mask",
			subnet: netip.MustParsePrefix("fd00:1:2:3::/64"),
			want:   netip.MustParseAddr("fd00:1:2:3:ffff:ffff:ffff:ffff"),
		},
		{
			name:   "IPv6 odd mask",
			subnet: netip.MustParsePrefix("fd00:1:2:3::/57"),
			want:   netip.MustParseAddr("fd00:1:2:007f:ffff:ffff:ffff:ffff"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BroadcastAddress(tt.subnet)
			if (err != nil) != tt.wantErr {
				t.Errorf("BroadcastAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BroadcastAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPToAddr(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want netip.Addr
	}{
		{
			name: "IPv4",
			ip:   "192.168.2.2",
			want: netip.MustParseAddr("192.168.2.2"),
		},
		{
			name: "IPv6",
			ip:   "2001:db8::2",
			want: netip.MustParseAddr("2001:db8::2"),
		},
		{
			name: "IPv4 in IPv6",
			ip:   "::ffff:192.168.0.1",
			want: netip.MustParseAddr("192.168.0.1"),
		},
		{
			name: "invalid",
			ip:   "invalid_ip",
			want: netip.Addr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if got := IPToAddr(ip); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPToAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
