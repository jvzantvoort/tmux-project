[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/designed-in-etch-a-sketch.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)


# tmux-project

**tmux-project** creates, maintains, archives, and removes profiles used in
combination with the tmux command.

This tool allows users to list and use different tmux-based
profiles. Together with bash (or other shell-based) profiles, you
can easily maintain multiple sessions and project environments.

Each project can have its own configuration, working directory, git
repositories, and setup actions, making it handy for ticket-based
workflows, multi-repository projects, and context switching between
different development environments.

## Motivation

This project was created to help manage multiple ticket-based
projects, each with its own repositories and environment. Instead of
manually tracking project state, `tmux-project` automates session
and environment management, making it easy to resume work on any
project.

## Features

- Manage tmux and shell profiles for multiple projects
- Archive and restore project environments
- List, create, edit, and remove project sessions
- Integrate with your shell for quick project switching

# Usage

* Check [shell](SHELL.md) for shell integration
* Check [commands](COMMANDS.md) for commands

# Functionality

## Project Structure

Each project consists of:

1. **Configuration File**: `${HOME}/.tmux.d/<project>.yaml` - Project metadata and settings
2. **Environment File**: `${HOME}/.tmux.d/<project>.env` - Project-specific environment variables
3. **Tmux Config**: `${HOME}/.tmux.d/<project>.rc` - Tmux session configuration
4. **Working Directory**: Location where project files are stored (defined in project type)

## Project Types

Project types are reusable templates stored in `${HOME}/.tmux-project/`. Each project type defines:

- **Directory pattern**: Where projects of this type are created
- **Name pattern**: Regex pattern for valid project names
- **Git repositories**: Repositories to clone on project creation
- **Setup actions**: Commands to run during project initialization
- **Target files**: Templates for environment and tmux configuration files

### Example Project Type Structure

```yaml
projecttype: "ticket"
directory: "${HOME}/projects/${NAME}"
pattern: "^[A-Z]+-[0-9]+$"
root: "directory"
repos:
  - url: "https://github.com/example/repo.git"
    destination: "repo"
    branch: "main"
setupactions:
  - "echo 'Project initialized'"
targets:
  - name: "env"
    destination: "${HOME}/.tmux.d/${NAME}.env"
    mode: "0644"
  - name: "tmux"
    destination: "${HOME}/.tmux.d/${NAME}.rc"
    mode: "0644"
```

## Configuration Files

### Environment File (`<project>.env`)
Shell script sourced when resuming a project. Can contain:
- Environment variables
- PATH modifications
- Aliases and functions
- Project-specific settings

### Tmux Config File (`<project>.rc`)
Tmux configuration loaded when starting the project session. Can define:
- Window layouts
- Pane arrangements
- Key bindings
- Status bar customizations

## Targets

| Target                              | Description                           |
|:------------------------------------|:--------------------------------------|
| `${HOME}/.tmux.d/<project>.yaml`    | Project configuration (YAML)          |
| `${HOME}/.tmux.d/<project>.env`     | Shell environment file                |
| `${HOME}/.tmux.d/<project>.rc`      | Tmux session configuration            |
| `${HOME}/.tmux-project/`            | Project type templates directory      |
| `PROJECTS`                          | Base directory where projects reside  |

## Template Variables

The following variables can be used in project type templates and configurations:

| Variable          | Description                                  |
|:------------------|:---------------------------------------------|
| `${NAME}`         | Project name                                 |
| `${HOME}`         | User's home directory                        |
| `${HOMEDIR}`      | Same as HOME                                 |
| `${USER}`         | Current username                             |
| `${GOPATH}`       | Go workspace path                            |
| `${GOARCH}`       | Target architecture                          |
| `${GOOS}`         | Target operating system                      |

# Examples

## Creating a Ticket-Based Project

```sh
# Initialize the ticket project type (first time only)
tmux-project type init ticket

# Create a new project for ticket PROJ-123
tmux-project project create PROJ-123 -t ticket -d "Fix authentication bug"

# Resume work on the project
tmux-project resume PROJ-123
```

## Creating a Go Development Project

```sh
# Create project type for Go development
tmux-project type init golang

# Create a new Go project
tmux-project project create myapp -t golang -d "New Go application"

# Resume the project
tmux-project resume myapp
```

## Archiving and Removing a Project

```sh
# Archive a project before removal
tmux-project project remove oldproject --archive

# Or archive without removing
tmux-project project archive myproject -o ~/backups/myproject.tar.gz
```

# Development

## Project Structure

```
tmux-project/
├── cmd/
│   ├── tmux-project/    # Main CLI application
│   ├── resume/          # Resume command helper
│   └── proj_info/       # Project info utilities
├── project/             # Project management core logic
├── projecttype/         # Project type handling
├── config/              # Configuration utilities
├── tmux/                # Tmux integration
├── git/                 # Git operations
├── archive/             # Archive creation
├── messages/            # Embedded message templates
├── utils/               # Common utilities
├── errors/              # Custom error types
├── update/              # Self-update functionality
└── version/             # Version information
```

## Building from Source

```sh
# Build the main binary
go build -o tmux-project ./cmd/tmux-project/

# Run tests
go test ./...

# Build all commands
./build.sh
```

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

# Requirements

- **Go**: 1.23.0 or later
- **tmux**: For session management features
- **git**: For repository cloning (optional)

# Dependencies

Key dependencies include:
- `github.com/spf13/cobra`: CLI framework
- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/sirupsen/logrus`: Logging
- `gopkg.in/yaml.v2`: YAML parsing

# See Also
- [CHANGELOG.rst](CHANGELOG.rst)
- [contrib/README.md](contrib/README.md)

---

For more details, see the inline documentation and comments in the source code.
