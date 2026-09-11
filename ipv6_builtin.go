// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

package addrto6

import (
	"encoding/binary"
	"net/netip"
	"unique"
)

// The types below mirror the internals of net/netip.
type (
	uint128 struct{ hi, lo uint64 }

	addrDetail struct {
		isV6   bool
		zoneV6 string
	}
)

// Addr represents an IPv4 or IPv6 address (with or without a scoped addressing
// zone), similar to [net.IP] or [net.IPAddr].
type Addr struct {
	addr uint128
	z    unique.Handle[addrDetail]
}

var (
	z4    = unique.Make(addrDetail{})
	z6noz = unique.Make(addrDetail{isV6: true})
)

// Is4 reports whether ip is an IPv4 address.
//
// It returns false for IPv4-mapped IPv6 addresses.
func (ip Addr) Is4() bool {
	return ip.z == z4
}

// To6 maps an IPv4 address to an IPv4-mapped IPv6 address. It returns an IPv6
// address unmodified. This version needs to be built into Go's standard library
// to be useful.
func (ip Addr) To6() Addr {
	if ip.Is4() {
		ip.z = z6noz
	}
	return ip
}

// toNetipAddr copies an Addr back into a netip.Addr.
func (ip Addr) toNetipAddr() netip.Addr {
	if ip.z == (unique.Handle[addrDetail]{}) {
		return netip.Addr{}
	}
	var b [16]byte
	binary.BigEndian.PutUint64(b[:8], ip.addr.hi)
	binary.BigEndian.PutUint64(b[8:], ip.addr.lo)
	if ip.z == z4 {
		return netip.AddrFrom16(b).Unmap()
	}
	return netip.AddrFrom16(b).WithZone(ip.z.Value().zoneV6)
}

// String returns the string form of the IP address ip. The real netip formats
// the address by hand, this version asks netip to do it.
func (ip Addr) String() string {
	return ip.toNetipAddr().String()
}
