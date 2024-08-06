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
## archive

Creates a TAR archive of a project.


```
Usage:
  tmux-project archive <project> [flags]

Flags:
  -a, --archivename string   Archive file
  -h, --help                 help for archive

Global Flags:
  -v, --verbose   Verbose logging
```


## create

Create a new project


```
Usage:
  tmux-project create <project> [flags]

Flags:
  -d, --description string   Description of the project
  -h, --help                 help for create
  -t, --type string          Type of project (default "default")

Global Flags:
  -v, --verbose   Verbose logging
```


## edit

Edit the config of a project


```
Usage:
  tmux-project edit <project> [flags]

Flags:
  -h, --help   help for edit

Global Flags:
  -v, --verbose   Verbose logging
```


## init



```
Usage:
  tmux-project init <projecttype> [flags]

Flags:
  -f, --force   Force
  -h, --help    help for init

Global Flags:
  -v, --verbose   Verbose logging
```


## list

List the available sessions


```
Usage:
  tmux-project list [flags]

Flags:
  -f, --full   Print full
  -h, --help   help for list

Global Flags:
  -v, --verbose   Verbose logging
```


## listfiles



```
Usage:
  tmux-project listfiles <project> [flags]

Flags:
  -h, --help   help for listfiles

Global Flags:
  -v, --verbose   Verbose logging
```


## projectinit



```
```


## resume

resume a session


```
Usage:
  tmux-project resume <project> [flags]

Flags:
  -h, --help   help for resume

Global Flags:
  -v, --verbose   Verbose logging
```


## shell

Provides a way to integrate tmux-project into shell by executing:

  eval "$(tmux-project shell)"

(don't forget the quotes)


```
Usage:
  tmux-project shell [<shell>] [flags]

Flags:
  -h, --help   help for shell

Global Flags:
  -v, --verbose   Verbose logging
```



# Functionality

## Targets

| Target                              | Description                 |
|:------------------------------------|:----------------------------|
| ```${HOME}/.tmux.d/<project>.env``` | environment file            |
| ```${HOME}/.tmux.d/<project>.rc```  | tmux configuration          |
| ```PROJECTS```                      | location projects are setup |
