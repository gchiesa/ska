---
title: Partial Sections
description: Keep specific file sections centrally managed while allowing customization elsewhere.
---

SKA's partial sections feature lets you maintain only specific portions of files while allowing developers to freely customize the rest.

## How It Works

Use special markers in your template files to define centrally managed sections:

- `ska-start:<identifier>` - Marks the beginning of a managed section
- `ska-end` - Marks the end of a managed section

Everything **inside** these markers is controlled by the upstream template. Everything **outside** can be freely edited by developers.

## Example: YAML Configuration

```yaml
# codecov.yml - Only key4 section is managed by SKA
---
root:
  key1: value1
  key2: custom-value  # User can modify this
  key3: value3
  # ska-start:key4
  key4:
    subkey: "{{.appName}}"
    subkey2: standardized-value
  # ska-end
  key5: value5  # User can modify this too
```

When `ska update` runs:
- The `key4` section is refreshed from the upstream template
- All other keys retain their local values

## Marker Syntax

The markers work as comments in any file type:

### YAML
```yaml
# ska-start:section-name
managed: content
# ska-end
```

### JavaScript/Go/Java
```javascript
// ska-start:config
const config = { /* managed */ };
// ska-end
```

### HTML/XML
```html
<!-- ska-start:header -->
<header>Managed content</header>
<!-- ska-end -->
```

### Shell/Python
```bash
# ska-start:env-setup
export PATH="$PATH:/managed/path"
# ska-end
```

## Identifier Names

The identifier after `ska-start:` should be:

- **Descriptive**: `ska-start:ci-step`, `ska-start:codecov-config`
- **Unique within the file**: Each section needs a distinct name
- **Stable**: Don't rename frequently as it affects update matching

## Use Cases

### CI/CD Pipelines

Keep build steps standardized while allowing custom stages:

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  # ska-start:standard-checks
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run linter
        run: make lint
  # ska-end

  # Teams can add custom jobs below
  custom-tests:
    runs-on: ubuntu-latest
    steps:
      - name: Custom test suite
        run: ./run-custom-tests.sh
```

### Configuration Files

Standardize security settings while allowing application-specific config:

```yaml
# app-config.yml
app:
  name: "{{.appName}}"

  # ska-start:security
  security:
    csrf: enabled
    xss-protection: enabled
    content-security-policy: strict
  # ska-end

  # Application-specific settings
  features:
    beta-mode: true
    custom-feature: enabled
```

### Docker Files

Maintain base image standards:

```dockerfile
# ska-start:base
FROM node:20-alpine
WORKDIR /app
RUN apk add --no-cache dumb-init
# ska-end

# Custom build steps
COPY package*.json ./
RUN npm ci --only=production

COPY . .
CMD ["dumb-init", "node", "server.js"]
```

## Best Practices

:::tip[Partial Section Tips]
- **Keep sections small and focused** - Easier to maintain and less merge friction
- **Use clear identifiers** - `ska-start:security-headers` is better than `ska-start:section1`
- **Document the contract** - Let developers know which sections are managed
- **Don't nest sections** - Keep the structure flat and simple
:::

## Behavior on Update

When you run `ska update`:

1. SKA identifies all `ska-start/ska-end` blocks in both template and local files
2. Managed sections are refreshed with the latest upstream content
3. Content outside managed sections is preserved exactly as-is
4. New template variables are applied within managed sections
