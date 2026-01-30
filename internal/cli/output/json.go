package output

import (
	"encoding/json"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

// JSONFormatter implements Formatter interface for JSON output
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// FormatTasks formats tasks as JSON
func (f *JSONFormatter) FormatTasks(tasks []domain.Task, options TaskFormatOptions) string {
	output := map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	}
	return f.marshal(output)
}

// FormatProjects formats projects as JSON
func (f *JSONFormatter) FormatProjects(projects []domain.Project, options ProjectFormatOptions) string {
	output := map[string]interface{}{
		"projects": projects,
		"count":    len(projects),
	}
	return f.marshal(output)
}

// FormatTags formats tags as JSON
func (f *JSONFormatter) FormatTags(tags []domain.Tag, options TagFormatOptions) string {
	output := map[string]interface{}{
		"tags":  tags,
		"count": len(tags),
	}
	return f.marshal(output)
}

// FormatTask formats a single task as JSON
func (f *JSONFormatter) FormatTask(task domain.Task) string {
	output := map[string]interface{}{
		"task": task,
	}
	return f.marshal(output)
}

// FormatProject formats a single project as JSON
func (f *JSONFormatter) FormatProject(project domain.Project) string {
	output := map[string]interface{}{
		"project": project,
	}
	return f.marshal(output)
}

// FormatTag formats a single tag as JSON
func (f *JSONFormatter) FormatTag(tag domain.Tag) string {
	output := map[string]interface{}{
		"tag": tag,
	}
	return f.marshal(output)
}

// FormatError formats an error as JSON
func (f *JSONFormatter) FormatError(err error) string {
	output := map[string]interface{}{
		"error": err.Error(),
	}

	// Check if it's a LazyFocusError - need to import the errors package
	// This is done via type assertion to avoid import cycle
	type lazyFocusError interface {
		error
		ExitCode() int
		Suggestion() string
	}

	if lfErr, ok := err.(lazyFocusError); ok {
		output["code"] = lfErr.ExitCode()
		if suggestion := lfErr.Suggestion(); suggestion != "" {
			output["suggestion"] = suggestion
		}
	}

	return f.marshal(output)
}

// FormatCreatedTask formats a newly created task as JSON
func (f *JSONFormatter) FormatCreatedTask(task domain.Task) string {
	output := map[string]interface{}{
		"success": true,
		"task":    task,
	}
	return f.marshal(output)
}

// FormatModifiedTask formats a modified task as JSON
func (f *JSONFormatter) FormatModifiedTask(task domain.Task) string {
	output := map[string]interface{}{
		"success": true,
		"task":    task,
	}
	return f.marshal(output)
}

// FormatCompletedTask formats a completed task operation result as JSON
func (f *JSONFormatter) FormatCompletedTask(result domain.OperationResult) string {
	output := map[string]interface{}{
		"success": result.Success,
		"id":      result.ID,
		"message": result.Message,
	}
	return f.marshal(output)
}

// FormatDeletedTask formats a deleted task operation result as JSON
func (f *JSONFormatter) FormatDeletedTask(result domain.OperationResult) string {
	output := map[string]interface{}{
		"success": result.Success,
		"id":      result.ID,
		"message": result.Message,
	}
	return f.marshal(output)
}

// marshal converts data to indented JSON string
func (f *JSONFormatter) marshal(data interface{}) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		// Fallback to error JSON if marshaling fails
		return `{"error": "failed to marshal JSON"}`
	}
	return string(bytes)
}
