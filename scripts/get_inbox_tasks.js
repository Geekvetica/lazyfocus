(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;
    const inboxTasks = doc.inboxTasks;

    const tasks = [];

    for (let i = 0; i < inboxTasks.length; i++) {
      const task = inboxTasks[i];

      // Extract tag names from task tags
      const taskTags = task.tags;
      const tags = [];
      for (let j = 0; j < taskTags.length; j++) {
        tags.push(taskTags[j].name());
      }

      // Convert dates to ISO 8601 format or null
      const dueDate = task.dueDate();
      const deferDate = task.deferDate();
      const completedDate = task.completionDate();

      tasks.push({
        id: task.id(),
        name: task.name(),
        note: task.note() || "",
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
