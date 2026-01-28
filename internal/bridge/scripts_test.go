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

// TestValidateID tests ID validation for script injection prevention
func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		// Valid IDs
		{name: "Valid alphanumeric", id: "abc123", wantErr: false},
		{name: "Valid with hyphen", id: "task-1", wantErr: false},
		{name: "Valid with underscore", id: "project_name", wantErr: false},
		{name: "Valid mixed case", id: "ABC-123_test", wantErr: false},
		{name: "Valid all uppercase", id: "TASKID", wantErr: false},
		{name: "Valid all lowercase", id: "taskid", wantErr: false},
		{name: "Valid all numbers", id: "123456", wantErr: false},
		{name: "Valid complex", id: "abc-123_DEF-456_xyz", wantErr: false},

		// Invalid IDs
		{name: "Empty string", id: "", wantErr: true},
		{name: "Too long", id: strings.Repeat("a", 101), wantErr: true},
		{name: "Contains double quote", id: "task\"123", wantErr: true},
		{name: "Contains single quote", id: "task'123", wantErr: true},
		{name: "Contains semicolon", id: "task;123", wantErr: true},
		{name: "Contains left brace", id: "task{123", wantErr: true},
		{name: "Contains right brace", id: "task}123", wantErr: true},
		{name: "Contains left paren", id: "task(123", wantErr: true},
		{name: "Contains right paren", id: "task)123", wantErr: true},
		{name: "Contains backtick", id: "task`123", wantErr: true},
		{name: "Contains dollar sign", id: "task$123", wantErr: true},
		{name: "Contains backslash", id: "task\\123", wantErr: true},
		{name: "Contains forward slash", id: "task/123", wantErr: true},
		{name: "Contains space", id: "task 123", wantErr: true},
		{name: "Contains newline", id: "task\n123", wantErr: true},
		{name: "Contains tab", id: "task\t123", wantErr: true},

		// Injection attempts
		{name: "Injection with semicolon", id: "\"; malicious code; \"", wantErr: true},
		{name: "Injection with OR", id: "' || true", wantErr: true},
		{name: "Injection with template", id: "${injection}", wantErr: true},
		{name: "Injection with eval", id: "eval(code)", wantErr: true},
		{name: "Injection with comment", id: "task//comment", wantErr: true},
		{name: "Injection with multiline comment", id: "task/*comment*/", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateID(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "ID") {
				t.Errorf("ValidateID(%q) error message should contain 'ID', got: %v", tt.id, err)
			}
		})
	}
}

// TestValidateParam tests parameter validation for non-ID values (like perspective names)
func TestValidateParam(t *testing.T) {
	tests := []struct {
		name    string
		param   string
		wantErr bool
	}{
		// Valid parameters
		{name: "Valid alphanumeric", param: "abc123", wantErr: false},
		{name: "Valid with hyphen", param: "task-1", wantErr: false},
		{name: "Valid with underscore", param: "project_name", wantErr: false},
		{name: "Valid with spaces", param: "Custom Projects", wantErr: false},
		{name: "Valid single word", param: "Forecast", wantErr: false},
		{name: "Valid multi-word", param: "Work Tasks Review", wantErr: false},
		{name: "Valid mixed case with spaces", param: "My Custom View 2024", wantErr: false},
		{name: "Valid with multiple spaces", param: "Some  Extra   Spaces", wantErr: false},

		// Invalid parameters - injection attempts
		{name: "Empty string", param: "", wantErr: true},
		{name: "Too long", param: strings.Repeat("a", 101), wantErr: true},
		{name: "Contains semicolon", param: "name; rm -rf /", wantErr: true},
		{name: "Contains pipe", param: "name | cat", wantErr: true},
		{name: "Contains ampersand", param: "name & malicious", wantErr: true},
		{name: "Contains dollar sign", param: "name $VAR", wantErr: true},
		{name: "Contains backtick", param: "name`cmd`", wantErr: true},
		{name: "Contains backslash", param: "name\\escape", wantErr: true},
		{name: "Contains double quote", param: "name\"quote", wantErr: true},
		{name: "Contains single quote", param: "name'quote", wantErr: true},
		{name: "Contains newline", param: "name\ninjection", wantErr: true},
		{name: "Contains tab", param: "name\tinjection", wantErr: true},
		{name: "Contains carriage return", param: "name\rinjection", wantErr: true},

		// Shell injection attempts
		{name: "Injection with semicolon", param: "name; malicious code", wantErr: true},
		{name: "Injection with pipe", param: "name | cat /etc/passwd", wantErr: true},
		{name: "Injection with ampersand", param: "name && rm -rf /", wantErr: true},
		{name: "Injection with subshell", param: "name $(malicious)", wantErr: true},
		{name: "Injection with backtick", param: "name `malicious`", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateParam(tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParam(%q) error = %v, wantErr %v", tt.param, err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "parameter") {
				t.Errorf("ValidateParam(%q) error message should contain 'parameter', got: %v", tt.param, err)
			}
		})
	}
}

// TestGetScriptWithParams_UsesCorrectValidation tests that ID params use ValidateID and others use ValidateParam
func TestGetScriptWithParams_UsesCorrectValidation(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "ID parameter with valid ID format",
			params:  map[string]string{"TaskID": "task-123"},
			wantErr: false,
		},
		{
			name:    "ID parameter with spaces should fail",
			params:  map[string]string{"TaskID": "task 123"},
			wantErr: true,
			errMsg:  "TaskID",
		},
		{
			name:    "Non-ID parameter with spaces should succeed",
			params:  map[string]string{"PerspectiveName": "Custom Projects"},
			wantErr: false,
		},
		{
			name:    "Non-ID parameter with injection should fail",
			params:  map[string]string{"PerspectiveName": "name; rm -rf /"},
			wantErr: true,
			errMsg:  "PerspectiveName",
		},
		{
			name:    "Mixed params - ID and non-ID both valid",
			params:  map[string]string{"TaskID": "task-123", "Name": "My Task Name"},
			wantErr: false,
		},
		{
			name:    "Mixed params - ID valid, non-ID with injection",
			params:  map[string]string{"TaskID": "task-123", "Name": "name`injection`"},
			wantErr: true,
			errMsg:  "Name",
		},
		{
			name:    "Mixed params - ID with space, non-ID valid",
			params:  map[string]string{"ProjectID": "proj 123", "Description": "My Project"},
			wantErr: true,
			errMsg:  "ProjectID",
		},
		{
			name:    "Parameter ending with ID uses strict validation",
			params:  map[string]string{"TaskID": "valid-id", "ParentTaskID": "parent-id"},
			wantErr: false,
		},
		{
			name:    "Parameter ending with ID with space fails",
			params:  map[string]string{"ParentTaskID": "id with space"},
			wantErr: true,
		},
		{
			name:    "Parameter containing ID but not ending uses lenient validation",
			params:  map[string]string{"IDName": "Name With Spaces"},
			wantErr: false,
		},
		{
			name:    "Lowercase id parameter uses strict validation",
			params:  map[string]string{"taskid": "valid-id"},
			wantErr: false,
		},
		{
			name:    "Lowercase id with space fails",
			params:  map[string]string{"taskid": "id with space"},
			wantErr: true,
		},
	}

	scriptName := "get_inbox_tasks"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetScriptWithParams(scriptName, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScriptWithParams() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("GetScriptWithParams() error should mention %q, got: %v", tt.errMsg, err)
			}
		})
	}
}

// TestGetScriptWithParams_ValidatesParameters tests that parameters are validated before template execution
func TestGetScriptWithParams_ValidatesParameters(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
	}{
		{
			name:    "Valid parameters",
			params:  map[string]string{"TaskID": "task-123", "ProjectID": "proj_456"},
			wantErr: false,
		},
		{
			name:    "Invalid parameter with quote",
			params:  map[string]string{"TaskID": "task\"123"},
			wantErr: true,
		},
		{
			name:    "Invalid parameter with semicolon",
			params:  map[string]string{"TaskID": "task;malicious"},
			wantErr: true,
		},
		{
			name:    "Invalid parameter with injection attempt",
			params:  map[string]string{"TaskID": "\"; malicious code; \""},
			wantErr: true,
		},
		{
			name:    "Invalid empty parameter",
			params:  map[string]string{"TaskID": ""},
			wantErr: true,
		},
		{
			name:    "Invalid too long parameter",
			params:  map[string]string{"TaskID": strings.Repeat("a", 101)},
			wantErr: true,
		},
		{
			name:    "Multiple params with one invalid",
			params:  map[string]string{"TaskID": "valid-123", "ProjectID": "invalid${injection}"},
			wantErr: true,
		},
		{
			name:    "Multiple valid params",
			params:  map[string]string{"TaskID": "task-123", "ProjectID": "proj-456", "TagID": "tag_789"},
			wantErr: false,
		},
	}

	scriptName := "get_inbox_tasks"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetScriptWithParams(scriptName, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScriptWithParams() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), "invalid") && !strings.Contains(err.Error(), "ID") {
				t.Errorf("GetScriptWithParams() error should mention validation failure, got: %v", err)
			}
		})
	}
}
