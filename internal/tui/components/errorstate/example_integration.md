# Error State Component - Integration Guide

## Overview

The `errorstate` component provides a reusable modal for displaying errors with optional retry functionality. It follows the same patterns as other overlay components like `confirm` and `taskdetail`.

## Features

- Display error messages in a centered modal overlay
- Optional retry command for retryable errors
- Keyboard shortcuts: `r` for retry (if available), `Enter`/`Esc` to dismiss
- Consistent styling with other TUI components
- Auto-hides after user interaction

## Basic Usage

### 1. Add to Model

```go
import "github.com/pwojciechowski/lazyfocus/internal/tui/components/errorstate"

type Model struct {
    // ... other fields
    errorState errorstate.Model
}
```

### 2. Initialize in Constructor

```go
func NewApp(svc service.OmniFocusService) Model {
    styles := tui.DefaultStyles()

    return Model{
        // ... other fields
        errorState: errorstate.NewWithStyles(styles),
    }
}
```

### 3. Handle Window Resize

```go
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
    m.width = msg.Width
    m.height = msg.Height

    // Update error state dimensions
    m.errorState = m.errorState.SetSize(msg.Width, msg.Height)

    // ... update other overlays
    return m, nil
}
```

### 4. Handle Error Messages

Replace the existing `ErrorMsg` handler:

```go
// Before:
if msg, ok := msg.(tui.ErrorMsg); ok {
    m.err = msg.Err
    return m, nil
}

// After:
if msg, ok := msg.(tui.ErrorOccurredMsg); ok {
    m.errorState = m.errorState.Show(msg.Error, msg.RetryCmd)
    return m, nil
}

// Handle legacy ErrorMsg for backward compatibility
if msg, ok := msg.(tui.ErrorMsg); ok {
    m.errorState = m.errorState.Show(msg.Err, nil)
    return m, nil
}
```

### 5. Add to Overlay Handling

In `handleOverlays()`, add error state at the appropriate priority:

```go
func (m Model) handleOverlays(msg tea.Msg) (Model, tea.Cmd, bool) {
    // 1. Error state (highest - should block all other interactions)
    if m.errorState.IsVisible() {
        var cmd tea.Cmd
        m.errorState, cmd = m.errorState.Update(msg)
        return m, cmd, true
    }

    // 2. Confirm modal
    if m.confirmModal.IsVisible() {
        // ...
    }

    // ... other overlays
}
```

### 6. Handle Error Dismissed Message

```go
func (m Model) handleCustomMessages(msg tea.Msg) (Model, tea.Cmd, bool) {
    // Handle error dismissed
    if _, ok := msg.(errorstate.ErrorDismissedMsg); ok {
        // Optionally clear legacy error field
        m.err = nil
        return m, nil, true
    }

    // ... other message handlers
}
```

### 7. Render in View

Add error state overlay at the highest priority:

```go
func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }

    // Render current view
    view := m.getCurrentView()

    // Layer overlays (lowest to highest priority)
    // ... other overlays

    // Error state (highest priority - rendered last)
    if m.errorState.IsVisible() {
        view = m.layerOverlay(view, m.errorState.View())
    }

    return view
}
```

## Emitting Retryable Errors

When an operation fails and can be retried:

```go
func (m Model) loadTasks() tea.Cmd {
    return func() tea.Msg {
        tasks, err := m.service.GetTasks()
        if err != nil {
            // Return retryable error
            return tui.ErrorOccurredMsg{
                Error:     err,
                Retryable: true,
                RetryCmd:  m.loadTasks(), // Same command to retry
            }
        }
        return tui.TasksLoadedMsg{Tasks: tasks}
    }
}
```

## Emitting Non-Retryable Errors

For validation errors or other non-retryable issues:

```go
func (m Model) validateAndSave() tea.Cmd {
    return func() tea.Msg {
        if err := validate(); err != nil {
            // Return non-retryable error
            return tui.ErrorOccurredMsg{
                Error:     err,
                Retryable: false,
                RetryCmd:  nil,
            }
        }
        // ... save logic
    }
}
```

## Migration Path

### Phase 1: Add Component (Done)
- Component is implemented and tested
- Message types added to `internal/tui/messages.go`

### Phase 2: Integrate into App (Optional)
- Add `errorState` field to `internal/app/app.go` Model
- Update message handlers
- Update overlay rendering
- Keep legacy `m.err` field for backward compatibility initially

### Phase 3: Update Components (Optional)
- Replace `tui.ErrorMsg` with `tui.ErrorOccurredMsg` in commands
- Add retry logic where appropriate
- Components: quickadd, taskedit, inbox, projects, tags, forecast, review

### Phase 4: Cleanup (Future)
- Remove legacy `m.err` field
- Remove `tui.ErrorMsg` type (or mark deprecated)

## Testing

The component includes comprehensive tests covering:
- Visibility state management
- Retry functionality
- Dismiss actions (Enter/Esc)
- View rendering with and without retry
- Non-visible state handling

Run tests:
```bash
go test ./internal/tui/components/errorstate/... -v
```

## Design Decisions

1. **Separate from ErrorMsg**: Created new `ErrorOccurredMsg` to support retry functionality without breaking existing code

2. **Optional Retry**: `RetryCmd` can be nil for non-retryable errors, providing flexibility

3. **Modal Overlay**: Follows the same pattern as `confirm` component for consistency

4. **High Priority**: Error state should be highest priority overlay to prevent user actions during error display

5. **Auto-dismiss**: Errors disappear after user interaction, keeping the TUI responsive
