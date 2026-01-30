package command

import (
	"testing"
)

func TestParse_SimpleCommand(t *testing.T) {
	p := NewParser()

	cmd, err := p.Parse("quit")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Name != "quit" {
		t.Errorf("name = %q, want %q", cmd.Name, "quit")
	}
	if len(cmd.Args) != 0 {
		t.Errorf("args = %v, want empty", cmd.Args)
	}
}

func TestParse_WithColonPrefix(t *testing.T) {
	p := NewParser()

	cmd, err := p.Parse(":quit")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Name != "quit" {
		t.Errorf("name = %q, want %q", cmd.Name, "quit")
	}
}

func TestParse_Alias(t *testing.T) {
	p := NewParser()

	cmd, err := p.Parse("q")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Name != "quit" {
		t.Errorf("name = %q, want %q", cmd.Name, "quit")
	}
}

func TestParse_WithArgs(t *testing.T) {
	p := NewParser()

	cmd, err := p.Parse("add Buy milk")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Name != "add" {
		t.Errorf("name = %q, want %q", cmd.Name, "add")
	}
	if len(cmd.Args) != 2 {
		t.Errorf("args count = %d, want 2", len(cmd.Args))
	}
}

func TestParse_WithQuotedArgs(t *testing.T) {
	p := NewParser()

	cmd, err := p.Parse(`add "Buy milk and eggs"`)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmd.Args) != 1 {
		t.Errorf("args count = %d, want 1", len(cmd.Args))
	}
	if cmd.Args[0] != "Buy milk and eggs" {
		t.Errorf("args[0] = %q, want %q", cmd.Args[0], "Buy milk and eggs")
	}
}

func TestParse_UnknownCommand(t *testing.T) {
	p := NewParser()

	_, err := p.Parse("unknown")

	if err == nil {
		t.Error("expected error for unknown command")
	}
}

func TestParse_EmptyCommand(t *testing.T) {
	p := NewParser()

	_, err := p.Parse("")

	if err == nil {
		t.Error("expected error for empty command")
	}
}

func TestGetCompletions(t *testing.T) {
	p := NewParser()

	completions := p.GetCompletions("qu")

	found := false
	for _, c := range completions {
		if c == "quit" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("completions should include 'quit', got %v", completions)
	}
}

func TestGetCompletions_WithColon(t *testing.T) {
	p := NewParser()

	completions := p.GetCompletions(":q")

	found := false
	for _, c := range completions {
		if c == "quit" || c == "q" {
			found = true
			break
		}
	}
	if !found {
		t.Error("completions should include 'quit' or 'q'")
	}
}
