// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

package stringtable

import (
	"time"

	"github.com/woozymasta/pofile"
)

const (
	// DefaultContentHashHeader is the default PO header for content hash.
	DefaultContentHashHeader = "X-Content-Hash"

	// DefaultSourceHashHeader is the default PO header for source hash.
	DefaultSourceHashHeader = "X-CSV-Hash"
)

// HeaderOptions controls optional build/header updates for PO/POT catalogs.
type HeaderOptions struct {
	// ProjectVersion sets Project-Id-Version when non-empty.
	ProjectVersion string `json:"project_version,omitempty" yaml:"project_version,omitempty"`

	// Generator sets X-Generator when non-empty.
	Generator string `json:"generator,omitempty" yaml:"generator,omitempty"`

	// SourceHash writes source hash to source hash header when non-empty.
	SourceHash string `json:"source_hash,omitempty" yaml:"source_hash,omitempty"`

	// ContentHashHeader overrides content hash header name.
	ContentHashHeader string `json:"content_hash_header,omitempty" yaml:"content_hash_header,omitempty"`

	// SourceHashHeader overrides source hash header name.
	SourceHashHeader string `json:"source_hash_header,omitempty" yaml:"source_hash_header,omitempty"`

	// WriteStandardHeaders writes MIME and language headers when missing.
	WriteStandardHeaders bool `json:"write_standard_headers,omitempty" yaml:"write_standard_headers,omitempty"`

	// WriteHash writes content hash into hash header.
	WriteHash bool `json:"write_hash,omitempty" yaml:"write_hash,omitempty"`

	// WriteDateOnChange updates creation/revision date only when content changed.
	WriteDateOnChange bool `json:"write_date_on_change,omitempty" yaml:"write_date_on_change,omitempty"`

	// Template switches date header to POT-Creation-Date (otherwise PO-Revision-Date).
	Template bool `json:"template,omitempty" yaml:"template,omitempty"`
}

// UpdateBuildHeaders updates optional build headers and returns whether
// effective catalog content changed (by content hash).
func UpdateBuildHeaders(catalog *pofile.Catalog, options HeaderOptions) bool {
	if catalog == nil {
		return false
	}

	hashHeader := options.ContentHashHeader
	if hashHeader == "" {
		hashHeader = DefaultContentHashHeader
	}

	sourceHashHeader := options.SourceHashHeader
	if sourceHashHeader == "" {
		sourceHashHeader = DefaultSourceHashHeader
	}

	if options.ProjectVersion != "" {
		catalog.SetHeader("Project-Id-Version", options.ProjectVersion)
	}
	if options.Generator != "" {
		catalog.SetHeader("X-Generator", options.Generator)
	}
	if options.WriteStandardHeaders {
		applyStandardHeaders(catalog)
	}
	if options.SourceHash != "" {
		catalog.SetHeader(sourceHashHeader, options.SourceHash)
	}

	oldHash := catalog.Header(hashHeader)
	newHash := catalog.ContentHash()
	changed := oldHash == "" || oldHash != newHash

	if options.WriteDateOnChange && changed {
		now := time.Now().Format("2006-01-02 15:04-0700")
		if options.Template {
			catalog.SetHeader("POT-Creation-Date", now)
		} else {
			catalog.SetHeader("PO-Revision-Date", now)
		}
	}
	if options.WriteHash {
		catalog.SetHeader(hashHeader, newHash)
	}

	return changed
}

// applyStandardHeaders applies common gettext headers when they are missing.
func applyStandardHeaders(catalog *pofile.Catalog) {
	if catalog.Header("Language") == "" && catalog.Language != "" {
		catalog.SetHeader("Language", catalog.Language)
	}
	if catalog.Header("MIME-Version") == "" {
		catalog.SetHeader("MIME-Version", "1.0")
	}
	if catalog.Header("Content-Type") == "" {
		catalog.SetHeader("Content-Type", "text/plain; charset=UTF-8")
	}
	if catalog.Header("Content-Transfer-Encoding") == "" {
		catalog.SetHeader("Content-Transfer-Encoding", "8bit")
	}
}
