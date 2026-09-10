// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

package addrto6

import "net/netip"

// AddrTo6Safe maps an IPv4 address to an IPv4-mapped IPv6 address. It returns
// an IPv6 address unmodified. This is the safest but slowest version.
func AddrTo6Safe(ip netip.Addr) netip.Addr {
	if ip.Is4() {
		return netip.AddrFrom16(ip.As16())
	}
	return ip
}
