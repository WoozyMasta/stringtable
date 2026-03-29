// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/woozymasta/lintkit/lint"
)

const (
	// lintAnalysisRunValueKey stores parsed CSV lint analysis cache.
	lintAnalysisRunValueKey = "stringtable.lint.analysis"
)

// csvLintIssue stores one precomputed lint finding for one stable code.
type csvLintIssue struct {
	// Message is user-facing lint finding text.
	Message string

	// Severity stores normalized lint severity for this finding.
	Severity lint.Severity

	// Line stores 1-based source line.
	Line int

	// Column stores 1-based source column.
	Column int
}

// csvLintAnalysis stores parsed CSV rows and precomputed rule issues.
type csvLintAnalysis struct {
	// IssuesByCode stores precomputed findings grouped by stable lint code.
	IssuesByCode map[lint.Code][]csvLintIssue

	// Rows stores parsed CSV data records after header.
	Rows []csvQuotedRecord

	// Header stores parsed first CSV record.
	Header csvQuotedRecord
}

// csvQuotedCell stores parsed cell value with quote origin metadata.
type csvQuotedCell struct {
	// Value stores parsed logical cell value.
	Value string

	// Quoted reports whether source cell token started as quoted.
	Quoted bool
}

// csvQuotedRecord stores parsed CSV record and source line.
type csvQuotedRecord struct {
	// Cells stores parsed logical record cells.
	Cells []csvQuotedCell

	// Line stores 1-based source line where record starts.
	Line int
}

// getCSVLintAnalysis returns cached or freshly built CSV lint analysis.
func getCSVLintAnalysis(run *lint.RunContext) (*csvLintAnalysis, error) {
	if run == nil {
		return nil, nil
	}

	if cached, ok := lint.GetRunValue[*csvLintAnalysis](run, lintAnalysisRunValueKey); ok {
		return cached, nil
	}

	content, err := resolveLintContent(run)
	if err != nil {
		return nil, err
	}

	analysis, err := analyzeCSVLint(content)
	if err != nil {
		return nil, err
	}

	lint.SetRunValue(run, lintAnalysisRunValueKey, analysis)
	return analysis, nil
}

// resolveLintContent returns target content from in-memory bytes or filesystem.
func resolveLintContent(run *lint.RunContext) ([]byte, error) {
	if run == nil {
		return nil, nil
	}

	if len(run.Content) > 0 {
		return run.Content, nil
	}

	path := strings.TrimSpace(run.TargetPath)
	if path == "" {
		return nil, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read target content: %w", err)
	}

	return content, nil
}

// analyzeCSVLint parses CSV and precomputes diagnostics for fixed rules.
func analyzeCSVLint(content []byte) (*csvLintAnalysis, error) {
	records, err := parseCSVQuotedRecords(content)
	if err != nil {
		return nil, err
	}

	analysis := &csvLintAnalysis{
		IssuesByCode: make(map[lint.Code][]csvLintIssue, 6),
	}
	if len(records) == 0 {
		return analysis, nil
	}

	analysis.Header = records[0]
	if len(records) > 1 {
		analysis.Rows = records[1:]
	}

	headerNames := headerValues(analysis.Header)
	addEmptyHeaderLanguageIssues(analysis, headerNames)
	addDuplicateHeaderLanguageIssues(analysis, headerNames)
	addMissingDefaultLanguageIssues(analysis, headerNames)
	addUnknownHeaderLanguageIssues(analysis, headerNames)
	addUnquotedColumnIssues(analysis, records)
	addRowColumnMismatchIssues(analysis)
	addDuplicateKeyIssues(analysis)
	addKeyTrimMismatchIssues(analysis)
	addOriginalEmptyIssues(analysis, headerNames)
	addTranslationEmptyIssues(analysis, headerNames)

	return analysis, nil
}

// addEmptyHeaderLanguageIssues reports empty language names in header columns.
func addEmptyHeaderLanguageIssues(analysis *csvLintAnalysis, header []string) {
	if len(header) < 3 {
		return
	}

	for index := 2; index < len(header); index++ {
		if strings.TrimSpace(header[index]) != "" {
			continue
		}

		analysis.IssuesByCode[CodeLintEmptyHeaderLanguage] = append(
			analysis.IssuesByCode[CodeLintEmptyHeaderLanguage],
			csvLintIssue{
				Severity: lint.SeverityError,
				Message:  "header language column name is empty",
				Line:     analysis.Header.Line,
				Column:   index + 1,
			},
		)
	}
}

// addDuplicateHeaderLanguageIssues reports duplicate language columns in header.
func addDuplicateHeaderLanguageIssues(analysis *csvLintAnalysis, header []string) {
	if len(header) < 3 {
		return
	}

	seen := make(map[string]int, len(header)-2)
	for index := 2; index < len(header); index++ {
		name := strings.TrimSpace(header[index])
		if name == "" {
			continue
		}

		if firstColumn, ok := seen[name]; ok {
			analysis.IssuesByCode[CodeLintDuplicateHeaderLanguage] = append(
				analysis.IssuesByCode[CodeLintDuplicateHeaderLanguage],
				csvLintIssue{
					Severity: lint.SeverityError,
					Message: "duplicate language column " +
						strconv.Quote(name) +
						", first at column " + strconv.Itoa(firstColumn),
					Line:   analysis.Header.Line,
					Column: index + 1,
				},
			)
			continue
		}

		seen[name] = index + 1
	}
}

// addMissingDefaultLanguageIssues reports missing default language columns.
func addMissingDefaultLanguageIssues(analysis *csvLintAnalysis, header []string) {
	if len(header) < 2 {
		return
	}

	present := make(map[string]struct{}, len(header))
	for index := 2; index < len(header); index++ {
		language := strings.TrimSpace(header[index])
		if language == "" {
			continue
		}

		present[language] = struct{}{}
	}

	missing := make([]string, 0, len(DefaultLanguages))
	for index := range DefaultLanguages {
		language := DefaultLanguages[index]
		if _, ok := present[language]; ok {
			continue
		}

		missing = append(missing, language)
	}

	if len(missing) == 0 {
		return
	}

	analysis.IssuesByCode[CodeLintMissingDefaultLanguages] = append(
		analysis.IssuesByCode[CodeLintMissingDefaultLanguages],
		csvLintIssue{
			Severity: lint.SeverityWarning,
			Message: "header is missing default language columns: " +
				strings.Join(missing, ", "),
			Line:   analysis.Header.Line,
			Column: 1,
		},
	)
}

// addUnknownHeaderLanguageIssues reports unknown language columns in header.
func addUnknownHeaderLanguageIssues(analysis *csvLintAnalysis, header []string) {
	if len(header) < 2 {
		return
	}

	for index := 2; index < len(header); index++ {
		name := strings.TrimSpace(header[index])
		if containsString(DefaultLanguages, name) {
			continue
		}

		analysis.IssuesByCode[CodeLintUnknownHeaderLanguage] = append(
			analysis.IssuesByCode[CodeLintUnknownHeaderLanguage],
			csvLintIssue{
				Severity: lint.SeverityError,
				Message:  "unknown language column " + strconv.Quote(name),
				Line:     analysis.Header.Line,
				Column:   index + 1,
			},
		)
	}
}

// addUnquotedColumnIssues reports unquoted non-key columns across records.
func addUnquotedColumnIssues(
	analysis *csvLintAnalysis,
	records []csvQuotedRecord,
) {
	for recordIndex := range records {
		record := records[recordIndex]
		for columnIndex := 1; columnIndex < len(record.Cells); columnIndex++ {
			if record.Cells[columnIndex].Quoted {
				continue
			}

			analysis.IssuesByCode[CodeLintUnquotedNonKeyColumn] = append(
				analysis.IssuesByCode[CodeLintUnquotedNonKeyColumn],
				csvLintIssue{
					Severity: lint.SeverityWarning,
					Message: "column " + strconv.Itoa(columnIndex+1) +
						" should be quoted",
					Line:   record.Line,
					Column: columnIndex + 1,
				},
			)
		}
	}
}

// addRowColumnMismatchIssues reports data row/header column count mismatch.
func addRowColumnMismatchIssues(analysis *csvLintAnalysis) {
	headerColumns := len(analysis.Header.Cells)
	if headerColumns == 0 {
		return
	}

	for rowIndex := range analysis.Rows {
		record := analysis.Rows[rowIndex]
		if len(record.Cells) == headerColumns {
			continue
		}

		analysis.IssuesByCode[CodeLintRowColumnCountMismatch] = append(
			analysis.IssuesByCode[CodeLintRowColumnCountMismatch],
			csvLintIssue{
				Severity: lint.SeverityError,
				Message: "row has " + strconv.Itoa(len(record.Cells)) +
					" columns, expected " + strconv.Itoa(headerColumns),
				Line:   record.Line,
				Column: 1,
			},
		)
	}
}

// addDuplicateKeyIssues reports duplicate key values across data rows.
func addDuplicateKeyIssues(analysis *csvLintAnalysis) {
	seen := make(map[string]int, len(analysis.Rows))
	for rowIndex := range analysis.Rows {
		record := analysis.Rows[rowIndex]
		if len(record.Cells) == 0 {
			continue
		}

		key := record.Cells[0].Value
		if key == "" {
			continue
		}

		if firstLine, ok := seen[key]; ok {
			analysis.IssuesByCode[CodeLintDuplicateKey] = append(
				analysis.IssuesByCode[CodeLintDuplicateKey],
				csvLintIssue{
					Severity: lint.SeverityError,
					Message: "duplicate translation key " +
						strconv.Quote(key) + ", first at line " +
						strconv.Itoa(firstLine),
					Line:   record.Line,
					Column: 1,
				},
			)
			continue
		}

		seen[key] = record.Line
	}
}

// addKeyTrimMismatchIssues reports key values with surrounding spaces.
func addKeyTrimMismatchIssues(analysis *csvLintAnalysis) {
	for rowIndex := range analysis.Rows {
		record := analysis.Rows[rowIndex]
		if len(record.Cells) == 0 {
			continue
		}

		key := record.Cells[0].Value
		trimmed := strings.TrimSpace(key)
		if trimmed == key || trimmed == "" {
			continue
		}

		analysis.IssuesByCode[CodeLintKeyTrimMismatch] = append(
			analysis.IssuesByCode[CodeLintKeyTrimMismatch],
			csvLintIssue{
				Severity: lint.SeverityWarning,
				Message: "translation key has surrounding spaces, normalized value is " +
					strconv.Quote(trimmed),
				Line:   record.Line,
				Column: 1,
			},
		)
	}
}

// addOriginalEmptyIssues reports empty required original column values.
func addOriginalEmptyIssues(analysis *csvLintAnalysis, header []string) {
	if len(header) < 2 {
		return
	}

	headerColumns := len(analysis.Header.Cells)
	for rowIndex := range analysis.Rows {
		record := analysis.Rows[rowIndex]
		if len(record.Cells) != headerColumns {
			continue
		}

		if strings.TrimSpace(record.Cells[1].Value) != "" {
			continue
		}

		analysis.IssuesByCode[CodeLintOriginalEmpty] = append(
			analysis.IssuesByCode[CodeLintOriginalEmpty],
			csvLintIssue{
				Severity: lint.SeverityError,
				Message:  "original column value must be non-empty",
				Line:     record.Line,
				Column:   2,
			},
		)
	}
}

// addTranslationEmptyIssues reports empty translation values after original.
func addTranslationEmptyIssues(analysis *csvLintAnalysis, header []string) {
	if len(header) < 3 {
		return
	}

	headerColumns := len(analysis.Header.Cells)
	for rowIndex := range analysis.Rows {
		record := analysis.Rows[rowIndex]
		if len(record.Cells) != headerColumns {
			continue
		}

		for columnIndex := 2; columnIndex < len(record.Cells); columnIndex++ {
			if strings.TrimSpace(record.Cells[columnIndex].Value) != "" {
				continue
			}

			columnName := strings.TrimSpace(header[columnIndex])
			analysis.IssuesByCode[CodeLintTranslationEmpty] = append(
				analysis.IssuesByCode[CodeLintTranslationEmpty],
				csvLintIssue{
					Severity: lint.SeverityWarning,
					Message: "translation column " +
						strconv.Quote(columnName) + " is empty",
					Line:   record.Line,
					Column: columnIndex + 1,
				},
			)
		}
	}
}

// keyPatternIssues reports key format issues using current rule options.
func keyPatternIssues(
	analysis *csvLintAnalysis,
	run *lint.RunContext,
) ([]csvLintIssue, error) {
	if analysis == nil {
		return nil, nil
	}

	headerColumns := len(analysis.Header.Cells)
	if headerColumns == 0 {
		return nil, nil
	}

	pattern, err := resolveKeyPattern(run)
	if err != nil {
		return nil, err
	}

	out := make([]csvLintIssue, 0, 4)
	for rowIndex := range analysis.Rows {
		record := analysis.Rows[rowIndex]
		if len(record.Cells) != headerColumns || len(record.Cells) == 0 {
			continue
		}

		key := record.Cells[0].Value
		if pattern.MatchString(key) {
			continue
		}

		out = append(out, csvLintIssue{
			Severity: lint.SeverityError,
			Message: "translation key " + strconv.Quote(key) +
				" does not match pattern " + strconv.Quote(pattern.String()),
			Line:   record.Line,
			Column: 1,
		})
	}

	return out, nil
}

// resolveKeyPattern returns compiled regexp from rule options or defaults.
func resolveKeyPattern(run *lint.RunContext) (*regexp.Regexp, error) {
	pattern := DefaultKeyPattern

	if options, ok := lint.GetCurrentRuleOptions[KeyPatternRuleOptions](run); ok {
		if strings.TrimSpace(options.Pattern) != "" {
			pattern = strings.TrimSpace(options.Pattern)
		}
	}

	if optionsMap, ok := lint.GetCurrentRuleOptions[map[string]any](run); ok {
		if value, exists := optionsMap["pattern"]; exists {
			if text, ok := value.(string); ok && strings.TrimSpace(text) != "" {
				pattern = strings.TrimSpace(text)
			}
		}
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKeyPattern, err)
	}

	return compiled, nil
}

// parseCSVQuotedRecords parses CSV text preserving quote token origin flags.
func parseCSVQuotedRecords(content []byte) ([]csvQuotedRecord, error) {
	if len(content) == 0 {
		return nil, nil
	}

	records := make([]csvQuotedRecord, 0, 16)
	record := csvQuotedRecord{
		Line:  1,
		Cells: make([]csvQuotedCell, 0, 8),
	}

	cell := make([]byte, 0, 32)
	quoted := false
	cellStarted := false
	inQuotes := false
	line := 1

	flushCell := func() {
		record.Cells = append(record.Cells, csvQuotedCell{
			Value:  string(cell),
			Quoted: quoted,
		})
		cell = cell[:0]
		quoted = false
		cellStarted = false
	}

	flushRecord := func() {
		if len(record.Cells) == 1 &&
			record.Cells[0].Value == "" &&
			!record.Cells[0].Quoted {
			record = csvQuotedRecord{
				Line:  line,
				Cells: make([]csvQuotedCell, 0, 8),
			}
			return
		}

		records = append(records, record)
		record = csvQuotedRecord{
			Line:  line,
			Cells: make([]csvQuotedCell, 0, 8),
		}
	}

	for index := 0; index < len(content); index++ {
		ch := content[index]
		if inQuotes {
			if ch == '"' {
				if index+1 < len(content) && content[index+1] == '"' {
					cell = append(cell, '"')
					cellStarted = true
					index++
					continue
				}

				inQuotes = false
				continue
			}

			cell = append(cell, ch)
			cellStarted = true
			continue
		}

		switch ch {
		case '"':
			if !cellStarted && len(cell) == 0 {
				inQuotes = true
				quoted = true
				continue
			}

			cell = append(cell, ch)
			cellStarted = true
		case ',':
			flushCell()
		case '\r', '\n':
			flushCell()
			flushRecord()
			if ch == '\r' && index+1 < len(content) && content[index+1] == '\n' {
				index++
			}

			line++
			record.Line = line
		default:
			cell = append(cell, ch)
			cellStarted = true
		}
	}

	if inQuotes {
		return nil, fmt.Errorf("parse csv: unterminated quoted value at line %d", line)
	}

	if len(cell) > 0 || quoted || cellStarted || len(record.Cells) > 0 {
		flushCell()
		flushRecord()
	}

	return records, nil
}

// headerValues converts header record to plain string slice.
func headerValues(header csvQuotedRecord) []string {
	out := make([]string, len(header.Cells))
	for index := range header.Cells {
		out[index] = strings.TrimSpace(header.Cells[index].Value)
	}

	return out
}

// containsString reports whether values slice contains target token.
func containsString(values []string, target string) bool {
	return slices.Contains(values, target)
}
