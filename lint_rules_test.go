// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/woozymasta/lintkit/lint"
	"github.com/woozymasta/lintkit/linting"
	"github.com/woozymasta/lintkit/linttest"
	"github.com/woozymasta/pofile"
)

func TestRegisterLintRulesNilRegistrar(t *testing.T) {
	t.Parallel()

	if err := RegisterLintRules(nil); !errors.Is(err, ErrNilLintRuleRegistrar) {
		t.Fatalf(
			"RegisterLintRules(nil) error=%v, want ErrNilLintRuleRegistrar",
			err,
		)
	}
}

func TestLintRulesProviderIncludesPofileRules(t *testing.T) {
	t.Parallel()

	engine, err := linting.NewEngineWithProviders(LintRulesProvider{})
	if err != nil {
		t.Fatalf("NewEngineWithProviders() error: %v", err)
	}

	rules := engine.Rules()
	if len(rules) == 0 {
		t.Fatal("len(engine.Rules())=0, want >0")
	}

	if !hasRuleID(rules, LintRuleID(CodeLintInvalidKeyPattern)) {
		t.Fatalf("missing stringtable rule %q", LintRuleID(CodeLintInvalidKeyPattern))
	}

	if !hasRuleID(rules, pofile.LintRuleID(pofile.CodeLintEntryWithoutID)) {
		t.Fatalf(
			"missing nested pofile rule %q",
			pofile.LintRuleID(pofile.CodeLintEntryWithoutID),
		)
	}
}

func TestRegisterLintRulesByScope(t *testing.T) {
	t.Parallel()

	engine := linting.NewEngine()
	if err := RegisterLintRulesByScope(engine, string(StageLint)); err != nil {
		t.Fatalf("RegisterLintRulesByScope() error: %v", err)
	}

	rules := engine.Rules()
	if len(rules) == 0 {
		t.Fatal("len(engine.Rules())=0, want >0")
	}

	if !hasRuleID(rules, LintRuleID(CodeLintInvalidKeyPattern)) {
		t.Fatalf(
			"missing stringtable scope rule %q",
			LintRuleID(CodeLintInvalidKeyPattern),
		)
	}
}

func TestRegisterLintRulesByStage(t *testing.T) {
	t.Parallel()

	engine := linting.NewEngine()
	if err := RegisterLintRulesByStage(engine, StageLint); err != nil {
		t.Fatalf("RegisterLintRulesByStage() error: %v", err)
	}

	rules := engine.Rules()
	if len(rules) == 0 {
		t.Fatal("len(engine.Rules())=0, want >0")
	}

	if !hasRuleID(rules, LintRuleID(CodeLintInvalidKeyPattern)) {
		t.Fatalf(
			"missing stringtable stage rule %q",
			LintRuleID(CodeLintInvalidKeyPattern),
		)
	}
}

func TestLintRuleSpecsMatchCatalog(t *testing.T) {
	t.Parallel()

	linttest.AssertCatalogContract(
		t,
		LintModule,
		DiagnosticCatalog(),
		LintRuleSpecs(),
		LintRuleID,
	)
}

func TestStringtableLintRules(t *testing.T) {
	t.Parallel()

	content := strings.Join([]string{
		`Language,original,english,foo`,
		`"BAD KEY","","","ok"`,
		`"GOOD","ok"`,
	}, "\n")

	result := runStringtableLint(t, content, nil)
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintMissingDefaultLanguages))
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintUnknownHeaderLanguage))
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintUnquotedNonKeyColumn))
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintInvalidKeyPattern))
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintOriginalEmpty))
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintTranslationEmpty))
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintRowColumnCountMismatch))
}

func TestStringtableLintHeaderAndKeyHygieneRules(t *testing.T) {
	t.Parallel()

	content := strings.Join([]string{
		`Language,original,english,,english`,
		`" KEY ","base","ok","","ok2"`,
		`"KEY","base2","ok","","ok3"`,
		`"KEY","base3","ok","","ok4"`,
	}, "\n")

	result := runStringtableLint(t, content, nil)
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintEmptyHeaderLanguage))
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintDuplicateHeaderLanguage))
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintKeyTrimMismatch))
	assertHasDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintDuplicateKey))
}

func TestStringtableLintKeyPatternOptions(t *testing.T) {
	t.Parallel()

	content := strings.Join([]string{
		buildFullHeaderLine(),
		buildFullDataLine("UI.BAD", "Hello"),
	}, "\n")

	policy := linting.RunPolicy{
		Rules: map[string]linting.RuleSettings{
			LintRuleID(CodeLintInvalidKeyPattern): {
				Options: KeyPatternRuleOptions{Pattern: "^[-._A-Za-z0-9]+$"},
			},
		},
	}

	result := runStringtableLint(t, content, &policy)
	assertNoDiagnosticCode(t, result.Diagnostics, lintCodeToken(t, CodeLintInvalidKeyPattern))
}

func runStringtableLint(
	t *testing.T,
	content string,
	policy *linting.RunPolicy,
) linting.RunResult {
	t.Helper()

	engine, err := linting.NewEngineWithProviders(LintRulesProvider{})
	if err != nil {
		t.Fatalf("NewEngineWithProviders() error: %v", err)
	}

	options := &linting.RunOptions{Policy: policy}
	result, err := engine.Run(context.Background(), lint.RunContext{
		TargetPath: "stringtable.csv",
		TargetKind: LintFileKindCSV,
		Content:    []byte(content),
	}, options)
	if err != nil {
		t.Fatalf("engine.Run() error: %v", err)
	}

	return result
}

func hasRuleID(rules []lint.RuleSpec, target string) bool {
	for index := range rules {
		if rules[index].ID == target {
			return true
		}
	}

	return false
}

func assertHasDiagnosticCode(
	t *testing.T,
	diagnostics []lint.Diagnostic,
	code string,
) {
	t.Helper()

	for index := range diagnostics {
		if diagnostics[index].Code == code {
			return
		}
	}

	t.Fatalf("missing diagnostic code %q", code)
}

func assertNoDiagnosticCode(
	t *testing.T,
	diagnostics []lint.Diagnostic,
	code string,
) {
	t.Helper()

	for index := range diagnostics {
		if diagnostics[index].Code == code {
			t.Fatalf("unexpected diagnostic code %q", code)
		}
	}
}

func buildFullHeaderLine() string {
	cells := make([]string, 0, len(DefaultLanguages)+2)
	cells = append(cells, `"Language"`, `"original"`)
	for index := range DefaultLanguages {
		cells = append(cells, `"`+DefaultLanguages[index]+`"`)
	}

	return strings.Join(cells, ",")
}

func buildFullDataLine(key string, original string) string {
	cells := make([]string, 0, len(DefaultLanguages)+2)
	cells = append(cells, `"`+key+`"`, `"`+original+`"`)
	for range DefaultLanguages {
		cells = append(cells, `"ok"`)
	}

	return strings.Join(cells, ",")
}

func lintCodeToken(t *testing.T, code lint.Code) string {
	t.Helper()

	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		t.Fatalf("getDiagnosticCodeCatalog() error: %v", err)
	}

	spec, ok := catalog.ByCode(code)
	if !ok {
		t.Fatalf("catalog.ByCode(%d)=false", code)
	}

	return catalog.RuleSpec(spec).Code
}
