// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

package addrto6

import (
	"net/netip"
	"unsafe"
)

// addrProxy has the same layout as netip.Addr. TestNetIPAddrStructure checks
// this is still true.
type addrProxy struct {
	addr [2]uint64      // netip.uint128
	z    unsafe.Pointer // unique.Handle[netip.addrDetail]
}

var (
	anyIPv6    = netip.IPv6Unspecified()
	netipZ6noz = (*addrProxy)(unsafe.Pointer(&anyIPv6)).z
)

// AddrTo6Unsafe maps an IPv4 address to an IPv4-mapped IPv6 address. It returns an
// IPv6 address unmodified. netip already stores an IPv4 address as
// ::ffff:a.b.c.d, so only the family marker has to change. This is unsafe, but
// there is a test to ensure netip.Addr is like we expect. Copying a
// unique.Handle bypasses the unique package bookkeeping, but netipZ6noz lives for
// the whole program.
func AddrTo6Unsafe(ip netip.Addr) netip.Addr {
	if !ip.Is4() {
		return ip
	}
	p := *(*addrProxy)(unsafe.Pointer(&ip))
	p.z = netipZ6noz
	return *(*netip.Addr)(unsafe.Pointer(&p))
}
