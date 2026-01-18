---
title: Installation
description: How to install SKA on your system.
---

SKA is distributed via [Homebrew](https://brew.sh/), making installation quick and simple.

## Prerequisites

- macOS or Linux
- [Homebrew](https://brew.sh/) package manager installed

## Install via Homebrew

Run the following commands to install SKA:

```bash
brew tap gchiesa/ska
brew install ska
```

## Verify Installation

After installation, verify SKA is working correctly:

```bash
ska --version
```

You should see the version number printed to the console.

## Getting Help

SKA includes built-in help for all commands:

```bash
# General help
ska --help

# Help for specific commands
ska create --help
ska update --help
ska config --help
```

## Next Steps

Now that SKA is installed, head to the [Quick Start](/ska/getting-started/quick-start/) guide to create your first scaffolded project.
