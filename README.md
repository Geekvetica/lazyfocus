# LazyFocus

A CLI and TUI tool for interacting with OmniFocus on macOS via Omni Automation.

## Overview

LazyFocus (`lf`) provides quick terminal access to OmniFocus, serving two audiences:

- **Humans** — Fast terminal access to tasks with readable output
- **AI Agents** — Structured JSON interface for LLMs to query and manipulate tasks

## Requirements

- macOS (Omni Automation requires macOS)
- OmniFocus 3 or 4
- Go 1.21+ (for building from source)

## Installation

### From Source

```bash
git clone https://github.com/pwojciechowski/lazyfocus.git
cd lazyfocus
go build -o lazyfocus ./cmd/lazyfocus
```

### Add to PATH (optional)

```bash
# Move to a directory in your PATH
mv lazyfocus /usr/local/bin/

# Or create a symlink
ln -s $(pwd)/lazyfocus /usr/local/bin/lf
```

## Usage

### View Inbox Tasks

```bash
./lazyfocus
```

Output:
```
LazyFocus - Inbox Tasks
==================================================
Found 3 task(s):

☐ Buy groceries                    📅 Today
  Note: Remember milk
  Tags: errands

☐ Review PR #142                   🚩
  Tags: work, code-review

☑ Completed task
  Completed: Jan 27
```

### First Run

On first run, macOS will prompt for Automation permission. Grant access to allow LazyFocus to communicate with OmniFocus.

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

# Integration tests (requires OmniFocus running)
go test -tags=integration ./internal/bridge/...
```

### Project Structure

```
lazyfocus/
├── cmd/lazyfocus/main.go      # Application entry point
├── internal/
│   ├── bridge/                # OmniFocus communication layer
│   │   ├── executor.go        # osascript execution
│   │   ├── parser.go          # JSON response parsing
│   │   ├── scripts.go         # Embedded JS scripts
│   │   └── scripts/           # Omni Automation scripts
│   └── domain/                # Domain models
│       ├── task.go
│       ├── project.go
│       └── tag.go
└── scripts/                   # Reference JXA scripts
```

## Roadmap

- [x] Phase 1: Foundation & Bridge Layer
- [ ] Phase 2: CLI Commands (Read Operations)
- [ ] Phase 3: CLI Commands (Write Operations)
- [ ] Phase 4: TUI - Basic Structure
- [ ] Phase 5: TUI - Full Implementation
- [ ] Phase 6: Polish & Distribution

## License

MIT
