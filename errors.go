// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import "errors"

var (
	// ErrNilTable indicates nil table argument.
	ErrNilTable = errors.New("table is nil")

	// ErrNilCatalog indicates nil PO catalog argument.
	ErrNilCatalog = errors.New("catalog is nil")

	// ErrInvalidHeader indicates invalid CSV header shape.
	ErrInvalidHeader = errors.New("invalid stringtable csv header")

	// ErrDuplicateKey indicates duplicate key in CSV rows.
	ErrDuplicateKey = errors.New("duplicate key")
)
