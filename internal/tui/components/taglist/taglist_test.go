package taglist

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()

	m := New(styles, keys)

	if len(m.tags) != 0 {
		t.Errorf("new model should have 0 tags, got %d", len(m.tags))
	}
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
	if !m.empty {
		t.Error("should be empty initially")
	}
}

func TestSetTags(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	tags := []domain.Tag{
		{ID: "t1", Name: "Tag 1"},
		{ID: "t2", Name: "Tag 2"},
	}
	counts := map[string]int{"t1": 5, "t2": 3}

	m = m.SetTags(tags, counts)

	if len(m.tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(m.tags))
	}
	if m.empty {
		t.Error("should not be empty after setting tags")
	}
	if m.tags[0].Count != 5 {
		t.Errorf("tag 1 count = %d, want 5", m.tags[0].Count)
	}
	if m.tags[1].Count != 3 {
		t.Errorf("tag 2 count = %d, want 3", m.tags[1].Count)
	}
}

func TestNavigationDown(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	tags := []domain.Tag{
		{ID: "t1", Name: "Tag 1"},
		{ID: "t2", Name: "Tag 2"},
		{ID: "t3", Name: "Tag 3"},
	}
	counts := map[string]int{"t1": 5, "t2": 3, "t3": 1}
	m = m.SetTags(tags, counts)

	// Move down
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.cursor)
	}

	// Move down again
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 2 {
		t.Errorf("cursor = %d, want 2", m.cursor)
	}

	// Wrap around
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 (wrapped)", m.cursor)
	}
}

func TestNavigationUp(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	tags := []domain.Tag{
		{ID: "t1", Name: "Tag 1"},
		{ID: "t2", Name: "Tag 2"},
		{ID: "t3", Name: "Tag 3"},
	}
	counts := map[string]int{"t1": 5, "t2": 3, "t3": 1}
	m = m.SetTags(tags, counts)

	// Wrap around to bottom
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 2 {
		t.Errorf("cursor = %d, want 2 (wrapped)", m.cursor)
	}

	// Move up
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.cursor)
	}
}

func TestSelectedTag(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	tags := []domain.Tag{
		{ID: "t1", Name: "Tag 1"},
		{ID: "t2", Name: "Tag 2"},
	}
	counts := map[string]int{"t1": 5, "t2": 3}
	m = m.SetTags(tags, counts)

	// First tag selected
	tag := m.SelectedTag()
	if tag == nil {
		t.Fatal("expected selected tag, got nil")
	}
	if tag.ID != "t1" {
		t.Errorf("selected tag ID = %q, want %q", tag.ID, "t1")
	}

	// Move down and check
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	tag = m.SelectedTag()
	if tag == nil {
		t.Fatal("expected selected tag, got nil")
	}
	if tag.ID != "t2" {
		t.Errorf("selected tag ID = %q, want %q", tag.ID, "t2")
	}
}

func TestSelectedTag_Empty(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	tag := m.SelectedTag()
	if tag != nil {
		t.Error("expected nil for empty list")
	}
}

func TestViewLoading(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)
	m = m.SetLoading(true)

	view := m.View()
	if !strings.Contains(view, "Loading") {
		t.Error("loading view should contain 'Loading'")
	}
}

func TestViewEmpty(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	view := m.View()
	if !strings.Contains(view, "No tags") {
		t.Error("empty view should contain 'No tags'")
	}
}

func TestViewTags(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	tags := []domain.Tag{
		{ID: "t1", Name: "Tag 1"},
	}
	counts := map[string]int{"t1": 5}
	m = m.SetTags(tags, counts)

	view := m.View()
	if !strings.Contains(view, "Tag 1") {
		t.Error("view should contain tag name")
	}
	if !strings.Contains(view, "(5)") {
		t.Error("view should contain task count")
	}
}

func TestHierarchicalTags(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	childTag := domain.Tag{ID: "t2", Name: "Child Tag"}
	parentTag := domain.Tag{ID: "t1", Name: "Parent Tag", Children: []domain.Tag{childTag}}

	tags := []domain.Tag{parentTag}
	counts := map[string]int{"t1": 5, "t2": 3}
	m = m.SetTags(tags, counts)

	// Should have 2 tags in flattened list
	if len(m.tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(m.tags))
	}

	// First tag should have depth 0
	if m.tags[0].Depth != 0 {
		t.Errorf("parent tag depth = %d, want 0", m.tags[0].Depth)
	}

	// Second tag should have depth 1
	if m.tags[1].Depth != 1 {
		t.Errorf("child tag depth = %d, want 1", m.tags[1].Depth)
	}

	// Child should be "Child Tag"
	if m.tags[1].Tag.Name != "Child Tag" {
		t.Errorf("child tag name = %q, want %q", m.tags[1].Tag.Name, "Child Tag")
	}
}

func TestCursorClamp(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	tags := []domain.Tag{
		{ID: "t1", Name: "Tag 1"},
		{ID: "t2", Name: "Tag 2"},
		{ID: "t3", Name: "Tag 3"},
	}
	counts := map[string]int{"t1": 5, "t2": 3, "t3": 1}
	m = m.SetTags(tags, counts)

	// Set cursor beyond bounds
	m.cursor = 10
	m = m.SetTags(tags, counts)

	if m.cursor != 2 {
		t.Errorf("cursor = %d, want 2 (clamped to last item)", m.cursor)
	}
}
