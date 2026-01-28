(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;
    const allProjects = doc.flattenedProjects;

    // Template parameter for status filter: "active", "on-hold", "completed", "dropped", "all"
    const statusFilter = "{{.Status}}";

    const projects = [];

    for (let i = 0; i < allProjects.length; i++) {
      const project = allProjects[i];

      // Determine project status
      let projectStatus = "active";
      if (project.completed()) {
        projectStatus = "completed";
      } else if (project.dropped()) {
        projectStatus = "dropped";
      } else if (project.status() === "on hold") {
        projectStatus = "on-hold";
      }

      // Apply status filter
      if (statusFilter !== "all" && statusFilter !== "" && statusFilter !== projectStatus) {
        continue;
      }

      projects.push({
        id: project.id(),
        name: project.name(),
        status: projectStatus,
        note: project.note() || ""
      });
    }

    return JSON.stringify({ projects: projects }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
