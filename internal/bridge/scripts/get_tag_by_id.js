(() => {
  try {
    const app = Application("OmniFocus");
    app.includeStandardAdditions = true;

    // Check if OmniFocus is running
    if (!app.running()) {
      return JSON.stringify({ error: "OmniFocus is not running" });
    }

    const doc = app.defaultDocument;
    const tagID = "{{.TagID}}";

    // Find the tag by ID by traversing the hierarchy
    // This also allows us to track parent IDs
    let targetTag = null;
    let targetParentID = "";

    function findTag(tag, parentID) {
      if (tag.id() === tagID) {
        targetTag = tag;
        targetParentID = parentID;
        return true;
      }

      const childTags = tag.tags;
      for (let i = 0; i < childTags.length; i++) {
        if (findTag(childTags[i], tag.id())) {
          return true;
        }
      }
      return false;
    }

    const topLevelTags = doc.tags;
    for (let i = 0; i < topLevelTags.length; i++) {
      if (findTag(topLevelTags[i], "")) {
        break;
      }
    }

    if (!targetTag) {
      return JSON.stringify({ error: "Tag not found" });
    }

    // Get child tags
    const childTags = targetTag.tags;
    const children = [];

    for (let i = 0; i < childTags.length; i++) {
      children.push({
        id: childTags[i].id(),
        name: childTags[i].name()
      });
    }

    const tag = {
      id: targetTag.id(),
      name: targetTag.name(),
      parentID: targetParentID
    };

    if (children.length > 0) {
      tag.children = children;
    }

    return JSON.stringify({ tag: tag }, null, 2);

  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
})();
