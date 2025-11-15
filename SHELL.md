# Shell Integration
To integrate tmux-project into your shell, add the following to your
bashrc:

```sh
eval "$(tmux-project shell)"
```

or to your zshrc:

```sh
eval "$(tmux-project shell zsh)"
```

This will bootstrap your environment file for the project. Through
the following it also allows for sourcing of per project custom
environment files:

```sh
if [ -f "$HOME/.tmux.d/${SESSIONNAME}.env" ]
then
  source  "$HOME/.tmux.d/${SESSIONNAME}.env"
fi
```

## builtin functions

### pcd

``pcd`` allows you to change directories based on relative paths and
aliases.

Given the following layout:

```
PROJDIR/
├── build
│   └── containers
├── diagnostics
└── work
    ├── ansible
    └── terraform
```

Executing ``pcd work/ansible`` would land you in
``PROJDIR/work/ansible`` etc.

Adding a configuration file ``PROJDIR/.pcd.rc`` allows you to add
aliases:

Given:

```
containers;PROJDIR/build/containers
terraform;PROJDIR/work/terraform
```

Executing ``pcd containers`` will end you up in
``PROJDIR/build/containers``.
