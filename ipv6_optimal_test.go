// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

//go:build !race

package addrto6

import (
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"golang.org/x/arch/x86/x86asm"
)

// instructionCount disassembles the compiled body of fn and counts the
// instructions. The padding between functions is not counted.
func instructionCount(t *testing.T, fn any) int {
	t.Helper()
	pc := reflect.ValueOf(fn).Pointer()

	// Walk forward until we leave fn to get its size
	var size uintptr
	for {
		owner := runtime.FuncForPC(pc + size)
		if owner == nil || owner.Entry() != pc {
			break
		}
		size++
	}

	code := unsafe.Slice((*byte)(unsafe.Pointer(pc)), size)
	count, padding := 0, 0
	for len(code) > 0 {
		inst, err := x86asm.Decode(code, 64)
		if err != nil {
			t.Fatalf("x86asm.Decode() error:\n%+v", err)
		}
		code = code[inst.Len:]
		if inst.Op == x86asm.INT {
			// Functions are padded with INT3. Only count them if more code
			// follows.
			padding++
			continue
		}
		count += padding + 1
		padding = 0
	}
	return count
}

// TestAddrTo6Optimal checks AddrTo6Unsafe and AddrTo6Safe compile to as few
// instructions as a builtin To6 method. This is version-dependent.
func TestAddrTo6Optimal(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skipf("no disassembler for %s", runtime.GOARCH)
	}
	if testing.CoverMode() != "" {
		t.Skip("instrumented build")
	}

	builtin := instructionCount(t, Addr.To6)
	withUnsafe := instructionCount(t, AddrTo6Unsafe)
	withoutUnsafe := instructionCount(t, AddrTo6Safe)
	t.Logf("instructions: builtin %d, unsafe %d, safe %d",
		builtin, withUnsafe, withoutUnsafe)

	if unsafeIsOptimal {
		AssertEqual(t, withUnsafe, builtin, "AddrTo6Unsafe() instruction count")
	} else if withUnsafe >= withoutUnsafe {
		t.Errorf("AddrTo6Unsafe() instruction count: %d, expected less than %d",
			withUnsafe, withoutUnsafe)
	}
	if safeIsOptimal {
		AssertEqual(t, withoutUnsafe, builtin, "AddrTo6Safe() instruction count")
	} else if withoutUnsafe <= builtin {
		t.Errorf("AddrTo6Safe() instruction count: %d, expected more than %d",
			withoutUnsafe, builtin)
	}
	if !goVersionChecked {
		t.Errorf("%s is newer than the last checked Go version: AddrTo6Safe may be optimal now",
			runtime.Version())
	}
}
