// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import "testing"

const benchmarkCSV = "" +
	"Language,original,english,russian,german\n" +
	"UI_OK,OK,OK,Ок,OK\n" +
	"UI_CANCEL,Cancel,Cancel,Отмена,Abbrechen\n" +
	"UI_EXIT,Exit,Exit,Выход,Beenden\n"

// BenchmarkParseCSV benchmarks read/parse flow.
func BenchmarkParseCSV(b *testing.B) {
	data := []byte(benchmarkCSV)

	b.ReportAllocs()
	for b.Loop() {
		_, err := ParseCSV(data)
		if err != nil {
			b.Fatalf("ParseCSV error: %v", err)
		}
	}
}

// BenchmarkFormatCSV benchmarks write/format flow.
func BenchmarkFormatCSV(b *testing.B) {
	table, err := ParseCSV([]byte(benchmarkCSV))
	if err != nil {
		b.Fatalf("ParseCSV setup error: %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		_, err := FormatCSV(table)
		if err != nil {
			b.Fatalf("FormatCSV error: %v", err)
		}
	}
}

// BenchmarkUpdateCatalogFromTable benchmarks top-level preprocess flow.
func BenchmarkUpdateCatalogFromTable(b *testing.B) {
	table, err := ParseCSV([]byte(benchmarkCSV))
	if err != nil {
		b.Fatalf("ParseCSV setup error: %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		_, err := UpdateCatalogFromTable(table, "russian", nil)
		if err != nil {
			b.Fatalf("UpdateCatalogFromTable error: %v", err)
		}
	}
}
