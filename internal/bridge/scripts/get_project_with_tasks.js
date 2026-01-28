(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;
    const projectID = "{{.ProjectID}}";

    // Find the project by ID
    const allProjects = doc.flattenedProjects;
    let targetProject = null;

    for (let i = 0; i < allProjects.length; i++) {
      if (allProjects[i].id() === projectID) {
        targetProject = allProjects[i];
        break;
      }
    }

    if (!targetProject) {
      return JSON.stringify({ error: "Project not found" });
    }

    // Determine project status
    let projectStatus = "active";
    if (targetProject.completed()) {
      projectStatus = "completed";
    } else if (targetProject.dropped()) {
      projectStatus = "dropped";
    } else if (targetProject.status() === "on hold") {
      projectStatus = "on-hold";
    }

    // Get all tasks in the project
    const projectTasks = targetProject.flattenedTasks;
    const tasks = [];

    for (let i = 0; i < projectTasks.length; i++) {
      const task = projectTasks[i];

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
        projectID: targetProject.id(),
        projectName: targetProject.name(),
        tags: tags,
        dueDate: dueDate ? dueDate.toISOString() : null,
        deferDate: deferDate ? deferDate.toISOString() : null,
        flagged: task.flagged(),
        completed: task.completed(),
        completedDate: completedDate ? completedDate.toISOString() : null
      });
    }

    const project = {
      id: targetProject.id(),
      name: targetProject.name(),
      status: projectStatus,
      note: targetProject.note() || "",
      tasks: tasks
    };

    return JSON.stringify({ project: project }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
