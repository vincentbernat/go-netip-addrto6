// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

package addrto6

import (
	"net/netip"
	"reflect"
	"testing"
	"unsafe"
)

type fieldLayout struct {
	Offset uintptr
	Size   uintptr
	Align  int
}

func layoutOf(t reflect.Type) []fieldLayout {
	layout := make([]fieldLayout, t.NumField())
	for i := range layout {
		field := t.Field(i)
		layout[i] = fieldLayout{field.Offset, field.Type.Size(), field.Type.Align()}
	}
	return layout
}

// TestNetIPAddrStructure checks the layout of netip.Addr matches what we expect.
func TestNetIPAddrStructure(t *testing.T) {
	type field struct {
		Name   string
		Type   string
		Offset uintptr
		Size   uintptr
	}
	addrType := reflect.TypeFor[netip.Addr]()
	got := make([]field, addrType.NumField())
	for i := range got {
		f := addrType.Field(i)
		got[i] = field{f.Name, f.Type.String(), f.Offset, f.Type.Size()}
	}
	expected := []field{
		{"addr", "netip.uint128", 0, 16},
		{"z", "unique.Handle[net/netip.addrDetail]", 16, 8},
	}
	if !AssertEqual(t, got, expected, "netip.Addr fields") {
		return
	}
	AssertEqual(t, addrType.Size(), uintptr(24), "netip.Addr size")

	// addr is an uint128 made of two uint64.
	uint128Type := addrType.Field(0).Type
	for i := range uint128Type.NumField() {
		f := uint128Type.Field(i)
		AssertEqual(t, f.Type.Kind(), reflect.Uint64, "netip.uint128 field %q kind", f.Name)
	}
	// z is a pointer
	zType := addrType.Field(1).Type
	if !AssertEqual(t, zType.NumField(), 1, "unique.Handle field count") {
		return
	}
	AssertEqual(t, zType.Field(0).Type.Kind(), reflect.Pointer, "unique.Handle field kind")
}

// TestAddrProxy checks addrProxy matches netip.Addr structure
func TestAddrProxy(t *testing.T) {
	addrType := reflect.TypeFor[netip.Addr]()
	proxyType := reflect.TypeFor[addrProxy]()

	AssertEqual(t, proxyType.Size(), addrType.Size(), "addrProxy size")
	AssertEqual(t, layoutOf(proxyType), layoutOf(addrType), "addrProxy layout")
	AssertEqual(t, proxyType.Field(1).Type.Kind(), reflect.UnsafePointer, "addrProxy z kind")

	// netipZ6noz must be the handle netip gives to an IPv6 address without a zone.
	ipv6 := netip.MustParseAddr("2001:db8::1")
	AssertEqual(t, (*addrProxy)(unsafe.Pointer(&ipv6)).z, netipZ6noz, "netipZ6noz")
}
