#!/bin/bash
C_SCRIPTPATH="$(readlink -f "$0")"
C_SCRIPTNAME="$(basename "$C_SCRIPTPATH" .sh)"
C_SCRIPTDIR="$(dirname "${C_SCRIPTPATH}")"
C_TOPDIR="$(git rev-parse --show-toplevel)"

readonly C_SCRIPTPATH
readonly C_SCRIPTNAME
readonly C_SCRIPTDIR
readonly C_TOPDIR

C_HEADER="$(cat <<-ENDHEADER
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

* ``HOME/.tmux.d/<project>.rc``, the tmux configuration used
  for this.
* ``HOME/.tmux.d/<project>.env``, the bash configuration
  sourced when resuming.
* ``PROJECSTDIR`` the location where projects are checked out.

For the longest time I had only one type of project to work on and the original
client/organization specific solution I wrote in Python covered this neatly.
However others recently came. Different ticket name format, different archive,
etc.. And instead of re-writing my python thing I instead opted for a golang
based approach. Why?  Because I'm shit at golang, it's the Christmas holiday
and I have nothing better to do.


# Synopsis

ENDHEADER
)"

C_FOOTER="$(cat <<-ENDFOOTER

# Functionality

## Targets

| Target                              | Description                 |
|:------------------------------------|:----------------------------|
| \`\`\`\${HOME}/.tmux.d/<project>.env\`\`\` | environment file            |
| \`\`\`\${HOME}/.tmux.d/<project>.rc\`\`\`  | tmux configuration          |
| \`\`\`PROJECTS\`\`\`                      | location projects are setup |

ENDFOOTER
)"

readonly C_HEADER
readonly C_FOOTER



function pathmunge()
{
  [ -d "$1" ] || return

  if echo "$PATH" | grep -E -q "(^|:)$1($|:)"
  then
    return
  fi

  if [ "$2" = "after" ]
  then
      PATH=$PATH:$1
  else
      PATH=$1:$PATH
  fi
}

function logging()
{
  local priority="$1"; shift
  logger -p "${C_FACILITY}.${priority}" -i -s -t "${C_SCRIPTNAME}" -- "${priority} $*"
}
function logging_err()   { logging "err" "$@";   }
function logging_info()  { logging "info" "$@";  }
function logging_warn()  { logging "warn" "$@";  }
function logging_debug() { logging "debug" "$@"; }
function die()           { script_exit "$1" 1;   }
function script_exit()
{
  local string="$1"
  local retv="${2:-0}"
  if [ "$retv" = "0" ]
  then
    logging_info "$string"
  else
    logging_err "$string"
  fi
  exit "$retv"
}

function print_option()
{
  local option="$1"
  printf "## %s\n\n" "${option}"

  cat "messages/long/${option}"

  printf "\n\n\`\`\`\n"
  $COMMAND "${option}" --help 2>/dev/null | sed -n '/Usage/,$p'
  printf "\`\`\`\n\n\n"

}

function print_sub_option()
{
  local subcommand="$1"
  local option="$2"
  printf "### %s\n\n" "${option}"

  cat "messages/long/${subcommand}/${option}"

  printf "\n\n\`\`\`\n"
  $COMMAND "${subcommand}" "${option}" --help 2>/dev/null | sed -n '/Usage/,$p'
  printf "\`\`\`\n\n\n"

}

function list_options()
{

  pushd "${C_TOPDIR}" >/dev/null 2>&1 || return
  find messages/ -maxdepth 2 -mindepth 2 \( \
       -path '*/short/*'  \
    -o -path '*/use/*' \
    -o -path '*/long/*' \) -type f -printf "%f\n"|sort|uniq -c | \
    while read -r num option
    do
      [[ "${num}" == 3 ]] || continue
      [[ "${option}" == "root" ]] && continue
      echo "${option}"
    done
}

function list_sub_options()
{
  local subcommand="$1"
  pushd "${C_TOPDIR}" >/dev/null 2>&1 || return

  find messages/  \( \
       -path "*/short/${subcommand}/*"  \
    -o -path "*/use/${subcommand}/*" \
    -o -path "*/long/${subcommand}/*" \) -type f -printf "%f\n"|sort|uniq -c | \
    while read -r num option
    do
      [[ "${num}" == 3 ]] || continue
      [[ "${option}" == "root" ]] && continue
      echo "${option}"
    done

}

function main()
{

  COMMAND="${C_TOPDIR}/build/$(go env GOOS)/$(go env GOARCH)/tmux-project"

  echo "${C_HEADER}"

  list_options | while read -r option
  do
    print_option "${option}" -h
  done

  find messages -mindepth 3 -maxdepth 3 -type f -name root -printf "%h\n" | \
    cut -d/ -f 3|sort -u | while read -r subcommand
  do
    printf "## %s\n\n" "$subcommand"
    list_sub_options "${subcommand}" | while read -r option
    do
      print_sub_option "$subcommand" "$option"
    done
  done

  echo "${C_FOOTER}"


}

#------------------------------------------------------------------------------#
#                                    Main                                      #
#------------------------------------------------------------------------------#

main > "${C_TOPDIR}/README.md"

#------------------------------------------------------------------------------#
#                                  The End                                     #
#------------------------------------------------------------------------------#
