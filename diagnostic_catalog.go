// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import "github.com/woozymasta/lintkit/lint"

// KeyPatternRuleOptions defines options for translation key pattern rule.
type KeyPatternRuleOptions struct {
	// Pattern is regexp used to validate key token format.
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty" jsonschema:"example=^[-_A-Za-z0-9]+$,example=^[A-Z0-9_]+$"`
}

// csvCodeSpec returns code spec with normalized file kind override.
func csvCodeSpec(spec lint.CodeSpec) lint.CodeSpec {
	return lint.WithCodeRule(spec, lint.CodeRuleOverride{
		FileKinds: []lint.FileKind{LintFileKindCSV},
	})
}

// csvCodeSpecDescription returns code spec with user-facing description.
func csvCodeSpecDescription(spec lint.CodeSpec, description string) lint.CodeSpec {
	spec.Description = description
	return csvCodeSpec(spec)
}

// diagnosticCatalog stores stable diagnostics metadata table.
var diagnosticCatalog = []lint.CodeSpec{
	csvCodeSpecDescription(
		lint.ErrorCodeSpec(
			CodeLintEmptyHeaderLanguage,
			StageLint,
			"header language column name must be non-empty",
		),
		"Every translation column after `Language,original` must have "+
			"a non-empty language name.",
	),
	csvCodeSpecDescription(
		lint.ErrorCodeSpec(
			CodeLintDuplicateHeaderLanguage,
			StageLint,
			"header language columns must be unique",
		),
		"Each language column name should appear once. Duplicate names make "+
			"column mapping ambiguous for parsers and exporters.",
	),
	csvCodeSpecDescription(
		lint.ErrorCodeSpec(
			CodeLintUnknownHeaderLanguage,
			StageLint,
			"header contains unknown language name",
		),
		"Language columns after `Language,original` must use supported DayZ "+
			"language names. Rename unsupported columns or update conversion flow.",
	),
	csvCodeSpecDescription(
		lint.WarningCodeSpec(
			CodeLintMissingDefaultLanguages,
			StageLint,
			"header is missing required default language columns",
		),
		"Add missing DayZ default language columns to keep expected export order "+
			"and avoid incomplete localization coverage.",
	),
	csvCodeSpecDescription(
		lint.ErrorCodeSpec(
			CodeLintRowColumnCountMismatch,
			StageLint,
			"row column count must match header column count",
		),
		"Every data row must have the same number of cells as the header; "+
			"extra or missing cells shift translations into wrong languages.",
	),
	csvCodeSpecDescription(
		lint.ErrorCodeSpec(
			CodeLintDuplicateKey,
			StageLint,
			"translation key must be unique",
		),
		"Duplicate `Language` keys create ambiguous lookup and merge behavior. "+
			"Keep exactly one row per key.",
	),
	csvCodeSpecDescription(
		lint.WarningCodeSpec(
			CodeLintKeyTrimMismatch,
			StageLint,
			"translation key should not have leading or trailing spaces",
		),
		"Keys with surrounding spaces are hard to spot and may behave as "+
			"different tokens than visually similar trimmed keys.",
	),
	csvCodeSpecDescription(
		lint.WithCodeOptions(
			lint.ErrorCodeSpec(
				CodeLintInvalidKeyPattern,
				StageLint,
				"translation key must match configured token pattern",
			),
			KeyPatternRuleOptions{Pattern: DefaultKeyPattern},
		),
		"Keep keys machine-safe and deterministic. Change `pattern` option only "+
			"when project naming rules intentionally differ.",
	),
	csvCodeSpecDescription(
		lint.ErrorCodeSpec(
			CodeLintOriginalEmpty,
			StageLint,
			"`original` column must be non-empty",
		),
		"`original` stores the source text. Empty source usually means broken "+
			"authoring or accidental row damage.",
	),
	csvCodeSpecDescription(
		lint.WarningCodeSpec(
			CodeLintTranslationEmpty,
			StageLint,
			"translation columns should be non-empty",
		),
		"Empty translation is allowed but likely unfinished localization. "+
			"Fill value or intentionally suppress this warning in your workflow.",
	),
	csvCodeSpecDescription(
		lint.WarningCodeSpec(
			CodeLintUnquotedNonKeyColumn,
			StageLint,
			"non-key cells should be quoted",
		),
		"Quote `original` and translation cells to keep CSV stable when text "+
			"contains commas, quotes, or leading/trailing spaces.",
	),
}
