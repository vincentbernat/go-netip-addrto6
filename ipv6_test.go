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

// MustParseAddr parses s as an Addr and panics on error. The builtin type has
// no parser, so netip parses the address and the bytes are copied over.
func MustParseAddr(s string) Addr {
	parsed := netip.MustParseAddr(s)
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
	parsers := map[reflect.Type]func(string) reflect.Value{
		reflect.TypeFor[netip.Addr](): func(s string) reflect.Value {
			return reflect.ValueOf(netip.MustParseAddr(s))
		},
		reflect.TypeFor[Addr](): func(s string) reflect.Value {
			return reflect.ValueOf(MustParseAddr(s))
		},
	}
	// Each implementation works on its own address type, so the parser is
	// looked up from the type the function takes.
	implementations := []struct {
		name    string
		addrTo6 any
	}{
		{"safe", AddrTo6Safe},
		{"unsafe", AddrTo6Unsafe},
		{"builtin", Addr.To6},
	}
	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			addrTo6 := reflect.ValueOf(implementation.addrTo6)
			addrType := addrTo6.Type().In(0)
			parse := func(s string) reflect.Value {
				if s == "" {
					return reflect.Zero(addrType)
				}
				return parsers[addrType](s)
			}
			for _, tc := range cases {
				input, want := parse(tc.input), parse(tc.output)
				got := addrTo6.Call([]reflect.Value{input})[0]
				AssertEqual(t, got.Interface(), want.Interface(),
					"AddrTo6(%s)", input)
				AssertEqual(t,
					got.Interface().(fmt.Stringer).String(),
					want.Interface().(fmt.Stringer).String(),
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
