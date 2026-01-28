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

    const projects = [];

    for (let i = 0; i < allProjects.length; i++) {
      const project = allProjects[i];
      const status = project.status();

      // Filter to only active projects
      if (status === "active") {
        projects.push({
          id: project.id(),
          name: project.name(),
          status: status,
          note: project.note() || ""
        });
      }
    }

    return JSON.stringify({ projects: projects }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
