// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.28

package addrto6

// Go 1.28 and later were never checked. A newer compiler may turn AddrTo6Safe
// into the same code as a builtin To6 and make the unsafe version useless.
// Let's hope this will be true!
const (
	unsafeIsOptimal  = true
	goVersionChecked = false
)
