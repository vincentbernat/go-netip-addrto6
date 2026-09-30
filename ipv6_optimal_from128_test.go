// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.28

package addrto6

// From Go 1.28, AddrTo6Safe is as short as a builtin To6, like AddrTo6Unsafe.
// This was only checked with a patched development version.
const (
	unsafeIsOptimal  = true
	safeIsOptimal    = true
	goVersionChecked = true
)
