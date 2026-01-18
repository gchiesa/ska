---
title: Partial File Management
description: Keep only specific sections of files centrally managed while allowing developer customization.
---

This guide shows how to:

- Keep only a specific section of a file centrally managed by SKA while letting developers customize the rest
- Use a blueprint living in a subfolder of a remote repository (template catalog)
- After the initial scaffolding, pin/switch the template to a new upstream version and onboard changes safely

## Why Partial Updates?

Partial updates allow you to:

- **Enforce standards** in a small, named block (e.g., a CI step or a config stanza)
- **Preserve developer changes** outside that block across updates

## How It Works: ska-start / ska-end

In your blueprint template, wrap the centrally managed portion with markers:

- Start with a comment line containing `ska-start:<identifier>`
- End with a comment line containing `ska-end`

### Example (YAML)

Only the `key4` section is centrally managed; everything else remains user-editable:

```yaml
# config.yml - SKA only manages the key4 section
---
root:
  key1: value1
  key2: {{.notManaged}}
  key3: value3
  # ska-start:key4
  key4:
    subkey: "{{.appName}}"
    subkey2: value2
  # ska-end
  key5: value5
```

:::note
- The identifier (`key4` above) is a friendly name you choose to reference the block
- These markers can be used in many file types (YAML, JSON-with-comments, scripts, Dockerfiles, etc.) as comments
:::

## Create from a Template

If your organization keeps multiple templates in one repo, reference the subfolder in the blueprint URL:

```bash
ska create \
  --blueprint https://github.com/your-org/templates-repo//templates/service-rest@v1.2.0 \
  --output ./my-rest-service
```

### What Happens

1. SKA fetches the template from the specified subfolder
2. It renders files; only the sections inside `ska-start/ska-end` are centrally managed
3. The project gains a `.ska-config/` folder that records the upstream and your captured variables

## Editing After Scaffolding

- **Outside managed blocks**: Edit freely. Changes are preserved across updates.
- **Inside managed blocks**: Controlled by the upstream template. Refreshed on update.

:::tip
Keep your local customizations outside the managed blocks to avoid conflicts.
:::

## Updating and Pinning to a New Version

You can onboard upstream changes (including updates to the managed block) and pin to a new version:

### Step 1: Update the Upstream Reference

Change the upstream reference to a new tag or branch. For example, switch from `v1.2.0` to `v1.3.0`:

```
github.com/your-org/templates-repo//templates/service-rest@v1.3.0
```

You can either:
- Update your SKA configuration to the new ref, or
- Re-run create in a new location with the new ref to compare outputs

### Step 2: Run the Update

From your project root:

```bash
ska update --path .
```

### Step 3: Review Changes

- SKA updates only the content within managed blocks
- Your edits outside those blocks remain untouched

If you maintain multiple SKA configurations in the same project root, use the named configuration:

```bash
ska update --path . --name service-rest
```

## Practical Example Flow

**Day 0**: Create a project from a subfolder template at a known tag:

```bash
ska create \
  --blueprint https://github.com/your-org/templates-repo//templates/service-rest@v1.2.0 \
  --output ./my-rest-service
```

**Day 1**: Developers customize files outside SKA-managed blocks.

**Day 15**: Central team releases `v1.3.0` adjusting the managed block (e.g., a linter config).

1. Update the upstream ref to `@v1.3.0` in your project's SKA config
2. Run the update:

```bash
cd ./my-rest-service
ska update --path .
```

**Result**: The managed block content updates to v1.3.0; custom changes outside the block remain intact.

## Best Practices

:::tip[Partial Management Tips]
- **Prefer tags** (e.g., `@v1.3.0`) for reproducibility; use branches (`@main`) for fast iteration
- **Keep managed blocks small** and focused to minimize merge friction
- **Use clear identifiers** in ska-start lines (`ska-start:ci-step`, `ska-start:codecov`, etc.)
- **Clean working tree** before switching refs or updating to make diffs clear
:::
