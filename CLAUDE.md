# LazyFocus - Project Instructions

## Project Overview

LazyFocus (`lazyfocus` / `lf`) is a CLI and TUI tool that interfaces with OmniFocus on macOS using Omni Automation (JavaScript for Automation). It serves two audiences:

1. **Humans** — quick terminal access to OmniFocus with readable output
2. **AI Agents** — structured JSON interface for LLMs to query and manipulate tasks

## Technical Stack

- **Language:** Go
- **CLI Framework:** Cobra (`github.com/spf13/cobra`)
- **TUI Framework:** Bubble Tea (`github.com/charmbracelet/bubbletea`)
- **UI Components:** Bubbles (`github.com/charmbracelet/bubbles`)
- **Styling:** Lip Gloss (`github.com/charmbracelet/lipgloss`)
- **OmniFocus Interface:** Omni Automation via `osascript -l JavaScript`
- **Output Formats:** Human-readable (default), JSON (`--json` flag)

## Project Structure

```
lazyfocus/
├── cmd/lazyfocus/main.go          # Single entrypoint
├── internal/
│   ├── bridge/                    # Omni Automation execution layer
│   │   ├── executor.go            # osascript wrapper
│   │   ├── scripts.go             # Embedded JS scripts
│   │   └── parser.go              # JSON response parsing
│   ├── domain/                    # Shared domain models
│   │   ├── task.go
│   │   ├── project.go
│   │   ├── tag.go
│   │   └── perspective.go
│   ├── cli/                       # Cobra command implementations
│   │   ├── root.go
│   │   ├── tasks.go
│   │   ├── projects.go
│   │   ├── add.go
│   │   ├── complete.go
│   │   ├── modify.go
│   │   └── output.go              # Human vs JSON formatting
│   └── tui/                       # Bubble Tea TUI
│       ├── app.go                 # Main model, orchestration
│       ├── keys.go                # Keybinding definitions
│       ├── styles.go              # Lip Gloss styles
│       ├── components/            # Reusable UI components
│       └── views/                 # Screen implementations
└── scripts/                       # Raw Omni Automation JS (reference/testing)
```

## Subagent Selection

When implementing features, use the appropriate subagents:

| Task Type | Subagent |
|-----------|----------|
| Go code (services, domain, CLI) | `go-tdd-expert` |
| TUI components (Bubble Tea) | `go-bubbletea-expert` |
| Omni Automation scripts | `omni-automation-expert` |
| Architecture decisions | `solutions-architect` |
| Codebase exploration | `Explore` |

## Development Phases

### Phase 1: Foundation & Bridge Layer
- Go module setup with dependencies
- Omni Automation bridge (execute JS, parse responses)
- Domain models (Task, Project, Tag)
- Verify OmniFocus communication

### Phase 2: CLI Commands (Read Operations)
- Cobra CLI structure
- Read-only commands: `tasks`, `projects`, `tags`, `show`, `perspective`
- Human and JSON output formatting
- Filtering and querying support

### Phase 3: CLI Commands (Write Operations)
- Task creation with natural syntax
- Task modification
- Completion/deletion
- Natural date parsing

### Phase 4: TUI - Basic Structure
- Bubble Tea application shell
- Basic navigation between views
- Task list component
- Quick add overlay

### Phase 5: TUI - Full Implementation
- All views (Inbox, Projects, Tags, Forecast, Review)
- Search/filter functionality
- All task actions within TUI
- Vim-style command mode

### Phase 6: Polish & Distribution
- Error handling and edge cases
- Configuration file support
- Shell completions
- Homebrew formula

## Go Development Standards

### Code Organization
- Keep packages focused on single responsibility
- Use `internal/` for non-exported packages
- Prefer composition over inheritance
- Make dependencies explicit via constructor injection

### Error Handling
- Return errors, don't panic
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Define sentinel errors for expected conditions
- Handle OmniFocus-specific errors gracefully:
  - OmniFocus not installed
  - OmniFocus not running
  - Invalid task/project IDs
  - Automation permission issues

**Common error scenarios:**
- **Empty task name:** "task name is required"
- **Task not found:** "task not found" (when using invalid task ID)
- **Project not found:** "failed to resolve project: project not found"
- **Invalid date format:** "invalid due date: unrecognized date format: xyz"
- **No modifications:** "no modifications specified" (modify command without flags)
- **Missing confirmation:** "confirmation required: use --force to delete" (delete without --force)

In JSON mode, errors return `{"error": "message"}` with appropriate exit codes.

### Testing
- Follow TDD: Red → Green → Refactor
- Use table-driven tests for multiple cases
- Mock external dependencies (osascript calls)
- Test both human and JSON output formatting
- Use `testify` for assertions if needed

### Naming Conventions
- Use idiomatic Go naming (camelCase for unexported, PascalCase for exported)
- Interfaces end with `-er` suffix when appropriate
- Keep names concise but descriptive

## Omni Automation Guidelines

### Script Execution
```go
// Execute via osascript
cmd := exec.Command("osascript", "-l", "JavaScript", "-e", script)
```

### Script Structure
```javascript
(() => {
  const app = Application("OmniFocus");
  const doc = app.defaultDocument;
  // ... operations ...
  return JSON.stringify(result);
})();
```

### Return Format
- Always return JSON from scripts
- Include error information in JSON when errors occur
- Use consistent date formats (ISO 8601)

## CLI Command Reference

### Read Commands

- `tasks` - List tasks with filtering support
- `projects` - List all projects
- `tags` - List all tags
- `show` - Show task details
- `perspective` - View custom perspectives

### Write Commands

#### `add` - Create a new task

**Natural Syntax:**
```bash
lazyfocus add "Buy milk #groceries due:tomorrow"
lazyfocus add "Review PR @Work due:friday !"
lazyfocus add "Meeting prep @\"Big Project\" defer:\"next monday\""
```

**Flags:**
```bash
lazyfocus add "Task name" --project Work --tag urgent --due tomorrow --flagged
lazyfocus add "Quick task" -p Work -t urgent -t followup -d friday -f
lazyfocus add "With note" --note "Additional details here"
```

**Available flags:**
- `-p, --project <name>` - Project name or ID
- `-t, --tag <name>` - Tags (repeatable)
- `-d, --due <date>` - Due date
- `--defer <date>` - Defer date
- `-f, --flagged` - Mark as flagged
- `-n, --note <text>` - Task note

Command-line flags override natural syntax when both are present.

#### `complete` - Mark tasks as complete

```bash
lazyfocus complete abc123
lazyfocus complete task1 task2 task3
lazyfocus complete abc123 --json
```

Accepts multiple task IDs. Continues processing even if some tasks fail.

#### `delete` - Delete tasks

```bash
lazyfocus delete abc123 --force
lazyfocus delete task1 task2 task3 --force
lazyfocus delete abc123 --json
```

**Flags:**
- `-f, --force` - Skip confirmation prompt

In JSON mode, confirmation is automatically skipped. Multiple task IDs supported.

#### `modify` - Modify existing task

```bash
lazyfocus modify task123 --name "Updated name"
lazyfocus modify task123 --due tomorrow --flagged true
lazyfocus modify task123 --add-tag urgent --remove-tag low
lazyfocus modify task123 --clear-due --clear-defer
lazyfocus modify task123 --project Work --note "New note"
```

**Available flags:**
- `--name <text>` - New task name
- `--note <text>` - New note
- `--project <name>` - Move to project (name or ID)
- `--add-tag <name>` - Add tag (repeatable)
- `--remove-tag <name>` - Remove tag (repeatable)
- `--due <date>` - Set due date
- `--defer <date>` - Set defer date
- `--flagged <true|false>` - Set flagged status
- `--clear-due` - Clear due date
- `--clear-defer` - Clear defer date

Requires at least one modification flag.

### Natural Syntax Guide

The `add` command supports natural language task input:

**Tags:** `#tagname` or `#"tag with spaces"`
```bash
lazyfocus add "Buy groceries #errands #shopping"
lazyfocus add "Team sync #\"project alpha\""
```

**Projects:** `@projectname` or `@"project with spaces"`
```bash
lazyfocus add "Review code @Work"
lazyfocus add "Planning meeting @\"Big Project\""
```

**Due dates:** `due:date` or `due:"date phrase"`
```bash
lazyfocus add "Submit report due:friday"
lazyfocus add "Call client due:\"next monday\""
```

**Defer dates:** `defer:date` or `defer:"date phrase"`
```bash
lazyfocus add "Review proposal defer:tomorrow"
lazyfocus add "Follow up defer:\"in 3 days\""
```

**Flagged:** `!` anywhere in input
```bash
lazyfocus add "Urgent task !"
lazyfocus add "! High priority item"
```

**Supported date formats:**
- Relative: `today`, `tomorrow`, `yesterday`
- Next occurrence: `next monday`, `next week`
- In N units: `in 3 days`, `in 2 weeks`
- ISO format: `2024-01-15`
- Month/day: `Jan 15`, `January 15 2024`

All dates without explicit times default to 5:00 PM local time.

## CLI Output Standards

### Human Output
```
INBOX (3 tasks)
───────────────────────────────────────
☐ Buy groceries                    📅 Today
  #errands
☐ Review PR #142                   🚩
  #work #code-review
```

### JSON Output

**List tasks:**
```json
{
  "tasks": [...],
  "count": 3
}
```

**Created task:**
```json
{
  "id": "abc123",
  "name": "Buy milk",
  "tags": ["groceries"],
  "due": "2024-01-16T17:00:00Z",
  "flagged": false
}
```

**Modified task:**
```json
{
  "id": "abc123",
  "name": "Updated task name",
  "project": "Work",
  "tags": ["urgent"],
  "due": "2024-01-20T17:00:00Z"
}
```

**Completed/deleted task:**
```json
{
  "id": "abc123",
  "name": "Task name",
  "completed": true
}
```

**Error response:**
```json
{
  "error": "task not found"
}
```

## TUI Standards

### Bubble Tea Patterns
- Keep Model immutable, return new Model from Update
- Use commands for async operations
- Separate components into their own files
- Use Lip Gloss for consistent styling

### Key Bindings
- `j/k` or `↑/↓` for navigation
- `h/l` or `←/→` for pane switching
- `1-5` for view switching
- `a` for quick add, `c` for complete, `e` for edit
- `q` to quit, `?` for help

### Component Interface
```go
type Component interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Component, tea.Cmd)
    View() string
}
```

## Testing Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/bridge/...

# Run with verbose output
go test -v ./...
```

## Build Commands

```bash
# Build binary
go build -o lazyfocus ./cmd/lazyfocus

# Run directly
go run ./cmd/lazyfocus

# Install locally
go install ./cmd/lazyfocus
```

## Platform Constraints

- **macOS only** — Omni Automation requires macOS
- **OmniFocus must be running** — Scripts execute via the running application
- **Automation permissions required** — First run triggers system permission prompt
- **Some features require OmniFocus Pro** — Custom perspectives, review functionality

## Agent-Friendly Design Requirements

- All commands must support `--json` flag
- Errors in JSON mode must return `{"error": "message"}`
- Support `--quiet` flag for scripting (exit codes only)
- Task/Project IDs must be stable and usable in subsequent commands
- Provide clear, parseable output for AI agent consumption

## CI Debugging Commands

When CI fails, use these commands to investigate:

```bash
# List recent CI runs
gh run list --limit 5

# View specific run details (shows jobs and annotations)
gh run view <run-id>

# View failed job logs
gh run view <run-id> --log-failed

# Watch a running workflow
gh run watch <run-id>
```
