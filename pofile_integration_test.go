// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/woozymasta/pofile"
)

func TestCatalogFromTableLanguage(t *testing.T) {
	t.Parallel()

	table := &Table{
		Languages: []string{"english", "russian"},
		Rows: []Row{
			{
				Key:      "UI_OK",
				Original: "OK",
				Translations: map[string]string{
					"english": "OK",
					"russian": "Ок",
				},
			},
		},
	}

	catalog, err := CatalogFromTableLanguage(table, "russian")
	if err != nil {
		t.Fatalf("CatalogFromTableLanguage error: %v", err)
	}
	if got := catalog.Header("Language"); got != "russian" {
		t.Fatalf("Language header = %q, want %q", got, "russian")
	}
	if got := catalog.Translation("UI_OK", "OK"); got != "Ок" {
		t.Fatalf("translation = %q, want %q", got, "Ок")
	}
}

func TestMergeCatalogLanguage(t *testing.T) {
	t.Parallel()

	table := &Table{
		Languages: []string{"english"},
		Rows: []Row{
			{
				Key:      "UI_OK",
				Original: "OK",
				Translations: map[string]string{
					"english": "OK",
				},
			},
			{
				Key:      "UI_CANCEL",
				Original: "Cancel",
				Translations: map[string]string{
					"english": "Cancel",
				},
			},
		},
	}

	catalog := pofile.NewCatalog()
	catalog.UpsertMessage("UI_OK", "OK", "")
	cancel := catalog.UpsertMessage("UI_CANCEL", "Cancel", "")
	cancel.Flags = []string{"notranslate"}

	err := MergeCatalogLanguage(table, "russian", catalog, MergeOptions{})
	if err != nil {
		t.Fatalf("MergeCatalogLanguage error: %v", err)
	}
	if got := table.Rows[0].Translations["russian"]; got != "OK" {
		t.Fatalf("row0 russian = %q, want %q", got, "OK")
	}
	if got := table.Rows[1].Translations["russian"]; got != "Cancel" {
		t.Fatalf("row1 russian = %q, want %q", got, "Cancel")
	}
}

func TestUpdateCatalogFromTablePreservesExistingMessageData(t *testing.T) {
	t.Parallel()

	table := &Table{
		Languages: []string{"russian"},
		Rows: []Row{
			{
				Key:      "UI_OK",
				Original: "OK",
				Translations: map[string]string{
					"russian": "Ок",
				},
			},
		},
	}

	existing := pofile.NewCatalog()
	message := existing.UpsertMessage("UI_OK", "OK", "Норм")
	message.Comments = []string{"# keep"}
	message.Flags = []string{"fuzzy"}
	message.References = []string{"ui.cpp:10"}

	updated, err := UpdateCatalogFromTable(table, "russian", existing)
	if err != nil {
		t.Fatalf("UpdateCatalogFromTable error: %v", err)
	}

	got := updated.FindMessage("UI_OK", "OK")
	if got == nil {
		t.Fatal("updated message is nil")
	}
	if got.TranslationAt(0) != "Норм" {
		t.Fatalf("translation = %q, want %q", got.TranslationAt(0), "Норм")
	}
	if len(got.Comments) != 1 || got.Comments[0] != "# keep" {
		t.Fatalf("comments = %#v, want [# keep]", got.Comments)
	}
	if len(got.Flags) != 1 || got.Flags[0] != "fuzzy" {
		t.Fatalf("flags = %#v, want [fuzzy]", got.Flags)
	}
}

func TestReadWriteCatalogSet(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	catalogs := map[string]*pofile.Catalog{
		"russian": func() *pofile.Catalog {
			catalog := pofile.NewCatalog()
			catalog.SetHeader("Language", "russian")
			catalog.UpsertMessage("UI_OK", "OK", "Ок")
			return catalog
		}(),
		"german": func() *pofile.Catalog {
			catalog := pofile.NewCatalog()
			catalog.SetHeader("Language", "german")
			catalog.UpsertMessage("UI_OK", "OK", "OK")
			return catalog
		}(),
	}

	if err := WriteCatalogSet(dir, catalogs); err != nil {
		t.Fatalf("WriteCatalogSet error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "russian.po")); err != nil {
		t.Fatalf("russian.po stat error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "german.po")); err != nil {
		t.Fatalf("german.po stat error: %v", err)
	}

	loaded, err := ReadCatalogSet(dir)
	if err != nil {
		t.Fatalf("ReadCatalogSet error: %v", err)
	}
	if got := loaded["russian"].Translation("UI_OK", "OK"); got != "Ок" {
		t.Fatalf("russian translation = %q, want %q", got, "Ок")
	}
	if got := loaded["german"].Translation("UI_OK", "OK"); got != "OK" {
		t.Fatalf("german translation = %q, want %q", got, "OK")
	}
}
