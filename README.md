# SKA

SKA is a "skaffolding" tool that allows you to expand folders based on local or remote blueprint folder structure.

Additionally, you can update your folder structure from the upstream blueprint to always onboard new changes that are
added centrally to the upstream.

The interface is design for immediate use and simplicity. SKA offers also a simple package to integrate its
functionality in other tools and frameworks.


## How To Use

SKA is very simple to use. Just run it with the upstream blueprint you want to expand, with the following command

```shell 
ska create -b https://github.com/gchiesa/ska-example-template -o ~/my-project
```

This will load the remote template, collect the required variable via Terminal UI and create the project for you. In
the root of the project you can find the SKA configuration file that will keep track from now on of state and
upstream reference.

When you want to update your project, just use:

```shell
cd ~/my-project
ska update -p ~/my-project
```

you can modify the variables or accept the current version. Your project will be updated based on latest changes in
the upstream blueprint.

More information are available with `ska --help`.

## Concepts and Templating 

### Upstream Blueprint

This is the typically centrally maintained template that everyone can use to expand its own folder structure. You
can specify the blueprint both locally and remotely with specific URIs:

* **file://** for local blueprints. E.g. `file:///Users/gchiesa/git/ska-example-template`
* **https://** for GitHub or GitLab blueprints. You can optionally pin a specific reference (tag or branch with the `@`
  symbol. E.g. `https://github.com/gchiesa/ska-example-template@v1.2.3`

### Update from upstream

Whenever the upstream template changes, you might want to onboard the changes yourself. This is natively supported
by SKA, by offering a simple `update` command with not additional arguments required.

See How To Use section for more information.


### Upstream Template Language

SKA currently fully supports **[Go Template framework][go-template]** for templating your upstream blueprint.

Moreover, in addition to Go Template you can use the extensions offered by [Sprig functions][sprig]

If you are not familiar with Go Template have a look to this [simple how-to][go-template-how-to].

[go-template]: https://pkg.go.dev/text/template

[sprig]: https://masterminds.github.io/sprig/

[go-template-how-to]: https://www.digitalocean.com/community/tutorials/how-to-use-templates-in-go#step-4-writing-a-template

### Templating with partial sections

SKA offers the capability to manage only part of files. This is generally useful to keep only a specific part of the 
file centrally managed by the upstream blueprint and let the user change the rest of the file. 

This is achievable with the named partials, a type of block you can use in the SKA template. 

For example, if we want to have only a section of a larger file managed by SKA we can use the approach below:

```yaml
# this is an example YAML and we want SKA to only manage the section `key4`
---
root:
  key1: value1
  key2: {{.notManaged}}
  key3: value3
  # ska-start:key/4
  key4:
    subkey: "{{.appName}}"
    subkey2: value2
  # ska-end
  key5: value5
```

by using the `ska-start:<blockName>` and `ska-end` tags, you instruct SKA to only process the block and leave the 
rest of the file as it is. 

This is quite useful for files where only a part should be centrally managed and the rest will be customized by the 
user, after the initial creation.

## Credits

SKA is made with 💙 by Giuseppe Chiesa as exercise to learn a bit better GoLang.