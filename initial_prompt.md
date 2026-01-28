# LazyFocus - Claude Code Project Prompt

## Project Overview

Build a CLI and TUI tool called `lazyfocus` (or `lf` for short) that interfaces with OmniFocus on macOS using Omni Automation (JavaScript for Automation). The tool serves two audiences:

1. **Humans** — quick terminal access to OmniFocus with readable output
2. **AI Agents** — structured JSON interface for LLMs to query and manipulate tasks

## Technical Decisions

- **Language:** Go
- **CLI Framework:** Cobra
- **TUI Framework:** Bubble Tea (with Bubbles, Lip Gloss)
- **OmniFocus Interface:** Omni Automation via `osascript -l JavaScript`
- **Output Formats:** Human-readable (default), JSON (`--json` flag)

## Project Structure

```
lazyfocus/
├── cmd/
│   └── lazyfocus/
│       └── main.go                 # Single entrypoint, TUI via `lf` or `lf tui`
├── internal/
│   ├── bridge/                     # Omni Automation execution layer
│   │   ├── executor.go             # osascript wrapper with error handling
│   │   ├── scripts.go              # Embedded JS scripts
│   │   └── parser.go               # JSON response parsing
│   ├── domain/                     # Shared domain models
│   │   ├── task.go
│   │   ├── project.go
│   │   ├── tag.go
│   │   └── perspective.go
│   ├── cli/                        # Cobra command implementations
│   │   ├── root.go
│   │   ├── tasks.go
│   │   ├── projects.go
│   │   ├── add.go
│   │   ├── complete.go
│   │   ├── modify.go
│   │   └── output.go               # Human vs JSON formatting
│   └── tui/                        # Bubble Tea TUI
│       ├── app.go                  # Main model, orchestration
│       ├── keys.go                 # Keybinding definitions
│       ├── styles.go               # Lip Gloss styles
│       ├── components/
│       │   ├── tasklist.go
│       │   ├── projecttree.go
│       │   ├── taglist.go
│       │   ├── taskdetail.go
│       │   ├── quickadd.go
│       │   └── help.go
│       └── views/
│           ├── inbox.go
│           ├── projects.go
│           ├── forecast.go
│           └── review.go
├── scripts/                        # Raw Omni Automation JS (for reference/testing)
│   ├── get_tasks.js
│   ├── get_projects.js
│   ├── add_task.js
│   └── complete_task.js
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Phase 1: Foundation & Bridge Layer

### Goals

- Set up Go module with dependencies
- Implement the Omni Automation bridge (execute JS, parse responses)
- Create basic domain models
- Verify OmniFocus communication works

### Tasks

1. **Initialize Go module**

   ```bash
   go mod init github.com/[user]/lazyfocus
   ```

2. **Add dependencies**

   - `github.com/spf13/cobra`
   - `github.com/charmbracelet/bubbletea`
   - `github.com/charmbracelet/bubbles`
   - `github.com/charmbracelet/lipgloss`

3. **Implement bridge/executor.go**

   - Function to execute Omni Automation JavaScript via osascript
   - Proper error handling (OmniFocus not running, script errors)
   - Timeout handling
   - Return structured results

4. **Create domain models**

   - `Task`: ID, Name, Note, Project, Tags, DueDate, DeferDate, Flagged, Completed, etc.
   - `Project`: ID, Name, Tasks, Status, Note, etc.
   - `Tag`: ID, Name, Parent, Children

5. **Write initial Omni Automation scripts**

   - `getTasks()` — fetch inbox tasks
   - `getProjects()` — fetch all projects
   - Script should return JSON that Go can parse

6. **Test the bridge**
   - Create a simple main.go that fetches and prints inbox tasks
   - Verify JSON parsing works correctly

### Omni Automation Notes

```javascript
// Example: Get inbox tasks as JSON
(() => {
  const app = Application("OmniFocus");
  const doc = app.defaultDocument;
  const inbox = doc.inboxTasks();

  const tasks = inbox.map((t) => ({
    id: t.id(),
    name: t.name(),
    note: t.note(),
    flagged: t.flagged(),
    dueDate: t.dueDate() ? t.dueDate().toISOString() : null,
    completed: t.completed(),
  }));

  return JSON.stringify(tasks);
})();
```

Execute with:

```bash
osascript -l JavaScript -e '<script>'
```

---

## Phase 2: CLI Commands (Read Operations)

### Goals

- Implement Cobra CLI structure
- Add read-only commands with human and JSON output
- Support filtering and querying

### Commands to Implement

```bash
# List tasks
lf tasks                      # Inbox tasks (default)
lf tasks --project "Work"     # Tasks in specific project
lf tasks --tag "waiting"      # Tasks with specific tag
lf tasks --flagged            # Flagged tasks only
lf tasks --due today          # Due today
lf tasks --due week           # Due this week
lf tasks --available          # Available tasks only
lf tasks --json               # JSON output for agents

# List projects
lf projects                   # All active projects
lf projects --folder "Work"   # Projects in folder
lf projects --stalled         # Stalled projects
lf projects --json

# List tags
lf tags                       # All tags
lf tags --unused              # Tags with no tasks
lf tags --json

# Show single item details
lf show <task-id>             # Full task details
lf show <project-id>          # Full project details

# Perspectives (if accessible via Omni Automation)
lf perspective "Forecast"     # Run built-in perspective
lf perspective "My Custom"    # Run custom perspective
```

### Output Formatting

Human output example:

```
INBOX (3 tasks)
───────────────────────────────────────
☐ Buy groceries                    📅 Today
  #errands
☐ Review PR #142                   🚩
  #work #code-review
☐ Call dentist
  #calls
```

JSON output example:

```json
{
  "tasks": [
    {
      "id": "abc123",
      "name": "Buy groceries",
      "project": null,
      "tags": ["errands"],
      "dueDate": "2025-01-28",
      "flagged": false
    }
  ],
  "count": 3
}
```

---

## Phase 3: CLI Commands (Write Operations)

### Goals

- Add task creation with natural syntax
- Implement task modification
- Add completion/deletion

### Commands to Implement

```bash
# Quick add (natural language-ish)
lf add "Buy milk"
lf add "Buy milk @errands"                    # With tag
lf add "Buy milk @errands #Shopping"          # With tag and project
lf add "Buy milk due:tomorrow"                # With due date
lf add "Buy milk defer:monday due:friday"     # With defer and due
lf add "Buy milk !"                           # Flagged
lf add "Buy milk //This is a note"            # With note

# Structured add
lf add --name "Buy milk" --project "Shopping" --tag "errands" --due "2025-02-01"

# Modify existing
lf modify <task-id> --name "New name"
lf modify <task-id> --due "tomorrow"
lf modify <task-id> --tag +urgent             # Add tag
lf modify <task-id> --tag -waiting            # Remove tag
lf modify <task-id> --project "Other Project"
lf modify <task-id> --flag                    # Set flagged
lf modify <task-id> --unflag                  # Remove flag

# Complete
lf complete <task-id>
lf complete <task-id> <task-id> ...           # Multiple

# Delete
lf delete <task-id>
lf delete <task-id> --force                   # Skip confirmation

# Defer/reschedule
lf defer <task-id> --to "next monday"
lf defer <task-id> --by "3 days"
```

### Natural Date Parsing

Support common patterns:

- `today`, `tomorrow`, `yesterday`
- `monday`, `next monday`, `last friday`
- `+3d`, `+1w`, `+2m` (relative)
- `2025-02-15` (ISO)
- `feb 15`, `15 feb` (natural)

Consider using a library like `github.com/tj/go-naturaldate` or similar.

---

## Phase 4: TUI - Basic Structure

### Goals

- Set up Bubble Tea application shell
- Implement basic navigation between views
- Create task list component

### Initial TUI Layout

```
┌─ LAZYFOCUS ───────────────────────────────────────────┐
│ [I]nbox  [P]rojects  [T]ags  [F]orecast  [R]eview     │
├───────────────────────┬───────────────────────────────┤
│ INBOX (12)            │ Buy groceries                 │
│                       │ ─────────────────────────────│
│ ● Buy groceries       │ Project: Shopping             │
│   Review PR #142      │ Tags: errands, home           │
│   Call dentist        │ Due: Today                    │
│   Plan vacation       │ Defer: -                      │
│   ...                 │ Flagged: No                   │
│                       │                               │
│                       │ Note:                         │
│                       │ Get milk, eggs, bread         │
│                       │                               │
├───────────────────────┴───────────────────────────────┤
│ a:add  e:edit  c:complete  d:defer  f:flag  ?:help    │
└───────────────────────────────────────────────────────┘
```

### Key Bindings

```
Navigation:
  j/k or ↑/↓     Navigate list
  h/l or ←/→     Switch panes
  1-5            Switch views (Inbox/Projects/Tags/Forecast/Review)
  Tab            Cycle focus
  Enter          Expand/drill down
  Esc            Back/close

Actions:
  a              Quick add task
  e              Edit selected task
  c              Complete selected task
  d              Defer task
  f              Toggle flag
  t              Add/remove tag
  p              Move to project
  /              Search/filter
  ?              Help

Global:
  q              Quit
  r              Refresh
  :              Command mode (vim-style)
```

### Components to Build

1. **tasklist.go** — Scrollable list with selection, supports filtering
2. **taskdetail.go** — Right pane showing full task info
3. **quickadd.go** — Text input overlay for adding tasks
4. **help.go** — Help overlay showing keybindings

---

## Phase 5: TUI - Full Implementation

### Goals

- Complete all views (Projects, Tags, Forecast, Review)
- Add search/filter functionality
- Implement all task actions within TUI
- Polish UI and error handling

### Views to Implement

1. **Inbox View** — Flat list of inbox tasks
2. **Projects View** — Tree view of folders → projects → tasks
3. **Tags View** — List of tags, selecting shows tagged tasks
4. **Forecast View** — Tasks grouped by date (today, tomorrow, this week, etc.)
5. **Review View** — Projects needing review (if OmniFocus Pro)

### Advanced Features

- **Search** — `/` to enter search mode, filter across all tasks
- **Bulk Actions** — Select multiple tasks (Space to toggle), then action applies to all
- **Vim-style Command Mode** — `:complete`, `:defer +3d`, `:move "Project Name"`
- **Custom Perspectives** — If Omni Automation supports it, allow switching perspectives

---

## Phase 6: Polish & Distribution

### Goals

- Error handling and edge cases
- Performance optimization (caching?)
- Documentation
- Release workflow

### Tasks

1. **Error Handling**

   - OmniFocus not installed
   - OmniFocus not running
   - Invalid task/project IDs
   - Permission issues (Automation permissions)

2. **Configuration File** (optional)

   - `~/.config/lazyfocus/config.yaml`
   - Default project for quick add
   - Custom keybindings
   - Theme/colors

3. **Shell Completions**

   - Bash/Zsh/Fish completions via Cobra

4. **Documentation**

   - README with installation, usage examples
   - GIF demos of TUI
   - Man page generation

5. **Release**
   - Makefile with build targets
   - Homebrew formula
   - GitHub releases with goreleaser

---

## Constraints & Considerations

### Omni Automation Limitations

- Only works when OmniFocus is running
- Some features may require OmniFocus Pro
- Custom perspectives may have limited API access
- Performance: each call spawns osascript process

### macOS Specific

- Requires macOS (Omni Automation is Apple-only)
- First run will trigger Automation permission prompt
- Consider adding check/instructions for granting permissions

### Agent-Friendly Design

- All commands should support `--json` for structured output
- Errors should also be JSON in JSON mode: `{"error": "message"}`
- Consider `--quiet` flag for scripting (exit codes only)
- IDs should be stable and usable in subsequent commands

---

## Getting Started Command

Start with Phase 1. Initialize the project and get a basic "fetch inbox tasks" working:

```bash
# Create project directory
mkdir lazyfocus && cd lazyfocus

# Initialize and set up basic structure
# Then verify we can talk to OmniFocus
```

Focus on getting the bridge layer solid before building CLI commands on top.

---

## Success Criteria

**Phase 1 Complete When:**

- `go run ./cmd/lazyfocus` successfully prints inbox tasks from OmniFocus
- JSON parsing works, domain models populated correctly

**Phase 2 Complete When:**

- All read commands work with both human and JSON output
- Filtering flags work correctly

**Phase 3 Complete When:**

- Can add, modify, complete, delete tasks from CLI
- Natural date parsing works

**Phase 4 Complete When:**

- TUI launches with task list visible
- Can navigate with keyboard
- Can complete a task from TUI

**Phase 5 Complete When:**

- All views implemented
- Search works
- All actions available in TUI

**Phase 6 Complete When:**

- `brew install lazyfocus` works
- README is comprehensive
- No crashes on edge cases

---

## Branding & Naming

- **Full name:** LazyFocus
- **Binary name:** `lazyfocus`
- **Short alias:** `lf` (create symlink or shell alias)
- **Tagline:** "A lazy way to manage OmniFocus from the terminal"
- **Style:** Follows the lazygit/lazydocker naming convention for TUI tools

### Trademark Notice (include in README)

```markdown
## Disclaimer

LazyFocus is not affiliated with, endorsed by, or sponsored by The Omni Group.
OmniFocus is a registered trademark of Omni Development, Inc.

This is an independent, open-source project that uses Omni Automation
(the official JavaScript automation API) to interface with OmniFocus.
```

### Shell Alias Setup (include in README)

```bash
# Add to ~/.zshrc or ~/.bashrc for short alias
alias lf="lazyfocus"
```
