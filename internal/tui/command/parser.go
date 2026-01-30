// Package command provides command parsing for the TUI.
package command

import (
	"fmt"
	"strings"
)

// Command represents a parsed command
type Command struct {
	Name string
	Args []string
}

// Def defines a command with its aliases
type Def struct {
	Name        string
	Aliases     []string
	Description string
	ArgsHint    string // e.g., "<task name>", "[project name]"
}

// Available commands
var commands = []Def{
	{Name: "quit", Aliases: []string{"q", "exit"}, Description: "Quit application"},
	{Name: "refresh", Aliases: []string{"w", "sync"}, Description: "Refresh current view"},
	{Name: "add", Aliases: []string{"a"}, Description: "Add new task", ArgsHint: "<task name>"},
	{Name: "complete", Aliases: []string{"done", "c"}, Description: "Complete selected task"},
	{Name: "delete", Aliases: []string{"del", "rm"}, Description: "Delete selected task"},
	{Name: "project", Aliases: []string{"p"}, Description: "Filter by project", ArgsHint: "<project name>"},
	{Name: "tag", Aliases: []string{"t"}, Description: "Filter by tag", ArgsHint: "<tag name>"},
	{Name: "due", Aliases: []string{}, Description: "Filter by due date", ArgsHint: "<today|tomorrow|week>"},
	{Name: "flagged", Aliases: []string{}, Description: "Show only flagged tasks"},
	{Name: "clear", Aliases: []string{"reset"}, Description: "Clear all filters"},
	{Name: "help", Aliases: []string{"?"}, Description: "Show available commands"},
}

// Parser parses command strings
type Parser struct {
	aliasMap map[string]string // alias -> canonical name
}

// NewParser creates a new command parser
func NewParser() *Parser {
	aliasMap := make(map[string]string)

	for _, cmd := range commands {
		aliasMap[cmd.Name] = cmd.Name
		for _, alias := range cmd.Aliases {
			aliasMap[alias] = cmd.Name
		}
	}

	return &Parser{aliasMap: aliasMap}
}

// Parse parses a command string
func (p *Parser) Parse(input string) (*Command, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty command")
	}

	// Remove leading : if present
	if strings.HasPrefix(input, ":") {
		input = strings.TrimPrefix(input, ":")
		input = strings.TrimSpace(input)
	}

	parts := strings.SplitN(input, " ", 2)
	cmdName := strings.ToLower(parts[0])

	// Resolve alias
	canonicalName, exists := p.aliasMap[cmdName]
	if !exists {
		return nil, fmt.Errorf("unknown command: %s", cmdName)
	}

	// Parse arguments
	var args []string
	if len(parts) > 1 {
		argsStr := strings.TrimSpace(parts[1])
		if argsStr != "" {
			args = p.parseArgs(argsStr)
		}
	}

	return &Command{
		Name: canonicalName,
		Args: args,
	}, nil
}

// parseArgs splits arguments, respecting quotes
func (p *Parser) parseArgs(argsStr string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, ch := range argsStr {
		if !inQuotes && (ch == '"' || ch == '\'') {
			inQuotes = true
			quoteChar = ch
		} else if inQuotes && ch == quoteChar {
			inQuotes = false
			quoteChar = 0
		} else if !inQuotes && ch == ' ' {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

// GetCompletions returns command completions for the given prefix
func (p *Parser) GetCompletions(prefix string) []string {
	prefix = strings.ToLower(prefix)
	prefix = strings.TrimPrefix(prefix, ":")

	var completions []string
	for alias := range p.aliasMap {
		if strings.HasPrefix(alias, prefix) {
			completions = append(completions, alias)
		}
	}
	return completions
}

// GetCommands returns all available command definitions
func GetCommands() []Def {
	return commands
}
