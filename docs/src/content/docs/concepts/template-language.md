---
title: Template Language
description: Learn about Go Templates and Sprig functions used in SKA blueprints.
---

SKA uses the **Go Template framework** for templating your upstream blueprints. Additionally, it includes all **Sprig functions** for extended functionality.

## Go Template Basics

Go templates use double curly braces for template expressions:

```go
Hello, {{.name}}!
```

### Variables

Access variables from your `.ska-upstream.yaml` inputs:

```go
// Simple variable
{{.appName}}

// With formatting
Application: {{.appName | upper}}
```

### Conditionals

```go
{{if .enableTests}}
test:
  enabled: true
{{end}}

{{if eq .environment "production"}}
replicas: 3
{{else}}
replicas: 1
{{end}}
```

### Loops

```go
{{range .services}}
- name: {{.name}}
  port: {{.port}}
{{end}}
```

### Pipelines

Chain functions using the pipe operator:

```go
{{.appName | lower | replace " " "-"}}
```

## Sprig Functions

SKA includes 100+ functions from the [Sprig library](https://masterminds.github.io/sprig/). Here are the most useful ones:

### String Functions

| Function | Description | Example |
|----------|-------------|---------|
| `upper` | Uppercase | `{{.name \| upper}}` → `MYAPP` |
| `lower` | Lowercase | `{{.name \| lower}}` → `myapp` |
| `title` | Title case | `{{.name \| title}}` → `Myapp` |
| `replace` | Replace string | `{{.name \| replace "-" "_"}}` |
| `trim` | Trim whitespace | `{{.name \| trim}}` |
| `quote` | Add quotes | `{{.name \| quote}}` → `"myapp"` |
| `default` | Default value | `{{.port \| default "8080"}}` |

### List Functions

```go
// First item
{{first .items}}

// Last item
{{last .items}}

// Join with separator
{{.tags | join ", "}}
```

### Date Functions

```go
// Current date
{{now | date "2006-01-02"}}

// Format timestamp
{{.createdAt | date "Jan 2, 2006"}}
```

### Encoding Functions

```go
// Base64 encode
{{.secret | b64enc}}

// JSON encode
{{.config | toJson}}

// YAML encode
{{.config | toYaml}}
```

## File and Folder Names

Template syntax works in file and folder names too:

```
src/
├── {{.appName}}/
│   ├── {{.appName}}_test.go
│   └── main.go
```

With `appName: "myservice"`, this becomes:

```
src/
├── myservice/
│   ├── myservice_test.go
│   └── main.go
```

## Jinja2-like Engine

For templates that heavily use Go code (where `{{` conflicts), SKA offers a Jinja2-like engine:

```bash
ska create --blueprint ./my-template --output ./project --engine jinja
```

This uses `{%` and `%}` delimiters instead:

```jinja
func main() {
    fmt.Println("Hello, {% .appName %}!")
}
```

:::note
The Jinja2 engine (powered by [Pongo2](https://www.schlachter.tech/solutions/pongo2-template-engine/)) has a more limited function set compared to Go templates with Sprig.
:::

## Whitespace Control

Control whitespace around template tags:

```go
// Trim leading whitespace
{{- .appName}}

// Trim trailing whitespace
{{.appName -}}

// Trim both
{{- .appName -}}
```

## Learning Resources

- [Go Template Documentation](https://pkg.go.dev/text/template)
- [Sprig Function Reference](https://masterminds.github.io/sprig/)
- [Go Template How-To](https://www.digitalocean.com/community/tutorials/how-to-use-templates-in-go)
