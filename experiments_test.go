// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

package addrto6

import (
	"encoding/binary"
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
