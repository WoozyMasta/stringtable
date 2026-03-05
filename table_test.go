// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"errors"
	"strings"
	"testing"
)

func TestParseCSVDuplicateKey(t *testing.T) {
	t.Parallel()

	input := "" +
		`Language,original,english` + "\n" +
		`UI_OK,OK,OK` + "\n" +
		`UI_OK,Another,Another` + "\n"

	_, err := ParseCSVReader(strings.NewReader(input))
	if err == nil {
		t.Fatal("ParseCSVReader error = nil, want duplicate key error")
	}
	if !errors.Is(err, ErrDuplicateKey) {
		t.Fatalf("error = %v, want ErrDuplicateKey", err)
	}
}

func TestFormatCSVWritesAllDefaultLanguagesByDefault(t *testing.T) {
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

	data, err := FormatCSV(table)
	if err != nil {
		t.Fatalf("FormatCSV error: %v", err)
	}

	header := strings.Split(strings.TrimSpace(string(data)), "\n")[0]
	for _, language := range DefaultLanguages {
		if !strings.Contains(header, ","+language) &&
			!strings.HasSuffix(header, ","+language+"\r") {
			t.Fatalf("header misses default language %q: %q", language, header)
		}
	}
}

func TestFormatCSVWithOptionsUseTableLanguagesOnly(t *testing.T) {
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

	data, err := FormatCSVWithOptions(table, &WriteOptions{
		UseTableLanguagesOnly: true,
	})
	if err != nil {
		t.Fatalf("FormatCSVWithOptions error: %v", err)
	}

	header := strings.Split(strings.TrimSpace(string(data)), "\n")[0]
	if strings.Contains(header, ",german") {
		t.Fatalf("header contains unexpected default language: %q", header)
	}
	if !strings.Contains(header, ",english") || !strings.Contains(header, ",russian") {
		t.Fatalf("header misses expected table languages: %q", header)
	}
}
