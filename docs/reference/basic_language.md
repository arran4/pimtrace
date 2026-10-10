# Basic Query Language Reference

The "Basic" parser allows you to construct queries by chaining operations linearly. Each operation receives data from the previous one, similar to Unix pipelines.

## Operation Grammar

A query consists of one or more operations executed sequentially:

`<operation> <arguments> [ <operation> <arguments> ... ]`

### Supported Operations

#### `filter` or `where`
Excludes records that do not match the specified condition.

**Syntax:** `filter <condition>` or `where <condition>`

The `<condition>` is a boolean expression evaluating to true or false.

#### `into table`
Transforms the data into a tabular format, keeping only the specified columns.

**Syntax:** `into table <column_expr> [ <column_expr> ... ]`

#### `into summary`
Groups the incoming data by the specified columns and allows for aggregate calculations. Used in conjunction with `calculate`.

**Syntax:** `into summary <group_by_expr> [ <group_by_expr> ... ] calculate <aggregate_expr> [ <aggregate_expr> ... ]`

#### `calculate`
Specifies the aggregate functions to apply to the groups created by `into summary`.

**Syntax:** `calculate <aggregate_expr> [ <aggregate_expr> ... ]`

#### `sort`
Sorts the output rows based on a specific expression.

**Syntax:** `sort <expression> [ asc | desc ]`

The default sort direction is `asc` (ascending). Use `desc` for descending order.

#### `limit`
Limits the number of output rows to the specified integer. Must be placed at the end of the query or after `sort`.

**Syntax:** `limit <N>` (where N is a positive integer)

---

## Expressions

Expressions represent data points from the underlying records, transformed values, or literal constants.

### Field References
*   `c.<ColumnName>` : Refers to a CSV or general table column (e.g., `c.Amount`).
*   `h.<HeaderName>` : Refers to a Mail header field (e.g., `h.From`).
*   `p.<PropertyName>` : Refers to an iCalendar property (e.g., `p.DTSTART`).

### Functions
Functions transform values or extract information. They are prefixed with `f.` and take arguments in brackets `[]` without spaces.
*   **Syntax:** `f.<function_name>[<arg1>,<arg2>]`
*   **Example:** `f.year[p.DTSTART]`

For a full list of available functions, see [functions.md](../../functions.md). Generated aggregate values (like `f.sum[c.Amount]`) produce columns named according to the formula, e.g., `c.sum-Amount`. Wait until output generation to reference them by these hyphenated names, or use `f.as[expr,name]` to rename.

### Literals
String and numeric constants must be prefixed with a dot (`.`).
*   **Example:** `.Meeting` or `.50`

---

## Filter Conditions & Operators

Filter conditions use comparison operators to evaluate data, and boolean operators to combine multiple conditions.

### Comparison Operators

| Operator | Description | Type Behavior |
| :--- | :--- | :--- |
| `eq` | Textual equality check. | Textual |
| `contains` | Case-sensitive substring check. | Textual |
| `icontains` | Case-insensitive substring check. | Textual |
| `gt` | Greater than. | Numeric or Date |
| `gte` | Greater than or equal to. | Numeric or Date |
| `lt` | Less than. | Numeric or Date |
| `lte` | Less than or equal to. | Numeric or Date |

### Typed Comparison & Coercion Rules

When using relational operators (`gt`, `gte`, `lt`, `lte`), PIMTrace performs typed numeric (integer/decimal) or chronological date comparisons.

*   Operands must successfully coerce to the expected type (number or date).
*   If an operand cannot be coerced (e.g., comparing text like "apple" or an empty string against a number), evaluation produces an explicit error, rather than silently falling back to a zero value or lexical ordering.
*   Textual operators (`eq`, `contains`, `icontains`) always compare values lexically as strings (e.g., `"01"` is not equal to `"1"` under `eq`).

### Null/Missing Field Behavior

When referring to a field (e.g., `c.MissingField` or `h.NonExistentHeader`) that does not exist in the record, it typically evaluates to an empty string `""`. However, applying numeric/date relational operators to a missing field will trigger a coercion error.

### Boolean Operators & Precedence

You can combine conditions using boolean operators.

*   **`not`**: Unary logical NOT. Binds most tightly.
*   **`and`**: Logical AND. Binds more tightly than `or`.
*   **`or`**: Logical OR. Binds least tightly.

**Precedence:** `not` > `and` > `or`.
For example, `A or B and C` is evaluated as `A or (B and C)`.

#### Grouping
Use parentheses `(` and `)` to override default precedence: `( A or B ) and C`.
*Note on Shell usage:* Unquoted parentheses trigger subshell syntax in Bash/Zsh. You must escape them with backslashes (e.g., `\(`, `\)`) or pass them quoted (e.g., `'('`, `')'`).

---

## Execution Environment

### Standard Stream Contract

PIMTrace CLI tools (`csvtrace`, `mailtrace`, `icaltrace`) strictly adhere to standard stream decoupling:
*   **`stdout` (Standard Output):** Reserved exclusively for the actual query results, output formats (like `table` or `csv`), or explicit help/version output. Safe to pipe to downstream tools like `grep` or `jq`.
*   **`stderr` (Standard Error):** Used for all diagnostic messages, progress indicators, warnings, parse errors, and runtime errors.

### Error and Exit Behavior

*   **Exit Code 0**: Success (query executed successfully, or help/version requested).
*   **Exit Code 1**: Runtime errors (e.g., failed to read input file, filter evaluation coercion error, division by zero).
*   **Exit Code 2**: Usage errors (e.g., invalid command-line flags, invalid basic query syntax).

### Functions Reference
A full list of available functions is maintained in `functions.md` at the root of the repository.

To regenerate this file and ensure it matches the current source code capabilities, run:
```bash
go run cmd/docs/genfunctionmd/main.go
```
