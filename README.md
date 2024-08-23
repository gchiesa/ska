# SKA

SKA is a "skaffolding" tool that allows you to expand folders based on local or remote blueprints (GitHub and GitLab 
hosted)

Additionally, you can update your folder structure from the upstream blueprint to onboard new changes that are
maintained centrally and added overtime to upstream.

SKA is designed for simplicity and quick usage. If you want to integrate SKA capabilities in your app or framework
you can leverage the package under `pkg/`.

## SKA in action

![demo](assets/demo.gif)

## Get Started

SKA is distributed via [Homebrew](https://brew.sh/), so you can install it with the commands below:

```shell
brew tap gchiesa/ska
brew install ska
```

## How To Use

SKA is very simple to use. Just run `create` with the upstream blueprint you want to use. See below example for
scaffolding a golang command line tool:

```shell 
ska create --blueprint https://github.com/gchiesa/ska-golang-cli-template@master --output ~/workspace/myNewApp
```

This will load the remote template, collect the required variable via Terminal UI and create the project for you.

In the root of the project you can find the SKA configuration file `.ska-config.yml ` that will keep track from now
on of state and upstream reference.

When you want to update your project, just use:

```shell
cd ~/workspace/myNewApp
ska update --path .
```

you can modify the variables or accept the current version. Your project will be updated based on latest changes in
the upstream blueprint.

More information are available with `ska --help`.


## Concepts and Templating

### Upstream Blueprint

This is the typically centrally maintained template that everyone can use to expand its own folder structure. You
can specify the blueprint both locally and remotely with specific URIs:

* **file://** for local blueprints. E.g. `file:///Users/gchiesa/git/ska-example-template`
* **https://** for GitHub or GitLab (soon) blueprints. You can optionally pin a specific reference (tag or branch with
  the `@`
  symbol. E.g. `https://github.com/gchiesa/ska-example-template@v1.2.3`

### Update from upstream

Whenever the upstream template changes, you might want to onboard the changes yourself. This is natively supported
by SKA, by offering a simple `update` command with not additional arguments required.

See How To Use section for more information.


### Upstream Template Language

SKA currently fully supports primarily **[Go Template framework][go-template]** for templating your upstream blueprint.

Moreover, in addition to Go Template you can use the extensions offered by [Sprig functions][sprig]

If you are not familiar with Go Template have a look to this [simple how-to][go-template-how-to].

_NOTE:_ SKA supports also a Jinja2 like engine (use the `--engine jinja` argument) that will make your life easier
when you need to deal with native go code. However, the function set is limited. The support is based on [Pongo2
project][pongo2].


[go-template]: https://pkg.go.dev/text/template

[sprig]: https://masterminds.github.io/sprig/

[go-template-how-to]: https://www.digitalocean.com/community/tutorials/how-to-use-templates-in-go#step-4-writing-a-template

[pongo2]: https://www.schlachter.tech/solutions/pongo2-template-engine/

### Templating with partial sections

SKA offers the capability to manage only part of files.  
This is generally useful to keep only a specific part of the file centrally managed by the upstream blueprint and
let the user change the rest of the file. See as example
the `codecov.yml` [here](https://github.com/gchiesa/ska-golang-cli-template/blob/master/codecov.yml)

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


### Template variable collection via Terminal UI

While developing the blueprint, you can configure SKA to collect the variables via Terminal UI.

To do this you need to create on the root of your blueprint folder a special file `.ska-upstream.yaml`. The file is a
simple YAML file with the following structure (see example with comments):

```yaml
# list of path that will be ignored in the blueprint 
# typically it's wise to ignore `.git` and `.idea` (and your preferred IDE) since you might use
# the blueprint not only from remote but also from your local filesystem 
ignorePaths:
- .git
- .idea

# this is the section that SKA will consume to generate the input form.
# each input supports the following information:
# 
# * placeholder: the name of the variable you used in the template blueprint
# * label: the label used in the form
# * help: the help inline to be presented in the form for that specific field
# * regexp: a validation regexp that will be used to specify the set of accepted characters (NOT A ENTIRE PATTERN VALIDATION)
# * default: the default value you might want to present in the form
inputs:
- placeholder: githubRepo
  label: Github Repository
  help: The url of the github repository e.g. https://github.com/org/repo
  regexp: "^[a-z0-9-/:.]*$"

# ... other inputs 
```

SKA will check if the file is present and if yes, then a Terminal UI interface will be used to collect the inputs.

The interactive Terminal UI can be disabled with `--non-interactive (-n)` command line argument.


### Example template

Check this example golang cli tool template: https://github.com/gchiesa/ska-golang-cli-template

## Contribute

You are more than welcome! please have a look to [CONTRIBUTING.md](CONTRIBUTING.md)

## Credits

SKA is made with 💙 by Giuseppe Chiesa as exercise to learn a bit better GoLang.