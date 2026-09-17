// SPDX-FileCopyrightText: 2026 Free Mobile
// SPDX-License-Identifier: BSD-3-Clause

//go:build !go1.26

package addrto6

// Up to Go 1.25, the unsafe version of AddrTo6 has not the same cost as the
// builtin one.
const (
	unsafeIsOptimal  = false
	goVersionChecked = true
)
