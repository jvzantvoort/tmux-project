[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/designed-in-etch-a-sketch.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)


# tmux-project

**tmux-project** creates, maintains, archives, and removes profiles used in
combination with the tmux command.

This tool allows users to list and use different tmux-based profiles. Together with bash (or other shell-based) profiles, you can easily maintain multiple sessions and project environments.

## Motivation

This project was created to help manage multiple ticket-based projects, each with its own repositories and environment. Instead of manually tracking project state, `tmux-project` automates session and environment management, making it easy to resume work on any project.

## Features
- Manage tmux and shell profiles for multiple projects
- Archive and restore project environments
- List, create, edit, and remove project sessions
- Integrate with your shell for quick project switching

# Usage

## Shell Integration
To integrate tmux-project into your shell, run:

```sh
eval "$(tmux-project shell)"
```

## Project Commands

### Archive
Create a TAR archive of a project:

```sh
tmux-project project archive <project> [flags]
```

### Create
Create a new project:

```sh
tmux-project project create <project> [flags]
```

### Edit
Edit the config of a project:

```sh
tmux-project project edit <project> [flags]
```

### List
List the available sessions:

```sh
tmux-project project list [flags]
```

### List Files
List files in a project:

```sh
tmux-project project listfiles <project> [flags]
```

### Remove
Remove a project (optionally archive first):

```sh
tmux-project project remove <project> [flags]
```

### Resume
Resume a session:

```sh
tmux-project project resume <project> [flags]
```

### Project Types
Initialize or list project types:

```sh
tmux-project type init <projecttype> [flags]
tmux-project type list [flags]
```

# Functionality

## Targets

| Target                              | Description                 |
|:------------------------------------|:----------------------------|
| `${HOME}/.tmux.d/<project>.env`     | environment file            |
| `${HOME}/.tmux.d/<project>.rc`      | tmux configuration          |
| `PROJECTS`                          | location projects are setup |

# See Also
- [CHANGELOG.rst](CHANGELOG.rst)
- [contrib/README.md](contrib/README.md)

---

For more details, see the inline documentation and comments in the source code.
