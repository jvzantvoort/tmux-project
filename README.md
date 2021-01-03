[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/designed-in-etch-a-sketch.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)


# tmux-project

**tmux-project** creates, maintains, archives and removes profiles
used by the
[resume](https://github.com/jvzantvoort/homebin/blob/master/bin/resume_tmux)
command. This command allows the user to list and use different tmux
based profiles. Together with bash (or other shell based) profiles
you can easily maintain multiple sessions.

# Synopsis

## Commands


* ``tmux-project archive``
* ``tmux-project create``
* ``tmux-project edit``
* ``tmux-project init``
* ``tmux-project list``

All commands are described in help.

# Functionality

## Targets

| Target                                   | Description                 |
|:-----------------------------------------|:----------------------------|
| ```${HOME}/.bash/tmux.d/<project>.env``` | environment file            |
| ```${HOME}/.bash/tmux.d/<project>.rc```  | tmux configuration          |
| ```PROJECTS```                           | location projects are setup |

## Use in bash

The following lines allow profiles to source created environment
files:

```
export SESSIONNAME=`tmux display-message -p '#S'`
if [ -f "$HOME/.bash/tmux.d/${SESSIONNAME}.env" ]
then
  source  "$HOME/.bash/tmux.d/${SESSIONNAME}.env"
fi
```
