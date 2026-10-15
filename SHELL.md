# Shell Integration

To integrate tmux-project into your shell, add the following to your
`~/.bashrc`:

```sh
eval "$(tmux-project shell)"
```

or to your `~/.zshrc`:

```sh
eval "$(tmux-project shell zsh)"
```

Outside tmux this does nothing. In a shell inside a tmux session it:

- sets `SESSIONNAME` to the name of the tmux session
- sets `HISTFILE` to `~/.bash_history_<session>` (bash) or
  `~/.zsh_history_<session>` (zsh), so every project keeps its own history
- sources `~/.tmux.d/<session>.env`, which sets `PROJDIR` and anything else
  the project needs
- sources `~/.tmux.d/<type>.shared` when it exists, for settings shared by
  all projects of a type; this needs `jq` to read the type of the project
- defines the `pcd` function

## Builtin functions

### pcd

`pcd` changes directory relative to the project directory, with
completion. Without an argument it goes to `PROJDIR`.

Given the following layout:

```text
PROJDIR/
├── build
│   └── containers
├── diagnostics
└── work
    ├── ansible
    └── terraform
```

Executing `pcd work/ansible` would land you in `PROJDIR/work/ansible` etc.

Adding a configuration file `PROJDIR/.pcd.rc` allows you to add aliases, one
`<alias>;<directory>` per line, where `PROJDIR` is replaced by the project
directory:

```text
containers;PROJDIR/build/containers
terraform;PROJDIR/work/terraform
```

Executing `pcd containers` will end you up in `PROJDIR/build/containers`.
