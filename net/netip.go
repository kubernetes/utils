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
	"net/netip"
)

// IPToAddr converts a net.IP to a netip.Addr
// if the net.IP is not valid it returns an empty netip.Addr{}
func IPToAddr(ip net.IP) netip.Addr {
	// https://pkg.go.dev/net/netip#AddrFromSlice can return an IPv4 in IPv6 format
	// so we have to check the IP family to return exactly the format that we want
	// address, _ := netip.AddrFromSlice(net.ParseIPSloppy(192.168.0.1)) returns
	// an address like ::ffff:192.168.0.1/32
	bytes := ip.To4()
	if bytes == nil {
		bytes = ip.To16()
	}
	// AddrFromSlice returns Addr{}, false if the input is invalid.
	address, _ := netip.AddrFromSlice(bytes)
	return address
}

// BroadcastAddress returns the broadcast address of the subnet
// The broadcast address is obtained by setting all the host bits
// in a subnet to 1.
// network 192.168.0.0/24 : subnet bits 24 host bits 32 - 24 = 8
// broadcast address 192.168.0.255
func BroadcastAddress(subnet netip.Prefix) (netip.Addr, error) {
	base := subnet.Masked().Addr()
	bytes := base.AsSlice()
	// get all the host bits from the subnet
	n := 8*len(bytes) - subnet.Bits()
	// set all the host bits to 1
	for i := len(bytes) - 1; i >= 0 && n > 0; i-- {
		if n >= 8 {
			bytes[i] = 0xff
			n -= 8
		} else {
			mask := ^uint8(0) >> (8 - n)
			bytes[i] |= mask
			break
		}
	}

	addr, ok := netip.AddrFromSlice(bytes)
	if !ok {
		return netip.Addr{}, fmt.Errorf("invalid address %v", bytes)
	}
	return addr, nil
}
