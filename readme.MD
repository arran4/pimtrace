# PIMTrace Tools

PIMTrace is a collection of command-line utilities designed to make querying, filtering, and summarizing structured data (like CSV, Email, and iCal) as easy as writing a simple sentence. It aims to bridge the gap between simple shell one-liners and complex query languages (like SQL), providing a "batteries included" experience for data exploration.

## Status

**Active Development / Beta.**
This project is currently being developed to solve specific personal needs but is open to contributions. It is stable enough for the documented use cases, but you may encounter edge cases or missing features.

## Goals

*   **Simplicity**: Provide a query language that feels natural and requires minimal syntax (no complex quoting or escaping).
*   **Utility**: Handle common PIM (Personal Information Management) formats like iCal, Mbox, and CSV out of the box.
*   **Shell Friendliness**: Work well with pipes and standard streams.

## Installation

### Pre-built Binaries

Pre-built binaries are available on the [Releases Page](https://github.com/arran4/pimtrace/releases).
1.  Download the package for your OS (Linux, macOS, Windows).
2.  Extract the archive.
3.  Place the binaries (`csvtrace`, `mailtrace`, `icaltrace`) in a directory included in your `$PATH`.

**Debian/Ubuntu Example:**
```bash
VERSION=$(curl -s https://api.github.com/repos/arran4/pimtrace/releases/latest | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
wget https://github.com/arran4/pimtrace/releases/download/v$VERSION/pimtrace_${VERSION}_linux_amd64.deb
sudo apt install ./pimtrace_${VERSION}_linux_amd64.deb
```

### Build from Source

Requires [Go](https://go.dev/) 1.25.3+.

```bash
git clone https://github.com/arran4/pimtrace.git
cd pimtrace
mkdir -p bin
go build -o bin/csvtrace ./cmd/csvtrace
go build -o bin/mailtrace ./cmd/mailtrace
go build -o bin/icaltrace ./cmd/icaltrace
```
This will build `csvtrace`, `mailtrace`, and `icaltrace` in the `bin/` directory.

## Usage Overview

PIMTrace provides three main tools, each targeting a specific data format:

| Tool | Description |
| :--- | :--- |
| `csvtrace` | For processing tabular data, specifically CSV files. |
| `mailtrace` | For querying email archives (Mbox files) or single email files. |
| `icaltrace` | For filtering and summarizing iCalendar (.ics) files. |

All tools follow a similar usage pattern:

```bash
toolname -input <file> -input-type <type> -parser basic -output-type <type> [QUERY]
```

*   `-input`: The source file (defaults to stdin `-`).
*   `-input-type <type>`: Specifies the input format. (Defaults to `list`, which prints available types and exits).
*   `-output-type <type>`: Specifies the output format. (Defaults to `list`, which prints available types and exits).
*   `-parser basic`: **Required.** Specifies the query parser to use.
*   `[QUERY]`: The sequence of operations to perform on the data.

## Exit Codes and Streams

PIMTrace tools follow a consistent contract for exit codes and standard streams:

### Exit Codes
*   **`0` (Success):** The query executed successfully. Explicit help (`-help`) or version (`-version`) requests also exit with `0`.
*   **`1` (Runtime Error):** An error occurred while reading input, parsing data, executing the query, or writing output.
*   **`2` (Usage Error):** Invalid flags, missing required arguments (like the query itself), or invalid query syntax.

### Streams
*   **`stdout` (Standard Output):** Reserved exclusively for the actual query results, or explicit help/version output. This ensures that the output is safe to pipe to other tools or redirect to files.
*   **`stderr` (Standard Error):** Used for all diagnostic messages, progress indicators, warnings, and error messages.


### Discovering Fields with "-describe"

If you have an unfamiliar file and don't know what columns or headers you can query, you can use the "-describe" flag. This discovers available fields, infer their data types by sampling up to the first 1000 records, and prints the exact query expression you need to use.

1. **Describe the file:**
   ```bash
   csvtrace -input unknown_data.csv -input-type csv -describe
   ```
   *Output:*
   ```text
   Note: Values are sampled up to a maximum of 1000 records to infer types.

   Component	Name      Expression   Observed Types   Matches (Sampled)
   ---------	----      ----------   --------------   -----------------
   table    	ID        c.ID         string           10/10
   table    	Name      c.Name       string           10/10
   table    	Balance   c.Balance    string           10/10
   ```

2. **Copy an expression from the output.** For example, "c.Balance".

3. **Use it in a normal query:**
   ```bash
   csvtrace -input unknown_data.csv -input-type csv -parser basic -output-type table \
     filter c.Balance gt .100 \
     into table c.Name c.Balance
   ```

You can also use this with standard input by passing "-input -":
```bash
cat data.csv | csvtrace -input - -input-type csv -describe
```

## Documentation

For comprehensive information on how to construct queries, use functions, and see examples for various file types, please consult our documentation:

*   **[Basic Query Language Reference](docs/reference/basic_language.md)**: A systematic guide to the query language syntax, operations, expressions, and type behaviors.
*   **[PIMTrace Cookbook](docs/cookbook/cookbook.md)**: A collection of task-oriented examples for `csvtrace`, `mailtrace`, and `icaltrace`, as well as pipeline integration and discovery workflows.
*   **[Functions Reference](functions.md)**: The definitive, auto-generated list of all supported functions (e.g., `f.sum`, `f.year`).

## Supported Formats

### Inputs
*   **CSV**: Comma Separated Values.
*   **Mbox**: Unix mailbox format (common export format for email).
*   **Mail**: Single email message files (`.eml`).
*   **iCal**: iCalendar files (`.ics`).

### Outputs
*   **Table**: ASCII table.
*   **CSV**: Good for piping into other tools.
*   **Plot**: Generate simple bar charts (e.g., `-output-type plot.bar -output chart.png`).
*   **List**: Prints available formats (this is the default behavior if `-input-type` or `-output-type` are omitted).

## FAQ

**Q: Why `*trace`?**
A: It's a nod to `dtrace`. While the functionality is different, the spirit of powerful, on-the-fly instrumentation and querying is similar.

**Q: Why do I need `-parser basic`?**
A: Currently, `basic` is the only implemented parser. It allows for simple, space-separated queries. We require the flag to ensure backward compatibility if/when a more complex parser is introduced.

**Q: Can it handle huge files?**
A: It depends. `csvtrace` and `mailtrace` (in mbox mode) stream data, so they are generally memory-efficient. However, sorting or summarizing requires holding the result set in memory.

## License

This project is licensed under the AGPL.
