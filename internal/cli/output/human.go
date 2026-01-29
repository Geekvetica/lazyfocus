package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

// HumanFormatter implements Formatter interface for human-readable output
type HumanFormatter struct{}

// NewHumanFormatter creates a new human-readable formatter
func NewHumanFormatter() *HumanFormatter {
	return &HumanFormatter{}
}

// FormatTasks formats tasks in a human-readable format
func (f *HumanFormatter) FormatTasks(tasks []domain.Task, options TaskFormatOptions) string {
	var b strings.Builder

	// Header
	taskCount := len(tasks)
	taskWord := "task"
	if taskCount != 1 {
		taskWord = "tasks"
	}
	b.WriteString(fmt.Sprintf("TASKS (%d %s)\n", taskCount, taskWord))
	b.WriteString(strings.Repeat("─", 50) + "\n")

	if taskCount == 0 {
		b.WriteString("No tasks found\n")
		return b.String()
	}

	// Tasks
	for i, task := range tasks {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(f.formatTaskLine(task, options))
	}

	return b.String()
}

// FormatProjects formats projects in a human-readable format
func (f *HumanFormatter) FormatProjects(projects []domain.Project, options ProjectFormatOptions) string {
	var b strings.Builder

	// Header
	projectCount := len(projects)
	projectWord := "project"
	if projectCount != 1 {
		projectWord = "projects"
	}
	b.WriteString(fmt.Sprintf("PROJECTS (%d %s)\n", projectCount, projectWord))
	b.WriteString(strings.Repeat("─", 50) + "\n")

	if projectCount == 0 {
		b.WriteString("No projects found\n")
		return b.String()
	}

	// Projects
	for i, project := range projects {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(f.formatProjectSection(project, options))
	}

	return b.String()
}

// FormatTags formats tags in a human-readable format
func (f *HumanFormatter) FormatTags(tags []domain.Tag, options TagFormatOptions) string {
	var b strings.Builder

	// Header
	tagCount := len(tags)
	tagWord := "tag"
	if tagCount != 1 {
		tagWord = "tags"
	}
	b.WriteString(fmt.Sprintf("TAGS (%d %s)\n", tagCount, tagWord))
	b.WriteString(strings.Repeat("─", 50) + "\n")

	if tagCount == 0 {
		b.WriteString("No tags found\n")
		return b.String()
	}

	// Tags
	for i, tag := range tags {
		if i > 0 {
			b.WriteString("\n")
		}
		if options.Flat {
			b.WriteString(f.formatTagFlat(tag))
		} else {
			b.WriteString(f.formatTagHierarchical(tag, 0))
		}
	}

	return b.String()
}

// FormatTask formats a single task
func (f *HumanFormatter) FormatTask(task domain.Task) string {
	return f.formatTaskLine(task, TaskFormatOptions{
		ShowProject: true,
		ShowTags:    true,
	})
}

// FormatProject formats a single project
func (f *HumanFormatter) FormatProject(project domain.Project) string {
	return f.formatProjectSection(project, ProjectFormatOptions{
		ShowNotes: true,
		ShowTasks: false,
	})
}

// FormatTag formats a single tag
func (f *HumanFormatter) FormatTag(tag domain.Tag) string {
	return f.formatTagFlat(tag)
}

// FormatError formats an error message
func (f *HumanFormatter) FormatError(err error) string {
	return fmt.Sprintf("Error: %s\n", err.Error())
}

// FormatCreatedTask formats a newly created task
func (f *HumanFormatter) FormatCreatedTask(task domain.Task) string {
	var b strings.Builder

	// Success header
	b.WriteString(fmt.Sprintf("✓ Created task: %s\n", task.ID))
	b.WriteString(fmt.Sprintf("  %s\n", task.Name))

	// Due date (if present)
	if task.DueDate != nil {
		b.WriteString(fmt.Sprintf("  Due: %s\n", formatDate(*task.DueDate)))
	}

	// Tags (if present)
	if len(task.Tags) > 0 {
		tagStr := make([]string, len(task.Tags))
		for i, tag := range task.Tags {
			tagStr[i] = "#" + tag
		}
		b.WriteString(fmt.Sprintf("  Tags: %s\n", strings.Join(tagStr, ", ")))
	}

	// Project (if present)
	if task.ProjectName != "" {
		b.WriteString(fmt.Sprintf("  Project: %s\n", task.ProjectName))
	}

	return b.String()
}

// FormatModifiedTask formats a modified task
func (f *HumanFormatter) FormatModifiedTask(task domain.Task) string {
	var b strings.Builder

	// Success header
	b.WriteString(fmt.Sprintf("✓ Modified task: %s\n", task.ID))
	b.WriteString(fmt.Sprintf("  %s\n", task.Name))

	// Due date (if present)
	if task.DueDate != nil {
		b.WriteString(fmt.Sprintf("  Due: %s\n", formatDate(*task.DueDate)))
	}

	// Flagged status (only show if flagged)
	if task.Flagged {
		b.WriteString("  Flagged: Yes\n")
	}

	return b.String()
}

// FormatCompletedTask formats a completed task operation result
func (f *HumanFormatter) FormatCompletedTask(result domain.OperationResult) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("✓ Completed: %s\n", result.ID))
	b.WriteString("  Task marked as complete\n")

	return b.String()
}

// FormatDeletedTask formats a deleted task operation result
func (f *HumanFormatter) FormatDeletedTask(result domain.OperationResult) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("✓ Deleted: %s\n", result.ID))
	b.WriteString("  Task moved to trash\n")

	return b.String()
}

// formatTaskLine formats a single task line with icons and details
func (f *HumanFormatter) formatTaskLine(task domain.Task, options TaskFormatOptions) string {
	var b strings.Builder

	// Status icon
	if task.Completed {
		b.WriteString("☑ ")
	} else {
		b.WriteString("☐ ")
	}

	// Task name
	b.WriteString(task.Name)

	// Flag icon
	if task.Flagged {
		b.WriteString(" 🚩")
	}

	// Due date
	if task.DueDate != nil {
		b.WriteString(fmt.Sprintf("   📅 %s", formatDate(*task.DueDate)))
	}

	b.WriteString("\n")

	// Note (indented)
	if task.Note != "" {
		b.WriteString(fmt.Sprintf("  Note: %s\n", task.Note))
	}

	// Project name (if enabled)
	if options.ShowProject && task.ProjectName != "" {
		b.WriteString(fmt.Sprintf("  Project: %s\n", task.ProjectName))
	}

	// Tags (if enabled)
	if options.ShowTags && len(task.Tags) > 0 {
		tagStr := make([]string, len(task.Tags))
		for i, tag := range task.Tags {
			tagStr[i] = "#" + tag
		}
		b.WriteString(fmt.Sprintf("  %s\n", strings.Join(tagStr, " ")))
	}

	return b.String()
}

// formatProjectSection formats a project with optional details
func (f *HumanFormatter) formatProjectSection(project domain.Project, options ProjectFormatOptions) string {
	var b strings.Builder

	// Project name and status
	b.WriteString(fmt.Sprintf("📁 %s (%s)\n", project.Name, project.Status))

	// Note (if enabled)
	if options.ShowNotes && project.Note != "" {
		b.WriteString(fmt.Sprintf("  Note: %s\n", project.Note))
	}

	// Tasks (if enabled)
	if options.ShowTasks && len(project.Tasks) > 0 {
		b.WriteString(fmt.Sprintf("  Tasks: %d\n", len(project.Tasks)))
		for _, task := range project.Tasks {
			taskLine := f.formatTaskLine(task, TaskFormatOptions{})
			// Indent task lines
			lines := strings.Split(strings.TrimRight(taskLine, "\n"), "\n")
			for _, line := range lines {
				b.WriteString("    " + line + "\n")
			}
		}
	}

	return b.String()
}

// formatTagHierarchical formats a tag with its children hierarchically
func (f *HumanFormatter) formatTagHierarchical(tag domain.Tag, indent int) string {
	var b strings.Builder

	// Tag name with indentation
	b.WriteString(strings.Repeat("  ", indent))
	b.WriteString(fmt.Sprintf("#%s\n", tag.Name))

	// Children
	for _, child := range tag.Children {
		b.WriteString(f.formatTagHierarchical(child, indent+1))
	}

	return b.String()
}

// formatTagFlat formats a tag and its children in a flat list
func (f *HumanFormatter) formatTagFlat(tag domain.Tag) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("#%s\n", tag.Name))

	for _, child := range tag.Children {
		b.WriteString(f.formatTagFlat(child))
	}

	return b.String()
}

// formatDate formats a time.Time into a human-readable string
func formatDate(t time.Time) string {
	now := time.Now()

	// Check if it's today
	if isSameDay(t, now) {
		return "Today"
	}

	// Check if it's tomorrow
	tomorrow := now.AddDate(0, 0, 1)
	if isSameDay(t, tomorrow) {
		return "Tomorrow"
	}

	// Check if it's yesterday
	yesterday := now.AddDate(0, 0, -1)
	if isSameDay(t, yesterday) {
		return "Yesterday"
	}

	// Check if it's within the same year
	if t.Year() == now.Year() {
		return t.Format("Jan 2")
	}

	return t.Format("Jan 2, 2006")
}

// isSameDay checks if two times are on the same calendar day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
