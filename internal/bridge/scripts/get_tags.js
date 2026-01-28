(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;
    const flat = "{{.Flat}}" === "true";

    if (flat) {
      // Return flat list of all tags
      // Note: parentID is not available in OmniFocus JXA API for tags
      // We build parent relationships by traversing the hierarchy
      const topLevelTags = doc.tags;
      const tags = [];

      function collectTags(tag, parentID) {
        tags.push({
          id: tag.id(),
          name: tag.name(),
          parentID: parentID
        });

        const childTags = tag.tags;
        for (let j = 0; j < childTags.length; j++) {
          collectTags(childTags[j], tag.id());
        }
      }

      for (let i = 0; i < topLevelTags.length; i++) {
        collectTags(topLevelTags[i], "");
      }

      return JSON.stringify({ tags: tags }, null, 2);
    } else {
      // Return hierarchical structure
      const topLevelTags = doc.tags;
      const tags = [];

      for (let i = 0; i < topLevelTags.length; i++) {
        tags.push(buildTagTree(topLevelTags[i]));
      }

      return JSON.stringify({ tags: tags }, null, 2);
    }

    // Helper function to build tag tree recursively
    function buildTagTree(tag) {
      const childTags = tag.tags;
      const children = [];

      for (let j = 0; j < childTags.length; j++) {
        children.push(buildTagTree(childTags[j]));
      }

      const result = {
        id: tag.id(),
        name: tag.name(),
        parentID: ""
      };

      if (children.length > 0) {
        result.children = children;
      }

      return result;
    }

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
