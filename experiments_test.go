// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

package addrto6

import (
	"encoding/binary"
	"math/rand/v2"
	"testing"
)

func foldingBEUint64(x uint64) uint64 {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], x)
	return binary.BigEndian.Uint64(b[:])
}

func foldingNativeUint64(x uint64) uint64 {
	var b [8]byte
	binary.NativeEndian.PutUint64(b[:], x)
	return binary.NativeEndian.Uint64(b[:])
}

func foldingNoAllocUint64(x uint64) uint64 {
	b0 := byte(x >> 56)
	b1 := byte(x >> 48)
	b2 := byte(x >> 40)
	b3 := byte(x >> 32)
	b4 := byte(x >> 24)
	b5 := byte(x >> 16)
	b6 := byte(x >> 8)
	b7 := byte(x)
	return uint64(b0)<<56 | uint64(b1)<<48 | uint64(b2)<<40 | uint64(b3)<<32 |
		uint64(b4)<<24 | uint64(b5)<<16 | uint64(b6)<<8 | uint64(b7)
}

func foldingNoAllocUint64Optimized(x uint64) uint64 {
	return x&0xff00000000000000 |
		uint64(byte(x>>48))<<48 | uint64(byte(x>>40))<<40 | uint64(byte(x>>32))<<32 |
		uint64(byte(x>>24))<<24 | uint64(byte(x>>16))<<16 | uint64(byte(x>>8))<<8 |
		uint64(byte(x))
}

func foldingBEUint32(x uint32) uint32 {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], x)
	return binary.BigEndian.Uint32(b[:])
}

func foldingNativeUint32(x uint32) uint32 {
	var b [4]byte
	binary.NativeEndian.PutUint32(b[:], x)
	return binary.NativeEndian.Uint32(b[:])
}

func foldingNoAllocUint32(x uint32) uint32 {
	b0 := byte(x >> 24)
	b1 := byte(x >> 16)
	b2 := byte(x >> 8)
	b3 := byte(x)
	return uint32(b0)<<24 | uint32(b1)<<16 | uint32(b2)<<8 | uint32(b3)
}

func foldingNoAllocUint32Optimized(x uint32) uint32 {
	return x&0xff000000 | uint32(byte(x>>16))<<16 | uint32(byte(x>>8))<<8 | uint32(byte(x))
}

func foldingNoAllocUint32Optimized2(x uint32) uint32 {
	return x&0xff000000 | x&0xff0000 | x&0xff00 | x&0xff
}

func foldingBEUint16(x uint16) uint16 {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], x)
	return binary.BigEndian.Uint16(b[:])
}

func foldingNativeUint16(x uint16) uint16 {
	var b [2]byte
	binary.NativeEndian.PutUint16(b[:], x)
	return binary.NativeEndian.Uint16(b[:])
}

func foldingNoAllocUint16(x uint16) uint16 {
	b0 := byte(x >> 8)
	b1 := byte(x)
	return uint16(b0)<<8 | uint16(b1)
}

func foldingNoAllocUint16Optimized(x uint16) uint16 {
	return x&0xff00 | uint16(byte(x))
}

func foldingNoAllocUint16FullyOptimized(x uint16) uint16 {
	return x
}

// Each folding function should give back its input. The tests check this on a
// few handpicked values and on random ones. They stop at the first mismatch to
// keep the output short.

func TestFoldingUint64(t *testing.T) {
	values := []uint64{
		0, 1, ^uint64(0), 1 << 32, 1 << 63, 0xffffffff,
		0x0123456789abcdef, 0xff00ff00ff00ff00, 0x00ff00ff00ff00ff,
	}
	r := rand.New(rand.NewPCG(1, 2))
	for range 1000 {
		values = append(values, r.Uint64())
	}
	implementations := []struct {
		name string
		fold func(uint64) uint64
	}{
		{"BE", foldingBEUint64},
		{"native", foldingNativeUint64},
		{"no alloc", foldingNoAllocUint64},
		{"no alloc optimized", foldingNoAllocUint64Optimized},
	}
	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			for _, x := range values {
				if !AssertEqual(t, implementation.fold(x), x, "fold(%#016x)", x) {
					return
				}
			}
		})
	}
}

func TestFoldingUint32(t *testing.T) {
	values := []uint32{
		0, 1, ^uint32(0), 1 << 31, 0xffff,
		0x01234567, 0xff00ff00, 0x00ff00ff,
	}
	r := rand.New(rand.NewPCG(1, 2))
	for range 1000 {
		values = append(values, r.Uint32())
	}
	implementations := []struct {
		name string
		fold func(uint32) uint32
	}{
		{"BE", foldingBEUint32},
		{"native", foldingNativeUint32},
		{"no alloc", foldingNoAllocUint32},
		{"no alloc optimized", foldingNoAllocUint32Optimized},
		{"no alloc tweaked", foldingNoAllocUint32Optimized2},
	}
	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			for _, x := range values {
				if !AssertEqual(t, implementation.fold(x), x, "fold(%#08x)", x) {
					return
				}
			}
		})
	}
}

func TestFoldingUint16(t *testing.T) {
	implementations := []struct {
		name string
		fold func(uint16) uint16
	}{
		{"BE", foldingBEUint16},
		{"native", foldingNativeUint16},
		{"no alloc", foldingNoAllocUint16},
		{"no alloc optimized", foldingNoAllocUint16Optimized},
		{"no alloc fully optimized", foldingNoAllocUint16FullyOptimized},
	}
	for _, implementation := range implementations {
		t.Run(implementation.name, func(t *testing.T) {
			// There are only 65536 values, so check them all.
			for i := range 1 << 16 {
				x := uint16(i)
				if !AssertEqual(t, implementation.fold(x), x, "fold(%#04x)", x) {
					return
				}
			}
		})
	}
}
