# USE CASE: Ignore Files After Initial Scaffolding

Some files should be generated once and then left under full developer control. SKA supports this via ignore lists that
are automatically set during the first scaffolding and honored on every later update.

This guide shows how to:

- Predefine which files are ignored after the first run.
- Use templates stored in subfolders of a remote repository.
- Update/pin the template version later without touching ignored files.
- Mirror the flow you can see in the demo cast: create → edit locally → update.

## When to use this

- Project-specific files that differ per team (e.g., local.env, README badges, example configs).
- Binary artifacts or lockfile you want to manage manually.
- Any content you don’t want future template updates to override.

## How it works

Your upstream template includes a .ska-upstream.yaml with a section that pre-populates the project’s local SKA config:

- `ignorePaths` in `.ska-upstream.yaml` → Controls which paths are ignored while reading the upstream template itself (
  e.g.,
  .git, .idea).
- `skaConfig.ignorePaths` in `.ska-upstream.yaml` → Populates the generated project’s ignore list so those files are NOT
  updated on future ska update runs.

Example upstream configuration (YAML):

```yaml
# Files/folders ignored when reading the upstream template repository
ignorePaths:
- .git
- .idea

# Files/folders that SKA should ignore in your generated project on future updates
skaConfig:
  ignorePaths:
  - "docs/changelog.md"
  - "local/*.env"
  - "configs/example-{{ .appName }}.yaml"
  - "*.local"
```

Notes:

- Patterns support globs (e.g., *.env, folder/*, **/path/**) and can be templated with variables collected at creation
  time.
- On the first creation, SKA writes these into the project’s internal SKA configuration, so future updates skip them.

## Create from a remote subfolder template

If your organization keeps multiple templates in one repository, target the specific subfolder and a version (
tag/branch):

```sh
ska create \
--blueprint github.com/your-org/templates-repo/templates/service-rest@v1.2.0 \
--output ./my-rest-service
```

What happens:

- SKA fetches and renders the template from templates/service-rest.
- It collects inputs in TUI (unless you pass --non-interactive).
- It seeds the project’s SKA configuration with the skaConfig.ignorePaths specified by the upstream.
- From now on, those paths are excluded from updates.

## After the initial scaffolding

- You can freely edit files that match the ignore list; SKA will not overwrite them.
- You can still update other files as the template evolves.

To update the project:

```sh
ska update --path .
```

If you manage multiple SKA configurations in the same root (e.g., multiple templates):

```sh
ska update --path . --name service-rest
```

## Pin or switch to a new template version

Ignored files remain untouched when you update or switch versions.

1. Change the upstream ref from, say, @v1.2.0 to @v1.3.0 (same subfolder path):
    ```
    github.com/your-org/templates-repo/templates/service-rest@v1.3.0
    ```
   Update the upstream reference in your project’s SKA config.

2. Apply the update:
   ```sh
   ska update --path .
   ```

Result:

- SKA applies upstream changes to non-ignored files.
- Files listed in the ignore list (e.g., docs/changelog.md, local/*.env) remain exactly as you customized them.

## Practical example patterns

Add these to `skaConfig.ignorePaths` in your upstream to control behavior in generated projects:

- Ignore a single file:
    - README.md
- Ignore a folder:
    - docs/*
- Ignore environment files:
    - .env
    - local/*.env
- Ignore generated example config tied to an input variable:
    - configs/example-{{ .appName }}.yaml
- Ignore editor or machine-local files:
    - .vscode/*
    - *.local

> [!TIP]
> - Keep the ignore list focused on truly project-owned files to avoid missing important upstream improvements.
> - Prefer tags (e.g., @v1.3.0) for reproducible scaffolding; use branches (@main) for fast iteration.
> - If a file was ignored but you later want SKA to manage it, remove it from the project’s ignore list and run ska
> update.
