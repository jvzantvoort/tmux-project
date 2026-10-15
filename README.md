[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/designed-in-etch-a-sketch.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)

# tmux-project

**tmux-project** creates, resumes, archives and removes project
workspaces, each with its own tmux session, shell environment and
working directory.

```sh
tmux-project project create -t golang -d "Fix the login timeout" JIRA-1234
resume JIRA-1234
```

## Reason

The reason I'm writing this thing.

At the moment I'm working in a project based on tickets. For each ticket I
re-checkout whatever repositories I need. Seems silly until you work on 4
projects at once and start to lose sight of things. In my case I wrote a
small wrapper in my bash profile that allows me to resume working on a
project by executing:

```sh
resume <projectname>
```

For the longest time I had only one type of project to work on and the
original client/organization specific solution I wrote in Python covered
this neatly. However others recently came. Different ticket name format,
different archive, etc.. And instead of re-writing my python thing I instead
opted for a golang based approach. Why? Because I'm shit at golang, it's the
Christmas holiday and I have nothing better to do.

## Requirements

- Linux or macOS
- tmux; 3.2 or newer for the project info popup (`prefix e`)
- bash or zsh for the shell integration
- git, when project types clone repositories
- optional: `jq`, to load settings shared by all projects of a type
- optional: [tmux-themepack](https://github.com/jimeh/tmux-themepack) in
  `~/.tmux.d/tmux-themepack`, used by the default tmux configuration

## Installation

Download and install the latest release for your platform, with checksum
verification:

```sh
curl -fsSLO https://raw.githubusercontent.com/jvzantvoort/tmux-project/main/tmux-project-update
bash tmux-project-update
```

As root the binaries go to `/usr/local/bin`, otherwise to `$GOBIN` or
`~/bin`. Alternatively download a `tmux-project.<version>.<os>.<arch>.tar.gz`
from the [releases](https://github.com/jvzantvoort/tmux-project/releases),
check it against its `.sha256` file and unpack it somewhere in your `PATH`.

Once installed, `tmux-project update` replaces the binaries with the latest
release.

A release contains three binaries:

| Binary         | Purpose                                              |
|:---------------|:-----------------------------------------------------|
| `tmux-project` | the main command                                     |
| `resume`       | shorthand for `tmux-project resume`                  |
| `proj_info`    | project overview, shown in a tmux popup (`prefix e`) |

## Quick start

1. Add the shell integration to `~/.bashrc` (see [SHELL.md](SHELL.md) for zsh):

   ```sh
   eval "$(tmux-project shell)"
   ```

2. Create a project type to start from:

   ```sh
   tmux-project type init golang
   ```

   and adjust `~/.tmux-project/golang/config.yml`, see
   [Project types](#project-types).

3. Create a project and start working on it:

   ```sh
   tmux-project project create -t golang -d "Fix the login timeout" JIRA-1234
   resume JIRA-1234
   ```

4. When you are done, archive and remove it:

   ```sh
   tmux-project project remove --archive JIRA-1234
   ```

All commands are described in [COMMANDS.md](COMMANDS.md).

## Files

| Path                                   | Description                                    |
|:---------------------------------------|:-----------------------------------------------|
| `~/.tmux.d/<project>.json`             | project definition and status                  |
| `~/.tmux.d/<project>.rc`               | tmux configuration of the project              |
| `~/.tmux.d/<project>.env`              | sourced by the shell integration in a session  |
| `~/.tmux.d/<project>.chapters.yml`     | grouping of directories for `proj_info`        |
| `~/.tmux.d/<type>.shared`              | sourced in every session of a project type     |
| `~/.tmux.d/common.rc`                  | tmux settings shared by all projects           |
| `~/.tmux-project/<type>/config.yml`    | project type definition                        |
| `~/.tmux-project/<type>/<template>`    | templates for the files of a new project       |

`common.rc` and `tmux_opt_source` are written to `~/.tmux.d` when they do
not exist yet; an existing `common.rc` is never overwritten.

## Project types

A project type describes how a project is set up. `tmux-project type init
<type>` creates `~/.tmux-project/<type>/` with a `config.yml` and the default
templates:

```yaml
---
type: golang
description: Go projects
# "directory" to create an empty directory, or a git URL to clone
root: "directory"
directory: "{{.HomeDir}}/{{.ProjectType}}/{{.Name}}"
# regular expression project names must match
pattern: ^[A-Z]+-[0-9]+$
# shell commands run in the project directory, in order
setupactions:
  - git init -q
  - mkdir -p docs
# repositories cloned into the project directory
repos:
  - url: git@github.com:example/service.git
    destination: service
    branch: main
# files created from the templates in this directory
targets:
  - name: default.env
    destination: "{{.Name}}.env"
    mode: 644
  - name: default.rc
    destination: "{{.Name}}.rc"
    mode: 644
  - name: default.chapters.yml
    destination: "{{.Name}}.chapters.yml"
    mode: 644
```

- `setupactions` run one after another through `sh -c`, so quoting, pipes
  and `&&` work. Creation stops at the first action that fails; the project
  is still saved so it can be inspected or removed.
- `repos` are cloned in parallel; `branch` is optional.
- A relative target `destination` is placed in `~/.tmux.d`; `mode` is octal.
- Project names must also be usable as tmux session names: no `/`, `\`,
  `.`, `:` or whitespace, and no leading `-`.

`directory`, `setupactions`, target destinations and template contents are
Go [templates](https://pkg.go.dev/text/template). Available fields:

| Field              | Value                                                |
|:-------------------|:-----------------------------------------------------|
| `{{.Name}}`        | project name                                         |
| `{{.Description}}` | project description                                  |
| `{{.ProjectType}}` | project type                                         |
| `{{.Directory}}`   | project directory (not in `directory` itself)        |
| `{{.HomeDir}}`     | home directory                                       |
| `{{.USER}}`        | user name                                            |
| `{{.GOOS}}`        | operating system, e.g. `linux`                       |
| `{{.GOARCH}}`      | architecture, e.g. `amd64`                           |
| `{{.GOPATH}}`      | Go path                                              |

## Development

```sh
go test -race ./...          # tests, they never touch your real ~/.tmux.d
./build.sh build             # build for this platform into build/
./build.sh check             # go vet, golangci-lint and staticcheck
./build.sh buildr linux arm64 /tmp/out   # what the release workflow runs
```

Changes are listed in [CHANGELOG.md](CHANGELOG.md), generated with
[git-cliff](https://git-cliff.org) from conventional commit messages.
