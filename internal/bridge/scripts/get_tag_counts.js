(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;
    const allTags = doc.flattenedTags;
    const counts = {};

    for (let i = 0; i < allTags.length; i++) {
      const tag = allTags[i];
      const tagID = tag.id();

      // Get tasks with this tag
      const tagTasks = tag.tasks;
      let incompleteCount = 0;

      // Count only incomplete tasks
      for (let j = 0; j < tagTasks.length; j++) {
        if (!tagTasks[j].completed()) {
          incompleteCount++;
        }
      }

      counts[tagID] = incompleteCount;
    }

    return JSON.stringify({ counts: counts }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
