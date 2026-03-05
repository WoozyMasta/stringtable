// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"strings"
	"testing"

	"github.com/woozymasta/pofile"
)

func TestParseCSVFixtures(t *testing.T) {
	t.Parallel()

	quoted := mustParseCSVFixture(t, "quoted.csv")
	if len(quoted.Rows) != 3 {
		t.Fatalf("quoted rows = %d, want 3", len(quoted.Rows))
	}
	if got := quoted.Rows[1].Original; got != `Say "Hello", user` {
		t.Fatalf("quoted original = %q, want %q", got, `Say "Hello", user`)
	}
	if got := quoted.Rows[1].Translations["russian"]; got != `Скажи "Привет", пользователь` {
		t.Fatalf("quoted russian = %q", got)
	}

	unquoted := mustParseCSVFixture(t, "unquoted.csv")
	if len(unquoted.Rows) != 2 {
		t.Fatalf("unquoted rows = %d, want 2", len(unquoted.Rows))
	}
	if got := unquoted.Rows[1].Translations["russian"]; got != "Выход" {
		t.Fatalf("unquoted russian = %q, want %q", got, "Выход")
	}
}

func TestParseCSVBrokenFixture(t *testing.T) {
	t.Parallel()

	_, err := ParseCSVFile(fixturePath("testdata", "csv", "broken.csv"))
	if err == nil {
		t.Fatal("ParseCSVFile(broken.csv) error = nil, want parse error")
	}
}

func TestRoundTripCSVToPOToCSV(t *testing.T) {
	t.Parallel()

	source := mustParseCSVFixture(t, "quoted.csv")
	formatted, err := FormatCSV(source)
	if err != nil {
		t.Fatalf("FormatCSV error: %v", err)
	}
	if !strings.Contains(string(formatted), "UI_HELLO") {
		t.Fatal("formatted csv misses expected key")
	}
	formattedTable, err := ParseCSV(formatted)
	if err != nil {
		t.Fatalf("ParseCSV(formatted) error: %v", err)
	}
	if len(formattedTable.Rows) != len(source.Rows) {
		t.Fatalf("formatted rows = %d, want %d", len(formattedTable.Rows), len(source.Rows))
	}

	catalog, err := CatalogFromTableLanguage(source, "russian")
	if err != nil {
		t.Fatalf("CatalogFromTableLanguage error: %v", err)
	}

	round := source.Clone()
	for index := range round.Rows {
		round.Rows[index].Translations["russian"] = ""
	}
	if err := MergeCatalogLanguage(round, "russian", catalog, MergeOptions{}); err != nil {
		t.Fatalf("MergeCatalogLanguage error: %v", err)
	}

	for index := range source.Rows {
		want := source.Rows[index].Translations["russian"]
		got := round.Rows[index].Translations["russian"]
		if got != want {
			t.Fatalf("row %d russian = %q, want %q", index, got, want)
		}
	}
}

func TestRoundTripPOToCSVToPO(t *testing.T) {
	t.Parallel()

	table := mustParseCSVFixture(t, "unquoted.csv")
	for index := range table.Rows {
		table.Rows[index].Translations["russian"] = ""
	}

	catalog, err := pofile.ParseFile(fixturePath("testdata", "po", "russian.po"))
	if err != nil {
		t.Fatalf("pofile.ParseFile error: %v", err)
	}
	if err := MergeCatalogLanguage(table, "russian", catalog, MergeOptions{}); err != nil {
		t.Fatalf("MergeCatalogLanguage error: %v", err)
	}

	round, err := CatalogFromTableLanguage(table, "russian")
	if err != nil {
		t.Fatalf("CatalogFromTableLanguage error: %v", err)
	}
	if got := round.Translation("UI_OK", "OK"); got != "Ок" {
		t.Fatalf("UI_OK = %q, want %q", got, "Ок")
	}
	if got := round.Translation("UI_EXIT", "Exit"); got != "Выход" {
		t.Fatalf("UI_EXIT = %q, want %q", got, "Выход")
	}
}
