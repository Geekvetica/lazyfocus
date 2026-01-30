# LazyFocus CLI Command Reference

Complete reference for all LazyFocus CLI commands, flags, and syntax.

## Table of Contents

- [Global Flags](#global-flags)
- [Exit Codes](#exit-codes)
- [Read Commands](#read-commands)
  - [tasks](#tasks)
  - [projects](#projects)
  - [tags](#tags)
  - [show](#show)
  - [perspective](#perspective)
- [Write Commands](#write-commands)
  - [add](#add)
  - [complete](#complete)
  - [delete](#delete)
  - [modify](#modify)
- [Utility Commands](#utility-commands)
  - [version](#version)
- [Natural Syntax Reference](#natural-syntax-reference)
- [Date Format Reference](#date-format-reference)

## Global Flags

These flags are available for all commands:

| Flag | Description | Default |
|------|-------------|---------|
| `--json` | Output in JSON format (machine-readable) | `false` |
| `--quiet` | Suppress all output, use exit codes only | `false` |
| `--timeout <duration>` | Timeout for OmniFocus operations (e.g., "30s", "1m") | `30s` |

### Examples

```bash
# Get JSON output
lazyfocus tasks --json

# Quiet mode (only exit codes)
lazyfocus tasks --quiet

# Set custom timeout
lazyfocus tasks --timeout 60s
```

## Exit Codes

LazyFocus uses the following exit codes:

| Code | Meaning |
|------|---------|
| `0` | Successful execution |
| `1` | General error (invalid arguments, missing flags) |
| `2` | OmniFocus not running or permission denied |
| `3` | Requested item not found (task, project, or tag) |

## Read Commands

### tasks

List tasks from OmniFocus with various filtering options.

**Usage:**
```bash
lazyfocus tasks [flags]
```

**Description:**

By default, shows inbox tasks. Use flags to filter by project, tag, due date, or other criteria.

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--inbox` | boolean | Show inbox tasks only (default behavior) |
| `--all` | boolean | Show all tasks (across all projects and inbox) |
| `--project <id>` | string | Filter by project ID |
| `--tag <id>` | string | Filter by tag ID |
| `--flagged` | boolean | Show flagged tasks only |
| `--due <date>` | string | Show tasks due on/before date (supports 'today', 'tomorrow', or YYYY-MM-DD) |
| `--completed` | boolean | Include completed tasks in output |

**Examples:**

```bash
# Show inbox tasks (default)
lazyfocus tasks
lazyfocus tasks --inbox

# Show all tasks
lazyfocus tasks --all

# Show all tasks including completed
lazyfocus tasks --all --completed

# Show flagged tasks
lazyfocus tasks --flagged

# Show tasks by project
lazyfocus tasks --project abc123

# Show tasks by tag
lazyfocus tasks --tag def456

# Show tasks due today or earlier
lazyfocus tasks --due today

# Show tasks due tomorrow or earlier
lazyfocus tasks --due tomorrow

# Show tasks due by specific date
lazyfocus tasks --due 2024-12-31

# Combine filters (tasks in project, due today, JSON output)
lazyfocus tasks --project abc123 --due today --json
```

**Human Output:**
```
INBOX (3 tasks)
───────────────────────────────────────
☐ Buy groceries                    📅 Today
  #errands
☐ Review PR #142                   🚩
  #work #code-review
☐ Call dentist
  @Personal
```

**JSON Output:**
```json
{
  "tasks": [
    {
      "id": "abc123",
      "name": "Buy groceries",
      "completed": false,
      "flagged": false,
      "dueDate": "2024-01-15T17:00:00Z",
      "tags": ["errands"],
      "project": null
    }
  ],
  "count": 3
}
```

---

### projects

List projects from OmniFocus with filtering options.

**Usage:**
```bash
lazyfocus projects [flags]
```

**Description:**

By default, shows active projects. Use flags to filter by status or include nested tasks.

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--status <status>` | string | Filter by status (active, on-hold, completed, dropped, all) | `active` |
| `--with-tasks` | boolean | Include nested tasks in output | `false` |

**Examples:**

```bash
# Show active projects (default)
lazyfocus projects
lazyfocus projects --status active

# Show all projects regardless of status
lazyfocus projects --status all

# Show on-hold projects
lazyfocus projects --status on-hold

# Show completed projects
lazyfocus projects --status completed

# Show dropped projects
lazyfocus projects --status dropped

# Show projects with their tasks
lazyfocus projects --with-tasks

# JSON output
lazyfocus projects --json
```

**Human Output:**
```
PROJECTS (5 active)
───────────────────────────────────────
📁 Work
   Due: Jan 30, 2024

📁 Personal

📁 Learning
   Defer: Feb 1, 2024
```

**JSON Output:**
```json
{
  "projects": [
    {
      "id": "proj123",
      "name": "Work",
      "status": "active",
      "dueDate": "2024-01-30T17:00:00Z",
      "deferDate": null,
      "note": "",
      "tasks": []
    }
  ],
  "count": 5
}
```

---

### tags

List tags from OmniFocus with optional hierarchy and task counts.

**Usage:**
```bash
lazyfocus tags [flags]
```

**Description:**

By default, shows tags with hierarchy. Use flags to show flat list or include task counts.

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--flat` | boolean | Show tags in flat list (no hierarchy) | `false` |
| `--with-counts` | boolean | Show task count per tag | `false` |

**Examples:**

```bash
# Show tags with hierarchy (default)
lazyfocus tags

# Show tags in flat list
lazyfocus tags --flat

# Show tags with task counts
lazyfocus tags --with-counts

# Combine flags
lazyfocus tags --flat --with-counts

# JSON output
lazyfocus tags --json
```

**Human Output (hierarchical):**
```
TAGS
───────────────────────────────────────
🏷  Work
  🏷  urgent
  🏷  meetings
🏷  Personal
  🏷  errands
  🏷  health
```

**Human Output (flat with counts):**
```
TAGS
───────────────────────────────────────
🏷  Work (12)
🏷  urgent (3)
🏷  meetings (5)
🏷  Personal (8)
🏷  errands (4)
```

**JSON Output:**
```json
{
  "tags": [
    {
      "id": "tag123",
      "name": "Work",
      "parentID": null,
      "children": ["tag456", "tag789"]
    },
    {
      "id": "tag456",
      "name": "urgent",
      "parentID": "tag123",
      "children": []
    }
  ],
  "count": 5
}
```

---

### show

Show detailed information for a specific item (task, project, or tag) by its ID.

**Usage:**
```bash
lazyfocus show <id> [flags]
```

**Description:**

The command will attempt to auto-detect the type of item unless you specify the type explicitly with the `--type` flag.

**Arguments:**

| Argument | Required | Description |
|----------|----------|-------------|
| `<id>` | Yes | The ID of the item to show |

**Flags:**

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--type <type>` | string | Item type: task, project, or tag | auto-detect |

**Examples:**

```bash
# Auto-detect type
lazyfocus show abc123

# Explicitly specify type
lazyfocus show abc123 --type task
lazyfocus show proj456 --type project
lazyfocus show tag789 --type tag

# JSON output
lazyfocus show abc123 --json
```

**Human Output (Task):**
```
TASK: Buy groceries
───────────────────────────────────────
ID:       abc123
Status:   Active
Flagged:  Yes
Due:      Today, Jan 15 at 5:00 PM
Defer:    -
Project:  Personal
Tags:     #errands, #shopping

Notes:
Need milk, eggs, and bread
```

**Human Output (Project):**
```
PROJECT: Work
───────────────────────────────────────
ID:       proj123
Status:   Active
Due:      Jan 30, 2024 at 5:00 PM
Tasks:    12 remaining

Notes:
Q1 work items
```

**JSON Output (Task):**
```json
{
  "id": "abc123",
  "name": "Buy groceries",
  "completed": false,
  "flagged": true,
  "dueDate": "2024-01-15T17:00:00Z",
  "deferDate": null,
  "note": "Need milk, eggs, and bread",
  "tags": ["errands", "shopping"],
  "project": "Personal"
}
```

**Error Cases:**

If the item is not found, returns exit code 3:

```bash
lazyfocus show nonexistent123
# Error: item not found: nonexistent123
```

---

### perspective

Show tasks from a named OmniFocus perspective.

**Usage:**
```bash
lazyfocus perspective <name>
```

**Description:**

View tasks from a custom OmniFocus perspective. Note that custom perspectives require OmniFocus Pro.

**Arguments:**

| Argument | Required | Description |
|----------|----------|-------------|
| `<name>` | Yes | The name of the perspective |

**Examples:**

```bash
# Show tasks from "Today" perspective
lazyfocus perspective Today

# Show tasks from custom perspective
lazyfocus perspective "Work Focus"

# JSON output
lazyfocus perspective Today --json
```

**Human Output:**
```
TODAY PERSPECTIVE (8 tasks)
───────────────────────────────────────
☐ Review PR #142                   🚩
  @Work #urgent
☐ Buy groceries
  @Personal #errands
☐ Call dentist                     📅 Today
  @Personal
```

**JSON Output:**
```json
{
  "tasks": [
    {
      "id": "abc123",
      "name": "Review PR #142",
      "flagged": true,
      "project": "Work",
      "tags": ["urgent"]
    }
  ],
  "count": 8
}
```

**Notes:**

- Perspective name must match exactly (case-sensitive)
- Built-in perspectives: Inbox, Projects, Tags, Forecast, Flagged, Review, Nearby, Search
- Custom perspectives require OmniFocus Pro subscription

---

## Write Commands

### add

Create a new task in OmniFocus with natural syntax or flags.

**Usage:**
```bash
lazyfocus add <task description> [flags]
```

**Description:**

Create tasks using either natural syntax embedded in the description or explicit command-line flags. Flags override natural syntax when both are present.

**Arguments:**

| Argument | Required | Description |
|----------|----------|-------------|
| `<task description>` | Yes | Task name with optional natural syntax markers |

**Flags:**

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--project <name>` | `-p` | string | Project name or ID |
| `--tag <name>` | `-t` | string | Tags (repeatable flag) |
| `--due <date>` | `-d` | string | Due date (see [Date Formats](#date-format-reference)) |
| `--defer <date>` | | string | Defer date (see [Date Formats](#date-format-reference)) |
| `--flagged` | `-f` | boolean | Mark as flagged |
| `--note <text>` | `-n` | string | Task note |

**Natural Syntax in Description:**

| Syntax | Description | Example |
|--------|-------------|---------|
| `#tagname` | Add tag | `#errands` |
| `#"tag with spaces"` | Add tag with spaces | `#"project alpha"` |
| `@projectname` | Set project | `@Work` |
| `@"project name"` | Project with spaces | `@"Big Project"` |
| `due:date` | Set due date | `due:tomorrow` |
| `due:"date phrase"` | Due with spaces | `due:"next monday"` |
| `defer:date` | Set defer date | `defer:friday` |
| `defer:"date phrase"` | Defer with spaces | `defer:"in 3 days"` |
| `!` | Mark as flagged | `!` (anywhere in text) |

**Examples:**

```bash
# Simple task
lazyfocus add "Buy milk"

# Natural syntax: tag
lazyfocus add "Buy milk #groceries"

# Natural syntax: multiple tags
lazyfocus add "Team meeting #work #meetings"

# Natural syntax: tag with spaces
lazyfocus add "Client review #\"project alpha\""

# Natural syntax: project
lazyfocus add "Review PR @Work"

# Natural syntax: project with spaces
lazyfocus add "Planning session @\"Big Project\""

# Natural syntax: due date
lazyfocus add "Submit report due:friday"
lazyfocus add "Call client due:tomorrow"
lazyfocus add "Annual review due:\"next monday\""

# Natural syntax: defer date
lazyfocus add "Follow up defer:\"in 3 days\""

# Natural syntax: flagged
lazyfocus add "Urgent task !"
lazyfocus add "! High priority"

# Natural syntax: combined
lazyfocus add "Review PR @Work #urgent due:tomorrow !"

# Using flags instead
lazyfocus add "Buy milk" --project Personal --tag groceries --due tomorrow

# Using flags (short form)
lazyfocus add "Review PR" -p Work -t urgent -t code-review -d friday -f

# Flags override natural syntax
lazyfocus add "Buy milk #errands" --tag groceries
# Result: tag is "groceries" (flag overrides #errands)

# With note
lazyfocus add "Meeting prep" --project Work --note "Prepare slides and agenda"

# Multiple tags with flag
lazyfocus add "Code review" --tag urgent --tag code-review --tag backend

# JSON output
lazyfocus add "Buy milk" --json
```

**Human Output:**
```
Created task: Buy groceries
ID: abc123
Tags: #errands
Due: Today, Jan 15 at 5:00 PM
```

**JSON Output:**
```json
{
  "id": "abc123",
  "name": "Buy milk",
  "completed": false,
  "flagged": false,
  "dueDate": "2024-01-16T17:00:00Z",
  "deferDate": null,
  "note": "",
  "tags": ["groceries"],
  "project": null
}
```

**Error Cases:**

```bash
# Empty task name
lazyfocus add ""
# Error: task name is required

# Invalid date format
lazyfocus add "Task" --due xyz
# Error: invalid due date: unrecognized date format: xyz

# Project not found
lazyfocus add "Task" --project NonExistent
# Error: failed to resolve project: project not found
```

**Important Notes:**

- **Tag Limitation:** Due to OmniFocus automation API constraints, only the first tag specified will be applied during task creation. Use `modify --add-tag` to add additional tags afterward. See [Notes and Limitations](#notes-and-limitations) section for details.
- Command-line flags always take precedence over natural syntax
- All dates without explicit times default to 5:00 PM local time

---

### complete

Mark one or more tasks as complete in OmniFocus.

**Usage:**
```bash
lazyfocus complete <task-id> [task-id...] [flags]
```

**Description:**

Mark tasks as complete. Accepts multiple task IDs. The command will attempt to complete all specified tasks, continuing even if some fail.

**Arguments:**

| Argument | Required | Description |
|----------|----------|-------------|
| `<task-id>` | Yes | One or more task IDs to complete |

**Examples:**

```bash
# Complete single task
lazyfocus complete abc123

# Complete multiple tasks
lazyfocus complete abc123 def456 ghi789

# JSON output
lazyfocus complete abc123 --json
```

**Human Output (single task):**
```
Completed: Buy groceries
ID: abc123
```

**Human Output (multiple tasks):**
```
Completed: Buy groceries
ID: abc123

Completed: Call dentist
ID: def456

Completed: Review PR
ID: ghi789
```

**JSON Output:**
```json
{
  "id": "abc123",
  "name": "Buy groceries",
  "completed": true
}
```

**Error Handling:**

If some tasks fail to complete, the command continues processing remaining tasks and shows errors for failed ones:

```bash
lazyfocus complete abc123 invalid456 def789

# Output:
# Completed: Buy groceries (abc123)
# Error: failed to complete invalid456: task not found
# Completed: Call dentist (def789)
```

---

### delete

Delete one or more tasks from OmniFocus.

**Usage:**
```bash
lazyfocus delete <task-id> [task-id...] [flags]
```

**Description:**

Delete tasks permanently. By default, requires confirmation. Use `--force` to skip confirmation. In JSON mode, confirmation is automatically skipped.

**Arguments:**

| Argument | Required | Description |
|----------|----------|-------------|
| `<task-id>` | Yes | One or more task IDs to delete |

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-f` | Skip confirmation prompt |

**Examples:**

```bash
# Delete single task (requires --force)
lazyfocus delete abc123 --force

# Delete multiple tasks
lazyfocus delete abc123 def456 ghi789 --force

# JSON mode (auto-skips confirmation)
lazyfocus delete abc123 --json
```

**Human Output:**
```
Deleted: Buy groceries
ID: abc123
```

**JSON Output:**
```json
{
  "id": "abc123",
  "name": "Buy groceries",
  "deleted": true
}
```

**Error Cases:**

```bash
# Missing --force flag
lazyfocus delete abc123
# Error: confirmation required: use --force to delete without confirmation

# Task not found
lazyfocus delete invalid123 --force
# Error: failed to delete invalid123: task not found
```

**Error Handling:**

Similar to `complete`, continues processing all tasks even if some fail:

```bash
lazyfocus delete abc123 invalid456 def789 --force

# Output:
# Deleted: Buy groceries (abc123)
# Error: failed to delete invalid456: task not found
# Deleted: Call dentist (def789)
```

---

### modify

Modify an existing task in OmniFocus.

**Usage:**
```bash
lazyfocus modify <task-id> [flags]
```

**Description:**

Modify properties of an existing task. Requires exactly one task ID and at least one modification flag.

**Arguments:**

| Argument | Required | Description |
|----------|----------|-------------|
| `<task-id>` | Yes | The ID of the task to modify |

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--name <text>` | string | New task name |
| `--note <text>` | string | New note (replaces existing note) |
| `--project <name>` | string | Move to project (name or ID) |
| `--add-tag <name>` | string | Add tag (repeatable) |
| `--remove-tag <name>` | string | Remove tag (repeatable) |
| `--due <date>` | string | Set due date (see [Date Formats](#date-format-reference)) |
| `--defer <date>` | string | Set defer date (see [Date Formats](#date-format-reference)) |
| `--flagged <bool>` | string | Set flagged status (true/false) |
| `--clear-due` | boolean | Clear due date |
| `--clear-defer` | boolean | Clear defer date |

**Examples:**

```bash
# Change task name
lazyfocus modify abc123 --name "New task name"

# Update note
lazyfocus modify abc123 --note "Updated task details"

# Move to different project
lazyfocus modify abc123 --project Work
lazyfocus modify abc123 --project "Big Project"

# Add tags
lazyfocus modify abc123 --add-tag urgent
lazyfocus modify abc123 --add-tag urgent --add-tag high-priority

# Remove tags
lazyfocus modify abc123 --remove-tag low-priority

# Set due date
lazyfocus modify abc123 --due tomorrow
lazyfocus modify abc123 --due "next monday"
lazyfocus modify abc123 --due 2024-12-31

# Set defer date
lazyfocus modify abc123 --defer friday
lazyfocus modify abc123 --defer "in 3 days"

# Set flagged status
lazyfocus modify abc123 --flagged true
lazyfocus modify abc123 --flagged false

# Clear dates
lazyfocus modify abc123 --clear-due
lazyfocus modify abc123 --clear-defer

# Multiple modifications at once
lazyfocus modify abc123 --name "Updated name" --due tomorrow --flagged true

# Combine adding and removing tags
lazyfocus modify abc123 --add-tag urgent --remove-tag low-priority

# JSON output
lazyfocus modify abc123 --name "New name" --json
```

**Human Output:**
```
Modified task: Buy groceries
ID: abc123
Due: Tomorrow, Jan 16 at 5:00 PM
Flagged: Yes
```

**JSON Output:**
```json
{
  "id": "abc123",
  "name": "Buy groceries",
  "completed": false,
  "flagged": true,
  "dueDate": "2024-01-16T17:00:00Z",
  "deferDate": null,
  "note": "",
  "tags": ["errands", "urgent"],
  "project": "Personal"
}
```

**Error Cases:**

```bash
# No modifications specified
lazyfocus modify abc123
# Error: no modifications specified

# Invalid task ID
lazyfocus modify invalid123 --name "New name"
# Error: failed to modify task: task not found

# Invalid date
lazyfocus modify abc123 --due xyz
# Error: invalid due date: unrecognized date format: xyz

# Invalid flagged value
lazyfocus modify abc123 --flagged maybe
# Error: invalid flagged value (use true/false)

# Project not found
lazyfocus modify abc123 --project NonExistent
# Error: failed to resolve project: project not found
```

**Important Notes:**

- **Tag Limitation:** Only one tag can be added per modify operation due to OmniFocus automation API constraints. Run multiple `modify` commands to add multiple tags. See [Notes and Limitations](#notes-and-limitations) for details.
- At least one modification flag is required
- All dates without explicit times default to 5:00 PM local time

---

## Utility Commands

### version

Print version information for LazyFocus.

**Usage:**
```bash
lazyfocus version
```

**Description:**

Display the current version of LazyFocus along with build information.

**Examples:**

```bash
lazyfocus version
```

**Output:**
```
lazyfocus version 0.1.0
Build date: 2024-01-15
Git commit: abc1234
```

**Notes:**

- Version information is set at build time using `-ldflags`
- Does not require OmniFocus to be running
- Does not support `--json` flag

---

## Natural Syntax Reference

The `add` command supports natural language syntax embedded directly in the task description.

### Tags

**Syntax:** `#tagname` or `#"tag with spaces"`

```bash
# Single tag
lazyfocus add "Buy milk #groceries"

# Multiple tags
lazyfocus add "Team meeting #work #meetings #planning"

# Tag with spaces (use quotes)
lazyfocus add "Sprint review #\"project alpha\""
```

### Projects

**Syntax:** `@projectname` or `@"project with spaces"`

```bash
# Simple project name
lazyfocus add "Code review @Work"

# Project with spaces (use quotes)
lazyfocus add "Planning meeting @\"Big Project\""
```

### Due Dates

**Syntax:** `due:date` or `due:"date with spaces"`

```bash
# Simple date
lazyfocus add "Submit report due:friday"
lazyfocus add "Pay bills due:tomorrow"

# Date phrase with spaces (use quotes)
lazyfocus add "Annual review due:\"next monday\""
lazyfocus add "Follow up due:\"in 2 weeks\""
```

See [Date Format Reference](#date-format-reference) for all supported date formats.

### Defer Dates

**Syntax:** `defer:date` or `defer:"date with spaces"`

```bash
# Simple defer date
lazyfocus add "Review proposal defer:tomorrow"

# Defer with phrase (use quotes)
lazyfocus add "Follow up defer:\"in 3 days\""
```

### Flagged Status

**Syntax:** `!` (can appear anywhere in the description)

```bash
# Flag at end
lazyfocus add "Urgent task !"

# Flag at beginning
lazyfocus add "! High priority item"

# Flag in middle
lazyfocus add "Review ! this ASAP"
```

### Combining Syntax

You can combine multiple natural syntax elements:

```bash
# Tags + project + due + flagged
lazyfocus add "Review PR #urgent #code-review @Work due:tomorrow !"

# Project + defer + tags
lazyfocus add "Follow up @Sales #client #follow-up defer:\"next week\""

# Everything
lazyfocus add "Sprint planning @\"Product Team\" #planning #meetings due:friday defer:thursday !"
```

### Precedence Rules

Command-line flags always override natural syntax:

```bash
# Natural syntax says #errands, but flag overrides to #groceries
lazyfocus add "Buy milk #errands" --tag groceries
# Result: Task has tag "groceries"

# Natural syntax says @Personal, but flag overrides to Work
lazyfocus add "Task @Personal" --project Work
# Result: Task is in project "Work"
```

---

## Date Format Reference

LazyFocus supports flexible date input across all date-related flags and natural syntax.

### Relative Dates

Simple relative date expressions:

| Format | Example | Description |
|--------|---------|-------------|
| `today` | `--due today` | Today at 5:00 PM |
| `tomorrow` | `--due tomorrow` | Tomorrow at 5:00 PM |
| `yesterday` | `--defer yesterday` | Yesterday at 5:00 PM |

```bash
lazyfocus add "Task" --due today
lazyfocus add "Task" --due tomorrow
lazyfocus add "Task" due:yesterday
```

### Next Weekday

Get the next occurrence of a specific weekday:

| Format | Example | Description |
|--------|---------|-------------|
| `next monday` | `--due "next monday"` | Next Monday at 5:00 PM |
| `next tuesday` | `--due "next tuesday"` | Next Tuesday at 5:00 PM |
| `next wednesday` | `--due "next wednesday"` | Next Wednesday at 5:00 PM |
| `next thursday` | `--due "next thursday"` | Next Thursday at 5:00 PM |
| `next friday` | `--due "next friday"` | Next Friday at 5:00 PM |
| `next saturday` | `--due "next saturday"` | Next Saturday at 5:00 PM |
| `next sunday` | `--due "next sunday"` | Next Sunday at 5:00 PM |

```bash
lazyfocus add "Task" --due "next monday"
lazyfocus add "Task" due:"next friday"
```

**Note:** If today is Monday and you say "next monday", it will be next week's Monday (7 days from now), not today.

### Relative Time Spans

Add a specific number of days or weeks from today:

| Format | Example | Description |
|--------|---------|-------------|
| `in N day(s)` | `in 3 days` | Three days from now |
| `in N week(s)` | `in 2 weeks` | Two weeks from now |
| `next week` | `next week` | Seven days from now |

```bash
lazyfocus add "Task" --due "in 3 days"
lazyfocus add "Task" --due "in 1 day"
lazyfocus add "Task" --due "in 2 weeks"
lazyfocus add "Task" --due "next week"

# Natural syntax
lazyfocus add "Task due:\"in 5 days\""
lazyfocus add "Task defer:\"in 1 week\""
```

### ISO Date Format

Standard ISO 8601 date format:

| Format | Example | Description |
|--------|---------|-------------|
| `YYYY-MM-DD` | `2024-12-31` | Specific date at 5:00 PM |

```bash
lazyfocus add "Task" --due 2024-12-31
lazyfocus add "Task" --due 2024-01-15
lazyfocus add "Task" due:2024-06-30
```

### Month and Day

Natural month/day format:

| Format | Example | Description |
|--------|---------|-------------|
| `Mon DD` | `Jan 15` | January 15 of current year |
| `Month DD` | `January 15` | January 15 of current year |
| `Mon DD YYYY` | `Jan 15 2024` | January 15, 2024 |
| `Month DD YYYY` | `January 15 2024` | January 15, 2024 |

**Supported month abbreviations:**
- Jan, Feb, Mar, Apr, May, Jun, Jul, Aug, Sep/Sept, Oct, Nov, Dec

**Supported full month names:**
- January, February, March, April, May, June, July, August, September, October, November, December

```bash
# Current year assumed
lazyfocus add "Task" --due "Jan 15"
lazyfocus add "Task" --due "January 15"

# Explicit year
lazyfocus add "Task" --due "Jan 15 2024"
lazyfocus add "Task" --due "January 15 2024"

# Natural syntax
lazyfocus add "Task due:\"Mar 20\""
lazyfocus add "Task defer:\"December 31 2024\""
```

### Default Time

**Important:** All dates without explicit times default to **5:00 PM (17:00)** local time.

```bash
# All of these result in 5:00 PM local time:
lazyfocus add "Task" --due today           # Today at 5:00 PM
lazyfocus add "Task" --due tomorrow        # Tomorrow at 5:00 PM
lazyfocus add "Task" --due "next friday"   # Next Friday at 5:00 PM
lazyfocus add "Task" --due 2024-12-31      # Dec 31, 2024 at 5:00 PM
lazyfocus add "Task" --due "Jan 15"        # Jan 15 at 5:00 PM
```

### Case Insensitivity

Date parsing is case-insensitive:

```bash
# All valid (case doesn't matter)
lazyfocus add "Task" --due TODAY
lazyfocus add "Task" --due Tomorrow
lazyfocus add "Task" --due "NEXT MONDAY"
lazyfocus add "Task" --due "JAN 15"
lazyfocus add "Task" --due "january 15 2024"
```

### Examples by Use Case

**This week:**
```bash
lazyfocus add "Task" --due today
lazyfocus add "Task" --due tomorrow
lazyfocus add "Task" --due "next friday"
```

**Next week:**
```bash
lazyfocus add "Task" --due "next week"
lazyfocus add "Task" --due "in 7 days"
lazyfocus add "Task" --due "next monday"
```

**Specific future dates:**
```bash
lazyfocus add "Task" --due "in 2 weeks"
lazyfocus add "Task" --due "Jan 30"
lazyfocus add "Task" --due 2024-12-31
```

**Combining with natural syntax:**
```bash
lazyfocus add "Review code due:tomorrow"
lazyfocus add "Client meeting due:\"next friday\""
lazyfocus add "Annual review due:\"Jan 30 2024\""
lazyfocus add "Follow up defer:\"in 3 days\""
```

### Error Handling

Unrecognized date formats will return an error:

```bash
lazyfocus add "Task" --due xyz
# Error: invalid due date: unrecognized date format: xyz

lazyfocus add "Task" --due "last monday"
# Error: invalid due date: unrecognized date format: last monday

lazyfocus add "Task" --due "15 Jan"  # Wrong order
# Error: invalid due date: unrecognized date format: 15 jan
```

---

## Notes and Limitations

### Tag Limitations

LazyFocus uses OmniFocus's JavaScript for Automation (JXA) API, which has limitations regarding tag handling:

**During Task Creation (`add` command):**
- Only the first tag specified (via `--tag` flag or `#tag` natural syntax) will be applied
- If multiple tags are provided, only the first will be set as the primary tag
- This is a limitation of the OmniFocus automation API, not LazyFocus

**During Task Modification (`modify` command):**
- The `--add-tag` flag can only add one tag per operation
- Multiple `--add-tag` flags can be provided, but only the first will be applied
- Run multiple `modify` commands to add multiple tags

**Workaround for Multiple Tags:**
```bash
# Create task with first tag
TASK_ID=$(lazyfocus add "Buy groceries #shopping" --json | jq -r '.task.id')

# Add additional tags with separate modify commands
lazyfocus modify "$TASK_ID" --add-tag weekly
lazyfocus modify "$TASK_ID" --add-tag errands
```

See [Troubleshooting: Tag Limitations](./troubleshooting.md#tag-limitations) for more details.

### OmniFocus Pro Requirements

**Custom Perspectives:** Viewing custom perspectives requires OmniFocus Pro subscription.

### Platform Requirements

- LazyFocus only works on macOS (OmniFocus requirement)
- OmniFocus must be running for CLI commands to work
- First run will trigger macOS automation permission prompt
- Some features require OmniFocus Pro (custom perspectives, review functionality)

### Performance

- Default timeout is 30 seconds for OmniFocus operations
- Use `--timeout` flag to adjust for larger databases or slower systems
- JSON output is generally faster for scripting than human-readable output

### Output Modes

**Human-readable output:**
- Uses icons, colors, and formatting for terminal display
- Best for interactive terminal use
- Output format may change between versions

**JSON output (`--json`):**
- Machine-readable structured data
- Stable format for scripting and integration
- Suitable for AI agents and automation
- Always produces valid JSON (even for errors)

**Quiet mode (`--quiet`):**
- Suppresses all output
- Only exit codes indicate success/failure
- Useful for scripts that don't need output
- Can be combined with any command

---

## See Also

- [LazyFocus README](../README.md) - Project overview and quick start guide
- [JSON Schemas](./json-schemas.md) - JSON response formats for AI agents
- [Troubleshooting Guide](./troubleshooting.md) - Common issues and solutions
- [Project Instructions (CLAUDE.md)](../CLAUDE.md) - Developer documentation

---

**Last Updated:** 2026-01-30
