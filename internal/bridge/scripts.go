package bridge

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
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
