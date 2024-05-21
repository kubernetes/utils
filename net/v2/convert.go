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
	"strings"
)

// AddrFromIP converts a net.IP to a netip.Addr. Given valid input this will always
// succeed; it will return the invalid netip.Addr on nil or garbage input.
//
// Use this rather than netip.AddrFromSlice(), which (despite the claims of its
// documentation) does not always do what you would expect if you pass it a net.IP.
func AddrFromIP(ip net.IP) netip.Addr {
	// Naively using netip.AddrFromSlice() gives unexpected results:
	//
	//   ip := net.ParseIP("1.2.3.4")
	//   addr, _ := netip.AddrFromSlice(ip)
	//   addr.String()  =>  "::ffff:1.2.3.4"
	//   addr.Is4()     =>  false
	//   addr.Is6()     =>  true
	//
	// This is because net.IP and netip.Addr have different ideas about how to handle
	// "IPv4-mapped IPv6" addresses, but netip.AddrFromSlice ignores that fact.
	//
	// In net.IP, parsing either "1.2.3.4" or "::ffff:1.2.3.4", will give you the
	// same result:
	//
	//   ip1 := net.ParseIP("1.2.3.4")
	//   []byte(ip1)   =>  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, 1, 2, 3, 4}
	//   ip1.String()  =>  "1.2.3.4"
	//   ip2 := net.ParseIP("::ffff:1.2.3.4")
	//   []byte(ip2)   =>  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, 1, 2, 3, 4}
	//   ip2.String()  =>  "1.2.3.4"
	//
	// net.IP normally stores IPv4 addresses as 16-byte IPv4-mapped IPv6 addresses,
	// but it hides that from the user, and it never stringifies an IPv4 IP to an
	// IPv4-mapped IPv6 form, even if that was the format you started with.
	//
	// net.IP *can* represent IPv4 addresses in a 4-byte format, but this is treated
	// as completly equivalent to the 16-byte representation:
	//
	//   ip4 := ip1.To4()
	//   []byte(ip4)     =>  []byte{1, 2, 3, 4}
	//   ip4.String()    =>  "1.2.3.4"
	//   ip1.Equal(ip4)  =>  true
	//
	// netip.Addr, on the other hand, treats "plain" IPv4 and IPv4-mapped IPv6 as two
	// completely separate things:
	//
	//   a1 := netip.MustParseAddr("1.2.3.4")
	//   a2 := netip.MustParseAddr("::ffff:1.2.3.4")
	//   a1.String()  =>  "1.2.3.4"
	//   a2.String()  =>  "::ffff:1.2.3.4"
	//   a1 == a2     =>  false
	//
	// which would be fine, except that netip.AddrFromSlice breaks net.IP's normal
	// semantics by converting the 4-byte and 16-byte net.IP forms to different
	// netip.Addr values, giving the confusing results above.
	//
	// In order to correctly convert an IPv4 address from net.IP to netip.Addr, you
	// need to either call .To4() on it before converting, or call .Unmap() on it
	// after converting. (The latter option is slightly simpler for us here because we
	// can just do it unconditionally, since it's a no-op in the IPv6 and invalid
	// cases).

	addr, _ := netip.AddrFromSlice(ip)
	return addr.Unmap()
}

// IPFromAddr converts a netip.Addr to a net.IP. Given valid input this will always
// succeed; it will return nil if addr is the invalid netip.Addr.
func IPFromAddr(addr netip.Addr) net.IP {
	// addr.AsSlice() returns:
	//   - a []byte of length 4 if addr is a normal IPv4 address
	//   - a []byte of length 16 if addr is an IPv6 address (including IPv4-mapped IPv6)
	//   - nil if addr is the zero Addr (which is the only other possibility)
	//
	// Any of those values can be correctly cast directly to a net.IP.
	//
	// Note that we don't bother to do any "cleanup" here like in the AddrFromIP case,
	// so converting a plain IPv4 netip.Addr to net.IP gives a different result than
	// converting an IPv4-mapped IPv6 netip.Addr:
	//
	//   ip1 := netutils.IPFromAddr(netip.MustParseAddr("1.2.3.4"))
	//   []byte(ip1)  =>  []byte{1, 2, 3, 4}
	//
	//   ip2 := netutils.IPFromAddr(netip.MustParseAddr("::ffff:1.2.3.4"))
	//   []byte(ip2)  =>  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, 1, 2, 3, 4}
	//
	// However, the net.IP API treats the two values as the same anyway, so it doesn't
	// matter.
	//
	//   ip1.String()    =>  "1.2.3.4"
	//   ip2.String()    =>  "1.2.3.4"
	//   ip2.Equal(ip1)  =>  true

	return net.IP(addr.AsSlice())
}

// IPFromInterfaceAddr can be used to extract the underlying IP address value (as a
// net.IP) from the return values of net.InterfaceAddrs(), net.Interface.Addrs(), or
// net.Interface.MulticastAddrs(). (net.Addr is also used in some other APIs, but this
// function should not be used on net.Addrs that are not "interface addresses".)
func IPFromInterfaceAddr(ifaddr net.Addr) net.IP {
	// On both Linux and Windows, the values returned from the "interface address"
	// methods are currently *net.IPNet for unicast addresses or *net.IPAddr for
	// multicast addresses.
	if ipnet, ok := ifaddr.(*net.IPNet); ok {
		return ipnet.IP
	} else if ipaddr, ok := ifaddr.(*net.IPAddr); ok {
		return ipaddr.IP
	}

	// Try to deal with other similar types... in particular, this is needed for
	// some existing unit tests...
	addrStr := ifaddr.String()
	// If it has a subnet length (like net.IPNet) or optional zone identifier (like
	// net.IPAddr), trim that away.
	if end := strings.IndexAny(addrStr, "/%"); end != -1 {
		addrStr = addrStr[:end]
	}
	// What's left is either an IP address, or something we can't parse.
	return ParseIPSloppy(addrStr)
}

// AddrFromInterfaceAddr can be used to extract the underlying IP address value (as a
// netip.Addr) from the return values of net.InterfaceAddrs(), net.Interface.Addrs(), or
// net.Interface.MulticastAddrs(). (net.Addr is also used in some other APIs, but this
// function should not be used on net.Addrs that are not "interface addresses".)
func AddrFromInterfaceAddr(ifaddr net.Addr) netip.Addr {
	return AddrFromIP(IPFromInterfaceAddr(ifaddr))
}

// PrefixFromIPNet converts a *net.IPNet to a netip.Prefix. Given valid input this will
// always succeed; it will return the invalid netip.Prefix on nil or garbage input.
func PrefixFromIPNet(ipnet *net.IPNet) netip.Prefix {
	if ipnet == nil {
		return netip.Prefix{}
	}

	addr := AddrFromIP(ipnet.IP)
	if !addr.IsValid() {
		return netip.Prefix{}
	}

	prefixLen, bits := ipnet.Mask.Size()
	if prefixLen == 0 && bits == 0 {
		// non-CIDR Mask representation; not representible as a netip.Prefix
		return netip.Prefix{}
	}
	if bits == 128 && addr.Is4() && (bits-prefixLen <= 32) {
		// In the same way that net.IP allows an IPv4 IP to be either 4 or 16
		// bytes (32 or 128 bits), *net.IPNet allows an IPv4 CIDR to have either a
		// 32-bit or a 128-bit mask. If the mask is 128 bits, we discard the
		// leftmost 96 bits.
		prefixLen -= 128 - 32
	} else if bits != addr.BitLen() {
		// invalid IPv4/IPv6 mix
		return netip.Prefix{}
	}

	return netip.PrefixFrom(addr, prefixLen)
}

// IPNetFromPrefix converts a netip.Prefix to a *net.IPNet. Given valid input this will
// always succeed; it will return nil if prefix is the invalid netip.Prefix or is
// otherwise invalid.
func IPNetFromPrefix(prefix netip.Prefix) *net.IPNet {
	addr := prefix.Addr()
	bits := prefix.Bits()
	if bits == -1 || !addr.IsValid() {
		return nil
	}
	addrLen := addr.BitLen()

	// (As with IPFromAddr, a plain IPv4 netip.Prefix and an equivalent IPv4-mapped
	// IPv6 netip.Prefix will get converted to distinct *net.IPNet values, but
	// *net.IPNet will treat them equivalently.)

	return &net.IPNet{
		IP:   IPFromAddr(addr),
		Mask: net.CIDRMask(bits, addrLen),
	}
}
