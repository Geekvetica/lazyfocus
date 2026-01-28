package bridge

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	scriptExtension = ".js"
	scriptsDir      = "scripts"
)

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

// GetScriptWithParams retrieves a script and replaces placeholders with provided values.
// Placeholders use the format {{.ParamName}} and are replaced using Go's text/template.
// If params is nil or empty, returns the original script unchanged.
// Returns error if script is not found or template parsing fails.
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
