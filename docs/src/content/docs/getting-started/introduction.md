---
title: Introduction
description: Learn what SKA is and how it can help you manage project scaffolding at scale.
---

**SKA** (Your Scaffolding Buddy) is a powerful templating tool that allows you to expand folders based on local or remote blueprint templates hosted on GitHub, GitLab, or your local filesystem.

## Why SKA?

Traditional scaffolding tools offer one-shot code generation. After the initial project is created, the code becomes autonomous with no way to centrally roll out new changes.

**SKA is different.** It maintains a connection between your scaffolded project and the upstream template, allowing you to:

- **Update projects** when templates evolve
- **Maintain consistency** across your organization
- **Keep sections synced** while allowing customization elsewhere

## Key Features

### Dynamic Forms for Data Capture

Each blueprint supports a `.ska-upstream.yaml` file that defines template variables. SKA uses this to dynamically create an interactive Terminal UI form to collect input data with:

- Input validation via regex patterns
- Default values
- Inline help text
- Label customization

### Central Management for Updates

Unlike traditional tools, SKA allows you to update scaffolded projects when the upstream template changes. Simply run:

```bash
ska update --path .
```

Your project will be updated with the latest template changes while preserving your customizations.

### Partial File Management

SKA supports managing only specific sections of files. Using `ska-start` and `ska-end` markers, you can:

- Keep certain sections centrally managed
- Allow developers to customize everything else
- Update managed sections without affecting local changes

## Who Is SKA For?

SKA is designed for **Platform Engineers** and teams who need to:

- Maintain multiple project templates
- Ensure consistency across projects
- Roll out template improvements over time
- Balance standardization with flexibility

## Next Steps

Ready to get started? Continue to the [Installation](/ska/getting-started/installation/) guide.
