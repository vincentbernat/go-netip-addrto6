// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.26 && !go1.28

package addrto6

// From Go 1.26 to Go 1.27, AddrTo6Unsafe is as short as a builtin To6 while
// AddrTo6Safe is much longer.
const (
	unsafeIsOptimal  = true
	goVersionChecked = true
)
