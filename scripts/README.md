# LazyFocus Omni Automation Scripts

This directory contains JavaScript for Automation (JXA) scripts that interface with OmniFocus on macOS.

## Scripts

### `get_inbox_tasks.js`

Fetches all tasks in the OmniFocus inbox (tasks not assigned to a project).

**Usage:**
```bash
osascript -l JavaScript scripts/get_inbox_tasks.js
```

**Output format:**
```json
{
  "tasks": [
    {
      "id": "abc123",
      "name": "Task name",
      "note": "Task notes",
      "tags": ["tag1", "tag2"],
      "dueDate": "2025-01-28T17:00:00.000Z",
      "deferDate": null,
      "flagged": false,
      "completed": false,
      "completedDate": null
    }
  ]
}
```

**Error format:**
```json
{
  "error": "OmniFocus is not running"
}
```

### `get_projects.js`

Fetches all active projects from OmniFocus.

**Usage:**
```bash
osascript -l JavaScript scripts/get_projects.js
```

**Output format:**
```json
{
  "projects": [
    {
      "id": "xyz789",
      "name": "Project Name",
      "status": "active",
      "note": "Project notes"
    }
  ]
}
```

**Error format:**
```json
{
  "error": "OmniFocus is not running"
}
```

## Testing Scripts

Test scripts directly from the command line:

```bash
# Test inbox tasks
osascript -l JavaScript scripts/get_inbox_tasks.js | jq

# Test projects
osascript -l JavaScript scripts/get_projects.js | jq

# Pipe to jq for pretty-printing (optional)
osascript -l JavaScript scripts/get_inbox_tasks.js | jq '.tasks[] | {name, dueDate}'
```

## Script Structure

All scripts follow this IIFE pattern:

```javascript
(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;

    // ... script logic ...

    return JSON.stringify(result, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
```

## Key OmniFocus API Notes

- Access tasks via `doc.inboxTasks` (property, not method)
- Access projects via `doc.flattenedProjects` (property, not method)
- Use `task.id()` to get unique identifier
- Date properties return JavaScript Date objects or null
- Convert dates to ISO 8601 using `date.toISOString()`
- Access arrays using index notation: `arr[i]`, not functional methods

## Requirements

- macOS
- OmniFocus 3 or 4
- Automation permissions granted to Terminal or the calling application
