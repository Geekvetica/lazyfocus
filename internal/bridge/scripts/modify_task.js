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
    const newName = "{{.Name}}";
    const newNote = "{{.Note}}";
    const projectID = "{{.ProjectID}}";
    const addTagsJSON = "{{.AddTags}}";
    const removeTagsJSON = "{{.RemoveTags}}";
    const dueDateStr = "{{.DueDate}}";
    const deferDateStr = "{{.DeferDate}}";
    const flaggedStr = "{{.Flagged}}";

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

    // Update name if provided
    if (newName) {
      targetTask.name = newName;
    }

    // Update note if provided
    if (newNote) {
      targetTask.note = newNote;
    }

    // Update flagged status if provided
    if (flaggedStr === "true") {
      targetTask.flagged = true;
    } else if (flaggedStr === "false") {
      targetTask.flagged = false;
    }

    // Update project if provided
    if (projectID) {
      if (projectID === "CLEAR") {
        // Move task to inbox by removing from project
        targetTask.assignedContainer = doc;
      } else {
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

        targetTask.assignedContainer = targetProject;
      }
    }

    // Update due date if provided
    if (dueDateStr) {
      if (dueDateStr === "CLEAR") {
        targetTask.dueDate = null;
      } else {
        const dueDate = new Date(dueDateStr);
        if (isNaN(dueDate.getTime())) {
          return JSON.stringify({ error: `Invalid due date format: ${dueDateStr}` });
        }
        targetTask.dueDate = dueDate;
      }
    }

    // Update defer date if provided
    if (deferDateStr) {
      if (deferDateStr === "CLEAR") {
        targetTask.deferDate = null;
      } else {
        const deferDate = new Date(deferDateStr);
        if (isNaN(deferDate.getTime())) {
          return JSON.stringify({ error: `Invalid defer date format: ${deferDateStr}` });
        }
        targetTask.deferDate = deferDate;
      }
    }

    // Add tags if specified
    // Note: Due to JXA/OmniFocus limitations, we can only set the primary tag
    // The tag must already exist in OmniFocus
    if (addTagsJSON && addTagsJSON !== "[]") {
      try {
        const tagNames = JSON.parse(addTagsJSON);
        if (tagNames.length > 0) {
          const tagName = tagNames[0]; // Only use first tag as primary tag

          // Find existing tag by name
          const existingTag = doc.flattenedTags.whose({name: tagName});

          if (existingTag.length > 0) {
            targetTask.primaryTag = existingTag[0];
          }
          // If tag doesn't exist, silently skip
        }
      } catch (e) {
        return JSON.stringify({ error: `Invalid add tags JSON: ${e.message}` });
      }
    }

    // Remove tags if specified (clear primary tag)
    if (removeTagsJSON && removeTagsJSON !== "[]") {
      try {
        const tagNames = JSON.parse(removeTagsJSON);
        // If removing the primary tag, clear it
        const currentPrimary = targetTask.primaryTag();
        if (currentPrimary && tagNames.includes(currentPrimary.name())) {
          targetTask.primaryTag = null;
        }
      } catch (e) {
        return JSON.stringify({ error: `Invalid remove tags JSON: ${e.message}` });
      }
    }

    // Retrieve the updated task to return full details
    const taskTags = targetTask.tags;
    const tags = [];
    for (let j = 0; j < taskTags.length; j++) {
      tags.push(taskTags[j].name());
    }

    const containingProject = targetTask.containingProject();
    const returnProjectID = containingProject ? containingProject.id() : "";
    const returnProjectName = containingProject ? containingProject.name() : "";

    const dueDate = targetTask.dueDate();
    const deferDate = targetTask.deferDate();
    const completedDate = targetTask.completionDate();

    const result = {
      id: targetTask.id(),
      name: targetTask.name(),
      note: targetTask.note() || "",
      projectID: returnProjectID,
      projectName: returnProjectName,
      tags: tags,
      dueDate: dueDate ? dueDate.toISOString() : null,
      deferDate: deferDate ? deferDate.toISOString() : null,
      flagged: targetTask.flagged(),
      completed: targetTask.completed(),
      completedDate: completedDate ? completedDate.toISOString() : null
    };

    return JSON.stringify({ task: result }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
