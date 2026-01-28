package bridge

import (
	"strings"
	"testing"
)

func TestGetScript_ExistingScript_ReturnsContent(t *testing.T) {
	script, err := GetScript("get_inbox_tasks")
	if err != nil {
		t.Fatalf("GetScript() error = %v, want nil", err)
	}

	if script == "" {
		t.Fatal("GetScript() returned empty string, want non-empty content")
	}

	// Verify it's valid JavaScript content
	if !strings.Contains(script, "Application(\"OmniFocus\")") {
		t.Errorf("GetScript() content doesn't contain expected OmniFocus code")
	}
}

func TestGetScript_NonExistentScript_ReturnsError(t *testing.T) {
	_, err := GetScript("nonexistent_script")
	if err == nil {
		t.Fatal("GetScript() error = nil, want error for non-existent script")
	}

	// Error message should be informative
	if !strings.Contains(err.Error(), "nonexistent_script") {
		t.Errorf("GetScript() error = %v, want error message to mention script name", err)
	}
}

func TestListScripts_ReturnsAllAvailableScripts(t *testing.T) {
	scripts := ListScripts()

	if len(scripts) == 0 {
		t.Fatal("ListScripts() returned empty list, want at least one script")
	}

	// Verify expected scripts are present
	expectedScripts := map[string]bool{
		"get_inbox_tasks": false,
		"get_projects":    false,
	}

	for _, script := range scripts {
		if _, exists := expectedScripts[script]; exists {
			expectedScripts[script] = true
		}
	}

	for name, found := range expectedScripts {
		if !found {
			t.Errorf("ListScripts() missing expected script: %s", name)
		}
	}

	// Verify scripts are returned without .js extension
	for _, script := range scripts {
		if strings.HasSuffix(script, ".js") {
			t.Errorf("ListScripts() returned script with .js extension: %s, want without extension", script)
		}
	}
}

func TestGetScript_AllListedScripts_CanBeRetrieved(t *testing.T) {
	scripts := ListScripts()

	for _, name := range scripts {
		t.Run(name, func(t *testing.T) {
			content, err := GetScript(name)
			if err != nil {
				t.Errorf("GetScript(%q) error = %v, want nil", name, err)
			}
			if content == "" {
				t.Errorf("GetScript(%q) returned empty content", name)
			}
		})
	}
}
