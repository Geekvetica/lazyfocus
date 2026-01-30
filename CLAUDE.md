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
│   ├── app/                       # Main TUI application
│   │   └── app.go                 # Root model, orchestration
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
│       ├── keys.go                # Keybinding definitions
│       ├── styles.go              # Lip Gloss styles
│       ├── messages.go            # Message types
│       ├── command/               # Vim-style command parsing
│       ├── filter/                # Search/filter state
│       ├── overlay/               # Overlay compositor
│       ├── components/            # Reusable UI components
│       │   ├── quickadd/          # Quick add task overlay
│       │   ├── taskdetail/        # Task detail view
│       │   ├── taskedit/          # Task editing overlay
│       │   ├── confirm/           # Confirmation modal
│       │   ├── searchinput/       # Search input
│       │   ├── commandinput/      # Command input
│       │   ├── tasklist/          # Task list display
│       │   ├── projectlist/       # Project list display
│       │   └── taglist/           # Tag list display
│       └── views/                 # Screen implementations
│           ├── inbox/             # Inbox view
│           ├── projects/          # Projects view
│           ├── tags/              # Tags view
│           ├── forecast/          # Forecast view
│           └── review/            # Review view
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

### Phase 1: Foundation & Bridge Layer ✅ COMPLETE
**Status:** Fully implemented and tested
- ✅ Go module setup with dependencies (Cobra, Bubble Tea, Lip Gloss)
- ✅ Omni Automation bridge (execute JS via osascript, parse JSON responses)
- ✅ Domain models (Task, Project, Tag with full field support)
- ✅ Verified OmniFocus communication and error handling

### Phase 2: CLI Commands (Read Operations) ✅ COMPLETE
**Status:** All read commands implemented with filtering
- ✅ Cobra CLI structure with root command
- ✅ `tasks` command with comprehensive filtering (inbox, all, project, tag, flagged, due, completed)
- ✅ `projects` command for listing all projects
- ✅ `tags` command for listing all tags
- ✅ `show` command for detailed task view
- ✅ Human and JSON output formatting
- ✅ `--json` flag support across all commands

### Phase 3: CLI Commands (Write Operations) ✅ COMPLETE
**Status:** Full CRUD operations with natural syntax parsing
- ✅ `add` command with natural syntax parsing (#tags, @projects, due:, defer:, !)
- ✅ `modify` command with granular field updates (add-tag, remove-tag, clear-due, etc.)
- ✅ `complete` command with multi-task support
- ✅ `delete` command with confirmation and force flag
- ✅ Natural date parsing (relative, next, in N units, ISO format)
- ⚠️ **Limitation:** Only one tag can be added during task creation via natural syntax due to OmniFocus API constraints. Use `modify --add-tag` for additional tags.

### Phase 4: TUI - Basic Structure ✅ COMPLETE
**Status:** Fully implemented with all core views and actions
- ✅ Bubble Tea application shell with model architecture
- ✅ Inbox view with task list rendering
- ✅ Quick Add overlay with natural syntax support
- ✅ Basic task navigation (j/k, Enter to view details)
- ✅ Task completion (c key)
- ✅ Task list component with formatting
- ✅ Help overlay (? key)
- ✅ Overlay compositor for proper layering

### Phase 5: TUI - Full Implementation ✅ COMPLETE
**Status:** Fully implemented with all views, actions, and advanced features
- ✅ All views: Projects (2), Tags (3), Forecast (4), Review (5)
- ✅ View switching via 1-5 keys
- ✅ Task detail view (Enter key) with full information display
- ✅ Task editing (e key) with tabbed form navigation
- ✅ Task deletion (d key) with confirmation modal
- ✅ Flag toggle (f key) for immediate flagging
- ✅ Search/filter functionality (/ key) with real-time filtering
- ✅ Vim-style command mode (: key) with command history and tab completion
- ✅ Projects view with drill-down navigation to project tasks
- ✅ Tags view with hierarchical display and drill-down
- ✅ Forecast view with tasks grouped by due date (Overdue, Today, Tomorrow, Week, Later)
- ✅ Review view for flagged tasks

### Phase 6: Polish & Distribution ⬚ NOT STARTED
**Status:** Planned for 1.0 release
- ⬚ Comprehensive error handling and edge cases
- ⬚ Configuration file support (~/.lazyfocus.yaml)
- ⬚ Shell completions (bash, zsh, fish)
- ⬚ Homebrew formula for easy installation
- ⬚ Release automation via GitHub Actions

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

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `task name is required` | Empty or whitespace-only task name | Provide a non-empty task name |
| `task not found` | Invalid or non-existent task ID | Verify task ID using `lazyfocus tasks` |
| `failed to resolve project: project not found` | Project name doesn't exist in OmniFocus | Check project name with `lazyfocus projects` |
| `invalid due date: unrecognized date format: xyz` | Date string not in supported format | Use relative (tomorrow), next (next monday), in (in 3 days), or ISO format |
| `no modifications specified` | `modify` command without any flags | Provide at least one modification flag |
| `confirmation required: use --force to delete` | `delete` command without `--force` | Add `--force` flag or use `--json` mode |
| `OmniFocus is not running` | OmniFocus application not launched | Launch OmniFocus before running commands |
| `automation permission denied` | Automation permission not granted | Allow Terminal/iTerm access in System Preferences > Security > Automation |
| `project requires single tag for natural syntax` | Multiple tags in natural syntax during `add` | Use only one `#tag` in natural syntax, or add via `modify --add-tag` |

**Exit codes:**
- `0` - Success
- `1` - General error (invalid arguments, missing flags)
- `2` - OmniFocus not running or permission denied
- `3` - Item not found (task, project, or tag)

In JSON mode, all errors return `{"error": "message"}` with appropriate exit codes.

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

### Known API Limitations

**Task Creation:**
- Only one tag can be assigned during task creation via the Omni Automation API
- Workaround: Create task with one tag, then use modify operations to add additional tags
- This affects both natural syntax parsing and flag-based task creation

**Task Modification:**
- Multiple tags can be added/removed in a single modify operation
- Date clearing requires explicit operations (cannot set dates to null in creation)

**Performance:**
- Each `osascript` call has ~100-200ms overhead
- Batch operations when possible
- Cache project and tag lists for validation

**Error Handling:**
- OmniFocus must be running for any operation
- Some operations fail silently if OmniFocus is in an inconsistent state
- Always validate IDs before operations when possible

## CLI Command Reference

### Read Commands

#### `tasks` - List tasks with filtering support

**Basic usage:**
```bash
lazyfocus tasks              # Show inbox (default)
lazyfocus tasks --all        # Show all tasks
lazyfocus tasks --json       # JSON output
```

**Filtering flags:**
- `--inbox` - Show inbox tasks only (default)
- `--all` - Show all incomplete tasks
- `--project <name>` - Filter by project name or ID
- `--tag <name>` - Filter by tag name
- `--flagged` - Show only flagged tasks
- `--due <date>` - Show tasks due on or before date
- `--completed` - Show completed tasks instead of incomplete

**Examples:**
```bash
lazyfocus tasks --project Work
lazyfocus tasks --tag urgent --flagged
lazyfocus tasks --due today
lazyfocus tasks --completed --project "Personal"
```

#### `projects` - List all projects

```bash
lazyfocus projects
lazyfocus projects --json
```

Lists all projects in OmniFocus with task counts.

#### `tags` - List all tags

```bash
lazyfocus tags
lazyfocus tags --json
```

Lists all tags in OmniFocus.

#### `show` - Show task details

```bash
lazyfocus show <task-id>
lazyfocus show abc123 --json
```

Displays detailed information about a specific task including name, project, tags, dates, notes, and completion status.

#### `perspective` - View custom perspectives

**Status:** Planned for future implementation (requires OmniFocus Pro)

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
- `-t, --tag <name>` - Tag name (only one tag supported during creation)
- `-d, --due <date>` - Due date
- `--defer <date>` - Defer date
- `-f, --flagged` - Mark as flagged
- `-n, --note <text>` - Task note

Command-line flags override natural syntax when both are present.

**Known Limitation:**
- Only **one tag** can be added during task creation due to OmniFocus Automation API constraints
- Workaround: Use `lazyfocus modify <task-id> --add-tag <name>` to add additional tags after creation
- Natural syntax with multiple `#tags` will only apply the first tag found

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
lazyfocus add "Buy groceries #errands"
lazyfocus add "Team sync #\"project alpha\""
```

⚠️ **Note:** Only the first tag in natural syntax will be applied due to OmniFocus API limitations. Use `modify --add-tag` for additional tags.

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

### Current Implementation Status

**✅ All Features Implemented:**

**Views:**
- Inbox view (key `1`) - Task list with completion status
- Projects view (key `2`) - Project list with drill-down to tasks
- Tags view (key `3`) - Hierarchical tag list with drill-down
- Forecast view (key `4`) - Tasks grouped by due date
- Review view (key `5`) - Flagged tasks for quick review

**Overlays:**
- Quick Add (`a`) - Natural syntax task creation
- Task Detail (`Enter`) - Full task information with actions
- Task Edit (`e`) - Tabbed form for modifying tasks
- Delete Confirmation (`d`) - Confirmation modal for destructive actions
- Search Input (`/`) - Real-time task filtering
- Command Input (`:`) - Vim-style command mode
- Help (`?`) - Keyboard shortcuts reference

**Task Actions:**
- Complete (`c`) - Mark task as complete
- Delete (`d`) - Delete with confirmation
- Edit (`e`) - Open edit overlay
- Flag (`f`) - Toggle flagged status

### Bubble Tea Patterns
- Keep Model immutable, return new Model from Update
- Use commands for async operations (task loading, task operations)
- Separate components into their own files
- Use Lip Gloss for consistent styling
- Handle errors gracefully with user-visible messages

### Key Bindings

**Navigation:**
- `j` or `↓` - Move down in list
- `k` or `↑` - Move up in list
- `Enter` - View task details / drill-down into project or tag
- `h` or `Esc` - Go back from drill-down view
- `1-5` - Switch between views (Inbox, Projects, Tags, Forecast, Review)

**Task Actions:**
- `a` - Open Quick Add overlay
- `c` - Complete selected task
- `d` - Delete selected task (with confirmation)
- `e` - Edit selected task
- `f` - Toggle flag on selected task

**Search & Commands:**
- `/` - Open search input (real-time filtering)
- `:` - Open command input (vim-style commands)

**General:**
- `?` - Toggle help overlay
- `q` or `Ctrl+C` - Quit application

### Vim-Style Commands

Available commands (all support aliases):
- `:quit` / `:q` / `:exit` - Quit application
- `:refresh` / `:w` / `:sync` - Refresh current view
- `:add` / `:a` `<task>` - Add new task
- `:complete` / `:done` / `:c` - Complete selected task
- `:delete` / `:del` / `:rm` - Delete selected task
- `:project` / `:p` `<name>` - Filter by project
- `:tag` / `:t` `<name>` - Filter by tag
- `:due` `<today|tomorrow|week|overdue>` - Filter by due date
- `:flagged` - Show only flagged tasks
- `:clear` / `:reset` - Clear all filters
- `:help` / `:?` - Show help

### Component Interface
```go
type Component interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Component, tea.Cmd)
    View() string
}
```

### TUI Architecture

The TUI follows a component-based architecture:

- **Main Model (`internal/app/app.go`)**: Root application state and orchestration
- **Views** (`internal/tui/views/`): Inbox, Projects, Tags, Forecast, Review
- **Components** (`internal/tui/components/`):
  - `quickadd` - Quick Add overlay with natural syntax
  - `taskdetail` - Task detail view overlay
  - `taskedit` - Task editing overlay with tabbed form
  - `confirm` - Reusable confirmation modal
  - `searchinput` - Search input with real-time filtering
  - `commandinput` - Vim-style command input
  - `tasklist` - Reusable task list display
  - `projectlist` - Project list display
  - `taglist` - Hierarchical tag list display
- **Filter State** (`internal/tui/filter/`): Search and filter state management
- **Command Parser** (`internal/tui/command/`): Vim-style command parsing
- **Message Passing**: Custom messages for async operations (TasksLoadedMsg, TaskCompletedMsg, etc.)
- **Overlay Compositor** (`internal/tui/overlay/`): Character-level overlay compositing

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

## Linting Commands

**CRITICAL: Always run lint checks before considering any implementation complete.**

### Install golangci-lint (one-time setup)

```bash
# Via Homebrew (recommended)
brew install golangci-lint

# Or via go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Run Lint Checks

```bash
# Run all linters (uses .golangci.yml config)
golangci-lint run

# Run with increased timeout for larger changes
golangci-lint run --timeout=3m

# Run on specific packages
golangci-lint run ./internal/app/...

# Show all issues including those from cache
golangci-lint run --new=false
```

### Auto-Fix Issues

```bash
# Auto-fix formatting issues
gofmt -w .

# Auto-fix imports (add missing, remove unused)
goimports -w .

# Some linters support --fix flag
golangci-lint run --fix
```

### Configuration Compatibility Warning

**CRITICAL: Do NOT modify `.golangci.yml` without verifying CI compatibility.**

The CI uses golangci-lint **v1.64.8** (pinned in `.github/workflows/ci.yml`). The configuration file uses v1 format, which is NOT compatible with golangci-lint v2.

**v2-only features that will BREAK CI:**
- `version: 2` at the root level
- `formatters:` section (v2 separates formatters from linters)
- `output.formats:` as an object (v1 requires array format)

If you have golangci-lint v2 installed locally, you may see config errors when running `golangci-lint config verify`. This is expected - the config is designed for the CI's v1.64.8 version.

### Common Lint Issues and Fixes

| Issue | Linter | Fix |
|-------|--------|-----|
| Cyclomatic complexity > 15 | gocyclo | Extract helper methods to reduce function complexity |
| Name stuttering (e.g., `pkg.PkgType`) | revive | Rename type to avoid package prefix repetition |
| Unused variable | unused | Remove or use the variable, use `_` for intentionally unused |
| Missing comment on exported type | revive | Add `// TypeName ...` comment above the type |
| Deprecated API usage | staticcheck | Replace with recommended alternative |
| Formatting issues | gofmt | Run `gofmt -w .` |

### Pre-Commit Checklist

Before considering any implementation complete, run:

```bash
# 1. Format code
gofmt -w .

# 2. Run linter
golangci-lint run --timeout=3m

# 3. Run tests
go test ./...

# 4. Build
go build ./...
```

All four commands must pass with no errors.

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

## Documentation References

LazyFocus maintains comprehensive documentation for different audiences:

- **[`docs/commands.md`](/Users/pwojciechowski/_dev/lazyfocus/docs/commands.md)** - Complete CLI command reference with examples and all flags
- **[`docs/json-schemas.md`](/Users/pwojciechowski/_dev/lazyfocus/docs/json-schemas.md)** - JSON schemas for AI agents and programmatic access
- **[`docs/troubleshooting.md`](/Users/pwojciechowski/_dev/lazyfocus/docs/troubleshooting.md)** - Common issues, error messages, and solutions
- **[`README.md`](/Users/pwojciechowski/_dev/lazyfocus/README.md)** - User-facing project overview and quick start guide

When working on the codebase:
- Update `docs/commands.md` when adding/modifying CLI commands
- Update `docs/json-schemas.md` when changing JSON output structure
- Update `docs/troubleshooting.md` when adding new error messages or handling edge cases
- Keep `CLAUDE.md` (this file) synchronized with implementation status

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
