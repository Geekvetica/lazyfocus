# Error State Component

A reusable Bubble Tea component for displaying errors in a modal overlay with optional retry functionality.

## Features

- **Modal Error Display**: Shows errors in a centered overlay that blocks interaction
- **Retry Support**: Optional retry command for retryable operations
- **Keyboard Navigation**:
  - `r` - Retry (if available)
  - `Enter` or `Esc` - Dismiss error
- **Consistent Styling**: Uses project's Lip Gloss styles for visual consistency
- **Auto-dismiss**: Automatically hides after user interaction

## API

### Types

```go
// Model represents the error state component
type Model struct { ... }

// ErrorDismissedMsg indicates the user dismissed an error
type ErrorDismissedMsg struct{}

// Styles for the error state
type Styles struct {
    Container lipgloss.Style
    Title     lipgloss.Style
    Message   lipgloss.Style
    Hint      lipgloss.Style
}
```

### Methods

```go
// New creates a new error state with default styles
func New() Model

// NewWithStyles creates a new error state with custom styles
func NewWithStyles(styles *tui.Styles) Model

// Show displays an error with optional retry command
func (m Model) Show(err error, retryCmd tea.Cmd) Model

// Hide hides the error state
func (m Model) Hide() Model

// IsVisible returns whether the error is visible
func (m Model) IsVisible() bool

// SetSize updates the dimensions
func (m Model) SetSize(width, height int) Model

// Init initializes the component (returns nil)
func (m Model) Init() tea.Cmd

// Update handles messages (key presses)
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)

// View renders the error modal
func (m Model) View() string
```

## Usage Example

### Basic Error (Non-retryable)

```go
// Show a validation error
m.errorState = m.errorState.Show(
    errors.New("Task name cannot be empty"),
    nil, // No retry command
)
```

### Retryable Error

```go
// Show a network error with retry option
m.errorState = m.errorState.Show(
    errors.New("Failed to connect to OmniFocus"),
    m.loadTasks(), // Command to retry
)
```

### Integration Pattern

```go
type AppModel struct {
    errorState errorstate.Model
    // ... other fields
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle error state first (highest priority)
    if m.errorState.IsVisible() {
        var cmd tea.Cmd
        m.errorState, cmd = m.errorState.Update(msg)
        return m, cmd
    }

    // Handle error occurred messages
    if errMsg, ok := msg.(tui.ErrorOccurredMsg); ok {
        m.errorState = m.errorState.Show(errMsg.Error, errMsg.RetryCmd)
        return m, nil
    }

    // Handle error dismissed
    if _, ok := msg.(errorstate.ErrorDismissedMsg); ok {
        // Optional: perform cleanup
        return m, nil
    }

    // ... handle other messages
}

func (m AppModel) View() string {
    view := m.renderMainView()

    // Layer error state overlay (render last for highest priority)
    if m.errorState.IsVisible() {
        view = m.layerOverlay(view, m.errorState.View())
    }

    return view
}
```

## Message Types

The component works with message types defined in `internal/tui/messages.go`:

```go
// ErrorOccurredMsg indicates an error occurred (app → component)
type ErrorOccurredMsg struct {
    Error     error
    Retryable bool
    RetryCmd  tea.Cmd
}

// ErrorDismissedMsg indicates user dismissed error (component → app)
type ErrorDismissedMsg struct{}
```

## Visual Example

```
╭──────────────────────────────╮
│                              │
│            Error             │
│                              │
│  Failed to connect to        │
│  OmniFocus. Please ensure    │
│  the application is running. │
│                              │
│    [r] Retry  [Enter/Esc]    │
│           Dismiss            │
│                              │
╰──────────────────────────────╯
```

## Testing

Comprehensive test coverage (97.8%) including:
- Visibility state management
- Retry functionality for retryable errors
- Dismiss actions (Enter/Esc keys)
- View rendering with/without retry option
- Keyboard input handling when hidden

Run tests:
```bash
go test ./internal/tui/components/errorstate/... -v -cover
```

## Design Decisions

1. **Follows Confirm Pattern**: Similar API to the existing `confirm` component for consistency
2. **Optional Retry**: Flexible - supports both retryable and non-retryable errors
3. **Blocking Modal**: High priority overlay that blocks user interaction
4. **Self-contained**: No external dependencies except standard TUI packages
5. **Message-based Communication**: Uses Bubble Tea message passing for state changes

## Implementation Status

- ✅ Component implementation
- ✅ Comprehensive test coverage
- ✅ Message types added to `internal/tui/messages.go`
- ✅ Integration documentation
- ⬚ Integration into main app (optional - can be done incrementally)

See `example_integration.md` for detailed integration guide.
