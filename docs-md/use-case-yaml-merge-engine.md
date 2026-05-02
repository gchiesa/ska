# USE CASE: YAML-Aware Updates with the `yaml-merge` Engine

This guide shows how to:

- Keep only specific keys of a YAML file centrally managed while letting developers freely extend the rest.
- Prevent upstream updates from wiping out keys that are not known to the blueprint.
- Adopt an existing, previously-unmanaged YAML file into SKA using a one-shot replacement directive, after which every
  subsequent update uses the YAML merge engine.

## Why the default text engine is not always enough

The default SKA engine replaces the content of a managed block verbatim on every update. That is exactly what you want
for most files — markup, scripts, Dockerfiles — where the template owns the whole block.

For YAML configuration files the story is different: teams often need to add service-specific keys *inside* a section
that is partially governed by a central platform team. With a text-based replacement those additions disappear on the
next `ska update`.

The `yaml-merge` engine solves this: it computes a **structural diff** between the upstream template and the current
destination block, then applies only the keys that are present in the template. Keys that exist only in the destination
are left untouched.

## How it works

Add the `[engine:yaml-merge]` bracket modifier to the `ska-start` tag in your **blueprint** template:

```yaml
# ska-start[engine:yaml-merge]:platform-settings
platform:
  logging:
    level: "{{.logLevel}}"
    format: json
  tracing:
    enabled: true
    samplingRate: "{{.tracingSamplingRate}}"
# ska-end
```

> [!NOTE]
> The modifier is only meaningful in the **blueprint**. SKA writes a plain `# ska-start:<id>` / `# ska-end` block
> into the destination — no modifier is stored there.

When `ska update` runs:

1. The upstream template is rendered with the captured variables.
2. If a `# ska-start:platform-settings` block already exists in the destination, SKA merges the rendered output
   **into** the existing block instead of replacing it wholesale.
3. Keys in the rendered template override the corresponding destination values.
4. Keys present only in the destination are preserved unchanged.
5. SKA logs every key path it patched at `DEBUG` level, e.g.:
   ```
   DEBUG  [engine:yaml-merge]  [patch]  [platform.tracing.samplingRate]
   ```

> [!TIP]
> On the **first** scaffolding (no block exists yet) the yaml-merge engine behaves like the default engine — the
> block is created wholesale. Merge semantics apply from the **second** `ska update` onward.

---

## Scenario: service runtime configuration

### Blueprint layout

```
templates-repo/
  service-runtime/
    .ska-upstream.yaml
    app-config.yaml       ← centrally managed YAML template
    main.go
    ...
```

### `app-config.yaml` in the blueprint

The platform team owns the `platform` section. Everything else can be freely extended by the service team.

```yaml
# Managed by the platform team — do not edit the platform section by hand.
# ska-start[engine:yaml-merge]:platform-settings
platform:
  logging:
    level: "{{.logLevel}}"
    format: json
  tracing:
    enabled: true
    samplingRate: "{{.tracingSamplingRate}}"
  security:
    tlsMinVersion: "{{.tlsMinVersion}}"
# ska-end

# Service-team section — add your own keys freely below this line.
service: { }
```

### After initial scaffolding

`ska create` expands the template. The destination `app-config.yaml` in the project looks like this (variables replaced,
markers written without engine modifier):

```yaml
# Managed by the platform team — do not edit the platform section by hand.
# ska-start:platform-settings
platform:
  logging:
    level: info
    format: json
  tracing:
    enabled: true
    samplingRate: "0.05"
  security:
    tlsMinVersion: "TLSv1.2"
# ska-end

# Service-team section — add your own keys freely below this line.
service: { }
```

### Developer adds service-specific settings

The team extends the file with their own keys — both inside and outside the managed block:

```yaml
# Managed by the platform team — do not edit the platform section by hand.
# ska-start:platform-settings
platform:
  logging:
    level: info
    format: json
  tracing:
    enabled: true
    samplingRate: "0.05"
  security:
    tlsMinVersion: "TLSv1.2"
  # team added: custom exporter endpoint
  metrics:
    exporterURL: "https://metrics.internal/push"
# ska-end

# Service-team section — add your own keys freely below this line.
service:
  name: payment-processor
  port: 8080
  db:
    host: postgres.internal
    maxConns: 20
```

### Platform team releases a new template version

The blueprint is updated: `tracingSamplingRate` bumped to `0.1` and a new key `security.mTLS` is added. The service team
runs:

```sh
ska update --path .
```

**Result** — only the keys managed by the upstream template are touched:

```yaml
# Managed by the platform team — do not edit the platform section by hand.
# ska-start:platform-settings
platform:
  logging:
    level: info
    format: json
  tracing:
    enabled: true
    samplingRate: "0.1"       # ← updated from upstream
  security:
    tlsMinVersion: "TLSv1.2"
    mTLS: true                # ← new key from upstream
  # team added: custom exporter endpoint
  metrics:
    exporterURL: "https://metrics.internal/push"   # ← preserved
# ska-end

# Service-team section — add your own keys freely below this line.
service:
  name: payment-processor
  port: 8080
  db:
    host: postgres.internal
    maxConns: 20
```

`metrics.exporterURL` — added by the team inside the managed block — survives the update because the yaml-merge engine
does not know about it and therefore does not touch it.

---

## Scenario: adopting an existing YAML file

Sometimes a project already has a YAML configuration file that was written by hand before SKA was introduced. The file
has no `ska-start` / `ska-end` markers, so SKA does not touch it during normal updates.

The solution is a **one-shot adoption** using the `ska-replace-match` directive in combination with the
`yaml-merge` engine modifier. On the first run, `ska-replace-match` replaces the **entire file content** with a freshly
rendered managed block. From that point on, the block is present and every subsequent update uses the
`yaml-merge` engine — incremental and non-destructive.

> [!WARNING]
> `ska-replace-match` with a whole-file regex is **destructive on first use**: it discards the original
> hand-crafted content and replaces it with the template output. Back up anything you need to keep before
> running the adoption update, then re-add it manually inside (or outside) the managed block afterwards.

### The existing file (in the project, not managed by SKA)

```yaml
# app-config.yaml — hand-crafted, no SKA markers yet
platform:
  logging:
    level: warn
  tracing:
    enabled: false
  security:
    tlsMinVersion: "TLSv1.1"
```

### Blueprint template that adopts the file

Add the `ska-replace-match` directive with the regex `(?s).*` — the `(?s)` flag makes `.` match newlines, so
`.*` greedily matches the entire file content and replaces it with the managed block.

```yaml
# ska-start[engine:yaml-merge]:platform-settings + ska-replace-match:(?s).*
platform:
  logging:
    level: "{{.logLevel}}"
    format: json
  tracing:
    enabled: true
    samplingRate: "{{.tracingSamplingRate}}"
  security:
    tlsMinVersion: "{{.tlsMinVersion}}"
# ska-end
```

> [!NOTE]
> `ska-replace-match` is a one-shot directive: it fires **only when no managed block with that identifier is
> present** in the destination. Once `# ska-start:platform-settings` exists, subsequent updates use the
> yaml-merge engine and the replace directive is silently skipped.

### First `ska update` — adoption run

SKA finds no `# ska-start:platform-settings` block in `app-config.yaml`. The `ska-replace-match:(?s).*`
directive fires: `(?s).*` matches the entire file in one pass, and the whole content is replaced with the rendered
managed block:

```yaml
# ska-start:platform-settings
platform:
  logging:
    level: info
    format: json
  tracing:
    enabled: true
    samplingRate: "0.05"
  security:
    tlsMinVersion: "TLSv1.2"
# ska-end
```

The original hand-crafted content is gone. The file is now clean and fully under SKA management.

### Team extends the managed block after adoption

After the adoption the service team adds their own keys inside the block:

```yaml
# ska-start:platform-settings
platform:
  logging:
    level: info
    format: json
  tracing:
    enabled: true
    samplingRate: "0.05"
  security:
    tlsMinVersion: "TLSv1.2"
  # team added after adoption
  metrics:
    exporterURL: "https://metrics.internal/push"
# ska-end
```

### Every subsequent `ska update` — yaml-merge takes over

From the second update on, the `# ska-start:platform-settings` block already exists. SKA applies the yaml-merge engine:
only the keys present in the upstream template are updated; anything the team has added inside the block is preserved.

For example, when the platform team bumps `samplingRate` to `0.1` and adds `security.mTLS`:

```yaml
# ska-start:platform-settings
platform:
  logging:
    level: info
    format: json
  tracing:
    enabled: true
    samplingRate: "0.1"      # ← updated from upstream
  security:
    tlsMinVersion: "TLSv1.2"
    mTLS: true               # ← new key from upstream
  # team added after adoption
  metrics:
    exporterURL: "https://metrics.internal/push"   # ← preserved
# ska-end
```

---

## Quick reference

### Tag syntax

| Use case                                                           | Blueprint tag                                                        |
|--------------------------------------------------------------------|----------------------------------------------------------------------|
| Default text engine (entire block replaced)                        | `# ska-start:my-block`                                               |
| YAML merge engine (structural key merge)                           | `# ska-start[engine:yaml-merge]:my-block`                            |
| One-shot adoption: replace entire file, then yaml-merge on updates | `# ska-start[engine:yaml-merge]:my-block + ska-replace-match:(?s).*` |

### Behaviour matrix

| Situation                             | Default engine                     | `yaml-merge` engine |
|---------------------------------------|------------------------------------|---------------------|
| Block does not exist yet              | Adopt directive fires (or nothing) | Same as default     |
| Block exists, key in template         | Value overwritten                  | Value overwritten   |
| Block exists, key only in destination | **Key deleted**                    | **Key preserved**   |
| Block exists, new key in template     | Key inserted                       | Key inserted        |
| YAML comments in destination          | Lost on update                     | Preserved           |

> [!TIP]
> **When to prefer each engine**
> - Use the **default engine** for non-YAML files, or for YAML blocks where the upstream template owns *all* keys
    > and developers should never add extra ones.
> - Use the **yaml-merge engine** whenever a YAML block is shared between central governance and local
    > customisation — typical for configuration files, Helm values, Kubernetes manifests, and similar.
