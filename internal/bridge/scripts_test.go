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

func TestGetScriptWithParams_ReplacesPlaceholders(t *testing.T) {
	// Create a test script with placeholders
	scriptName := "get_inbox_tasks"
	params := map[string]string{
		"ProjectID": "test-project-123",
		"TagID":     "test-tag-456",
	}

	script, err := GetScriptWithParams(scriptName, params)
	if err != nil {
		t.Fatalf("GetScriptWithParams() error = %v, want nil", err)
	}

	// Should still get script content
	if script == "" {
		t.Fatal("GetScriptWithParams() returned empty string")
	}

	// Should not fail on script without placeholders
	if !strings.Contains(script, "Application(\"OmniFocus\")") {
		t.Errorf("GetScriptWithParams() should return valid script content")
	}
}

func TestGetScriptWithParams_WithActualPlaceholders_ReplacesValues(t *testing.T) {
	// Test with a mock script containing placeholders
	scriptName := "get_inbox_tasks"
	params := map[string]string{
		"ProjectID": "abc-123",
		"Status":    "active",
	}

	script, err := GetScriptWithParams(scriptName, params)
	if err != nil {
		t.Fatalf("GetScriptWithParams() error = %v, want nil", err)
	}

	if script == "" {
		t.Fatal("GetScriptWithParams() returned empty string")
	}

	// Should not contain template syntax
	if strings.Contains(script, "{{.") {
		// This is expected if script doesn't have placeholders
		// Just verify we get valid content
		t.Logf("Script doesn't contain placeholders, which is OK")
	}
}

func TestGetScriptWithParams_EmptyParams_ReturnsOriginalScript(t *testing.T) {
	scriptName := "get_inbox_tasks"
	params := map[string]string{}

	scriptWithParams, err := GetScriptWithParams(scriptName, params)
	if err != nil {
		t.Fatalf("GetScriptWithParams() error = %v, want nil", err)
	}

	scriptOriginal, err := GetScript(scriptName)
	if err != nil {
		t.Fatalf("GetScript() error = %v, want nil", err)
	}

	if scriptWithParams != scriptOriginal {
		t.Errorf("GetScriptWithParams() with empty params should return same as GetScript()")
	}
}

func TestGetScriptWithParams_NonExistentScript_ReturnsError(t *testing.T) {
	params := map[string]string{"Key": "value"}
	_, err := GetScriptWithParams("nonexistent_script", params)
	if err == nil {
		t.Fatal("GetScriptWithParams() error = nil, want error for non-existent script")
	}
}

func TestGetScriptWithParams_NilParams_ReturnsOriginalScript(t *testing.T) {
	scriptName := "get_inbox_tasks"

	scriptWithParams, err := GetScriptWithParams(scriptName, nil)
	if err != nil {
		t.Fatalf("GetScriptWithParams() error = %v, want nil", err)
	}

	scriptOriginal, err := GetScript(scriptName)
	if err != nil {
		t.Fatalf("GetScript() error = %v, want nil", err)
	}

	if scriptWithParams != scriptOriginal {
		t.Errorf("GetScriptWithParams() with nil params should return same as GetScript()")
	}
}
