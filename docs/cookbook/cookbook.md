# PIMTrace Cookbook

This cookbook provides task-oriented examples for using PIMTrace with different file formats. The examples rely on standard test fixtures available in the repository.

## Table of Contents

- [CSVTrace Examples](#csvtrace-examples)
- [MailTrace Examples](#mailtrace-examples)
- [ICalTrace Examples](#icaltrace-examples)
- [Pipelines (stdin/stdout)](#pipelines-stdinstdout)
- [Discovery Workflow](#discovery-workflow)

---

## CSVTrace Examples

Imagine a CSV file named `expenses.csv` with columns: `Date`, `Category`, `Amount`.

### Select Columns
Extract specific columns into a table:
```bash
csvtrace -input testdata/expenses.csv -input-type csv -parser basic -output-type table \
  into table c.Date c.Amount
```

### Filter Rows
Show all large expenses (over 50):
```bash
csvtrace -input testdata/expenses.csv -input-type csv -parser basic -output-type table \
  filter c.Amount gt .50 \
  into table c.Category c.Amount
```

### Summarize
Calculate the total amount spent per category:
```bash
csvtrace -input testdata/expenses.csv -input-type csv -parser basic -output-type table \
  into summary c.Category calculate f.sum[c.Amount]
```

### Sort
List expenses sorted by amount in descending order:
```bash
csvtrace -input testdata/expenses.csv -input-type csv -parser basic -output-type table \
  sort c.Amount desc \
  into table c.Date c.Category c.Amount
```

### Amount/Date Ranges
Find expenses within a specific range:
```bash
csvtrace -input testdata/expenses.csv -input-type csv -parser basic -output-type table \
  filter c.Amount gt .10 and c.Amount lte .100 \
  into table c.Date c.Category c.Amount
```

---

## MailTrace Examples

Using an mbox file named `inbox.mbox`.

### Find Subjects/Senders
List the subjects and senders of all emails:
```bash
mailtrace -input testdata/inbox.mbox -input-type mbox -parser basic -output-type table \
  into table h.From h.Subject
```

### Top Senders
Find the top 10 senders in the inbox:
```bash
mailtrace -input testdata/inbox.mbox -input-type mbox -parser basic -output-type table \
  into summary h.From calculate f.count \
  sort c.count desc limit 10
```

### Year/Month Summaries
Count the number of emails received per year and month:
```bash
mailtrace -input testdata/inbox.mbox -input-type mbox -parser basic -output-type table \
  into summary f.year[h.Date] f.month[h.Date] calculate f.count
```

---

## ICalTrace Examples

Using an iCalendar file named `calendar.ics`.

### Filter Events
List all events containing "Meeting" in the summary:
```bash
icaltrace -input testdata/calendar.ics -input-type ical -parser basic -output-type table \
  filter p.SUMMARY icontains .Meeting \
  into table p.DTSTART p.SUMMARY
```

### Monthly Summaries
Count how many events occur per month:
```bash
icaltrace -input testdata/calendar.ics -input-type ical -parser basic -output-type table \
  into summary f.year[p.DTSTART] f.month[p.DTSTART] calculate f.count
```

### Timezone-Aware Examples
PIMTrace automatically respects local timezones defined in the iCalendar file for date-based functions:
```bash
icaltrace -input testdata/calendar.ics -input-type ical -parser basic -output-type table \
  into summary f.date[p.DTSTART] calculate f.count
```

### Weekday / Hour / Duration Examples
Group meetings by weekday and hour to find the busiest times:
```bash
icaltrace -input testdata/calendar.ics -input-type ical -parser basic -output-type table \
  into summary f.weekday[p.DTSTART] f.hour[p.DTSTART] calculate f.count \
  sort c.count desc limit 5
```

Find the 5 longest meetings using the `f.duration` function:
```bash
icaltrace -input testdata/calendar.ics -input-type ical -parser basic -output-type table \
  sort f.duration desc limit 5 \
  into table p.SUMMARY f.duration
```

Get total meeting time (in seconds) per day:
```bash
icaltrace -input testdata/calendar.ics -input-type ical -parser basic -output-type table \
  into summary f.date[p.DTSTART] calculate f.sum[f.duration] \
  sort c.date-DTSTART asc \
  into table c.date-DTSTART c.sum-duration
```

---

## Pipelines (stdin/stdout)

PIMTrace tools natively support standard streams for composability.

### Using Standard Input/Output
Use `-` as the input path to read from stdin, and pipe stdout to another command:
```bash
cat testdata/expenses.csv | csvtrace -input - -input-type csv -parser basic -output-type table \
  filter c.Amount gt .10 \
  into table c.Category c.Amount
```

### CSV Output
Output the results as CSV for downstream processing:
```bash
csvtrace -input testdata/expenses.csv -input-type csv -parser basic -output-type csv \
  into summary c.Category calculate f.sum[c.Amount]
```

---

## Discovery Workflow

If you have an unfamiliar file and don't know what columns or headers you can query, use the `-describe` flag. This discovers available fields and infers their data types by sampling up to the first 1000 records.

### Describe then Query
1. **Describe the file:**
   ```bash
   csvtrace -input testdata/expenses.csv -input-type csv -describe
   ```
2. **Copy an expression from the output.** For example, `c.Amount`.
3. **Use it in a normal query:**
   ```bash
   csvtrace -input testdata/expenses.csv -input-type csv -parser basic -output-type table \
     filter c.Amount gt .100 \
     into table c.Category c.Amount
   ```
