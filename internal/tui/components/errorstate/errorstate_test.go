package errorstate_test

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/errorstate"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	m := errorstate.New()

	assert.False(t, m.IsVisible(), "New error state should not be visible")
}

func TestShow(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("test error")
	retryCmd := func() tea.Msg { return nil }

	m = m.Show(testErr, retryCmd)

	assert.True(t, m.IsVisible(), "Error state should be visible after Show()")
}

func TestHide(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("test error")

	m = m.Show(testErr, nil)
	assert.True(t, m.IsVisible(), "Error should be visible after Show()")

	m = m.Hide()
	assert.False(t, m.IsVisible(), "Error should be hidden after Hide()")
}

func TestUpdate_RetryKeyWithRetryableError(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("test error")
	retryCmd := func() tea.Msg { return "retried" }

	m = m.Show(testErr, retryCmd)

	// Press 'r' key
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

	assert.False(t, newModel.IsVisible(), "Error should be hidden after retry")
	assert.NotNil(t, cmd, "Should return retry command")
}

func TestUpdate_RetryKeyWithoutRetryableError(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("test error")

	m = m.Show(testErr, nil)

	// Press 'r' key
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

	assert.True(t, newModel.IsVisible(), "Error should still be visible (not retryable)")
	assert.Nil(t, cmd, "Should not return retry command when not retryable")
}

func TestUpdate_EscapeKeyDismissesError(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("test error")

	m = m.Show(testErr, nil)

	// Press Esc key
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	assert.False(t, newModel.IsVisible(), "Error should be hidden after Esc")
	assert.NotNil(t, cmd, "Should return ErrorDismissedMsg")

	// Verify the command returns the correct message
	msg := cmd()
	_, ok := msg.(errorstate.ErrorDismissedMsg)
	assert.True(t, ok, "Should return ErrorDismissedMsg")
}

func TestUpdate_EnterKeyDismissesError(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("test error")

	m = m.Show(testErr, nil)

	// Press Enter key
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.False(t, newModel.IsVisible(), "Error should be hidden after Enter")
	assert.NotNil(t, cmd, "Should return ErrorDismissedMsg")

	// Verify the command returns the correct message
	msg := cmd()
	_, ok := msg.(errorstate.ErrorDismissedMsg)
	assert.True(t, ok, "Should return ErrorDismissedMsg")
}

func TestUpdate_IgnoresKeysWhenNotVisible(t *testing.T) {
	m := errorstate.New()

	// Try pressing 'r' when not visible
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

	assert.False(t, newModel.IsVisible(), "Error should remain hidden")
	assert.Nil(t, cmd, "Should not return any command when hidden")
}

func TestView_ShowsErrorMessage(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("connection timeout")

	m = m.Show(testErr, nil)

	view := m.View()

	assert.Contains(t, view, "connection timeout", "View should contain error message")
}

func TestView_EmptyWhenHidden(t *testing.T) {
	m := errorstate.New()

	view := m.View()

	assert.Empty(t, view, "View should be empty when hidden")
}

func TestView_ShowsRetryHintWhenRetryable(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("test error")
	retryCmd := func() tea.Msg { return nil }

	m = m.Show(testErr, retryCmd)

	view := m.View()

	assert.Contains(t, view, "r", "View should show retry hint when retryable")
	assert.Contains(t, view, "Retry", "View should mention retry option")
}

func TestView_HidesRetryHintWhenNotRetryable(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("test error")

	m = m.Show(testErr, nil)

	view := m.View()

	// Should not contain retry hint (checking lowercase to avoid false positives)
	assert.NotContains(t, view, "[r]", "View should not show retry key when not retryable")
}

func TestView_ShowsDismissHint(t *testing.T) {
	m := errorstate.New()
	testErr := errors.New("test error")

	m = m.Show(testErr, nil)

	view := m.View()

	assert.Contains(t, view, "Enter", "View should show dismiss hint")
	assert.Contains(t, view, "Esc", "View should show escape hint")
}

func TestSetSize(t *testing.T) {
	m := errorstate.New()

	m = m.SetSize(100, 50)

	// Size is set internally - we can't test it directly,
	// but we can verify the method doesn't panic
	assert.NotNil(t, m, "SetSize should return valid model")
}

func TestInit(t *testing.T) {
	m := errorstate.New()

	cmd := m.Init()

	assert.Nil(t, cmd, "Init should return nil command")
}
