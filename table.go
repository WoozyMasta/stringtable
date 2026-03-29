// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"maps"
	"os"
	"strings"
)

const (
	keyColumnName      = "Language"
	originalColumnName = "original"
)

// Table is a DayZ stringtable CSV model.
type Table struct {
	// Languages defines translation column order after key+original columns.
	Languages []string `json:"languages,omitempty" yaml:"languages,omitempty"`

	// Rows stores key/original and per-language values.
	Rows []Row `json:"rows,omitempty" yaml:"rows,omitempty"`
}

// WriteOptions controls CSV serialization behavior.
type WriteOptions struct {
	// UseTableLanguagesOnly writes only table header languages as-is.
	// When false, writer emits all default DayZ languages plus any extra
	// languages found in table data.
	UseTableLanguagesOnly bool `json:"use_table_languages_only,omitempty" yaml:"use_table_languages_only,omitempty"`
}

// Row is one stringtable.csv entry.
type Row struct {
	// Translations stores values per language column.
	Translations map[string]string `json:"translations,omitempty" yaml:"translations,omitempty"`

	// Key is value from "Language" column.
	Key string `json:"key" yaml:"key"`

	// Original is value from "original" column.
	Original string `json:"original" yaml:"original"`
}

// NewTable creates an empty table with language order.
func NewTable(languages []string) *Table {
	return &Table{
		Languages: normalizeLanguageOrder(languages),
		Rows:      make([]Row, 0),
	}
}

// ParseCSV parses stringtable CSV bytes.
func ParseCSV(data []byte) (*Table, error) {
	return ParseCSVReader(bytes.NewReader(data))
}

// ParseCSVReader parses stringtable CSV from reader.
func ParseCSVReader(reader io.Reader) (*Table, error) {
	records, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}
	if len(records) == 0 {
		return NewTable(nil), nil
	}

	header := records[0]
	if err := validateHeader(header); err != nil {
		return nil, err
	}

	languages := normalizeLanguageOrder(header[2:])
	table := NewTable(languages)
	seen := make(map[string]int, len(records)-1)

	for index, record := range records[1:] {
		if len(record) == 0 {
			continue
		}

		key := cellAt(record, 0)
		if key == "" {
			continue
		}
		if first, ok := seen[key]; ok {
			return nil, fmt.Errorf(
				"row %d key %q duplicates row %d: %w",
				index+2,
				key,
				first+2,
				ErrDuplicateKey,
			)
		}
		seen[key] = index + 1

		row := Row{
			Key:          key,
			Original:     cellAt(record, 1),
			Translations: make(map[string]string, len(languages)),
		}
		for langIndex, language := range languages {
			row.Translations[language] = cellAt(record, langIndex+2)
		}

		table.Rows = append(table.Rows, row)
	}

	return table, nil
}

// ParseCSVFile parses stringtable CSV file from disk.
func ParseCSVFile(path string) (*Table, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open csv file: %w", err)
	}
	defer func() { _ = file.Close() }()

	table, err := ParseCSVReader(file)
	if err != nil {
		return nil, fmt.Errorf("parse csv file: %w", err)
	}

	return table, nil
}

// FormatCSV serializes table into stringtable CSV bytes.
func FormatCSV(table *Table) ([]byte, error) {
	return FormatCSVWithOptions(table, nil)
}

// FormatCSVWithOptions serializes table into stringtable CSV bytes.
func FormatCSVWithOptions(table *Table, options *WriteOptions) ([]byte, error) {
	if table == nil {
		return nil, ErrNilTable
	}

	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	writer.UseCRLF = true

	languages := resolveWriteLanguages(table, options)
	header := make([]string, 0, len(languages)+2)
	header = append(header, keyColumnName, originalColumnName)
	header = append(header, languages...)
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("write csv header: %w", err)
	}

	for _, row := range table.Rows {
		record := make([]string, 0, len(languages)+2)
		record = append(record, row.Key, row.Original)
		for _, language := range languages {
			record = append(record, row.Translations[language])
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("write csv row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("flush csv writer: %w", err)
	}

	return []byte(builder.String()), nil
}

// WriteCSVFile writes table as stringtable CSV file.
func WriteCSVFile(path string, table *Table) error {
	return WriteCSVFileWithOptions(path, table, nil)
}

// WriteCSVFileWithOptions writes table as stringtable CSV file.
func WriteCSVFileWithOptions(path string, table *Table, options *WriteOptions) error {
	data, err := FormatCSVWithOptions(table, options)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write csv file: %w", err)
	}

	return nil
}

// Clone deep-copies table.
func (t *Table) Clone() *Table {
	if t == nil {
		return nil
	}

	out := NewTable(t.Languages)
	out.Rows = make([]Row, 0, len(t.Rows))
	for _, row := range t.Rows {
		copied := Row{
			Key:          row.Key,
			Original:     row.Original,
			Translations: maps.Clone(row.Translations),
		}
		out.Rows = append(out.Rows, copied)
	}

	return out
}

// EnsureLanguage appends missing language to column order.
func (t *Table) EnsureLanguage(language string) {
	if t == nil {
		return
	}
	if language == "" || ContainsLanguage(t.Languages, language) {
		return
	}

	t.Languages = append(t.Languages, language)
}

// Validate checks structural consistency.
func (t *Table) Validate() error {
	if t == nil {
		return ErrNilTable
	}

	seen := make(map[string]int, len(t.Rows))
	for index, row := range t.Rows {
		if row.Key == "" {
			continue
		}
		if first, ok := seen[row.Key]; ok {
			return fmt.Errorf(
				"rows[%d] key %q duplicates rows[%d]: %w",
				index,
				row.Key,
				first,
				ErrDuplicateKey,
			)
		}

		seen[row.Key] = index
	}

	return nil
}

// cellAt returns record value at index or empty string.
func cellAt(record []string, index int) string {
	if index < 0 || index >= len(record) {
		return ""
	}

	return record[index]
}

// validateHeader checks required stringtable header columns.
func validateHeader(header []string) error {
	if len(header) < 2 {
		return fmt.Errorf("header requires at least 2 columns: %w", ErrInvalidHeader)
	}
	if header[0] != keyColumnName || header[1] != originalColumnName {
		return fmt.Errorf(
			`header must start with %q,%q: %w`,
			keyColumnName,
			originalColumnName,
			ErrInvalidHeader,
		)
	}

	return nil
}

// normalizeLanguageOrder deduplicates languages while preserving first order.
func normalizeLanguageOrder(languages []string) []string {
	out := make([]string, 0, len(languages))
	seen := make(map[string]struct{}, len(languages))
	for _, language := range languages {
		trimmed := strings.TrimSpace(language)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}

		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}

	return out
}

// resolveWriteLanguages selects output language columns for serialization.
func resolveWriteLanguages(table *Table, options *WriteOptions) []string {
	if options != nil && options.UseTableLanguagesOnly {
		return normalizeLanguageOrder(table.Languages)
	}

	estimated := len(DefaultLanguages) + len(table.Languages)
	for rowIndex := range table.Rows {
		estimated += len(table.Rows[rowIndex].Translations)
	}

	out := make([]string, 0, estimated)
	seen := make(map[string]struct{}, estimated)
	for _, language := range DefaultLanguages {
		out = appendUniqueTrimmedLanguage(out, seen, language)
	}
	for _, language := range table.Languages {
		out = appendUniqueTrimmedLanguage(out, seen, language)
	}
	for _, row := range table.Rows {
		for language := range row.Translations {
			out = appendUniqueTrimmedLanguage(out, seen, language)
		}
	}

	return out
}

// appendUniqueTrimmedLanguage appends non-empty unique language names.
func appendUniqueTrimmedLanguage(
	out []string,
	seen map[string]struct{},
	language string,
) []string {
	trimmed := strings.TrimSpace(language)
	if trimmed == "" {
		return out
	}
	if _, ok := seen[trimmed]; ok {
		return out
	}

	seen[trimmed] = struct{}{}
	return append(out, trimmed)
}
