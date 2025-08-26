# USE CASE: Manage Partial File Updates

This guide shows how to:

- Keep only a specific section of a file centrally managed by SKA while letting developers customize the rest.
- Use a blueprint living in a subfolder of a remote repository (template catalog).
- After the initial scaffolding, pin/switch the template to a new upstream version and onboard changes safely.

## Why partial updates?

Partial updates allow you to:

- Enforce standards in a small, named block (e.g., a CI step or a config stanza).
- Preserve developer changes outside that block across updates.

## How it works: ska-start / ska-end

In your blueprint template, wrap the centrally managed portion with markers:

- Start with a comment line containing `ska-start:<identifier>`
- End with a comment line containing `ska-end`

Example (YAML) — only the key4 section is centrally managed; everything else remains user-editable:

```yaml
# this is an example YAML and we want SKA to only manage the section `key4`
---
root:
  key1: value1
  key2: { { .notManaged } }
  key3: value3
  # ska-start:key4
  key4:
    subkey: "{{ .appName }}"
    subkey2: value2
  # ska-end
  key5: value5
```

> [!NOTE]
> - The identifier (key4 above) is a friendly name you choose to reference the block.
> - These markers can be used in many file types (YAML, JSON-with-comments, scripts, Dockerfiles, etc.) as comments.

## Create from a template in a remote repo subfolder

If your organization keeps multiple templates in one repo, reference the subfolder in the blueprint URL and optionally
pin a ref (branch or tag):

- HTTPS form:
  ```sh
  ska create \
    --blueprint https://github.com/your-org/templates-repo/templates/service-rest@v1.2.0 \
    --output ./my-rest-service
  ```

What happens:

- SKA fetches the template from the specified subfolder.
- It renders files; only the sections inside ska-start/ska-end are centrally managed.
- The project gains a .ska-config/ folder that records the upstream and your captured variables.

## Editing after scaffolding

- You can freely edit anything outside the managed blocks. Those edits are preserved across updates.
- The contents inside `ska-start/ska-end` are controlled by the upstream template and variables; they will be
  refreshed on update.

> [!TIP]
> Keep your local customizations outside the managed blocks to avoid conflicts.

## Updating and pinning to a new version

You can onboard upstream changes (including updates to the managed block) and pin to a new version:

1. Change the upstream reference to a new tag or branch. For example, switch from v1.2.0 to v1.3.0 while keeping the
   same template subfolder path:
   ```
   github.com/your-org/templates-repo/templates/service-rest@v1.3.0
   ```

   You can either:

- Update your SKA configuration to the new ref, or
- Re-run create in a new location with the new ref if you want to compare outputs, then move the change to your existing
  project.

2. From your project root, run:

    ```sh
    ska update --path .
    ```

3. Review changes:

- SKA updates only the content within the managed blocks.
- Your edits outside those blocks remain untouched.

If you maintain multiple SKA configurations in the same project root (e.g., multiple templates applied to different
subfolders), use a named configuration:

```
ska update --path . --name service-rest
```

## Practical example flow

- Day 0: Create a project from a subfolder template at a known tag:
    ```sh
    ska create \
    --blueprint github.com/your-org/templates-repo/templates/service-rest@v1.2.0 \
    --output ./my-rest-service
    ```
  
- Day 1: Developers customize files outside ska-managed blocks.
 
- Day 15: Central team releases v1.3.0 adjusting the managed block (e.g., a linter config).
    - Update the upstream ref to @v1.3.0 in your project’s SKA config.
    - Run:
      ```
      cd ./my-rest-service
      ska update --path .
      ```
    - Result: The managed block content updates to v1.3.0; custom changes outside the block remain intact.

> [!TIP]
> **Tips and good practices**
> - Prefer tags (e.g., @v1.3.0) for reproducibility; use branches (e.g., @main) for fast iteration.
> - Keep your managed blocks as small and focused as possible to minimize merge friction.
> - Use clear identifiers in ska-start lines (ska-start:ci-step, ska-start:codecov, etc.) for maintainability.
> - Before switching refs or subfolders, run an update and ensure the working tree is clean to make diffs clear.
