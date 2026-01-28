(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;
    const taskID = "{{.TaskID}}";

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

    // Extract tag names from task tags
    const taskTags = targetTask.tags;
    const tags = [];
    for (let j = 0; j < taskTags.length; j++) {
      tags.push(taskTags[j].name());
    }

    // Get project info if task belongs to a project
    const containingProject = targetTask.containingProject();
    const projectID = containingProject ? containingProject.id() : "";
    const projectName = containingProject ? containingProject.name() : "";

    // Convert dates to ISO 8601 format or null
    const dueDate = targetTask.dueDate();
    const deferDate = targetTask.deferDate();
    const completedDate = targetTask.completionDate();

    const task = {
      id: targetTask.id(),
      name: targetTask.name(),
      note: targetTask.note() || "",
      projectID: projectID,
      projectName: projectName,
      tags: tags,
      dueDate: dueDate ? dueDate.toISOString() : null,
      deferDate: deferDate ? deferDate.toISOString() : null,
      flagged: targetTask.flagged(),
      completed: targetTask.completed(),
      completedDate: completedDate ? completedDate.toISOString() : null
    };

    return JSON.stringify({ task: task }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
