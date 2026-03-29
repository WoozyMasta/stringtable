<!-- Automatically generated file, do not modify! -->

# Lint Rules Registry

This document contains the current registry of lint rules.

Total rules: 11.

## stringtable

### csv

> Stringtable CSV structure and translation diagnostics.

Codes:
[STBL2001](#stbl2001),
[STBL2002](#stbl2002),
[STBL2003](#stbl2003),
[STBL2004](#stbl2004),
[STBL2005](#stbl2005),
[STBL2006](#stbl2006),
[STBL2007](#stbl2007),
[STBL2008](#stbl2008),
[STBL2009](#stbl2009),
[STBL2010](#stbl2010),
[STBL2011](#stbl2011),

#### `STBL2001`

Header language column name must be non-empty

> Every translation column after `Language,original` must have a non-empty
> language name.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.header-language-column-name-must-be-non-empty` |
| Scope | `csv` |
| Severity | `error` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

#### `STBL2002`

Header language columns must be unique

> Each language column name should appear once. Duplicate names make column
> mapping ambiguous for parsers and exporters.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.header-language-columns-must-be-unique` |
| Scope | `csv` |
| Severity | `error` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

#### `STBL2003`

Header contains unknown language name

> Language columns after `Language,original` must use supported DayZ language
> names. Rename unsupported columns or update conversion flow.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.header-contains-unknown-language-name` |
| Scope | `csv` |
| Severity | `error` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

#### `STBL2004`

Header is missing default languages

> Add missing DayZ default language columns to keep expected export order and
> avoid incomplete localization coverage.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.header-is-missing-default-languages` |
| Scope | `csv` |
| Severity | `warning` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

#### `STBL2005`

Row column count must match header column count

> Every data row must have the same number of cells as the header; extra or
> missing cells shift translations into wrong languages.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.row-column-count-must-match-header-column-count` |
| Scope | `csv` |
| Severity | `error` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

#### `STBL2006`

Translation key must be unique

> Duplicate `Language` keys create ambiguous lookup and merge behavior. Keep
> exactly one row per key.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.translation-key-must-be-unique` |
| Scope | `csv` |
| Severity | `error` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

#### `STBL2007`

Translation key has surrounding spaces

> Keys with surrounding spaces are hard to spot and may behave as different tokens
> than visually similar trimmed keys.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.translation-key-has-surrounding-spaces` |
| Scope | `csv` |
| Severity | `warning` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

#### `STBL2008`

Translation key must match key pattern

> Keep keys machine-safe and deterministic. Change `pattern` option only when
> project naming rules intentionally differ.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.translation-key-must-match-key-pattern` |
| Scope | `csv` |
| Severity | `error` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

Default options:
```json
{
  "pattern": "^[-_A-Za-z0-9]+$"
}
```

#### `STBL2009`

`Original` column must be non-empty

> `original` stores the source text. Empty source usually means broken authoring
> or accidental row damage.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.original-column-must-be-non-empty` |
| Scope | `csv` |
| Severity | `error` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

#### `STBL2010`

Translation columns should be non-empty

> Empty translation is allowed but likely unfinished localization. Fill value or
> intentionally suppress this warning in your workflow.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.translation-columns-should-be-non-empty` |
| Scope | `csv` |
| Severity | `warning` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

#### `STBL2011`

Non-key cells should be quoted

> Quote `original` and translation cells to keep CSV stable when text contains
> commas, quotes, or leading/trailing spaces.

| Field | Value |
| --- | --- |
| Rule ID | `stringtable.csv.non-key-cells-should-be-quoted` |
| Scope | `csv` |
| Severity | `warning` |
| Enabled | `true` (implicit) |
| File kinds | stringtable.csv |

---

> Generated with
> [lintkit](https://github.com/woozymasta/lintkit)
> version `dev`
> commit `unknown`

<!-- Automatically generated file, do not modify! -->
