---
title: Terminal UI
description: How SKA's interactive Terminal UI collects template variables.
---

SKA includes a dynamic Terminal UI that collects variable values when scaffolding or updating projects.

## How It Works

When a blueprint includes a `.ska-upstream.yaml` file with `inputs` defined, SKA automatically generates an interactive form in your terminal.

## Input Configuration

Define inputs in your blueprint's `.ska-upstream.yaml`:

```yaml
inputs:
  - placeholder: appName
    label: Application Name
    help: The name of your application (lowercase, no spaces)
    regexp: "^[a-z0-9-]*$"
    default: myapp

  - placeholder: author
    label: Author
    help: Your name or organization name

  - placeholder: port
    label: HTTP Port
    help: The port your service will listen on
    regexp: "^[0-9]+$"
    default: "8080"
```

## Input Properties

| Property | Required | Description |
|----------|----------|-------------|
| `placeholder` | Yes | Variable name used in templates (e.g., `{{.appName}}`) |
| `label` | Yes | Display label shown in the form |
| `help` | No | Inline help text for the field |
| `regexp` | No | Validation pattern for accepted characters |
| `default` | No | Pre-filled default value |

## Validation

The `regexp` property validates input as the user types:

```yaml
inputs:
  # Only lowercase letters, numbers, and hyphens
  - placeholder: projectSlug
    label: Project Slug
    regexp: "^[a-z0-9-]*$"

  # Semantic version format
  - placeholder: version
    label: Initial Version
    regexp: "^[0-9]+\\.[0-9]+\\.[0-9]+$"
    default: "1.0.0"

  # Email format
  - placeholder: email
    label: Contact Email
    regexp: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

:::note
The `regexp` validates individual characters, not the entire pattern. It determines which characters are accepted in the input field.
:::

## Non-Interactive Mode

For CI/CD pipelines or automation, skip the Terminal UI with `--non-interactive`:

```bash
ska create \
  --blueprint https://github.com/org/template@v1.0 \
  --output ./my-project \
  --non-interactive \
  -v appName=myservice \
  -v author="Platform Team" \
  -v port=3000
```

### Providing Variables

Use the `-v` flag (multiple times) to pass variables:

```bash
ska create \
  --blueprint ./my-template \
  --output ./project \
  -n \
  -v key1=value1 \
  -v key2=value2 \
  -v key3="value with spaces"
```

## Update Behavior

When running `ska update`, the Terminal UI:

1. Pre-fills fields with previously captured values
2. Allows you to modify any value
3. Applies the updated values to managed sections

```bash
# Interactive update - modify values as needed
ska update --path .

# Non-interactive update - keep existing values
ska update --path . --non-interactive

# Non-interactive with value overrides
ska update --path . -n -v port=9090
```

## Tips for Blueprint Authors

:::tip[Terminal UI Best Practices]
- **Order inputs logically** - Put the most important fields first
- **Provide sensible defaults** - Reduce friction for common cases
- **Write helpful help text** - Explain format expectations clearly
- **Use validation** - Prevent errors early with regexp patterns
- **Keep it minimal** - Only ask for what's truly needed
:::
