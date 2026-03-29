// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/woozymasta/lintkit/lint"
	"github.com/woozymasta/pofile"
)

// LintRulesProvider registers stringtable and nested pofile rules.
type LintRulesProvider struct{}

// RegisterRules adds provider-owned rules to target registrar.
func (provider LintRulesProvider) RegisterRules(
	registrar lint.RuleRegistrar,
) error {
	return RegisterLintRules(registrar)
}

// RegisterRulesByScope adds provider-owned rules filtered by scope tokens.
func (provider LintRulesProvider) RegisterRulesByScope(
	registrar lint.RuleRegistrar,
	scopes ...string,
) error {
	return RegisterLintRulesByScope(registrar, scopes...)
}

// RegisterRulesByStage adds provider-owned rules filtered by stage tokens.
func (provider LintRulesProvider) RegisterRulesByStage(
	registrar lint.RuleRegistrar,
	stages ...lint.Stage,
) error {
	return RegisterLintRulesByStage(registrar, stages...)
}

// RegisterLintRules registers stable stringtable and pofile rules.
func RegisterLintRules(registrar lint.RuleRegistrar) error {
	if registrar == nil {
		return ErrNilLintRuleRegistrar
	}

	return lint.RegisterRuleProviders(
		registrar,
		stringtableLintProvider{},
		pofile.LintRulesProvider{},
	)
}

// RegisterLintRulesByScope registers rules filtered by scope tokens.
func RegisterLintRulesByScope(
	registrar lint.RuleRegistrar,
	scopes ...string,
) error {
	if registrar == nil {
		return ErrNilLintRuleRegistrar
	}

	return lint.RegisterRuleProvidersByScope(
		registrar,
		scopes,
		stringtableLintProvider{},
		pofile.LintRulesProvider{},
	)
}

// RegisterLintRulesByStage registers rules filtered by stage tokens.
func RegisterLintRulesByStage(
	registrar lint.RuleRegistrar,
	stages ...lint.Stage,
) error {
	if registrar == nil {
		return ErrNilLintRuleRegistrar
	}

	return lint.RegisterRuleProvidersByStage(
		registrar,
		stages,
		stringtableLintProvider{},
		pofile.LintRulesProvider{},
	)
}

// stringtableLintProvider stores stringtable-owned rule registration behavior.
type stringtableLintProvider struct{}

// RegisterRules registers stringtable-owned rules into registrar.
func (provider stringtableLintProvider) RegisterRules(
	registrar lint.RuleRegistrar,
) error {
	if registrar == nil {
		return ErrNilLintRuleRegistrar
	}

	runners, err := stringtableLintRuleRunners()
	if err != nil {
		return err
	}

	return registrar.Register(runners...)
}

// RegisterRulesByScope registers stringtable rules filtered by scope tokens.
func (provider stringtableLintProvider) RegisterRulesByScope(
	registrar lint.RuleRegistrar,
	scopes ...string,
) error {
	if !containsStageScope(scopes, StageLint) {
		return nil
	}

	return provider.RegisterRules(registrar)
}

// RegisterRulesByStage registers stringtable rules filtered by stage tokens.
func (provider stringtableLintProvider) RegisterRulesByStage(
	registrar lint.RuleRegistrar,
	stages ...lint.Stage,
) error {
	if !containsStage(stages, StageLint) {
		return nil
	}

	return provider.RegisterRules(registrar)
}

// stringtableLintRuleRunners builds deterministic stringtable rule runners.
func stringtableLintRuleRunners() ([]lint.RuleRunner, error) {
	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		return nil, err
	}

	codes := []lint.Code{
		CodeLintEmptyHeaderLanguage,
		CodeLintDuplicateHeaderLanguage,
		CodeLintUnknownHeaderLanguage,
		CodeLintMissingDefaultLanguages,
		CodeLintRowColumnCountMismatch,
		CodeLintDuplicateKey,
		CodeLintKeyTrimMismatch,
		CodeLintInvalidKeyPattern,
		CodeLintOriginalEmpty,
		CodeLintTranslationEmpty,
		CodeLintUnquotedNonKeyColumn,
	}

	runners := make([]lint.RuleRunner, 0, len(codes))
	for index := range codes {
		code := codes[index]
		spec, ok := catalog.ByCode(code)
		if !ok {
			return nil, fmt.Errorf("missing code spec for code %d", code)
		}

		runners = append(runners, stringtableLintRuleRunner{
			Code: code,
			Spec: catalog.RuleSpec(spec),
		})
	}

	return runners, nil
}

// stringtableLintRuleRunner executes one stringtable rule by stable code.
type stringtableLintRuleRunner struct {
	// Spec stores stable metadata for this runner.
	Spec lint.RuleSpec

	// Code stores stable lint code for this runner.
	Code lint.Code
}

// RuleSpec returns stable metadata descriptor for current runner.
func (runner stringtableLintRuleRunner) RuleSpec() lint.RuleSpec {
	return runner.Spec
}

// Check runs one rule against parsed CSV lint analysis.
func (runner stringtableLintRuleRunner) Check(
	_ context.Context,
	run *lint.RunContext,
	emit lint.DiagnosticEmit,
) error {
	analysis, err := getCSVLintAnalysis(run)
	if err != nil {
		return err
	}

	switch runner.Code {
	case CodeLintInvalidKeyPattern:
		issues, err := keyPatternIssues(analysis, run)
		if err != nil {
			return err
		}

		for issueIndex := range issues {
			emit(issueToLintDiagnostic(
				issues[issueIndex],
				runner.Spec,
				run,
			))
		}

		return nil
	default:
		issues := analysis.IssuesByCode[runner.Code]
		for issueIndex := range issues {
			emit(issueToLintDiagnostic(
				issues[issueIndex],
				runner.Spec,
				run,
			))
		}

		return nil
	}
}

// issueToLintDiagnostic converts one precomputed issue into shared model.
func issueToLintDiagnostic(
	issue csvLintIssue,
	spec lint.RuleSpec,
	run *lint.RunContext,
) lint.Diagnostic {
	path := ""
	if run != nil {
		path = run.TargetPath
	}

	position := lint.Position{
		Line:   issue.Line,
		Column: issue.Column,
	}

	return lint.Diagnostic{
		RuleID:   spec.ID,
		Code:     spec.Code,
		Severity: issue.Severity,
		Message:  issue.Message,
		Path:     path,
		Start:    position,
		End:      position,
	}
}

// containsStageScope reports whether scopes include StageLint token.
func containsStageScope(scopes []string, stage lint.Stage) bool {
	if len(scopes) == 0 {
		return true
	}

	stageToken := strings.TrimSpace(string(stage))
	for index := range scopes {
		if strings.TrimSpace(scopes[index]) == stageToken {
			return true
		}
	}

	return false
}

// containsStage reports whether stages include given stage token.
func containsStage(stages []lint.Stage, stage lint.Stage) bool {
	if len(stages) == 0 {
		return true
	}

	return slices.Contains(stages, stage)
}
