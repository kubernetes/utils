/*
Copyright 2020 The Kubernetes Authors.

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
	"strconv"
)

// IPFamily refers to a specific family if not empty, i.e. "4" or "6"
type IPFamily string

// Constants refering to IPv4 and IPv6
const (
	IPv4 IPFamily = "4"
	IPv6          = "6"
)

// LocalPort describes a port on specific IP address and protocol
type LocalPort struct {
	// Description is the identity message of a given local port.
	Description string
	// IP is the IP address part of a given local port.
	// If this string is empty, the port binds to all local IP addresses.
	IP string
	// If IPFamily is not empty, the port binds only to addresses of this family
	IPFamily IPFamily
	// Port is the port part of a given local port.
	Port int
	// Protocol is the protocol part of a given local port.
	// The value is assumed to be lower-case. For example, "udp" not "UDP", "tcp" not "TCP".
	Protocol string
}

// NewLocalPort creates a new LocalPort struct
func NewLocalPort(desc, ip string, ipFamily IPFamily, port int, protocol string) (*LocalPort, error) {
	if protocol != "tcp" && protocol != "sctp" && protocol != "udp" {
		return nil, fmt.Errorf("Unsupported protocol %s", protocol)
	}
	if ipFamily != "" && ipFamily != "4" && ipFamily != "6" {
		return nil, fmt.Errorf("Invalid IP family %s", ipFamily)
	}
	if ip != "" {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return nil, fmt.Errorf("invalid ip address %s", ip)
		}
		asIPv4 := parsedIP.To4()
		if asIPv4 == nil && ipFamily == IPv4 || asIPv4 != nil && ipFamily == IPv6 {
			return nil, fmt.Errorf("ip address and family mismatch %s, %s", ip, ipFamily)
		}
	}
	return &LocalPort{Description: desc, IP: ip, IPFamily: ipFamily, Port: port, Protocol: protocol}, nil
}

func (lp *LocalPort) String() string {
	ipPort := net.JoinHostPort(lp.IP, strconv.Itoa(lp.Port))
	return fmt.Sprintf("%q (%s/%s%s)", lp.Description, ipPort, lp.Protocol, lp.IPFamily)
}

// Closeable is an interface around closing a port.
type Closeable interface {
	Close() error
}

// PortOpener is an interface around port opening/closing.
// Abstracted out for testing.
type PortOpener interface {
	OpenLocalPort(lp *LocalPort) (Closeable, error)
}

// listenPortOpener opens ports by calling bind() and listen().
type listenPortOpener struct{}

// OpenLocalPort holds the given local port open.
func (l *listenPortOpener) OpenLocalPort(lp *LocalPort) (Closeable, error) {
	return openLocalPort(lp)
}

func openLocalPort(lp *LocalPort) (Closeable, error) {
	var socket Closeable
	network := lp.Protocol + string(lp.IPFamily)
	hostPort := net.JoinHostPort(lp.IP, strconv.Itoa(lp.Port))
	switch lp.Protocol {
	case "tcp":
		listener, err := net.Listen(network, hostPort)
		if err != nil {
			return nil, err
		}
		socket = listener
	case "udp":
		addr, err := net.ResolveUDPAddr(network, hostPort)
		if err != nil {
			return nil, err
		}
		conn, err := net.ListenUDP(network, addr)
		if err != nil {
			return nil, err
		}
		socket = conn
	case "sctp":
		// SCTP ports are intentionally ignored, to ensure we don't cause the sctp
		// kernel module to be loaded, which breaks userspace SCTP support (and
		// may be considered a security risk by some administrators).
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown protocol %q", lp.Protocol)
	}
	return socket, nil
}
