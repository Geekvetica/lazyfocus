(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;
    const tagID = "{{.TagID}}";

    // Find the tag by ID
    const allTags = doc.flattenedTags;
    let targetTag = null;

    for (let i = 0; i < allTags.length; i++) {
      if (allTags[i].id() === tagID) {
        targetTag = allTags[i];
        break;
      }
    }

    if (!targetTag) {
      return JSON.stringify({ error: `Tag not found: ${tagID}` });
    }

    // Get all tasks and filter by tag
    const allTasks = doc.flattenedTasks;
    const tasks = [];

    for (let i = 0; i < allTasks.length; i++) {
      const task = allTasks[i];

      // Check if task has the target tag
      const taskTags = task.tags;
      let hasTag = false;
      for (let j = 0; j < taskTags.length; j++) {
        if (taskTags[j].id() === tagID) {
          hasTag = true;
          break;
        }
      }

      if (!hasTag) continue;

      // Extract all tag names
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
