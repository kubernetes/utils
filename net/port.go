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

import (
	"fmt"
	"net"
	"strconv"

	"k8s.io/klog"
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
	// For ports on node IPs, open the actual port and hold it, even though we
	// use ipvs or iptables to redirect traffic.
	// This ensures a) that it's safe to use that port and b) that (a) stays
	// true.  The risk is that some process on the node (e.g. sshd or kubelet)
	// is using a port and we give that same port out to a Service.  That would
	// be bad because ipvs or iptables would silently claim the traffic
	// but the process would never know.
	// NOTE: We should not need to have a real listen()ing socket - bind()
	// should be enough, but I can't figure out a way to e2e test without
	// it.  Tools like 'ss' and 'netstat' do not show sockets that are
	// bind()ed but not listen()ed, and at least the default debian netcat
	// has no way to avoid about 10 seconds of retries.
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
	default:
		return nil, fmt.Errorf("unknown protocol %q", lp.Protocol)
	}
	klog.V(2).Infof("Opened local port %s", lp.String())
	return socket, nil
}
