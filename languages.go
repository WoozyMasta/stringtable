// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"path/filepath"
	"slices"
	"strings"
)

// DefaultLanguages is the default DayZ language order.
var DefaultLanguages = []string{
	"english",
	"czech",
	"german",
	"russian",
	"polish",
	"hungarian",
	"italian",
	"spanish",
	"french",
	"chinese",
	"japanese",
	"portuguese",
	"chinesesimp",
}

// languageCodes maps DayZ language names to stable language codes.
var languageCodes = map[string]string{
	"english":     "en",
	"czech":       "cs",
	"german":      "de",
	"russian":     "ru",
	"polish":      "pl",
	"hungarian":   "hu",
	"italian":     "it",
	"spanish":     "es",
	"french":      "fr",
	"chinese":     "zh-Hant",
	"japanese":    "ja",
	"portuguese":  "pt",
	"chinesesimp": "zh-Hans",
}

// languageNamesByCode maps stable language codes back to DayZ language names.
var languageNamesByCode = map[string]string{
	"en":      "english",
	"cs":      "czech",
	"de":      "german",
	"ru":      "russian",
	"pl":      "polish",
	"hu":      "hungarian",
	"it":      "italian",
	"es":      "spanish",
	"fr":      "french",
	"zh-Hant": "chinese",
	"ja":      "japanese",
	"pt":      "portuguese",
	"zh-Hans": "chinesesimp",
}

// ParseLanguages parses comma-separated language list.
func ParseLanguages(value string) []string {
	if strings.TrimSpace(value) == "" {
		return slices.Clone(DefaultLanguages)
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		lang := strings.TrimSpace(part)
		if lang == "" {
			continue
		}
		if _, ok := seen[lang]; ok {
			continue
		}

		seen[lang] = struct{}{}
		out = append(out, lang)
	}

	return out
}

// LanguageCode returns language code for DayZ language name.
func LanguageCode(language string) (string, bool) {
	key := strings.ToLower(strings.TrimSpace(language))
	if key == "" {
		return "", false
	}

	code, ok := languageCodes[key]
	return code, ok
}

// LanguageNameFromCode returns DayZ language name for language code.
func LanguageNameFromCode(code string) (string, bool) {
	language, ok := languageNamesByCode[strings.TrimSpace(code)]
	return language, ok
}

// ContainsLanguage reports whether language exists in list.
func ContainsLanguage(list []string, language string) bool {
	return slices.Contains(list, language)
}

// ExtractLanguageName extracts language from "/path/lang.po" path.
func ExtractLanguageName(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// SelectLanguages returns deterministic language list using include/exclude.
func SelectLanguages(available, include, exclude []string) []string {
	base := include
	if len(base) == 0 {
		base = available
	}
	if len(base) == 0 {
		base = DefaultLanguages
	}

	excluded := make(map[string]struct{}, len(exclude))
	for _, language := range exclude {
		excluded[language] = struct{}{}
	}

	out := make([]string, 0, len(base))
	seen := make(map[string]struct{}, len(base))
	for _, language := range base {
		if language == "" {
			continue
		}
		if _, ok := seen[language]; ok {
			continue
		}
		if _, ok := excluded[language]; ok {
			continue
		}

		seen[language] = struct{}{}
		out = append(out, language)
	}

	return out
}
