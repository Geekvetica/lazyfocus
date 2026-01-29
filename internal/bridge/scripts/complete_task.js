(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;

    // Template parameters (filled by Go)
    const taskID = "{{.TaskID}}";

    if (!taskID) {
      return JSON.stringify({ error: "Task ID is required" });
    }

    // Find the task by ID
    const allTasks = doc.flattenedTasks;
    let targetTask = null;

    for (let i = 0; i < allTasks.length; i++) {
      if (allTasks[i].id() === taskID) {
        targetTask = allTasks[i];
        break;
      }
    }

    if (!targetTask) {
      return JSON.stringify({ error: `Task not found: ${taskID}` });
    }

    // Mark the task as complete
    targetTask.markComplete();

    const result = {
      success: true,
      id: taskID,
      message: "Task completed"
    };

    return JSON.stringify(result, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
