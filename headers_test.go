// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"testing"

	"github.com/woozymasta/pofile"
)

func TestUpdateBuildHeadersWriteHashAndDateOnChange(t *testing.T) {
	t.Parallel()

	catalog := pofile.NewCatalog()
	catalog.Language = "russian"
	catalog.UpsertMessage("UI_OK", "OK", "Ок")

	options := HeaderOptions{
		WriteStandardHeaders: true,
		WriteHash:            true,
		WriteDateOnChange:    true,
		Generator:            "stringtable-test",
	}

	firstChanged := UpdateBuildHeaders(catalog, options)
	if !firstChanged {
		t.Fatal("first update should report changed")
	}
	firstDate := catalog.Header("PO-Revision-Date")
	if firstDate == "" {
		t.Fatal("PO-Revision-Date should be set")
	}
	firstHash := catalog.Header(DefaultContentHashHeader)
	if firstHash == "" {
		t.Fatal("X-Content-Hash should be set")
	}
	if got := catalog.Header("Language"); got != "russian" {
		t.Fatalf("Language header = %q, want %q", got, "russian")
	}

	secondChanged := UpdateBuildHeaders(catalog, options)
	if secondChanged {
		t.Fatal("second update should report unchanged content")
	}
	if got := catalog.Header("PO-Revision-Date"); got != firstDate {
		t.Fatalf("PO-Revision-Date changed for unchanged content: %q vs %q", firstDate, got)
	}
	if got := catalog.Header(DefaultContentHashHeader); got != firstHash {
		t.Fatalf("X-Content-Hash changed for unchanged content: %q vs %q", firstHash, got)
	}
}

func TestUpdateBuildHeadersTemplateDate(t *testing.T) {
	t.Parallel()

	catalog := pofile.NewCatalog()
	catalog.UpsertMessage("UI_OK", "OK", "")

	changed := UpdateBuildHeaders(catalog, HeaderOptions{
		WriteHash:         true,
		WriteDateOnChange: true,
		Template:          true,
	})
	if !changed {
		t.Fatal("template update should report changed")
	}
	if catalog.Header("POT-Creation-Date") == "" {
		t.Fatal("POT-Creation-Date should be set")
	}
	if catalog.Header("PO-Revision-Date") != "" {
		t.Fatal("PO-Revision-Date should stay empty for template mode")
	}
}
