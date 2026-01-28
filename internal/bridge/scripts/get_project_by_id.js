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

    const project = {
      id: targetProject.id(),
      name: targetProject.name(),
      status: projectStatus,
      note: targetProject.note() || ""
    };

    return JSON.stringify({ project: project }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
