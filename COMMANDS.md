## Project Commands

### List Projects
List all available projects:

```sh
tmux-project list
# or simply
tmux-project project list [flags]
```

### Create a New Project
Create a new project from a project type:

```sh
tmux-project project create <project-name> [flags]
```

Flags:
- `-t, --type <type>`: Specify the project type (default: `default`)
- `-d, --description <desc>`: Project description

The project name is also the tmux session name and the base name of the
files in `~/.tmux.d`, so it must not contain `/`, `\`, `.`, `:` or
whitespace, and must not start with `-`. It must also match the `pattern` of
the project type.

The setup actions of the project type run one after another, in the project
directory, through `sh -c`; creation stops at the first action that fails.

### Resume a Project
Resume working on an existing project:

```sh
tmux-project project resume <project-name> [flags]

# or the shorthand
tmux-project resume <project-name>

# or even shorter
resume <project-name>
```

This command:
- Loads the project configuration
- Attaches to the tmux session, or creates it in the project directory
- Inside tmux, switches the current client to the session instead of
  nesting tmux

### Edit Project Configuration
Edit a project's configuration file:

```sh
tmux-project project edit <project-name> [flags]
```

Opens the project's YAML configuration in your default editor.

### Archive a Project
Create a tar.gz archive of a project:

```sh
tmux-project project archive <project-name> [flags]
```

Flags:
- `-o, --output <file>`: Specify output file path (default: `<project-name>.tar.gz`)

Creates a compressed archive containing all project files and symlinks.

### Remove a Project
Remove a project (with optional archiving):

```sh
tmux-project project remove <project-name> [flags]
```

Removes the project directory, its target files and its session files in
`~/.tmux.d`. The paths are listed and you are asked to confirm first. A
project whose tmux session is still running is not removed.

Flags:
- `-a, --archive`: Archive the project to `<directory>.tar.gz` before removal
- `-y, --yes`: Do not ask for confirmation

### List Project Files
List all files in a project directory:

```sh
tmux-project project listfiles <project-name> [flags]
```

## Project Type Commands

### List Available Project Types
View all available project type templates:

```sh
tmux-project type list [flags]
```

### Initialize a Project Type
Create a new project type template:

```sh
tmux-project type init <projecttype> [flags]
```

This creates a template configuration that can be customized for creating new projects.

## Interactive TUI

Launch an interactive terminal user interface for managing projects:

```sh
tmux-project tui
```

Provides a visual interface for browsing and managing projects.
