[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/designed-in-etch-a-sketch.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)


# tmux-project

**tmux-project** creates, maintains, archives and removes profiles used in
combination with the tmux command.

This command allows the user to list and use different tmux based profiles.
Together with bash (or other shell based) profiles you can easily maintain
multiple sessions.

# Reason

The reason I'm writing this thing.

At the moment I'm working in a project based on tickets. For each ticket I
re-checkout what ever repositories I need. Seems silly until you work on 4
projects at one and start to lose sight of things. In my case I wrote a small
wrapper in my bash profile that allows me to resume working on a project by
executing:

  resume <projectname>

This solution consists of a few distinct targets:

* HOME/.tmux.d/<project>.rc, the tmux configuration used
  for this.
* HOME/.tmux.d/<project>.env, the bash configuration
  sourced when resuming.
* PROJECSTDIR the location where projects are checked out.

For the longest time I had only one type of project to work on and the original
client/organization specific solution I wrote in Python covered this neatly.
However others recently came. Different ticket name format, different archive,
etc.. And instead of re-writing my python thing I instead opted for a golang
based approach. Why?  Because I'm shit at golang, it's the Christmas holiday
and I have nothing better to do.


# Synopsis
## shell

Provides a way to integrate tmux-project into shell by executing:

  eval "$(tmux-project shell)"

(don't forget the quotes)


```
```


## tui

Launch an interactive Terminal User Interface (TUI) for managing tmux-project.

The TUI provides a visual interface for:
  - Browsing and selecting projects
  - Editing project settings (name, directory, description, type)
  - Managing global configuration
  - Creating, archiving, and removing projects
  - Starting tmux sessions

Navigation:
  - Arrow keys / j,k: Move up/down
  - Enter: Select/Edit
  - Tab: Switch between panels
  - Esc: Go back
  - q: Quit

Examples:
  tmux-project tui          # Launch TUI
  tmux-project tui -v       # Launch with verbose logging


```
```


## project

### archive

Creates a TAR archive of a project.


```
```


### create

Create a new project


```
```


### edit

Edit the config of a project


```
```


### list

List the available sessions


```
```


### listfiles



```
```


### remove

remove a project


```
```


### resume

resume a session


```
```


## type

### init



```
```


### list

list project types


```
```



# Functionality

## Targets

| Target                              | Description                 |
|:------------------------------------|:----------------------------|
| ```${HOME}/.tmux.d/<project>.env``` | environment file            |
| ```${HOME}/.tmux.d/<project>.rc```  | tmux configuration          |
| ```PROJECTS```                      | location projects are setup |
