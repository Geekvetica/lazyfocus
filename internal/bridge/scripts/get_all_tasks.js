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
    const showCompleted = "{{.ShowCompleted}}" === "true";
    const flaggedOnly = "{{.FlaggedOnly}}" === "true";

    let allTasks = doc.flattenedTasks;
    const tasks = [];

    for (let i = 0; i < allTasks.length; i++) {
      const task = allTasks[i];

      // Skip completed unless requested
      if (!showCompleted && task.completed()) continue;

      // Skip if filtering for flagged only
      if (flaggedOnly && !task.flagged()) continue;

      // Extract tag names from task tags
      const taskTags = task.tags;
      const tags = [];
      for (let j = 0; j < taskTags.length; j++) {
        tags.push(taskTags[j].name());
      }

      // Get project info if task belongs to a project
      const containingProject = task.containingProject();
      const projectID = containingProject ? containingProject.id() : "";
      const projectName = containingProject ? containingProject.name() : "";

      // Convert dates to ISO 8601 format or null
      const dueDate = task.dueDate();
      const deferDate = task.deferDate();
      const completedDate = task.completionDate();

      tasks.push({
        id: task.id(),
        name: task.name(),
        note: task.note() || "",
        projectID: projectID,
        projectName: projectName,
        tags: tags,
        dueDate: dueDate ? dueDate.toISOString() : null,
        deferDate: deferDate ? deferDate.toISOString() : null,
        flagged: task.flagged(),
        completed: task.completed(),
        completedDate: completedDate ? completedDate.toISOString() : null
      });
    }

    return JSON.stringify({ tasks: tasks }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
