# Troubleshooting Guide

This guide covers common issues you might encounter when using LazyFocus and how to resolve them.

## Table of Contents

- [Automation Permission Issues](#automation-permission-issues)
- [OmniFocus Not Running](#omnifocus-not-running)
- [Item Not Found Errors](#item-not-found-errors)
- [Date Parsing Errors](#date-parsing-errors)
- [Modification Errors](#modification-errors)
- [Tag Limitations](#tag-limitations)
- [Timeout Issues](#timeout-issues)
- [OmniFocus Pro Requirements](#omnifocus-pro-requirements)
- [General Troubleshooting Tips](#general-troubleshooting-tips)

---

## Automation Permission Issues

### Symptom
Commands fail with permission-related errors on first run, or you see system prompts about automation access.

### Cause
macOS requires explicit permission for applications to control other applications via automation. LazyFocus uses Omni Automation (osascript) to communicate with OmniFocus, which requires this permission.

### Solution
1. When you first run LazyFocus, macOS will display a permission prompt
2. Click "OK" or "Allow" to grant permission
3. If you accidentally denied permission or need to change it later:
   - Open **System Settings** (or **System Preferences** on older macOS versions)
   - Go to **Privacy & Security** → **Automation**
   - Find your terminal application (Terminal, iTerm2, etc.)
   - Enable the checkbox next to **OmniFocus**

### Example
```bash
# First run - triggers permission prompt
lazyfocus tasks

# If permission was denied, you'll see an error like:
# osascript execution failed: Not authorized to send Apple events to OmniFocus
```

### Additional Notes
- Each terminal application requires separate permission (Terminal vs. iTerm2)
- If using LazyFocus through scripts or other tools, those applications also need permission
- Permission is persistent once granted

---

## OmniFocus Not Running

### Symptom
Commands fail immediately with "OmniFocus is not running" error.

### Cause
LazyFocus requires OmniFocus to be running to execute commands. The Omni Automation bridge communicates with the running OmniFocus application.

### Solution
1. Launch **OmniFocus** before running LazyFocus commands
2. OmniFocus can run in the background (doesn't need to be the active window)
3. Consider adding OmniFocus to your login items if you use LazyFocus frequently

### Example
```bash
# This will fail if OmniFocus is not running
lazyfocus add "Buy groceries"
# Error: OmniFocus is not running

# Solution: Launch OmniFocus first
open -a OmniFocus
lazyfocus add "Buy groceries"
# Success!
```

### JSON Mode
```bash
# In JSON mode, error is returned as JSON
lazyfocus tasks --json
# Output: {"error": "OmniFocus is not running"}
```

---

## Item Not Found Errors

### Symptom
Commands report "task not found", "project not found", or "tag not found" errors.

### Cause
- The provided ID is invalid or doesn't exist
- The item was deleted in OmniFocus
- Typo in the ID when copying/pasting
- Project or tag name doesn't match exactly (names are case-sensitive)

### Solution

#### For Task Not Found
```bash
# Error example
lazyfocus show abc123xyz
# Error: task not found

# Get valid task IDs
lazyfocus tasks --json | jq -r '.tasks[].id'

# Or use human output to see task IDs
lazyfocus tasks
```

#### For Project Not Found
```bash
# Error example
lazyfocus add "New task" --project "work"
# Error: failed to resolve project: project not found

# List all projects to find the correct name
lazyfocus projects

# Project names are case-sensitive
lazyfocus add "New task" --project "Work"  # Correct
```

#### For Tag Not Found
```bash
# List all available tags
lazyfocus tags

# Ensure exact name match (case-sensitive)
lazyfocus add "Task" --tag "urgent"  # Must match exactly
```

### Additional Notes
- IDs in OmniFocus are persistent and unique
- Names must match exactly (case-sensitive, including spaces)
- Use `--json` flag to get machine-readable IDs for scripting

---

## Date Parsing Errors

### Symptom
Commands fail with "invalid due date" or "unrecognized date format" errors.

### Cause
The date string provided doesn't match any of the supported formats.

### Solution

#### Supported Date Formats

**Relative dates:**
```bash
lazyfocus add "Task" --due today
lazyfocus add "Task" --due tomorrow
lazyfocus add "Task" --due yesterday
```

**Next occurrence:**
```bash
lazyfocus add "Task" --due "next monday"
lazyfocus add "Task" --due "next week"
```

**In N units:**
```bash
lazyfocus add "Task" --due "in 3 days"
lazyfocus add "Task" --due "in 2 weeks"
```

**ISO format:**
```bash
lazyfocus add "Task" --due "2024-01-15"
```

**Month/day:**
```bash
lazyfocus add "Task" --due "Jan 15"
lazyfocus add "Task" --due "January 15"
lazyfocus add "Task" --due "Jan 15 2024"
lazyfocus add "Task" --due "January 15 2024"
```

#### Error Examples
```bash
# These will fail
lazyfocus add "Task" --due "next month"
# Error: invalid due date: unrecognized date format: next month

lazyfocus add "Task" --due "15/01/2024"
# Error: invalid due date: unrecognized date format: 15/01/2024

lazyfocus add "Task" --due "Jan-15"
# Error: invalid due date: unrecognized date format: Jan-15
```

### Additional Notes
- All dates without explicit times default to 5:00 PM local time
- Date parsing is case-insensitive ("TODAY" works the same as "today")
- Quotes are required if the date contains spaces

---

## Modification Errors

### Symptom
Modify command fails with "no modifications specified" error.

### Cause
The `modify` command requires at least one modification flag.

### Solution
```bash
# Error: no flags provided
lazyfocus modify task123
# Error: no modifications specified

# Correct: provide at least one modification
lazyfocus modify task123 --name "Updated name"
lazyfocus modify task123 --due tomorrow
lazyfocus modify task123 --flagged true
```

#### Available Modification Flags
- `--name` - Change task name
- `--note` - Change task note
- `--project` - Move to project
- `--add-tag` - Add tag (repeatable)
- `--remove-tag` - Remove tag (repeatable)
- `--due` - Set due date
- `--defer` - Set defer date
- `--flagged` - Set flagged status (true/false)
- `--clear-due` - Clear due date
- `--clear-defer` - Clear defer date

### Example
```bash
# Multiple modifications at once
lazyfocus modify task123 \
  --name "New name" \
  --due tomorrow \
  --flagged true \
  --add-tag urgent
```

---

## Tag Limitations

### Symptom
When adding multiple tags using natural syntax (`#tag1 #tag2`), only the first tag is applied to the task.

### Cause
OmniFocus Omni Automation API limitation. Tasks can have multiple tags, but the automation API only sets the first tag specified (the "primary" tag).

### Solution

**Using Natural Syntax (only first tag works):**
```bash
# Only #groceries will be applied
lazyfocus add "Buy milk #groceries #shopping"
```

**Using Command-Line Flags (only first tag works):**
```bash
# Only "urgent" will be applied
lazyfocus add "Task" --tag urgent --tag followup
```

**Workaround: Use modify command to add additional tags:**
```bash
# Step 1: Create task with first tag
lazyfocus add "Buy milk #groceries"

# Step 2: Get the task ID from output
# Output: { "id": "abc123", ... }

# Step 3: Add additional tags
lazyfocus modify abc123 --add-tag shopping --add-tag weekly
```

### Additional Notes
- This is a limitation of the OmniFocus automation API, not LazyFocus
- The `--remove-tag` flag only removes the primary tag if it matches
- Future OmniFocus updates may improve tag handling in automation

---

## Timeout Issues

### Symptom
Commands hang and eventually fail with "script execution timed out" error, especially with large datasets.

### Cause
The default timeout (30 seconds) may be too short for operations that query large numbers of tasks or projects.

### Solution

#### Use the `--timeout` Flag
```bash
# Default timeout (30 seconds)
lazyfocus tasks

# Increase timeout to 60 seconds
lazyfocus tasks --timeout 60s

# Or use minutes
lazyfocus tasks --timeout 2m
```

#### Timeout Format
- Seconds: `30s`, `60s`
- Minutes: `1m`, `2m`
- Hours: `1h` (rarely needed)

### Example
```bash
# Querying all tasks in a large OmniFocus database
lazyfocus tasks --timeout 60s

# Getting project details with many tasks
lazyfocus show project123 --timeout 45s
```

### Additional Notes
- Only increase timeout if you're experiencing actual timeout issues
- Very long timeouts might indicate performance issues in OmniFocus itself
- Consider filtering results to reduce query time:
  ```bash
  lazyfocus tasks --flagged  # Faster than all tasks
  lazyfocus tasks --inbox    # Faster than all tasks
  ```

---

## OmniFocus Pro Requirements

### Symptom
Custom perspective commands fail or show unexpected behavior.

### Cause
Custom perspectives are a feature of **OmniFocus Pro**. The standard version of OmniFocus only includes built-in perspectives.

### Solution

#### If You Have OmniFocus Pro
```bash
# Access custom perspectives
lazyfocus perspective "My Custom View"
```

#### If You Have Standard OmniFocus
Use built-in commands instead of custom perspectives:
```bash
# Instead of custom perspectives, use built-in filters
lazyfocus tasks --inbox
lazyfocus tasks --flagged
lazyfocus tasks --project "Work"
lazyfocus tasks --tag "urgent"
```

### Checking Your OmniFocus Version
1. Open OmniFocus
2. Go to **OmniFocus** → **About OmniFocus**
3. Look for "Pro" in the version information

### Additional Notes
- All other LazyFocus features work with standard OmniFocus
- Custom perspectives are only available in OmniFocus Pro license
- Consider upgrading if you rely heavily on custom perspectives

---

## General Troubleshooting Tips

### Enable Verbose Output
Add `--json` flag to see raw output from OmniFocus:
```bash
lazyfocus tasks --json | jq .
```

### Check OmniFocus Directly
Verify the data exists in OmniFocus itself:
1. Open OmniFocus
2. Check if the task/project/tag exists
3. Note the exact name (case-sensitive)

### Test with Minimal Example
Start with the simplest possible command:
```bash
# Simplest read command
lazyfocus tasks --inbox

# Simplest write command
lazyfocus add "Test task"
```

### Check for Updates
Ensure you're running the latest version:
```bash
lazyfocus version
```

### Debugging Commands
```bash
# List all projects
lazyfocus projects --json

# List all tags
lazyfocus tags --json

# Get specific task details
lazyfocus show <task-id> --json
```

### Common Error Patterns

#### Error: "confirmation required"
```bash
# When deleting without --force flag
lazyfocus delete abc123
# Error: confirmation required: use --force to delete without confirmation

# Solution: add --force flag
lazyfocus delete abc123 --force
```

#### Error: "task name is required"
```bash
# When adding task with empty name
lazyfocus add ""
# Error: task name is required

# Solution: provide a task name
lazyfocus add "Valid task name"
```

#### Error: "invalid flagged value"
```bash
# When using invalid boolean value
lazyfocus modify abc123 --flagged yes
# Error: invalid flagged value (use true/false)

# Solution: use true or false
lazyfocus modify abc123 --flagged true
```

### Getting Help
```bash
# General help
lazyfocus --help

# Command-specific help
lazyfocus add --help
lazyfocus modify --help
lazyfocus tasks --help
```

### JSON Mode for Debugging
JSON mode provides structured error messages:
```bash
lazyfocus add "Task" --project "NonExistent" --json
# Output: {"error": "failed to resolve project: project not found"}
```

### Exit Codes
LazyFocus uses standard Unix exit codes for scripting and automation:
- `0` - Success
- `1` - General error (invalid arguments, missing flags)
- `2` - OmniFocus not running or permission denied
- `3` - Item not found (task, project, or tag)
- `4` - Validation error (invalid input data)
- `5` - Permission error (automation access denied)

Check exit codes in scripts:
  ```bash
  if lazyfocus add "Task"; then
    echo "Success"
  else
    echo "Failed with exit code $?"
  fi
  ```

Example handling specific exit codes:
  ```bash
  lazyfocus add "Task"
  case $? in
    0) echo "Success" ;;
    1) echo "Invalid arguments" ;;
    2) echo "OmniFocus not running" ;;
    3) echo "Item not found" ;;
    4) echo "Invalid input" ;;
    5) echo "Permission denied" ;;
    *) echo "Unknown error" ;;
  esac
  ```

---

## Still Having Issues?

If you continue to experience problems:

1. **Check OmniFocus automation permissions** - The most common issue
2. **Verify OmniFocus is running** - Required for all commands
3. **Test with simple commands first** - Isolate the problem
4. **Check case sensitivity** - Project and tag names must match exactly
5. **Review supported date formats** - Use formats from this guide
6. **Try increasing timeout** - For large databases
7. **Use JSON mode for debugging** - Get structured error messages

Remember: LazyFocus requires OmniFocus to be installed, running, and granted automation permissions. Most issues can be resolved by verifying these prerequisites.

---

## See Also

- [Commands Reference](./commands.md) - Complete CLI command documentation with all flags and examples
- [JSON Schemas](./json-schemas.md) - JSON response formats for programmatic access
- [LazyFocus README](../README.md) - Project overview and installation instructions

---

**Last Updated:** 2026-01-30
