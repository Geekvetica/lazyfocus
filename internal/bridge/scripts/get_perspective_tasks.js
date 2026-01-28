(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;

    // Template parameter (filled by Go)
    const perspectiveName = "{{.PerspectiveName}}";

    if (!perspectiveName || perspectiveName === "") {
      return JSON.stringify({ error: "Perspective name is required" });
    }

    const tasks = [];

    // Handle built-in perspectives with direct data access
    const normalizedName = perspectiveName.toLowerCase();

    switch (normalizedName) {
      case "inbox":
        collectInboxTasks(doc, tasks);
        break;

      case "flagged":
        collectFlaggedTasks(doc, tasks);
        break;

      case "forecast":
      case "due":
      case "due soon":
        collectDueTasks(doc, tasks);
        break;

      case "projects":
        collectProjectTasks(doc, tasks);
        break;

      default:
        // For custom perspectives, we cannot reliably access them via Automation API
        // Custom perspectives require OmniFocus Pro and GUI scripting which is unreliable
        return JSON.stringify({
          error: "Custom perspectives are not supported. Use built-in perspectives: inbox, flagged, forecast, projects"
        });
    }

    return JSON.stringify({ tasks: tasks }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }

  // Helper function to collect inbox tasks
  function collectInboxTasks(doc, tasks) {
    const inboxTasks = doc.inboxTasks;
    for (let i = 0; i < inboxTasks.length; i++) {
      const task = inboxTasks[i];
      if (!task.completed()) {
        addTaskToArray(task, tasks);
      }
    }
  }

  // Helper function to collect flagged tasks
  function collectFlaggedTasks(doc, tasks) {
    const allTasks = doc.flattenedTasks;
    for (let i = 0; i < allTasks.length; i++) {
      const task = allTasks[i];
      if (task.flagged() && !task.completed()) {
        addTaskToArray(task, tasks);
      }
    }
  }

  // Helper function to collect due tasks (next 7 days)
  function collectDueTasks(doc, tasks) {
    const allTasks = doc.flattenedTasks;
    const now = new Date();
    const sevenDaysFromNow = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);

    for (let i = 0; i < allTasks.length; i++) {
      const task = allTasks[i];
      if (task.completed()) continue;

      const dueDate = task.dueDate();
      if (dueDate && dueDate <= sevenDaysFromNow) {
        addTaskToArray(task, tasks);
      }
    }
  }

  // Helper function to collect all project tasks (non-inbox, non-completed)
  function collectProjectTasks(doc, tasks) {
    const allTasks = doc.flattenedTasks;
    for (let i = 0; i < allTasks.length; i++) {
      const task = allTasks[i];
      if (!task.completed() && task.containingProject()) {
        addTaskToArray(task, tasks);
      }
    }
  }

  // Helper function to format and add a task to the array
  function addTaskToArray(task, tasks) {
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
})();
