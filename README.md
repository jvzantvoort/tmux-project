[![forthebadge](https://forthebadge.com/images/badges/made-with-crayons.svg)](https://forthebadge.com)
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

## archive

Archive a project

```
tmux-project archive [-v] [-a|archivename <archivename>] -n | -projectname <projectname> 
  -a string
        Archive file
  -archivename string
        Archive file
  -n string
        Name of project
  -projectname string
        Name of project
  -v    Verbose logging
```

## create


Create a new project.

```
tmux-project create [-t|-projecttype <type>] -n <name> | -projectname <name> [-v]
  -n string
        Name of project
  -projectname string
        Name of project
  -projecttype string
        Type of project (default "default")
  -t string
        Type of project (default "default")
  -v    Verbose logging
```

## edit

Edit a projects environment and tmux configfile

```
tmux-project edit -n <projectname> [-v]
  -n string
        Name of project
  -projectname string
        Name of project
  -v    Verbose logging
```

## init

Initialize a new project type

```
tmux-project init [-v] -t | -projecttype <projecttype>
  -f    Force (re)creation
  -projecttype string
        Type of project (default "default")
  -t string
        Type of project (default "default")
  -v    Verbose logging
```

## list

The "list" command list projects currrently configured

```
tmux-project list [-projectname|-n <name>] [-v] [-v]
  -f    Print full
  -n string
        Name of project
  -projectname string
        Name of project
  -v    Verbose logging
```

## listfiles

The "listfiles" command lists the projects currrently in a project's
configuration.

```
tmux-project listfiles [-projectname|-n <name>] [-v]
  -n string
        Name of project
  -projectname string
        Name of project
  -v    Verbose logging
```

## shell

Allows tmux-project to be integrated in a shell. For example for bash add the
following to the profile (bash is the default).

```sh
eval "$(tmux-project shell)"
```

```
tmux-project shell [-s | -shellname <shell>]

  -s string
        Name of the shell profile to provide (default "bash")
  -shellname string
        Name of the shell profile to provide (default "bash")
  -v    Verbose logging
```

# Functionality

## Targets

| Target                              | Description                 |
|:------------------------------------|:----------------------------|
| ```${HOME}/.tmux.d/<project>.env``` | environment file            |
| ```${HOME}/.tmux.d/<project>.rc```  | tmux configuration          |
| ```PROJECTS```                      | location projects are setup |
