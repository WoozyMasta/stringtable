// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/stringtable

/*
Package stringtable provides DayZ stringtable.csv primitives and pofile
integration.

Use ParseCSVFile and WriteCSVFile for CSV workflows, ParseLanguages and
DefaultLanguages for language handling, and the pofile integration helpers to
convert between CSV rows and gettext PO catalogs. Use UpdateBuildHeaders to
optionally manage standard PO/POT build headers and content hash.

Package uses lintkit for CSV lint diagnostics.
*/
package stringtable
