//go:build integration

package bridge_test

import (
	"errors"
	"testing"

	"github.com/pwojciechowski/lazyfocus/internal/bridge"
)

func TestGetInboxTasks_Integration(t *testing.T) {
	// Create real executor
	executor := bridge.NewOSAScriptExecutor()

	// Load the get_inbox_tasks script
	script, err := bridge.GetScript("get_inbox_tasks")
	if err != nil {
		t.Fatalf("failed to load script: %v", err)
	}

	// Execute the script
	output, err := executor.Execute(script)

	// We expect either success or OmniFocus not running error
	if err != nil {
		if errors.Is(err, bridge.ErrOmniFocusNotRunning) {
			t.Skip("OmniFocus is not running - skipping integration test")
		}
		// Any other error is a failure
		t.Fatalf("unexpected error executing script: %v", err)
	}

	// Parse the results
	tasks, err := bridge.ParseTasks(output)
	if err != nil {
		if errors.Is(err, bridge.ErrOmniFocusNotRunning) {
			t.Skip("OmniFocus is not running - skipping integration test")
		}
		t.Fatalf("failed to parse tasks: %v", err)
	}

	// Verify we got a valid response (tasks slice should not be nil)
	if tasks == nil {
		t.Fatal("expected non-nil tasks slice")
	}

	// Log the number of tasks found
	t.Logf("Found %d inbox tasks", len(tasks))

	// If we got tasks, verify they have required fields
	for i, task := range tasks {
		if task.ID == "" {
			t.Errorf("task %d: missing ID", i)
		}
		if task.Name == "" {
			t.Errorf("task %d: missing Name", i)
		}
		t.Logf("Task %d: %s (ID: %s, Flagged: %v, Completed: %v)",
			i, task.Name, task.ID, task.Flagged, task.Completed)
	}
}

func TestGetProjects_Integration(t *testing.T) {
	// Create real executor
	executor := bridge.NewOSAScriptExecutor()

	// Load the get_projects script
	script, err := bridge.GetScript("get_projects")
	if err != nil {
		t.Fatalf("failed to load script: %v", err)
	}

	// Execute the script
	output, err := executor.Execute(script)

	// We expect either success or OmniFocus not running error
	if err != nil {
		if errors.Is(err, bridge.ErrOmniFocusNotRunning) {
			t.Skip("OmniFocus is not running - skipping integration test")
		}
		// Any other error is a failure
		t.Fatalf("unexpected error executing script: %v", err)
	}

	// Parse the results
	projects, err := bridge.ParseProjects(output)
	if err != nil {
		if errors.Is(err, bridge.ErrOmniFocusNotRunning) {
			t.Skip("OmniFocus is not running - skipping integration test")
		}
		t.Fatalf("failed to parse projects: %v", err)
	}

	// Verify we got a valid response (projects slice should not be nil)
	if projects == nil {
		t.Fatal("expected non-nil projects slice")
	}

	// Log the number of projects found
	t.Logf("Found %d projects", len(projects))

	// If we got projects, verify they have required fields
	for i, project := range projects {
		if project.ID == "" {
			t.Errorf("project %d: missing ID", i)
		}
		if project.Name == "" {
			t.Errorf("project %d: missing Name", i)
		}
		t.Logf("Project %d: %s (ID: %s, Status: %s)",
			i, project.Name, project.ID, project.Status)
	}
}

func TestEndToEnd_GetInboxTasksFlow(t *testing.T) {
	// This test verifies the complete flow:
	// 1. Create executor
	// 2. Load script
	// 3. Execute script
	// 4. Parse results

	// Step 1: Create executor
	executor := bridge.NewOSAScriptExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	// Step 2: Load script
	script, err := bridge.GetScript("get_inbox_tasks")
	if err != nil {
		t.Fatalf("step 2 failed - load script: %v", err)
	}
	if script == "" {
		t.Fatal("step 2 failed - script is empty")
	}

	// Step 3: Execute script
	output, err := executor.Execute(script)
	if err != nil {
		if errors.Is(err, bridge.ErrOmniFocusNotRunning) {
			t.Skip("OmniFocus is not running - skipping integration test")
		}
		t.Fatalf("step 3 failed - execute script: %v", err)
	}
	if output == "" {
		t.Fatal("step 3 failed - output is empty")
	}

	// Step 4: Parse results
	tasks, err := bridge.ParseTasks(output)
	if err != nil {
		if errors.Is(err, bridge.ErrOmniFocusNotRunning) {
			t.Skip("OmniFocus is not running - skipping integration test")
		}
		t.Fatalf("step 4 failed - parse tasks: %v", err)
	}

	// Success - log summary
	t.Logf("End-to-end test successful: retrieved %d tasks", len(tasks))
}
