# LazyFocus CLI Testing Checklist

This document provides a comprehensive test checklist for the LazyFocus CLI Phase 2 implementation.

## Build and Installation

```bash
# Build the binary
go build -o lazyfocus ./cmd/lazyfocus

# Verify build succeeded
./lazyfocus version
# Expected output: lazyfocus version 0.1.0
```

## Command Registration

Verify all commands are registered:

```bash
./lazyfocus --help
```

Expected commands:
- completion
- help
- perspective
- projects
- show
- tags
- tasks
- version

## Global Flags

All commands should support these global flags:
- `--json` - Output in JSON format
- `--quiet` - Suppress output, exit codes only
- `--timeout duration` - Timeout for OmniFocus operations (default 30s)

Test with each command:

```bash
./lazyfocus tasks --help
./lazyfocus projects --help
./lazyfocus tags --help
./lazyfocus show --help
./lazyfocus perspective --help
./lazyfocus version --help
```

## Version Command

```bash
# Test basic version
./lazyfocus version
# Expected: lazyfocus version 0.1.0

# Test with no args (should work)
./lazyfocus version

# Test with args (should fail)
./lazyfocus version extra-arg
# Expected: Error about no args accepted
```

## Tasks Command

```bash
# Test help
./lazyfocus tasks --help

# Test flags
./lazyfocus tasks --all
./lazyfocus tasks --inbox
./lazyfocus tasks --flagged
./lazyfocus tasks --project <project-id>
./lazyfocus tasks --tag <tag-id>
./lazyfocus tasks --due today
./lazyfocus tasks --due tomorrow
./lazyfocus tasks --due 2026-01-30
./lazyfocus tasks --completed

# Test JSON output
./lazyfocus tasks --json
# Expected: {"tasks": [...], "count": N}

# Test quiet mode
./lazyfocus tasks --quiet
# Expected: No output, only exit code
```

## Projects Command

```bash
# Test help
./lazyfocus projects --help

# Test flags
./lazyfocus projects
./lazyfocus projects --status active
./lazyfocus projects --status on-hold
./lazyfocus projects --status completed
./lazyfocus projects --status dropped
./lazyfocus projects --status all
./lazyfocus projects --with-tasks

# Test JSON output
./lazyfocus projects --json
# Expected: {"projects": [...], "count": N}

# Test quiet mode
./lazyfocus projects --quiet
# Expected: No output, only exit code
```

## Tags Command

```bash
# Test help
./lazyfocus tags --help

# Test flags
./lazyfocus tags
./lazyfocus tags --flat
./lazyfocus tags --with-counts

# Test JSON output
./lazyfocus tags --json
# Expected: {"tags": [...], "count": N}

# Test quiet mode
./lazyfocus tags --quiet
# Expected: No output, only exit code
```

## Show Command

```bash
# Test help
./lazyfocus show --help

# Test with task ID
./lazyfocus show <task-id>
./lazyfocus show <task-id> --type task
./lazyfocus show <task-id> --json
# Expected: {"task": {...}}

# Test with project ID
./lazyfocus show <project-id>
./lazyfocus show <project-id> --type project
./lazyfocus show <project-id> --json
# Expected: {"project": {...}}

# Test with tag ID
./lazyfocus show <tag-id>
./lazyfocus show <tag-id> --type tag
./lazyfocus show <tag-id> --json
# Expected: {"tag": {...}}

# Test with invalid ID
./lazyfocus show invalid-id
# Expected: Error exit code 3

# Test with invalid ID (JSON)
./lazyfocus show invalid-id --json
# Expected: {"error": "item not found: invalid-id"}

# Test quiet mode
./lazyfocus show <valid-id> --quiet
# Expected: No output, exit code 0
```

## Perspective Command

```bash
# Test help
./lazyfocus perspective --help

# Test with perspective name
./lazyfocus perspective "Forecast"
./lazyfocus perspective "Review"

# Test JSON output
./lazyfocus perspective "Forecast" --json
# Expected: {"tasks": [...], "count": N}

# Test quiet mode
./lazyfocus perspective "Forecast" --quiet
# Expected: No output, only exit code
```

## Error Handling

Test error scenarios:

```bash
# OmniFocus not running
# (Quit OmniFocus first)
./lazyfocus tasks
# Expected: Error message about OmniFocus not running
# Exit code: 2

./lazyfocus tasks --json
# Expected: {"error": "OmniFocus is not running"}
# Exit code: 2

# Item not found
./lazyfocus show nonexistent-id
# Expected: Error message about item not found
# Exit code: 3

./lazyfocus show nonexistent-id --json
# Expected: {"error": "item not found: nonexistent-id"}
# Exit code: 3

# Invalid flag
./lazyfocus tasks --invalid-flag
# Expected: Error about unknown flag
# Exit code: 1
```

## JSON Output Consistency

Verify JSON output format is consistent across commands:

- Collection commands return: `{"<items>": [...], "count": N}`
  - `tasks` → `{"tasks": [...], "count": N}`
  - `projects` → `{"projects": [...], "count": N}`
  - `tags` → `{"tags": [...], "count": N}`

- Single item commands return: `{"<item>": {...}}`
  - `show <task-id>` → `{"task": {...}}`
  - `show <project-id>` → `{"project": {...}}`
  - `show <tag-id>` → `{"tag": {...}}`

- Error commands return: `{"error": "message"}`

## Quiet Mode Consistency

Verify `--quiet` works for all commands:

```bash
./lazyfocus tasks --quiet
./lazyfocus projects --quiet
./lazyfocus tags --quiet
./lazyfocus show <id> --quiet
./lazyfocus perspective "Forecast" --quiet
```

All should produce no output and only return exit codes.

## Exit Codes

Verify correct exit codes:

- Success: 0
- General error: 1
- OmniFocus not running: 2
- Item not found: 3

## Integration Tests

Run all automated tests:

```bash
# Run all tests
go test ./...

# Run CLI tests specifically
go test ./internal/cli/... -v

# Run integration tests
go test ./internal/cli -run Integration -v
```

## Manual Testing Workflow

1. Start OmniFocus
2. Ensure you have some tasks in your inbox
3. Run each command with various flags
4. Verify output is readable and correct
5. Test with `--json` flag
6. Test with `--quiet` flag
7. Test error conditions
8. Verify exit codes

## AI Agent Testing

Test commands as an AI agent would use them:

```bash
# Get all inbox tasks as JSON
./lazyfocus tasks --json

# Get specific task details
TASK_ID=$(./lazyfocus tasks --json | jq -r '.tasks[0].id')
./lazyfocus show "$TASK_ID" --json

# Get project list
./lazyfocus projects --json

# Get specific project details
PROJECT_ID=$(./lazyfocus projects --json | jq -r '.projects[0].id')
./lazyfocus show "$PROJECT_ID" --json

# Get all tags
./lazyfocus tags --json
```

All commands should return valid JSON that can be parsed programmatically.

## Summary

Phase 2 CLI is complete when:
- [x] All commands are registered
- [x] All commands support --json, --quiet, --timeout flags
- [x] JSON output is consistent across all commands
- [x] Error output is consistent ({"error": "message"})
- [x] Exit codes are correct
- [x] Version command works
- [x] All tests pass
- [x] Build succeeds
- [x] Help text is consistent and clear
