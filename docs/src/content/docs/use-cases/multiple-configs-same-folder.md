---
title: Multiple Configs in Same Folder
description: Manage multiple SKA template configurations in a single project directory.
---

It's common to scaffold more than one SKA template into the same working directory (for example, app code + CI/pipeline template + infra snippets). In this scenario, SKA stores multiple configuration files under the local `.ska-config` directory—one per template—so each template can be updated independently.

## How Multiple Configurations Work

Each time you scaffold a template into the same root, SKA creates or updates an entry under `.ska-config/` dedicated to that template.

Every entry tracks:
- The upstream blueprint reference (including subfolder path and pinned ref like `@v1.2.3`)
- The variables captured during creation
- Behavior settings (e.g., ignore paths) carried over from the upstream

This enables you to:
- **Keep multiple templates side-by-side** in one project folder
- **Update each template independently** without affecting the others

## List Configurations

Use the `config list` command to see all template configurations present in the current directory:

```bash
ska config list
```

### What You Get

- The list of configuration names that identify each template's state within `.ska-config`
- A reference for which name to target in later commands (e.g., update, rename)

:::tip
If you've applied multiple templates, use this command to confirm the exact names before running updates.
:::

## Rename a Configuration

Give a configuration a clearer name (for example, rename a generic "default" to something meaningful like `ci-pipeline`):

```bash
ska config rename OLD_NAME NEW_NAME
```

### Why Rename?

- Make updates more explicit and safer, especially when multiple templates coexist
- Standardize naming across teams or CI pipelines

### Example

```bash
ska config rename default ci-pipeline
```

## Update a Specific Template

When multiple templates are present, specify which one to update via the named configuration:

```bash
ska update --path . --name <CONFIG_NAME>
```

### Examples

Update a CI/pipeline template only:

```bash
ska update --path . --name ci-pipeline
```

Update an application runtime template only:

```bash
ska update --path . --name app-runtime
```

### Non-Interactive Updates in CI

```bash
ska update --path . --name ci-pipeline --non-interactive -v key=value -v another=value
```

## Practical Workflow

### Step 1: Scaffold Multiple Templates

Apply multiple templates to the same folder:

```bash
# App template
ska create \
  --blueprint https://github.com/org/templates//app-base@v1.0 \
  --output ./my-project

# CI template
ska create \
  --blueprint https://github.com/org/templates//ci-pipeline@v2.0 \
  --output ./my-project
```

### Step 2: List Configurations

```bash
ska config list
```

Output:
```
Configurations in ./my-project:
  - default
  - ci-pipeline
```

### Step 3: Rename for Clarity

```bash
ska config rename default app-runtime
```

### Step 4: Update Independently

```bash
# Update just the CI template
ska update --path . --name ci-pipeline

# Update just the app template
ska update --path . --name app-runtime
```

## Best Practices

:::tip[Multi-Config Tips]
- **Choose descriptive names** (`app-runtime`, `ci-pipeline`, `infra-shared`) to make command intent obvious
- **Keep configuration names stable** so CI and scripts remain consistent
- **Clean working tree** before updating to review changes easily
- **Update one at a time** to clearly see which template caused which changes
:::

## Updating Upstream References

If a configuration's upstream reference needs to be pinned to a new version (e.g., `@v1.2.0` → `@v1.3.0`):

1. Update that configuration's upstream reference
2. Run the targeted update:

```bash
ska update --path . --name <CONFIG_NAME>
```
