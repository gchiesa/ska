# USE CASE: Manage Multiple SKA Configurations in the Same Folder

It’s common to scaffold more than one SKA template into the same working directory (for example, app code + CI/pipeline
template + infra snippets). In this scenario, SKA stores multiple configuration files under the local `.ska-config`
directory—one per template—so each template can be updated independently.

This document shows how to:

- Understand how multiple configurations are organized.
- List all SKA configurations present in the folder.
- Rename a specific configuration for clarity.
- Update a specific template by referencing its named configuration.

## How multiple configurations work

- Each time you scaffold a template into the same root, SKA creates or updates an entry under `.ska-config/` dedicated to
  that template.
- Every entry track:
    - The upstream blueprint reference (including subfolder path and pinned ref like @v1.2.3).
    - The variables which are captured during creation.
    - Any behavior settings (e.g., ignore paths) that are carried over from the upstream.

As a result, you can:

- Keep multiple templates side-by-side in one project folder.
- Update each template independently without affecting the others.

## List configurations

Use the config list command to see all template configurations present in the current directory:

```sh
ska config list
```

What you get:

- The list of configuration names that identify each template’s state within `.ska-config`.
- A handy reference for which name to target in later commands (e.g., update, rename).

Tip:

- If you’ve applied multiple templates, use this command to confirm the exact names before running updates.

## Rename a configuration

Give a configuration a clearer name (for example, rename a generic “default” to something meaningful like ci-pipeline):

```sh
ska config rename OLD_NAME NEW_NAME
```

Why rename?

- To make updates more explicit and safer, especially when multiple templates coexist in the same folder.
- To standardize naming across teams or CI pipelines.

## Update a specific template by name

When multiple templates are present, specify which one to update via the named configuration:

```sh
ska update --path . --name <CONFIG_NAME>
```

Examples:

- Update a CI/pipeline template only:
  ```sh
  ska update --path . --name ci-pipeline
  ```
- Update an application runtime template only:
  ```sh
  ska update --path . --name app-runtime
  ```

Notes:

- Use the long option --name to be explicit.
- You can pass variables during the update as needed:
  ```sh
  ska update --path . --name ci-pipeline 
  ```
- Run non-interactively in CI:
  ```sh
  ska update --path . --name ci-pipeline --non-interactive -v key=value -v another=value
  ```

## Practical workflow

1. Scaffold multiple templates into the same folder (e.g., an app template and a CI template).
2. Run:
   ```sh
   ska config list
   ```
   to see all configurations and their names.
3. Optionally, rename a configuration for clarity:
   ```sh
   ska config rename default ci-pipeline
   ```
4. Update a specific template by name whenever the upstream changes:
   ```sh
   ska update --path . --name ci-pipeline
   ```

> [!TIP]
> **Tips and good practices**
> - Choose descriptive names (app-runtime, ci-pipeline, infra-shared) to make command intent obvious.
> - Keep configuration names stable so CI and scripts remain consistent.
> - Before updating, ensure your working tree is clean to review changes easily.
> - If a configuration’s upstream reference needs to be pinned to a new version (e.g., @v1.2.0 → @v1.3.0), update that
>   configuration’s upstream reference and then run:
>   ```sh
>   ska update --path . --name <CONFIG_NAME>
>   ```
