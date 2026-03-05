# stringtable

`stringtable` is a Go module for DayZ `stringtable.csv` files.

It provides:

* CSV parse/read/write (`ParseCSV*`, `FormatCSV*`, `WriteCSVFile*`)
* table model with language columns (`Table`, `Row`)
* language helpers (`DefaultLanguages`, `ParseLanguages`)
* language selection helper (`SelectLanguages`)
* `pofile` integration for `.po/.pot` workflows
* merge/update helpers between CSV and PO catalogs
* PO directory helpers (`ReadCatalogSet`, `WriteCatalogSet`)
* optional build header management (`UpdateBuildHeaders`)

## Install

```bash
go get github.com/woozymasta/stringtable
```

## Quick Example

```go
package main

import (
    "log"

    "github.com/woozymasta/stringtable"
)

func main() {
    table, err := stringtable.ParseCSVFile("stringtable.csv")
    if err != nil {
        log.Fatal(err)
    }

    catalog, err := stringtable.CatalogFromTableLanguage(table, "russian")
    if err != nil {
        log.Fatal(err)
    }

    if err := stringtable.MergeCatalogLanguage(
        table,
        "russian",
        catalog,
        stringtable.MergeOptions{},
    ); err != nil {
        log.Fatal(err)
    }

    if err := stringtable.WriteCSVFile("stringtable_out.csv", table); err != nil {
        log.Fatal(err)
    }
}
```

By default, merge uses the original text as fallback for empty, missing, and
`notranslate` entries. Use `MergeOptions.Disable*` flags to turn off specific
fallback rules.

By default, CSV formatting writes all default DayZ language columns.
Use `WriteOptions.UseTableLanguagesOnly` to keep only current table columns.

`UpdateBuildHeaders` can write common headers, content hash, and update
`PO-Revision-Date` or `POT-Creation-Date` only when content actually changed.
