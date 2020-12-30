# tmux-project

**tmux-project** creates, maintains, archives and removes profiles
used by 

Create a tmux project.

In "https://github.com/jvzantvoort/homebin" the "resume" tool uses
two files to access different projects:

${HOME}/.bash/tmux.d/<project>.env | environment file
${HOME}/.bash/tmux.d/<project>.rc  | tmux configuration

It also uses another location namely the project directory.

tmux-project create -n <name> -t <type> 

tmux-project list [-t <type>]

tmux-project destroy -n <name>

tmux-project archive -n <name> [-a <archive dir>]

# Use

## Bash


```
export SESSIONNAME=`tmux display-message -p '#S'`
if [ -f "$HOME/.bash/tmux.d/${SESSIONNAME}.env" ]
then
  source  "$HOME/.bash/tmux.d/${SESSIONNAME}.env"
fi
```
