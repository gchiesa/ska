---
title: Quick Start
description: Create your first project with SKA in minutes.
---

This guide walks you through creating and updating a project using SKA.

## Create a New Project

Use the `create` command with a blueprint URL to scaffold a new project:

```bash
ska create \
  --blueprint https://github.com/gchiesa/ska-golang-cli-template@master \
  --output ~/workspace/myNewApp
```

### What Happens

1. SKA downloads the blueprint template
2. If the template defines inputs, an interactive form appears
3. SKA renders the template with your provided values
4. A `.ska-config/` folder is created to track state and upstream reference

## Interactive Form

When a blueprint includes a `.ska-upstream.yaml` file, SKA displays a Terminal UI form to collect variables:

- Fill in each field as prompted
- Use the provided defaults when appropriate
- Validation ensures correct input format

To skip the interactive form (useful in CI/CD), use the `--non-interactive` flag with variables:

```bash
ska create \
  --blueprint https://github.com/org/template@v1.0 \
  --output ./my-project \
  --non-interactive \
  -v appName=myapp \
  -v author="John Doe"
```

## Update Your Project

When the upstream template changes, update your project:

```bash
cd ~/workspace/myNewApp
ska update --path .
```

SKA will:
- Fetch the latest template version
- Show the interactive form (pre-filled with current values)
- Apply changes while respecting partial management markers

## Blueprint URL Formats

SKA supports multiple URL formats for blueprints:

| Format | Example |
|--------|---------|
| GitHub/GitLab HTTPS | `https://github.com/org/repo@v1.0` |
| With subfolder | `https://github.com/org/repo//templates/service@v1.0` |
| Local filesystem | `file:///path/to/template` |

:::tip
Pin templates to specific tags (e.g., `@v1.2.3`) for reproducible scaffolding. Use branches like `@main` only during development.
:::

## Example Templates

Try these example templates to see SKA in action:

- **Go CLI Template**: `https://github.com/gchiesa/ska-golang-cli-template@master`

## Next Steps

Learn more about SKA's core concepts:

- [Upstream Blueprints](/ska/concepts/upstream-blueprints/) - How templates are structured
- [Template Language](/ska/concepts/template-language/) - Go Template syntax and Sprig functions
- [Partial Sections](/ska/concepts/partial-sections/) - Managing specific file sections
