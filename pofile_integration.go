// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/woozymasta/pofile"
)

// MergeOptions controls PO-to-CSV translation merge behavior.
type MergeOptions struct {
	// DisableOriginalOnEmpty disables original fallback when msgstr is empty.
	DisableOriginalOnEmpty bool `json:"no_fallback_empty,omitempty" yaml:"no_fallback_empty,omitempty"`

	// DisableOriginalOnNoTranslate disables original fallback for "notranslate".
	DisableOriginalOnNoTranslate bool `json:"no_fallback_notranslate,omitempty" yaml:"no_fallback_notranslate,omitempty"`

	// DisableOriginalOnMissing disables original fallback when entry is missing.
	DisableOriginalOnMissing bool `json:"no_fallback_missing,omitempty" yaml:"no_fallback_missing,omitempty"`
}

// TemplateCatalogFromTable builds POT-like catalog from CSV table.
func TemplateCatalogFromTable(table *Table) (*pofile.Catalog, error) {
	if table == nil {
		return nil, ErrNilTable
	}
	if err := table.Validate(); err != nil {
		return nil, err
	}

	catalog := pofile.NewCatalog()
	for _, row := range table.Rows {
		if row.Key == "" {
			continue
		}

		catalog.UpsertMessage(row.Key, row.Original, "")
	}

	return catalog, nil
}

// CatalogFromTableLanguage builds language-specific PO catalog from CSV table.
func CatalogFromTableLanguage(table *Table, language string) (*pofile.Catalog, error) {
	if table == nil {
		return nil, ErrNilTable
	}
	if err := table.Validate(); err != nil {
		return nil, err
	}

	catalog := pofile.NewCatalog()
	if language != "" {
		catalog.Language = language
		catalog.SetHeader("Language", language)
	}
	for _, row := range table.Rows {
		if row.Key == "" {
			continue
		}

		catalog.UpsertMessage(row.Key, row.Original, row.Translations[language])
	}

	return catalog, nil
}

// CatalogSetFromTable builds PO catalogs for requested languages.
func CatalogSetFromTable(
	table *Table,
	languages []string,
) (map[string]*pofile.Catalog, error) {
	if table == nil {
		return nil, ErrNilTable
	}
	if err := table.Validate(); err != nil {
		return nil, err
	}

	selected := languages
	if len(selected) == 0 {
		selected = table.Languages
	}

	out := make(map[string]*pofile.Catalog, len(selected))
	for _, language := range selected {
		catalog, err := CatalogFromTableLanguage(table, language)
		if err != nil {
			return nil, fmt.Errorf("build catalog for %q: %w", language, err)
		}

		out[language] = catalog
	}

	return out, nil
}

// UpdateCatalogFromTable updates catalog entries using table rows.
func UpdateCatalogFromTable(
	table *Table,
	language string,
	existing *pofile.Catalog,
) (*pofile.Catalog, error) {
	if table == nil {
		return nil, ErrNilTable
	}
	if err := table.Validate(); err != nil {
		return nil, err
	}

	out := pofile.NewCatalog()
	if existing != nil {
		out.Headers = maps.Clone(existing.Headers)
		out.Language = existing.Language
	}
	if language != "" {
		out.Language = language
		out.SetHeader("Language", language)
	}

	if existing == nil {
		out.Messages = make([]*pofile.Message, 0, len(table.Rows))
		for _, row := range table.Rows {
			if row.Key == "" {
				continue
			}

			translation := row.Translations[language]
			message := &pofile.Message{
				Context:      row.Key,
				ID:           row.Original,
				Translations: map[int]string{0: translation},
			}
			out.Messages = append(out.Messages, message)
		}

		return out, nil
	}

	for _, row := range table.Rows {
		if row.Key == "" {
			continue
		}

		translation := row.Translations[language]
		source := existing.FindMessage(row.Key, row.Original)
		if source != nil {
			translation = source.TranslationAt(0)
		}

		target := out.UpsertMessage(row.Key, row.Original, translation)
		if source == nil {
			continue
		}

		target.Comments = slices.Clone(source.Comments)
		target.Flags = slices.Clone(source.Flags)
		target.References = slices.Clone(source.References)
		target.IDPlural = source.IDPlural
		target.Obsolete = source.Obsolete
		target.PreviousContext = source.PreviousContext
		target.PreviousID = source.PreviousID
		target.PreviousIDPlural = source.PreviousIDPlural
	}

	return out, nil
}

// MergeCatalogLanguage applies PO translations into one CSV table language.
func MergeCatalogLanguage(
	table *Table,
	language string,
	catalog *pofile.Catalog,
	options MergeOptions,
) error {
	if table == nil {
		return ErrNilTable
	}
	if catalog == nil {
		return ErrNilCatalog
	}
	if err := table.Validate(); err != nil {
		return err
	}

	table.EnsureLanguage(language)
	for rowIndex := range table.Rows {
		row := &table.Rows[rowIndex]
		if row.Translations == nil {
			row.Translations = make(map[string]string)
		}

		message := catalog.FindMessage(row.Key, row.Original)
		switch {
		case message == nil && !options.DisableOriginalOnMissing:
			row.Translations[language] = row.Original
		case message == nil:
			continue
		case !options.DisableOriginalOnNoTranslate &&
			message.HasFlag("notranslate"):
			row.Translations[language] = row.Original
		default:
			translation := message.TranslationAt(0)
			if translation == "" && !options.DisableOriginalOnEmpty {
				translation = row.Original
			}
			row.Translations[language] = translation
		}
	}

	return nil
}

// MergeCatalogSet applies multiple PO catalogs into table language columns.
func MergeCatalogSet(
	table *Table,
	catalogs map[string]*pofile.Catalog,
	options MergeOptions,
) error {
	if table == nil {
		return ErrNilTable
	}
	for language, catalog := range catalogs {
		if err := MergeCatalogLanguage(table, language, catalog, options); err != nil {
			return fmt.Errorf("merge language %q: %w", language, err)
		}
	}

	return nil
}

// ParseCatalogDir loads all *.po files from directory by language filename.
func ParseCatalogDir(dir string) (map[string]*pofile.Catalog, error) {
	return ReadCatalogSet(dir)
}

// ReadCatalogSet loads all *.po files from directory by language filename.
func ReadCatalogSet(dir string) (map[string]*pofile.Catalog, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read po directory: %w", err)
	}

	out := make(map[string]*pofile.Catalog)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".po") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		catalog, err := pofile.ParseFile(path)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", path, err)
		}

		language := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		out[language] = catalog
	}

	return out, nil
}

// WriteCatalogSet writes PO catalogs into "<language>.po" files in directory.
func WriteCatalogSet(dir string, catalogs map[string]*pofile.Catalog) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create po directory: %w", err)
	}

	languages := make([]string, 0, len(catalogs))
	for language := range catalogs {
		languages = append(languages, language)
	}
	sort.Strings(languages)

	for _, language := range languages {
		catalog := catalogs[language]
		path := filepath.Join(dir, language+".po")
		if err := pofile.WriteFile(path, catalog); err != nil {
			return fmt.Errorf("write %q: %w", path, err)
		}
	}

	return nil
}
