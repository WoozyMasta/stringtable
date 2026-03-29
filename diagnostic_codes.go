// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"github.com/woozymasta/lintkit/lint"
)

const (
	// LintModule is stable lint module namespace for stringtable rules.
	LintModule = "stringtable"

	// LintFileKindCSV is lint file kind token for stringtable CSV files.
	LintFileKindCSV lint.FileKind = "stringtable.csv"
)

const (
	// StageLint marks stringtable CSV lint diagnostics.
	StageLint lint.Stage = "csv"
)

const (
	// CodeLintEmptyHeaderLanguage reports empty language column name in header.
	CodeLintEmptyHeaderLanguage lint.Code = 2001

	// CodeLintDuplicateHeaderLanguage reports duplicate language columns in header.
	CodeLintDuplicateHeaderLanguage lint.Code = 2002

	// CodeLintUnknownHeaderLanguage reports unknown language column in header.
	CodeLintUnknownHeaderLanguage lint.Code = 2003

	// CodeLintMissingDefaultLanguages reports missing default language columns.
	CodeLintMissingDefaultLanguages lint.Code = 2004

	// CodeLintRowColumnCountMismatch reports row/header column count mismatch.
	CodeLintRowColumnCountMismatch lint.Code = 2005

	// CodeLintDuplicateKey reports duplicate key values in data rows.
	CodeLintDuplicateKey lint.Code = 2006

	// CodeLintKeyTrimMismatch reports key values with leading/trailing spaces.
	CodeLintKeyTrimMismatch lint.Code = 2007

	// CodeLintInvalidKeyPattern reports translation key token format mismatch.
	CodeLintInvalidKeyPattern lint.Code = 2008

	// CodeLintOriginalEmpty reports empty value in required original column.
	CodeLintOriginalEmpty lint.Code = 2009

	// CodeLintTranslationEmpty reports empty value in translation columns.
	CodeLintTranslationEmpty lint.Code = 2010

	// CodeLintUnquotedNonKeyColumn reports unquoted non-key CSV columns.
	CodeLintUnquotedNonKeyColumn lint.Code = 2011
)

const (
	// DefaultKeyPattern is default regexp for translation key values.
	DefaultKeyPattern = "^[-_A-Za-z0-9]+$"
)

var diagnosticCodeCatalogHandle = lint.NewCodeCatalogHandle(
	lint.CodeCatalogConfig{
		Module:            LintModule,
		CodePrefix:        "STBL",
		ModuleName:        "DayZ stringtable",
		ModuleDescription: "Lint rules for DayZ stringtable CSV files.",
		ScopeDescriptions: map[lint.Stage]string{
			StageLint: "Stringtable CSV structure and translation diagnostics.",
		},
	},
	diagnosticCatalog,
)

// getDiagnosticCodeCatalog returns lazy-initialized code catalog helper.
func getDiagnosticCodeCatalog() (lint.CodeCatalog, error) {
	return diagnosticCodeCatalogHandle.Catalog()
}

// DiagnosticCatalog returns stable diagnostics metadata list.
func DiagnosticCatalog() []lint.CodeSpec {
	return diagnosticCodeCatalogHandle.CodeSpecs()
}

// DiagnosticByCode returns diagnostic metadata for stable code.
func DiagnosticByCode(code lint.Code) (lint.CodeSpec, bool) {
	return diagnosticCodeCatalogHandle.ByCode(code)
}

// LintRuleID returns lint rule ID mapped from stable stringtable code.
func LintRuleID(code lint.Code) string {
	return diagnosticCodeCatalogHandle.RuleIDOrUnknown(code)
}

// LintRuleSpecs returns deterministic lint rule specs from diagnostics catalog.
func LintRuleSpecs() []lint.RuleSpec {
	return diagnosticCodeCatalogHandle.RuleSpecs()
}
