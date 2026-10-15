# Commands

All commands accept `-v, --verbose` for debug logging and `-h, --help`.

## Projects

### List projects

```sh
tmux-project list [-f]
tmux-project project list [-f]
```

Without flags the project names are printed, one per line. With `-f, --full`
a table shows type, status, whether the tmux session is active, description
and directory.

### Create a project

```sh
tmux-project project create [-t <type>] [-d <description>] <project>
```

Flags:

- `-t, --type <type>`: project type (default: `default`)
- `-d, --description <desc>`: description; asked for when left out

The project name is also the tmux session name and the base name of the
files in `~/.tmux.d`, so it must not contain `/`, `\`, `.`, `:` or
whitespace, and must not start with `-`. It must also match the `pattern` of
the project type.

Creating a project makes the project directory (or clones `root`), writes
the target files, clones the repositories and runs the setup actions of the
project type. The setup actions run one after another, in the project
directory, through `sh -c`; creation stops at the first action that fails.

### Resume a project

```sh
tmux-project resume <project>
tmux-project project resume <project>
resume <project>
```

This command:

- updates the last activity time of the project
- attaches to the tmux session of the project, or creates it in the project
  directory
- inside tmux, switches the current client to the session instead of
  nesting tmux

`resume ls` prints the same table as `tmux-project list -f`.

### Edit a project

```sh
tmux-project project edit <project>
```

Opens the `.rc` and `.env` files of the project in `$VISUAL`, `$EDITOR` or
vim. vim-like editors show them side by side.

### Archive a project

```sh
tmux-project project archive [-a <file>] <project>
```

Creates a gzipped tar archive of the project directory and the project's
files in `~/.tmux.d`.

Flags:

- `-a, --archivename <file>`: archive file (default: `<project
  directory>.tar.gz`)

### Remove a project

```sh
tmux-project project remove [-a] [-y] <project>
```

Removes the project directory, its target files and its session files in
`~/.tmux.d`. The paths are listed and you are asked to confirm first. A
project whose tmux session is still running is not removed.

Flags:

- `-a, --archive`: archive the project to `<directory>.tar.gz` before removal
- `-y, --yes`: do not ask for confirmation

### List project files

```sh
tmux-project project listfiles <project>
```

Prints the project directory and the files in `~/.tmux.d` that belong to
the project.

## Project types

### List project types

```sh
tmux-project type list
```

### Initialize a project type

```sh
tmux-project type init [-f] <type>
```

Creates `~/.tmux-project/<type>/` with a `config.yml` and the default
templates to customize, see [Project types](README.md#project-types). An
existing project type is only overwritten with `-f, --force`.

## Shell integration

```sh
eval "$(tmux-project shell [bash|zsh])"
```

Prints the shell code that loads the project environment in a tmux
session, see [SHELL.md](SHELL.md).

## Interactive TUI

```sh
tmux-project tui
```

Browse projects and edit their description, type and status. Press `?` in
the TUI for the keys.

## Updates and version

```sh
tmux-project update [-f]
tmux-project version
```

`update` downloads the latest release for this platform from GitHub,
verifies its sha256 checksum and replaces the installed binaries. It does
nothing when the latest version is already installed. Development builds
are only replaced with `-f, --force`.

## proj_info

```sh
proj_info [-s] [-d <depth>] [-v]
```

Shows the project of the current tmux session and the directories in it,
grouped by the chapters in `~/.tmux.d/<project>.chapters.yml`. The session
is taken from `$SESSIONNAME` or asked from tmux. `common.rc` binds it to
`prefix e`.

Flags:

- `-s`: short overview, closes after a key press or 10 seconds
- `-d <depth>`: maximum depth to search for directories (default: 1)
