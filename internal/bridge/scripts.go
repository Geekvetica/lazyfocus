package bridge

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	scriptExtension = ".js"
	scriptsDir      = "scripts"
	maxIDLength     = 100
)

var idPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

//go:embed scripts/*.js
var scriptsFS embed.FS

// GetScript retrieves a script by name (without .js extension).
// Returns the script content as string, or error if not found.
func GetScript(name string) (string, error) {
	filename := name + scriptExtension
	scriptPath := filepath.Join(scriptsDir, filename)

	content, err := scriptsFS.ReadFile(scriptPath)
	if err != nil {
		return "", fmt.Errorf("script not found: %s", name)
	}

	return string(content), nil
}

// ListScripts returns all available script names (without .js extension).
func ListScripts() []string {
	entries, err := fs.ReadDir(scriptsFS, scriptsDir)
	if err != nil {
		return []string{}
	}

	scripts := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), scriptExtension) {
			// Remove .js extension
			name := strings.TrimSuffix(entry.Name(), scriptExtension)
			scripts = append(scripts, name)
		}
	}

	return scripts
}

// ValidateID checks that an ID is safe for use in script templates.
// It prevents script injection by ensuring IDs only contain safe characters.
// Returns error if ID is empty, too long, or contains unsafe characters.
func ValidateID(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("ID cannot be empty")
	}
	if len(id) > maxIDLength {
		return fmt.Errorf("ID too long: max %d characters", maxIDLength)
	}
	if !idPattern.MatchString(id) {
		return fmt.Errorf("invalid ID format: only alphanumeric, hyphen, underscore allowed")
	}
	return nil
}

// GetScriptWithParams retrieves a script and replaces placeholders with provided values.
// Placeholders use the format {{.ParamName}} and are replaced using Go's text/template.
// If params is nil or empty, returns the original script unchanged.
// All parameter values are validated before template execution to prevent script injection.
// Returns error if script is not found, parameters are invalid, or template parsing fails.
func GetScriptWithParams(name string, params map[string]string) (string, error) {
	// Get the base script
	script, err := GetScript(name)
	if err != nil {
		return "", err
	}

	// If no params provided, return original script
	if len(params) == 0 {
		return script, nil
	}

	// Validate all parameter values before template execution
	for key, value := range params {
		if err := ValidateID(value); err != nil {
			return "", fmt.Errorf("invalid parameter %q: %w", key, err)
		}
	}

	// Parse script as template
	tmpl, err := template.New(name).Parse(script)
	if err != nil {
		return "", fmt.Errorf("failed to parse script template: %w", err)
	}

	// Execute template with params
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", fmt.Errorf("failed to execute script template: %w", err)
	}

	return buf.String(), nil
}
