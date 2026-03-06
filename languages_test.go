// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import "testing"

func TestParseLanguages(t *testing.T) {
	t.Parallel()

	languages := ParseLanguages("russian, english, russian, ,german")
	if len(languages) != 3 {
		t.Fatalf("len = %d, want 3", len(languages))
	}
	if languages[0] != "russian" || languages[1] != "english" || languages[2] != "german" {
		t.Fatalf("languages = %#v", languages)
	}
}

func TestExtractLanguageName(t *testing.T) {
	t.Parallel()

	got := ExtractLanguageName("l18n/russian.po")
	if got != "russian" {
		t.Fatalf("ExtractLanguageName = %q, want %q", got, "russian")
	}
}

func TestSelectLanguages(t *testing.T) {
	t.Parallel()

	available := []string{"english", "russian", "german"}
	include := []string{"german", "russian", "german"}
	exclude := []string{"russian"}

	got := SelectLanguages(available, include, exclude)
	if len(got) != 1 || got[0] != "german" {
		t.Fatalf("SelectLanguages = %#v, want %#v", got, []string{"german"})
	}
}

func TestSelectLanguagesFallback(t *testing.T) {
	t.Parallel()

	got := SelectLanguages(nil, nil, []string{"english"})
	if len(got) == 0 {
		t.Fatal("SelectLanguages returned empty fallback list")
	}
	if got[0] == "english" {
		t.Fatalf("first language = %q, expected excluded language to be removed", got[0])
	}
}

func TestLanguageCode(t *testing.T) {
	t.Parallel()

	code, ok := LanguageCode("russian")
	if !ok {
		t.Fatal("LanguageCode(russian) not found")
	}
	if code != "ru" {
		t.Fatalf("LanguageCode(russian) = %q, want %q", code, "ru")
	}
}

func TestLanguageNameFromCode(t *testing.T) {
	t.Parallel()

	name, ok := LanguageNameFromCode("zh-Hans")
	if !ok {
		t.Fatal("LanguageNameFromCode(zh-Hans) not found")
	}
	if name != "chinesesimp" {
		t.Fatalf("LanguageNameFromCode(zh-Hans) = %q, want %q", name, "chinesesimp")
	}
}
