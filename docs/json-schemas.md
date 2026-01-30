# JSON Schema Documentation for AI Agents

This document describes all JSON output formats for LazyFocus CLI commands when used with the `--json` flag. All responses use consistent structures to enable reliable parsing by AI agents and automation tools.

## Table of Contents

- [Overview](#overview)
- [Exit Codes](#exit-codes)
- [Domain Objects](#domain-objects)
  - [Task Object](#task-object)
  - [Project Object](#project-object)
  - [Tag Object](#tag-object)
  - [OperationResult Object](#operationresult-object)
- [Response Envelopes](#response-envelopes)
  - [List Response](#list-response)
  - [Single Item Response](#single-item-response)
  - [Operation Response](#operation-response)
  - [Error Response](#error-response)
- [Command Responses](#command-responses)
  - [tasks](#tasks)
  - [projects](#projects)
  - [tags](#tags)
  - [show](#show)
  - [add](#add)
  - [modify](#modify)
  - [complete](#complete)
  - [delete](#delete)

## Overview

All JSON output:
- Uses 2-space indentation
- Includes all fields with values
- Omits optional fields when they have no value (using `omitempty` JSON tag)
- Uses ISO 8601 format for all dates and timestamps
- Returns proper exit codes to indicate success or failure

## Exit Codes

LazyFocus uses the following exit codes:

| Code | Constant | Description |
|------|----------|-------------|
| 0 | `ExitSuccess` | Successful execution |
| 1 | `ExitGeneralError` | General error (invalid arguments, missing flags - see JSON error field for details) |
| 2 | `ExitOmniFocusNotRunning` | OmniFocus not running or permission denied |
| 3 | `ExitItemNotFound` | Requested item not found (task, project, or tag) |

For error scenarios, always check the JSON response for the `error` field which contains a human-readable error message.

**See [Commands Reference](./commands.md#exit-codes) for more details on exit codes.**

## Domain Objects

### Task Object

Represents a task in OmniFocus with the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier for the task |
| `name` | string | Yes | Task name/title |
| `note` | string | No | Optional note/description text |
| `projectId` | string | No | ID of the containing project |
| `projectName` | string | No | Name of the containing project |
| `tags` | string[] | No | Array of tag names assigned to the task |
| `dueDate` | string (ISO 8601) | No | Due date in ISO 8601 format (e.g., "2026-01-30T17:00:00Z") |
| `deferDate` | string (ISO 8601) | No | Defer date in ISO 8601 format |
| `flagged` | boolean | Yes | Whether the task is flagged (defaults to false) |
| `completed` | boolean | Yes | Whether the task is completed (defaults to false) |
| `completedDate` | string (ISO 8601) | No | Date when task was completed (only present if completed) |

#### Example Task Object

```json
{
  "id": "kGR3xMHww7P",
  "name": "Review PR #142",
  "note": "Check code style and test coverage",
  "projectId": "hKL4yNIxx8Q",
  "projectName": "Work",
  "tags": ["urgent", "code-review"],
  "dueDate": "2026-01-30T17:00:00Z",
  "deferDate": "2026-01-28T17:00:00Z",
  "flagged": true,
  "completed": false
}
```

#### Minimal Task Object

```json
{
  "id": "kGR3xMHww7P",
  "name": "Buy groceries",
  "flagged": false,
  "completed": false
}
```

### Project Object

Represents a project in OmniFocus with the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier for the project |
| `name` | string | Yes | Project name |
| `status` | string | Yes | Project status: "active", "on-hold", "completed", or "dropped" |
| `note` | string | No | Optional project note/description |
| `tasks` | Task[] | No | Array of tasks (only included in detailed views) |

#### Example Project Object

```json
{
  "id": "hKL4yNIxx8Q",
  "name": "Website Redesign",
  "status": "active",
  "note": "Complete redesign of company website",
  "tasks": [
    {
      "id": "kGR3xMHww7P",
      "name": "Design mockups",
      "flagged": false,
      "completed": false
    }
  ]
}
```

### Tag Object

Represents a tag in OmniFocus with the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier for the tag |
| `name` | string | Yes | Tag name |
| `parentId` | string | No | ID of parent tag (for nested tags) |
| `children` | Tag[] | No | Array of child tags (for hierarchical display) |

#### Example Tag Object

```json
{
  "id": "jPM5zOJyy9R",
  "name": "work",
  "children": [
    {
      "id": "kQN6AOKzz0S",
      "name": "meetings",
      "parentId": "jPM5zOJyy9R"
    }
  ]
}
```

### OperationResult Object

Represents the outcome of a write operation (complete, delete) with the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `success` | boolean | Yes | Whether the operation succeeded |
| `id` | string | Yes | ID of the affected task |
| `message` | string | Yes | Human-readable result message |

#### Example OperationResult Object

```json
{
  "success": true,
  "id": "kGR3xMHww7P",
  "message": "Task completed successfully"
}
```

## Response Envelopes

### List Response

Used for commands that return multiple items (`tasks`, `projects`, `tags`).

**Structure:**
```json
{
  "<items>": [<array of objects>],
  "count": <number>
}
```

Where `<items>` is `tasks`, `projects`, or `tags` depending on the command.

**Example:**
```json
{
  "tasks": [
    {
      "id": "kGR3xMHww7P",
      "name": "Buy groceries",
      "flagged": false,
      "completed": false
    },
    {
      "id": "lHS4zOJxx8Q",
      "name": "Review PR",
      "tags": ["urgent"],
      "dueDate": "2026-01-30T17:00:00Z",
      "flagged": true,
      "completed": false
    }
  ],
  "count": 2
}
```

### Single Item Response

Used for commands that return a single item in detail (`show`).

**Structure:**
```json
{
  "<item>": <object>
}
```

Where `<item>` is `task`, `project`, or `tag`.

**Example:**
```json
{
  "task": {
    "id": "kGR3xMHww7P",
    "name": "Review PR #142",
    "projectName": "Work",
    "tags": ["urgent", "code-review"],
    "dueDate": "2026-01-30T17:00:00Z",
    "flagged": true,
    "completed": false
  }
}
```

### Operation Response

Used for commands that create or modify items (`add`, `modify`, `complete`, `delete`).

**Structure for add/modify:**
```json
{
  "success": true,
  "task": <Task object>
}
```

**Structure for complete/delete:**
```json
{
  "success": <boolean>,
  "id": "<task-id>",
  "message": "<status message>"
}
```

**Example (add/modify):**
```json
{
  "success": true,
  "task": {
    "id": "kGR3xMHww7P",
    "name": "Buy groceries",
    "tags": ["shopping"],
    "dueDate": "2026-01-30T17:00:00Z",
    "flagged": false,
    "completed": false
  }
}
```

**Example (complete):**
```json
{
  "success": true,
  "id": "kGR3xMHww7P",
  "message": "Task completed successfully"
}
```

### Error Response

Used when any command fails. Always check the exit code and this response structure.

**Structure:**
```json
{
  "error": "<error message>"
}
```

**Common error messages:**
- `"task not found"` - Invalid task ID provided
- `"project not found"` - Invalid project name or ID
- `"failed to resolve project: project not found"` - Project specified in add/modify not found
- `"task name is required"` - Attempted to add task without name
- `"no modifications specified"` - Called modify without any modification flags
- `"OmniFocus is not running"` - OmniFocus application is not running
- `"invalid due date: unrecognized date format: xyz"` - Invalid date format provided
- `"confirmation required: use --force to delete"` - Delete command requires --force flag

**Example:**
```json
{
  "error": "task not found"
}
```

Exit code will be non-zero (typically 1 for general errors, 3 for not found errors).

## Command Responses

### tasks

Lists tasks with optional filtering.

**Command:**
```bash
lazyfocus tasks --json
lazyfocus tasks --flagged --json
lazyfocus tasks --project Work --json
```

**Response:**
```json
{
  "tasks": [
    {
      "id": "kGR3xMHww7P",
      "name": "Buy groceries",
      "tags": ["errands"],
      "dueDate": "2026-01-30T17:00:00Z",
      "flagged": false,
      "completed": false
    },
    {
      "id": "lHS4zOJxx8Q",
      "name": "Review PR #142",
      "projectName": "Work",
      "tags": ["urgent", "code-review"],
      "flagged": true,
      "completed": false
    }
  ],
  "count": 2
}
```

**Empty result:**
```json
{
  "tasks": [],
  "count": 0
}
```

### projects

Lists all projects.

**Command:**
```bash
lazyfocus projects --json
```

**Response:**
```json
{
  "projects": [
    {
      "id": "hKL4yNIxx8Q",
      "name": "Work",
      "status": "active"
    },
    {
      "id": "iLM5zOKyy9R",
      "name": "Personal",
      "status": "active",
      "note": "Personal tasks and goals"
    }
  ],
  "count": 2
}
```

### tags

Lists all tags.

**Command:**
```bash
lazyfocus tags --json
```

**Response:**
```json
{
  "tags": [
    {
      "id": "jPM5zOJyy9R",
      "name": "work",
      "children": [
        {
          "id": "kQN6AOKzz0S",
          "name": "meetings",
          "parentId": "jPM5zOJyy9R"
        }
      ]
    },
    {
      "id": "lRO7BPLAz1T",
      "name": "urgent"
    }
  ],
  "count": 2
}
```

### show

Shows detailed information about a single task.

**Command:**
```bash
lazyfocus show kGR3xMHww7P --json
```

**Response (success):**
```json
{
  "task": {
    "id": "kGR3xMHww7P",
    "name": "Review PR #142",
    "note": "Check code style and test coverage",
    "projectId": "hKL4yNIxx8Q",
    "projectName": "Work",
    "tags": ["urgent", "code-review"],
    "dueDate": "2026-01-30T17:00:00Z",
    "flagged": true,
    "completed": false
  }
}
```

**Response (not found):**
```json
{
  "error": "task not found"
}
```
Exit code: 3

### add

Creates a new task.

**Command:**
```bash
lazyfocus add "Buy groceries #shopping due:tomorrow" --json
lazyfocus add "Review PR" --project Work --tag urgent --due friday --json
```

**Response (success):**
```json
{
  "success": true,
  "task": {
    "id": "kGR3xMHww7P",
    "name": "Buy groceries",
    "tags": ["shopping"],
    "dueDate": "2026-01-31T17:00:00Z",
    "flagged": false,
    "completed": false
  }
}
```

**Response (error):**
```json
{
  "error": "task name is required"
}
```
Exit code: 1

**Response (project not found):**
```json
{
  "error": "failed to resolve project: project not found"
}
```
Exit code: 3

### modify

Modifies an existing task.

**Command:**
```bash
lazyfocus modify kGR3xMHww7P --name "Updated name" --flagged true --json
lazyfocus modify kGR3xMHww7P --due tomorrow --add-tag urgent --json
```

**Response (success):**
```json
{
  "success": true,
  "task": {
    "id": "kGR3xMHww7P",
    "name": "Updated name",
    "dueDate": "2026-01-31T17:00:00Z",
    "tags": ["urgent"],
    "flagged": true,
    "completed": false
  }
}
```

**Response (task not found):**
```json
{
  "error": "task not found"
}
```
Exit code: 3

**Response (no modifications):**
```json
{
  "error": "no modifications specified"
}
```
Exit code: 1

### complete

Marks one or more tasks as complete.

**Command:**
```bash
lazyfocus complete kGR3xMHww7P --json
lazyfocus complete task1 task2 task3 --json
```

**Response (single task success):**
```json
{
  "success": true,
  "id": "kGR3xMHww7P",
  "message": "Task completed successfully"
}
```

**Response (multiple tasks):**

When completing multiple tasks, LazyFocus outputs one JSON object per line (JSONL format):

```json
{"success": true, "id": "task1", "message": "Task completed successfully"}
{"success": true, "id": "task2", "message": "Task completed successfully"}
{"success": false, "id": "task3", "message": "task not found"}
```

Each line is a valid JSON object representing the result for one task. Parse line by line.

**Response (task not found):**
```json
{
  "success": false,
  "id": "invalid-id",
  "message": "task not found"
}
```
Exit code: 0 (command itself succeeded, check individual success fields)

### delete

Deletes one or more tasks.

**Command:**
```bash
lazyfocus delete kGR3xMHww7P --force --json
lazyfocus delete task1 task2 --force --json
```

**Response (single task success):**
```json
{
  "success": true,
  "id": "kGR3xMHww7P",
  "message": "Task deleted successfully"
}
```

**Response (multiple tasks):**

Similar to `complete`, outputs JSONL format (one JSON object per line):

```json
{"success": true, "id": "task1", "message": "Task deleted successfully"}
{"success": true, "id": "task2", "message": "Task deleted successfully"}
{"success": false, "id": "task3", "message": "task not found"}
```

**Response (task not found):**
```json
{
  "success": false,
  "id": "invalid-id",
  "message": "task not found"
}
```
Exit code: 0 (command itself succeeded, check individual success fields)

**Response (missing --force in non-JSON mode):**

Note: When using `--json` flag, the `--force` requirement is automatically bypassed. However, if somehow required:

```json
{
  "error": "confirmation required: use --force to delete"
}
```
Exit code: 1

## Date Format

All dates and timestamps use ISO 8601 format with timezone information:

**Format:** `YYYY-MM-DDTHH:MM:SSZ`

**Examples:**
- `2026-01-30T17:00:00Z` - January 30, 2026 at 5:00 PM UTC
- `2026-02-15T09:30:00Z` - February 15, 2026 at 9:30 AM UTC

Dates without explicit times default to 5:00 PM (17:00) local time when set via natural language parsing.

## Parsing Guidelines for AI Agents

1. **Always check the exit code first** - Exit code 0 means success, non-zero means error
2. **For list commands** - Parse the array field (`tasks`, `projects`, `tags`) and verify `count` matches array length
3. **For single item commands** - Extract the nested object (`task`, `project`, `tag`)
4. **For write operations** - Check `success` field and use `task` object or `id`/`message` as appropriate
5. **For errors** - Parse the `error` field for the human-readable message
6. **For multiple task operations** - Parse line by line (JSONL format), check each `success` field independently
7. **For optional fields** - Handle missing fields gracefully (they are omitted when null/empty)
8. **For dates** - Parse as ISO 8601 timestamps, handle null values for unset dates
9. **For boolean fields** - `flagged` and `completed` are always present (default false)

## Example Parsing Pseudocode

```
function parseTasksResponse(jsonString, exitCode):
  if exitCode != 0:
    error = parseJSON(jsonString).error
    return Error(error)

  data = parseJSON(jsonString)
  if data.tasks exists:
    return data.tasks  // Array of Task objects
  else:
    return Error("Invalid response format")

function parseAddResponse(jsonString, exitCode):
  if exitCode != 0:
    error = parseJSON(jsonString).error
    return Error(error)

  data = parseJSON(jsonString)
  if data.success and data.task exists:
    return data.task  // Task object with new ID
  else:
    return Error("Task creation failed")

function parseCompleteResponse(jsonString, exitCode):
  results = []
  for each line in jsonString.split("\n"):
    if line is not empty:
      result = parseJSON(line)
      results.append({
        id: result.id,
        success: result.success,
        message: result.message
      })
  return results
```

## Versioning

This schema documentation is for LazyFocus CLI version 0.1.0+. Future versions may add new fields but will maintain backward compatibility for existing fields. Parsers should ignore unknown fields gracefully.

---

## See Also

- [Commands Reference](./commands.md) - Complete CLI command documentation
- [Troubleshooting Guide](./troubleshooting.md) - Common issues and solutions
- [LazyFocus README](../README.md) - Project overview and quick start

---

**Last Updated:** 2026-01-30
