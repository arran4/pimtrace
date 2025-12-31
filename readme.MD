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
VERSION=<latest release tag>
wget https://github.com/arran4/pimtrace/releases/download/$VERSION/pimtrace_${VERSION}_linux_amd64.deb
sudo apt install ./pimtrace_${VERSION}_linux_amd64.deb
```

### Build from Source

Requires [Go](https://go.dev/) 1.24+.

```bash
git clone https://github.com/arran4/pimtrace.git
cd pimtrace
go build ./cmd/...
```
This will build `csvtrace`, `mailtrace`, and `icaltrace` in the current directory.

## Usage Overview

PIMTrace provides three main tools, each targeting a specific data format:

| Tool | Description |
| :--- | :--- |
| `csvtrace` | For processing tabular data, specifically CSV files. |
| `mailtrace` | For querying email archives (Mbox files) or single email files. |
| `icaltrace` | For filtering and summarizing iCalendar (.ics) files. |

All tools follow a similar usage pattern:

```bash
toolname -input <file> -parser basic [QUERY]
```

*   `-input`: The source file (defaults to stdin `-`).
*   `-parser basic`: **Required.** Specifies the query parser to use.
*   `[QUERY]`: The sequence of operations to perform on the data.

## The Query Language

The "Basic" parser allows you to chain operations together linearly. It reads like a sentence: "Filter this, then put it into a table with these columns, then sort by that."

### Core Operations

*   **`filter <condition>`**: Excludes records that do not match the condition.
*   **`sort <expression>`**: Sorts the results based on the given expression.
*   **`into table <columns...>`**: Transforms the data into a table with the specified columns.
*   **`into summary <columns...> calculate <aggregates...>`**: Groups data by the specified columns and calculates aggregate statistics (like counts or sums).

### Filter Conditions

Conditions are constructed using **operators**:

| Operator | Description | Example |
| :--- | :--- | :--- |
| `eq` | Equality check. | `filter h.From eq .admin@example.com` |
| `contains` | Case-sensitive substring check. | `filter h.Subject contains .Urgent` |
| `icontains` | Case-insensitive substring check. | `filter h.Subject icontains .urgent` |
| `not` | Negates a condition. | `filter not h.Status eq .Closed` |

### Expressions & Functions

You can refer to data fields and transform them using functions:

*   **Columns/Headers**:
    *   `c.columnName` (CSV/Table)
    *   `h.HeaderName` (Mail)
    *   `p.PropertyName` (iCal)
*   **Literals**: Strings starting with a dot (e.g., `.gmail`) are treated as text values.
*   **Functions**: Prefixed with `f.` (e.g., `f.count`, `f.year[c.date]`).

**Note:** Functions do not support spaces between arguments. Use `f.as[h.subject,.Title]`, not `f.as[h.subject, .Title]`.

**Common Functions:**

| Function | Description |
| :--- | :--- |
| `f.count` | Counts the number of items in a group. |
| `f.sum[field]` | Sums the values of a numeric field. |
| `f.year[date]` | Extracts the year from a date string or timestamp. |
| `f.month[date]` | Extracts the month from a date string or timestamp. |
| `f.as[expr,name]` | Renames a column (e.g., `f.as[h.subject,.Title]`). |

See [functions.md](functions.md) for a complete list.

---

## Examples

### 1. CSVTrace: Analyze Expenses & Visualize

Imagine a CSV `expenses.csv` with columns: `Date`, `Category`, `Amount`.

**Task:** Show me all "Food" expenses, showing just the Date and Amount.
```bash
csvtrace -input expenses.csv -parser basic \
  filter c.Category icontains .Food \
  into table c.Date c.Amount
```

**Task:** How much did I spend on each category?
```bash
csvtrace -input expenses.csv -parser basic \
  into summary c.Category calculate f.sum[c.Amount]
```
*Output:*
```text
+----------+------------+
| CATEGORY | SUM_AMOUNT |
+----------+------------+
| Food     |     150.50 |
| Transport|      45.00 |
| Utilities|     120.00 |
+----------+------------+
```

**Task:** Plot a bar chart of spending per category.
```bash
csvtrace -input expenses.csv -parser basic \
  into summary c.Category calculate f.sum[c.Amount] \
  -output-type plot.bar -output expenses_plot.png
```
*Output: Generates `expenses_plot.png` with a bar chart.*

### 2. MailTrace: Summarize Inbox

**Task:** Find all emails from "Amazon" and show the Subject and Date.
```bash
mailtrace -input inbox.mbox -input-type mbox -parser basic \
  filter h.From icontains .Amazon \
  into table h.Date h.Subject
```

**Task:** Who are the top senders in my inbox?
```bash
mailtrace -input inbox.mbox -input-type mbox -parser basic \
  into summary h.From calculate f.count \
  sort f.count
```

### 3. ICalTrace: Calendar Stats

**Task:** List all events containing "Meeting" in the summary.
```bash
icaltrace -input calendar.ics -parser basic \
  filter p.SUMMARY icontains .Meeting \
  into table p.DTSTART p.SUMMARY
```

**Task:** Count how many events I have per month.
```bash
icaltrace -input calendar.ics -parser basic \
  into summary f.year[p.DTSTART] f.month[p.DTSTART] calculate f.count
```
*Output:*
```text
+------+-------+-------+
| YEAR | MONTH | COUNT |
+------+-------+-------+
| 2023 |    10 |    15 |
| 2023 |    11 |    22 |
| 2023 |    12 |     5 |
+------+-------+-------+
```

## Supported Formats

### Inputs
*   **CSV**: Comma Separated Values.
*   **Mbox**: Unix mailbox format (common export format for email).
*   **Mail**: Single email message files (`.eml`).
*   **iCal**: iCalendar files (`.ics`).

### Outputs
*   **Table**: ASCII table (default for human reading).
*   **CSV**: Good for piping into other tools.
*   **Plot**: Generate simple bar charts (e.g., `-output-type plot.bar -output chart.png`).

## FAQ

**Q: Why `*trace`?**
A: It's a nod to `dtrace`. While the functionality is different, the spirit of powerful, on-the-fly instrumentation and querying is similar.

**Q: Why do I need `-parser basic`?**
A: Currently, `basic` is the only implemented parser. It allows for simple, space-separated queries. We require the flag to ensure backward compatibility if/when a more complex parser is introduced.

**Q: Can it handle huge files?**
A: It depends. `csvtrace` and `mailtrace` (in mbox mode) stream data, so they are generally memory-efficient. However, sorting or summarizing requires holding the result set in memory.

## License

This project is licensed under the AGPL.
