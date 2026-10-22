# Agent guidance

## CLI and data boundaries
- Treat input parsing, query evaluation and output formatting as separate behaviours. Preserve existing user-visible stream usage and documented CLI exit codes; inspect the README's "Exit Codes and Streams" section before modifying commands.
- Avoid hidden process state and synthetic XDG setups in ordinary logic tests. Inject readers, writers, clocks or data values where they clarify the scenario. Use temporary files and process streams intentionally when validating actual CLI or OS integration.
- Do not silently drop invalid input records or change strict/best-effort semantics. Characterise the intended behaviour and coordinate with issue #56 if implementation is needed.
- Compare meaningful parsed results and output bytes, including rejected input, instead of testing only that a fixture loads.

## Test and review workflow
- Prefer self-explanatory named CSV/mail/iCalendar fixtures over opaque encoded literals. Leave integration tests that assert real filesystem or subprocess behaviours in place.
- Reconcile issue #82 before proposing more filesystem fixture refactoring.
- Keep README usage/format descriptions aligned with changed observable CLI behaviour.

A coding agent must inspect source and tests, report untracked findings to issue #88, run the relevant Go checks, and keep implementation separate from this documentation-only adoption PR.
## Engineering reference and follow-up

- Adopt only relevant guidance from https://arran4.github.io/blog/post/2026/055-suggested-coding-standards-for-repositories-and-agents/ and https://arran4.github.io/blog/post/2026/020-optional-dependency-injection-via-type-switched-variadic-args/ .
- This documentation is not evidence of repository-wide compliance. A coding-agent pass must report checked files, observed gaps, tests/CI and existing issue overlap to https://github.com/arran4/pimtrace/issues/88 before any broad implementation.
- Do not merge implementation automatically; request review and follow repository-specific instructions.
