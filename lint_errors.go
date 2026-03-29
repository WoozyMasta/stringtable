// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"errors"

	"github.com/woozymasta/lintkit/lint"
)

var (
	// ErrNilLintRuleRegistrar indicates nil lint rule registrar in registration.
	ErrNilLintRuleRegistrar = lint.ErrNilRuleRegistrar

	// ErrInvalidKeyPattern reports invalid key pattern regexp in lint options.
	ErrInvalidKeyPattern = errors.New("invalid lint key pattern")
)
