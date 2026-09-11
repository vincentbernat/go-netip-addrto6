// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

package addrto6

import (
	"encoding/binary"
	"fmt"
	"net/netip"
	"reflect"
	"testing"
	"unique"
)

// AssertEqual fails the test when got and want differ. It returns true when
// they are equal, so a caller can stop early.
func AssertEqual[T any](t *testing.T, got, want T, format string, args ...any) bool {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		return true
	}
	t.Errorf("%s (-got, +want):\n-%+v\n+%+v", fmt.Sprintf(format, args...), got, want)
	return false
}

// FromNetipAddr copies a netip.Addr into an Addr.
func FromNetipAddr(parsed netip.Addr) Addr {
	if !parsed.IsValid() {
		return Addr{}
	}
	b := parsed.As16()
	ip := Addr{addr: uint128{
		binary.BigEndian.Uint64(b[:8]),
		binary.BigEndian.Uint64(b[8:]),
	}}
	switch {
	case parsed.Is4():
		ip.z = z4
	case parsed.Zone() == "":
		ip.z = z6noz
	default:
		ip.z = unique.Make(addrDetail{isV6: true, zoneV6: parsed.Zone()})
	}
	return ip
}

// AddrTo6Builtin wraps the builtin To6 to take and return a netip.Addr, like
// the other implementations.
func AddrTo6Builtin(ip netip.Addr) netip.Addr {
	return FromNetipAddr(ip).To6().toNetipAddr()
}

// MustParseAddr parses s as an Addr and panics on error. The builtin type has
// no parser, so netip parses the address and the bytes are copied over.
func MustParseAddr(s string) Addr {
	return FromNetipAddr(netip.MustParseAddr(s))
}

// TestMustParseAddr checks the bits taken from netip give back the same address.
func TestMustParseAddr(t *testing.T) {
	cases := []string{
		"0.0.0.0",
		"192.168.1.1",
		"255.255.255.255",
		"::",
		"::ffff:192.168.1.1",
		"2a01:db8::1",
		"fe80::1%eth0",
	}
	for _, tc := range cases {
		got := MustParseAddr(tc)
		want := netip.MustParseAddr(tc)
		AssertEqual(t, got.String(), want.String(), "MustParseAddr(%q).String()", tc)
		AssertEqual(t, got.Is4(), want.Is4(), "MustParseAddr(%q).Is4()", tc)
	}
}

func TestAddrTo6(t *testing.T) {
	// An empty string stands for the zero value, an invalid address.
	cases := []struct {
		input  string
		output string
	}{
		{"", ""},
		{"0.0.0.0", "::ffff:0.0.0.0"},
		{"192.168.1.1", "::ffff:192.168.1.1"},
		{"255.255.255.255", "::ffff:255.255.255.255"},
		{"::ffff:192.168.1.1", "::ffff:192.168.1.1"},
		{"::", "::"},
		{"2a01:db8::1", "2a01:db8::1"},
		{"fe80::1%eth0", "fe80::1%eth0"},
	}
	parse := func(s string) netip.Addr {
		if s == "" {
			return netip.Addr{}
		}
		return netip.MustParseAddr(s)
	}
	implementations := []struct {
		name    string
		addrTo6 func(netip.Addr) netip.Addr
	}{
		{"safe", AddrTo6Safe},
		{"unsafe", AddrTo6Unsafe},
		{"builtin", AddrTo6Builtin},
	}
	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			for _, tc := range cases {
				input, want := parse(tc.input), parse(tc.output)
				got := implementation.addrTo6(input)
				AssertEqual(t, got, want, "AddrTo6(%s)", input)
				AssertEqual(t, got.String(), want.String(),
					"AddrTo6(%s).String()", input)
			}
		})
	}
}

// The inputs are package-level and public on purpose. With a local, the
// compiler optimizes part of the work by putting it outside b.Loop().
var (
	BenchIPv4        = netip.MustParseAddr("192.168.1.1")
	BenchBuiltinIPv4 = MustParseAddr("192.168.1.1")
)

func BenchmarkAddrTo6(b *testing.B) {
	b.Run("safe", func(b *testing.B) {
		for b.Loop() {
			_ = AddrTo6Safe(BenchIPv4)
		}
	})
	b.Run("unsafe", func(b *testing.B) {
		for b.Loop() {
			_ = AddrTo6Unsafe(BenchIPv4)
		}
	})
	b.Run("builtin", func(b *testing.B) {
		for b.Loop() {
			_ = BenchBuiltinIPv4.To6()
		}
	})
	b.Run("do nothing", func(b *testing.B) {
		for b.Loop() {
		}
	})
}
