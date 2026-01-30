# Task Detail Component Integration Guide

## Overview

The task detail component (`internal/tui/components/taskdetail`) provides a modal overlay that displays comprehensive task information when a task is selected.

## Component Architecture

The component follows the Bubble Tea component pattern used in the codebase:

- **Model**: Stores task reference, visibility state, viewport for scrolling
- **Messages**: Emits action request messages (Edit, Complete, Delete, Flag)
- **View**: Renders as an overlay with header, scrollable content, and action hints
- **Lifecycle**: Show() → Update() → Hide() pattern

## Messages Emitted

The component emits these messages that the parent (app.go) should handle:

```go
taskdetail.CloseMsg              // User pressed Escape
taskdetail.EditRequestedMsg      // User pressed 'e'
taskdetail.CompleteRequestedMsg  // User pressed 'c'
taskdetail.DeleteRequestedMsg    // User pressed 'd'
taskdetail.FlagRequestedMsg      // User pressed 'f'
```

## Integration Steps

### 1. Add to Model (app.go)

```go
type Model struct {
    // ... existing fields ...
    taskDetail taskdetail.Model
}
```

### 2. Initialize in NewApp()

```go
func NewApp(svc service.OmniFocusService) Model {
    // ... existing code ...
    return Model{
        // ... existing fields ...
        taskDetail: taskdetail.New(styles, keys),
    }
}
```

### 3. Handle WindowSizeMsg

```go
if msg, ok := msg.(tea.WindowSizeMsg); ok {
    // ... existing resize logic ...
    m.taskDetail = m.taskDetail.SetSize(msg.Width, msg.Height)
    // ...
}
```

### 4. Handle Enter Key to Show Detail

In the global key handler section:

```go
if keyMsg, ok := msg.(tea.KeyMsg); ok {
    // ... existing keys ...

    // Show task detail on Enter
    if keyMsg.Type == tea.KeyEnter {
        task := m.getSelectedTask()
        if task != nil {
            m.taskDetail = m.taskDetail.Show(task)
        }
        return m, nil
    }
}
```

### 5. Delegate to Task Detail When Visible

```go
// After confirm modal but before quick add
if m.taskDetail.IsVisible() {
    var cmd tea.Cmd
    m.taskDetail, cmd = m.taskDetail.Update(msg)
    return m, cmd
}
```

### 6. Handle Task Detail Messages

```go
// Handle task detail close
if _, ok := msg.(taskdetail.CloseMsg); ok {
    m.taskDetail = m.taskDetail.Hide()
    return m, nil
}

// Handle edit request
if req, ok := msg.(taskdetail.EditRequestedMsg); ok {
    m.taskDetail = m.taskDetail.Hide()
    // TODO: Show edit overlay when implemented
    return m, nil
}

// Handle complete request
if req, ok := msg.(taskdetail.CompleteRequestedMsg); ok {
    m.taskDetail = m.taskDetail.Hide()
    return m, m.completeTask(req.TaskID)
}

// Handle delete request
if req, ok := msg.(taskdetail.DeleteRequestedMsg); ok {
    m.taskDetail = m.taskDetail.Hide()
    ctx := DeleteContext{TaskID: req.TaskID, TaskName: req.TaskName}
    m.confirmModal = m.confirmModal.ShowWithContext(
        "Delete Task",
        fmt.Sprintf("Delete \"%s\"?", req.TaskName),
        ctx,
    )
    return m, nil
}

// Handle flag request
if req, ok := msg.(taskdetail.FlagRequestedMsg); ok {
    // Keep detail view open, just update the flag
    return m, m.setTaskFlag(req.TaskID, req.Flagged)
}
```

### 7. Add Helper Function

```go
// setTaskFlag creates a command to set a task's flag status
func (m Model) setTaskFlag(taskID string, flagged bool) tea.Cmd {
    return func() tea.Msg {
        mod := domain.TaskModification{
            Flagged: &flagged,
        }
        result, err := m.service.ModifyTask(taskID, mod)
        if err != nil {
            return tui.ErrorMsg{Err: err}
        }
        return tui.TaskModifiedMsg{Task: *result}
    }
}
```

### 8. Overlay in View()

```go
func (m Model) View() string {
    // ... existing view code ...

    // Overlay task detail if visible (after quick add, before confirm)
    if m.taskDetail.IsVisible() {
        detailView := m.taskDetail.View()
        view = m.layerOverlay(view, detailView)
    }

    // ... rest of overlays ...
}
```

## Priority Order for Overlays

From highest to lowest priority:
1. Help overlay
2. Confirm modal
3. Task detail (NEW)
4. Quick add
5. Base view

## Key Bindings Used

- `Escape` - Close detail view
- `e` - Edit task (emits EditRequestedMsg)
- `c` - Complete task
- `d` - Delete task
- `f` - Toggle flag
- `j/k` or arrow keys - Scroll content

## Features

- Displays all task information (name, project, tags, dates, notes)
- Color-coded due dates (overdue in red, today in yellow)
- Flag icon for flagged tasks
- Checkbox icon shows completion status
- Scrollable viewport for long notes
- User-friendly date formatting

## Testing

Run tests with:
```bash
go test ./internal/tui/components/taskdetail/...
```

All tests verify:
- Component lifecycle (Show/Hide/IsVisible)
- Message emission for all actions
- Key handling
- View rendering with task data
- Visibility handling
