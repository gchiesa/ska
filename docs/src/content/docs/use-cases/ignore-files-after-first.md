---
title: Ignore Files After First Run
description: Generate files once and exclude them from future template updates.
---

Some files should be generated once and then left under full developer control. SKA supports this via ignore lists that are automatically set during the first scaffolding and honored on every later update.

## When to Use This

- **Project-specific files** that differ per team (e.g., `local.env`, README badges, example configs)
- **Binary artifacts or lockfiles** you want to manage manually
- **Any content** you don't want future template updates to override

## How It Works

Your upstream template includes a `.ska-upstream.yaml` with configuration that pre-populates the project's local SKA config:

| Section | Purpose |
|---------|---------|
| `ignorePaths` | Paths ignored when reading the upstream template itself (e.g., `.git`, `.idea`) |
| `skaConfig.ignorePaths` | Paths that will NOT be updated on future `ska update` runs |

### Example Configuration

```yaml
# .ska-upstream.yaml in your blueprint

# Files/folders ignored when reading the upstream template repository
ignorePaths:
  - .git
  - .idea

# Files/folders that SKA should ignore in generated projects on future updates
skaConfig:
  ignorePaths:
    - "docs/changelog.md"
    - "local/*.env"
    - "configs/example-{{.appName}}.yaml"
    - "*.local"
```

### Pattern Support

- **Globs**: `*.env`, `folder/*`, `**/path/**`
- **Templated paths**: `configs/example-{{.appName}}.yaml`

On the first creation, SKA writes these into the project's internal configuration, so future updates skip them.

## Create from a Template

```bash
ska create \
  --blueprint https://github.com/your-org/templates-repo//templates/service-rest@v1.2.0 \
  --output ./my-rest-service
```

### What Happens

1. SKA fetches and renders the template
2. It collects inputs in TUI (unless you pass `--non-interactive`)
3. It seeds the project's SKA configuration with the `skaConfig.ignorePaths` specified by the upstream
4. From now on, those paths are excluded from updates

## After the Initial Scaffolding

- **Ignored files**: Edit freely; SKA will not overwrite them
- **Other files**: Still updated as the template evolves

To update the project:

```bash
ska update --path .
```

If you manage multiple SKA configurations in the same root:

```bash
ska update --path . --name service-rest
```

## Pin or Switch Template Versions

Ignored files remain untouched when you update or switch versions:

1. Change the upstream ref (e.g., `@v1.2.0` → `@v1.3.0`)
2. Apply the update:

```bash
ska update --path .
```

**Result**:
- SKA applies upstream changes to non-ignored files
- Files in the ignore list remain exactly as you customized them

## Common Ignore Patterns

Add these to `skaConfig.ignorePaths` in your upstream to control behavior in generated projects:

| Pattern | What It Ignores |
|---------|-----------------|
| `README.md` | Single file |
| `docs/*` | Entire folder |
| `.env` | Environment file |
| `local/*.env` | Environment files in local folder |
| `configs/example-{{.appName}}.yaml` | Templated filename |
| `.vscode/*` | Editor settings |
| `*.local` | Files with .local extension |

## Re-enabling Management

If a file was ignored but you later want SKA to manage it:

1. Remove it from the project's ignore list (in `.ska-config/`)
2. Run `ska update`

## Best Practices

:::tip[Ignore File Tips]
- **Keep the ignore list focused** on truly project-owned files to avoid missing important upstream improvements
- **Use globs** for flexibility (`local/*.env` vs listing each file)
- **Template dynamic names** when file names include variable values
- **Document ignored files** so teams know which files are "theirs"
:::
