# LazyFocus

A powerful CLI and TUI tool for interacting with OmniFocus on macOS. LazyFocus brings the full power of OmniFocus to your terminal, with both human-friendly interfaces and structured JSON output for AI agents.

## Overview

LazyFocus (`lazyfocus` or `lf`) provides seamless terminal access to OmniFocus via Omni Automation, serving two distinct audiences:

- **Humans** — Fast, intuitive CLI commands and an interactive TUI for managing tasks
- **AI Agents** — Structured JSON interface enabling LLMs to query and manipulate OmniFocus data

## Features

- **10 CLI Commands** — Complete task management from the command line
- **Natural Syntax** — Create tasks with intuitive notation: `#tags`, `@project`, `due:tomorrow`, `!`
- **Interactive TUI** — Full-screen terminal interface with Vim-style navigation (Inbox view currently available)
- **JSON Output** — Every command supports `--json` for programmatic access
- **Flexible Filtering** — Query tasks by project, tag, due date, flagged status, and more
- **Custom Perspectives** — Access OmniFocus Pro perspectives from the terminal

## Requirements

- **macOS** — Omni Automation requires macOS
- **OmniFocus 3 or 4** — Must be installed and running
- **Go 1.21+** — For building from source

## Installation

### From Source

```bash
git clone https://github.com/pwojciechowski/lazyfocus.git
cd lazyfocus
go build -o lazyfocus ./cmd/lazyfocus
```

### Install to PATH

```bash
# Option 1: Move binary to /usr/local/bin
mv lazyfocus /usr/local/bin/

# Option 2: Use go install (installs to $GOPATH/bin)
go install ./cmd/lazyfocus

# Option 3: Create symbolic link as 'lf'
ln -s $(pwd)/lazyfocus /usr/local/bin/lf
```

### First Run

On first run, macOS will prompt for Automation permission. Grant access to allow LazyFocus to communicate with OmniFocus.

## Quick Start

### Launch the TUI

```bash
lazyfocus        # or just: lf
```

The TUI provides an interactive inbox view with keyboard navigation. Press `a` to quick-add tasks, `j/k` to navigate, `q` to quit.

**See the [Commands Reference](docs/commands.md) for complete CLI documentation.**

### List Inbox Tasks

```bash
lazyfocus tasks --inbox
```

Output:
```
INBOX (3 tasks)
───────────────────────────────────────
☐ Buy groceries                    📅 Today
  #errands
☐ Review PR #142                   🚩
  #work #code-review
☐ Meeting prep
  @Big Project
```

### Create a Task with Natural Syntax

```bash
lazyfocus add "Buy milk #groceries due:tomorrow"
lazyfocus add "Review PR @Work due:friday !"
lazyfocus add "Meeting prep @\"Big Project\" defer:\"next monday\""
```

### Get JSON Output (for AI Agents)

```bash
lazyfocus tasks --inbox --json
```

Output:
```json
{
  "tasks": [
    {
      "id": "abc123",
      "name": "Buy groceries",
      "tags": ["errands"],
      "due": "2024-01-30T17:00:00Z",
      "flagged": false,
      "completed": false
    }
  ],
  "count": 1
}
```

## CLI Commands

Quick reference for all commands. See [docs/commands.md](docs/commands.md) for complete documentation.

### Read Operations

#### `tasks` - List and filter tasks

```bash
# View inbox
lazyfocus tasks --inbox

# View all incomplete tasks
lazyfocus tasks --all

# Filter by project
lazyfocus tasks --project Work

# Filter by tag
lazyfocus tasks --tag urgent

# Show only flagged tasks
lazyfocus tasks --flagged

# Show tasks due soon
lazyfocus tasks --due

# Include completed tasks
lazyfocus tasks --completed

# Combine filters
lazyfocus tasks --project Work --tag urgent --flagged
```

**Available flags:**
- `--inbox` - Show inbox tasks only
- `--all` - Show all incomplete tasks
- `--project <name>` - Filter by project name
- `--tag <name>` - Filter by tag name
- `--flagged` - Show only flagged tasks
- `--due` - Show tasks with due dates
- `--completed` - Include completed tasks

#### `projects` - List all projects

```bash
lazyfocus projects
lazyfocus projects --json
```

#### `tags` - List all tags

```bash
lazyfocus tags
lazyfocus tags --json
```

#### `show` - Display item details

```bash
lazyfocus show abc123
lazyfocus show abc123 --json
```

#### `perspective` - View custom perspective

```bash
lazyfocus perspective "Review"
lazyfocus perspective "Weekly Planning" --json
```

**Note:** Requires OmniFocus Pro

### Write Operations

#### `add` - Create new tasks

**Using Natural Syntax:**

LazyFocus supports intuitive notation for creating tasks quickly:

```bash
# Tags: #tagname or #"tag with spaces"
lazyfocus add "Buy groceries #errands #shopping"
lazyfocus add "Team sync #\"project alpha\""

# Projects: @projectname or @"project with spaces"
lazyfocus add "Review code @Work"
lazyfocus add "Planning meeting @\"Big Project\""

# Due dates: due:date or due:"date phrase"
lazyfocus add "Submit report due:friday"
lazyfocus add "Call client due:\"next monday\""

# Defer dates: defer:date or defer:"date phrase"
lazyfocus add "Review proposal defer:tomorrow"
lazyfocus add "Follow up defer:\"in 3 days\""

# Flagged: ! anywhere in input
lazyfocus add "Urgent task !"
lazyfocus add "! High priority item"

# Combine notation
lazyfocus add "Code review @Work #urgent due:tomorrow !"
```

**Supported date formats:**
- Relative: `today`, `tomorrow`, `yesterday`
- Next occurrence: `next monday`, `next week`
- In N units: `in 3 days`, `in 2 weeks`
- ISO format: `2024-01-15`
- Month/day: `Jan 15`, `January 15 2024`

All dates default to 5:00 PM local time unless specified.

**Using Flags:**

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

**Important Notes:**
- **Tag Limitation:** Due to OmniFocus automation API constraints, only the first tag specified will be applied during task creation. Use `modify --add-tag` to add additional tags afterward.
- Command-line flags override natural syntax when both are present.
- See [Commands Reference](docs/commands.md) for complete documentation and [Date Format Reference](docs/commands.md#date-format-reference) for all supported date formats.

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
- `-f, --force` - Skip confirmation prompt (required)

In JSON mode, confirmation is automatically skipped. Multiple task IDs supported.

#### `modify` - Update existing tasks

```bash
# Change task name
lazyfocus modify task123 --name "Updated name"

# Update dates and flagged status
lazyfocus modify task123 --due tomorrow --flagged true

# Manage tags
lazyfocus modify task123 --add-tag urgent --remove-tag low

# Clear dates
lazyfocus modify task123 --clear-due --clear-defer

# Move to project and update note
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

#### `version` - Show version information

```bash
lazyfocus version
```

### Global Flags

All commands support these global flags:

- `--json` - Output in JSON format (for AI agents)
- `--quiet` - Suppress output, use exit codes only
- `--timeout <duration>` - Set execution timeout (default: 30s)

## TUI (Terminal User Interface)

Launch the interactive TUI by running `lazyfocus` or `lf` without any subcommand.

### Current Features

- **Inbox View** - Browse all inbox tasks with keyboard navigation
- **Quick Add** - Press `a` to open the quick-add overlay
- **Task List Navigation** - Vim-style `j/k` or arrow keys

### Key Bindings

| Key | Action |
|-----|--------|
| `j` or `↓` | Move down |
| `k` or `↑` | Move up |
| `a` | Quick add task |
| `q` | Quit |
| `?` | Help (coming soon) |

### Coming Soon

- Projects, Tags, Forecast, and Review views
- Task completion (`c`) and editing (`e`) from TUI
- Search and filter functionality
- Vim-style command mode (`:`)
- View switching (`1-5` keys)

## For AI Agents

LazyFocus is designed with AI agents in mind. Every command supports the `--json` flag for structured output.

### JSON Output Format

**List commands** (tasks, projects, tags):
```json
{
  "tasks": [...],
  "count": 3
}
```

**Create operations** (add):
```json
{
  "id": "abc123",
  "name": "Buy milk",
  "tags": ["groceries"],
  "due": "2024-01-16T17:00:00Z",
  "flagged": false
}
```

**Modify operations** (modify):
```json
{
  "id": "abc123",
  "name": "Updated task name",
  "project": "Work",
  "tags": ["urgent"],
  "due": "2024-01-20T17:00:00Z"
}
```

**Complete/delete operations**:
```json
{
  "id": "abc123",
  "name": "Task name",
  "completed": true
}
```

**Error responses**:
```json
{
  "error": "task not found"
}
```

### Exit Codes

- `0` - Success
- `1` - General error
- `2` - OmniFocus not running or permission denied
- `3` - Task/project/tag not found

**See [JSON Schemas](docs/json-schemas.md) for detailed JSON response formats.**

### Scripting Example

```bash
#!/bin/bash

# Create a task and capture its ID
TASK_JSON=$(lazyfocus add "Follow up" --project Work --json)
TASK_ID=$(echo "$TASK_JSON" | jq -r '.id')

# Mark it as flagged
lazyfocus modify "$TASK_ID" --flagged true --quiet

# Complete it later
lazyfocus complete "$TASK_ID" --quiet
```

## Known Limitations

1. **Single tag limitation** - When creating tasks, only one tag can be added at a time via the OmniFocus Automation API. Use `modify --add-tag` to add additional tags after creation.

2. **OmniFocus Pro required** - Custom perspectives feature requires OmniFocus Pro subscription.

3. **macOS only** - Omni Automation is exclusive to macOS. LazyFocus will not work on Linux or Windows.

4. **OmniFocus must be running** - All commands require OmniFocus to be open and running. LazyFocus communicates via the running application.

5. **Automation permissions** - First run requires granting Automation permission in macOS System Settings.

**See [Troubleshooting Guide](docs/troubleshooting.md) for common issues and solutions.**

## Project Structure

```
lazyfocus/
├── cmd/lazyfocus/main.go          # Application entry point
├── internal/
│   ├── bridge/                    # OmniFocus communication layer
│   │   ├── executor.go            # osascript execution wrapper
│   │   ├── parser.go              # JSON response parsing
│   │   ├── scripts.go             # Embedded Omni Automation scripts
│   │   └── scripts/               # Raw .js script files
│   ├── domain/                    # Shared domain models
│   │   ├── task.go
│   │   ├── project.go
│   │   └── tag.go
│   ├── cli/                       # Cobra CLI command implementations
│   │   ├── root.go
│   │   ├── tasks.go
│   │   ├── projects.go
│   │   ├── add.go
│   │   ├── complete.go
│   │   ├── modify.go
│   │   ├── delete.go
│   │   └── output.go              # Human vs JSON formatting
│   └── tui/                       # Bubble Tea TUI
│       ├── app.go                 # Main application model
│       ├── keys.go                # Keybinding definitions
│       ├── styles.go              # Lip Gloss styling
│       ├── components/            # Reusable UI components
│       └── views/                 # Screen implementations
└── scripts/                       # Reference JXA scripts for testing
```

## Development

### Prerequisites

- Go 1.21+
- OmniFocus installed and running (for integration tests)

### Build

```bash
go build -o lazyfocus ./cmd/lazyfocus
```

### Run Tests

```bash
# Unit tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Integration tests (requires OmniFocus running)
go test -tags=integration ./internal/bridge/...
```

### Run Locally

```bash
go run ./cmd/lazyfocus tasks --inbox
```

## Roadmap

### Phase 1: Foundation & Bridge Layer ✅ COMPLETE
- [x] Go module setup with Cobra, Bubble Tea, Lip Gloss
- [x] Omni Automation bridge (execute JavaScript, parse JSON)
- [x] Domain models (Task, Project, Tag)
- [x] OmniFocus communication verified

### Phase 2: CLI Commands (Read Operations) ✅ COMPLETE
- [x] Cobra CLI structure
- [x] `tasks` - List/filter tasks
- [x] `projects` - List projects
- [x] `tags` - List tags
- [x] `show` - Show item details
- [x] `perspective` - View custom perspectives
- [x] Human and JSON output formatting

### Phase 3: CLI Commands (Write Operations) ✅ COMPLETE
- [x] `add` - Create tasks with natural syntax
- [x] `complete` - Mark tasks complete
- [x] `delete` - Delete tasks
- [x] `modify` - Update tasks
- [x] Natural date parsing
- [x] `version` - Show version

### Phase 4: TUI - Basic Structure 🚧 IN PROGRESS
- [x] Bubble Tea application shell
- [x] Inbox view
- [x] Quick Add overlay
- [x] Basic navigation (j/k, arrows)
- [ ] Task completion from TUI
- [ ] Task editing from TUI

### Phase 5: TUI - Full Implementation ⬚ NOT STARTED
- [ ] Projects view
- [ ] Tags view
- [ ] Forecast view
- [ ] Review view
- [ ] Search/filter functionality
- [ ] All task actions within TUI
- [ ] Vim-style command mode (`:`)
- [ ] View switching (1-5 keys)
- [ ] Help overlay (`?`)

### Phase 6: Polish & Distribution ⬚ NOT STARTED
- [ ] Comprehensive error handling
- [ ] Configuration file support (`~/.lazyfocusrc`)
- [ ] Shell completions (bash, zsh, fish)
- [ ] Homebrew formula
- [ ] GitHub releases with binaries
- [ ] Documentation website

## Contributing

Contributions are welcome! This project follows Test-Driven Development (TDD) principles and Kent Beck's methodologies.

### Development Workflow

1. Write a failing test
2. Implement minimum code to pass
3. Refactor with tests passing
4. Run all tests before submitting

### Code Standards

- Follow Go idioms and conventions
- Use table-driven tests
- Mock external dependencies (OmniFocus calls)
- Test both human and JSON output formats
- Keep functions focused and composable

## License

MIT

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- TUI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Integrates with [OmniFocus](https://www.omnigroup.com/omnifocus) via Omni Automation
