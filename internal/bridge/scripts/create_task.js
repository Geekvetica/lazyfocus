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
    const taskName = "{{.Name}}";
    const taskNote = "{{.Note}}";
    const projectID = "{{.ProjectID}}";
    const tagsJSON = "{{.Tags}}";
    const dueDateStr = "{{.DueDate}}";
    const deferDateStr = "{{.DeferDate}}";
    const flaggedStr = "{{.Flagged}}";

    if (!taskName) {
      return JSON.stringify({ error: "Task name is required" });
    }

    // Create task properties object
    const taskProps = {
      name: taskName
    };

    if (taskNote) {
      taskProps.note = taskNote;
    }

    if (flaggedStr === "true") {
      taskProps.flagged = true;
    } else if (flaggedStr === "false") {
      taskProps.flagged = false;
    }

    // Parse and set due date
    if (dueDateStr) {
      const dueDate = new Date(dueDateStr);
      if (isNaN(dueDate.getTime())) {
        return JSON.stringify({ error: `Invalid due date format: ${dueDateStr}` });
      }
      taskProps.dueDate = dueDate;
    }

    // Parse and set defer date
    if (deferDateStr) {
      const deferDate = new Date(deferDateStr);
      if (isNaN(deferDate.getTime())) {
        return JSON.stringify({ error: `Invalid defer date format: ${deferDateStr}` });
      }
      taskProps.deferDate = deferDate;
    }

    // Create the task
    const newTask = app.Task(taskProps);

    // Add task to project or inbox
    if (projectID) {
      // Find project by ID
      const allProjects = doc.flattenedProjects;
      let targetProject = null;

      for (let i = 0; i < allProjects.length; i++) {
        if (allProjects[i].id() === projectID) {
          targetProject = allProjects[i];
          break;
        }
      }

      if (!targetProject) {
        return JSON.stringify({ error: `Project not found: ${projectID}` });
      }

      targetProject.tasks.push(newTask);
    } else {
      // Add to inbox
      doc.inboxTasks.push(newTask);
    }

    // Add tags if specified
    // Note: Due to JXA/OmniFocus limitations, only the first tag (primary tag) is supported
    // The tag must already exist in OmniFocus
    if (tagsJSON && tagsJSON !== "[]") {
      try {
        const tagNames = JSON.parse(tagsJSON);
        if (tagNames.length > 0) {
          const tagName = tagNames[0]; // Only use first tag as primary tag

          // Find existing tag by name
          const existingTag = doc.flattenedTags.whose({name: tagName});

          if (existingTag.length > 0) {
            newTask.primaryTag = existingTag[0];
          }
          // If tag doesn't exist, silently skip (don't create new tags via automation)
        }
      } catch (e) {
        return JSON.stringify({ error: `Invalid tags JSON: ${e.message}` });
      }
    }

    // Retrieve the created task to return full details
    const taskTags = newTask.tags;
    const tags = [];
    for (let j = 0; j < taskTags.length; j++) {
      tags.push(taskTags[j].name());
    }

    const containingProject = newTask.containingProject();
    const returnProjectID = containingProject ? containingProject.id() : "";
    const returnProjectName = containingProject ? containingProject.name() : "";

    const dueDate = newTask.dueDate();
    const deferDate = newTask.deferDate();

    const result = {
      id: newTask.id(),
      name: newTask.name(),
      note: newTask.note() || "",
      projectID: returnProjectID,
      projectName: returnProjectName,
      tags: tags,
      dueDate: dueDate ? dueDate.toISOString() : null,
      deferDate: deferDate ? deferDate.toISOString() : null,
      flagged: newTask.flagged(),
      completed: newTask.completed()
    };

    return JSON.stringify({ task: result }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
