---
title: Multiple Template Subfolders
description: Use SKA with a repository containing multiple templates in subfolders.
---

This guide explains how to use SKA when a remote repository hosts multiple templates in subfolders (for example, one repo acting as a catalog of templates).

## When to Use This

- Your organization keeps many templates in a single repository, each under its own folder
- You want to scaffold from a specific subfolder path and optionally pin a branch or tag

## Blueprint URL Format

Reference a subfolder inside a remote Git repository by including the path after the repository name:

```
https://github.com/ORG/REPO//path/to/template@ref
https://gitlab.com/ORG/REPO//path/to/template@ref
```

:::note
- The `ref` after `@` can be a branch (e.g., `@main`), a tag (e.g., `@v1.2.3`), or omitted to use the default branch
- Use double slashes (`//`) to separate the repo from the subfolder path
:::

## Example Repository Layout

Suppose your template catalog looks like this:

```
templates-repo/
└── templates/
    ├── service-rest/
    ├── service-grpc/
    ├── library-utility/
    ├── common/
    └── fragments/
```

Each folder under `templates/` is a separate blueprint you can target.

## Create from a Subfolder Template

Scaffold a new project by pointing `--blueprint` to the subfolder:

```bash
ska create \
  --blueprint https://github.com/your-org/templates-repo//templates/service-rest@v1.2.0 \
  --output ./my-rest-service
```

### What Happens

1. SKA downloads only what it needs from the repo
2. It loads the template from the specified subfolder
3. If the template defines inputs, SKA shows the interactive form (unless you use `--non-interactive`)

## Pinning and Updating

| Strategy | Example | Use Case |
|----------|---------|----------|
| Pin to tag | `@v1.2.0` | Production, reproducibility |
| Pin to branch | `@main` or `@develop` | Fast iteration |
| No pin | (omit `@`) | Use default branch |

After creation, your project contains a `.ska-config/` folder that stores the upstream reference. To bring in updates from the same upstream path and ref:

```bash
cd ./my-rest-service
ska update --path .
```

## Switching Template Versions

To update to a new template version:

1. The upstream reference is stored in `.ska-config/`
2. Update the reference to the new version (e.g., `@v1.2.0` → `@v1.3.0`)
3. Run the update command:

```bash
ska update --path .
```

## Multiple Templates, Same Project

You can scaffold multiple templates into the same directory. Each creates its own entry in `.ska-config/`:

```bash
# First template
ska create \
  --blueprint https://github.com/org/templates//app-base@v1.0 \
  --output ./my-project

# Second template (same output directory)
ska create \
  --blueprint https://github.com/org/templates//ci-pipeline@v2.0 \
  --output ./my-project
```

Manage them independently:

```bash
# List all configurations
ska config list

# Update specific template
ska update --path . --name ci-pipeline
```

## Best Practices

:::tip[Working with Template Catalogs]
- **Use consistent folder naming** across your template catalog
- **Pin versions for production** use; use branches during development
- **Document each template** with a README in its subfolder
- **Use semantic versioning** for clear upgrade paths
:::
