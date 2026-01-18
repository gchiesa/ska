---
title: Upstream Blueprints
description: Understanding how SKA blueprint templates are structured and referenced.
---

An **upstream blueprint** is a centrally maintained template that anyone can use to scaffold their own project structure. Blueprints can be hosted locally or remotely.

## Blueprint URL Formats

SKA supports multiple URI schemes for referencing blueprints:

### Local Blueprints

Use the `file://` scheme for templates on your local filesystem:

```bash
ska create --blueprint file:///Users/gchiesa/git/ska-example-template --output ./my-project
```

### Remote Blueprints (GitHub/GitLab)

Use HTTPS URLs for Git-hosted blueprints:

```bash
ska create --blueprint https://github.com/gchiesa/ska-golang-cli-template@master --output ./my-project
```

### Pinning Versions

Append `@ref` to pin a specific branch or tag:

| Reference Type | Example |
|----------------|---------|
| Tag | `https://github.com/org/repo@v1.2.3` |
| Branch | `https://github.com/org/repo@main` |
| Default branch | `https://github.com/org/repo` (no `@`) |

### Subfolder Templates

If a repository contains multiple templates in subfolders, reference them with a double slash:

```
https://github.com/org/repo//path/to/template@v1.0
```

For example:
```bash
ska create \
  --blueprint https://github.com/your-org/templates-repo//templates/service-rest@v1.2.0 \
  --output ./my-rest-service
```

## Blueprint Structure

A typical blueprint contains:

```
my-blueprint/
├── .ska-upstream.yaml    # SKA configuration (optional but recommended)
├── README.md             # Can contain template variables
├── src/
│   └── {{.appName}}/     # Folders can use template syntax
│       └── main.go       # Files can contain Go template syntax
└── ...
```

## The `.ska-upstream.yaml` File

This configuration file tells SKA how to process the blueprint:

```yaml
# Paths to ignore when reading the blueprint
ignorePaths:
  - .git
  - .idea

# Configuration to seed in generated projects
skaConfig:
  ignorePaths:
    - "*.local"
    - "docs/changelog.md"

# Input variable definitions for the Terminal UI form
inputs:
  - placeholder: appName
    label: Application Name
    help: The name of your application (lowercase, no spaces)
    regexp: "^[a-z0-9-]*$"
    default: myapp

  - placeholder: author
    label: Author Name
    help: Your name or organization
```

### Input Field Properties

| Property | Description |
|----------|-------------|
| `placeholder` | Variable name used in templates |
| `label` | Display label in the form |
| `help` | Inline help text shown to users |
| `regexp` | Validation pattern for accepted characters |
| `default` | Pre-filled default value |

## Project Configuration

After scaffolding, SKA creates a `.ska-config/` folder in your project containing:

- The upstream blueprint reference
- Captured variable values
- Ignore paths and other settings

This configuration enables the `ska update` command to work correctly.

## Best Practices

:::tip[Blueprint Design Tips]
- **Use semantic versioning** for tags to communicate breaking changes
- **Keep inputs minimal** - only ask for what's truly needed
- **Provide sensible defaults** to speed up the scaffolding process
- **Document your blueprint** with a clear README
- **Test updates** by scaffolding, making changes, then updating
:::
