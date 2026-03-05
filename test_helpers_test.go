// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"path/filepath"
	"runtime"
	"testing"
)

// fixturePath resolves test fixture path relative to this package.
func fixturePath(parts ...string) string {
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)

	all := make([]string, 0, len(parts)+1)
	all = append(all, base)
	all = append(all, parts...)

	return filepath.Join(all...)
}

// mustParseCSVFixture parses fixture CSV file or fails test.
func mustParseCSVFixture(t *testing.T, name string) *Table {
	t.Helper()

	table, err := ParseCSVFile(fixturePath("testdata", "csv", name))
	if err != nil {
		t.Fatalf("ParseCSVFile(%q) error: %v", name, err)
	}

	return table
}
