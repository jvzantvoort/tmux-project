Changelog
=========


tmux-project-0.3.1 (2024-07-20)
-------------------------------

New
~~~
- ProjectType is now a template variable.

Fix
~~~
- Fail before using an existing folder.
- Viper doesnot want to be called twice.
- Cleanup old crud.
- Move project related function to own subdir.


tmux-project-0.3.0 (2024-07-20)
-------------------------------

New
~~~
- Merge previouse development.

Fix
~~~
- Staticcheck fixes.
- Rename commands.
- Move .bash/tmux.d to .tmux.d.
- Replace subcommand by cobra.
- Remove go-bindata in favor of embed.

Other
~~~~~
- Merge pull request #2 from jvzantvoort/feature/new-build.
- Add build.sh.
- New workflows for releasing.


tmux-project-0.2.3 (2023-01-24)
-------------------------------
- Add activity and extend visible list.


tmux-project-0.2.2 (2023-01-06)
-------------------------------
- Updates on go install command.


tmux-project-0.2.1 (2023-01-06)
-------------------------------
- Update go to 1.19.


tmux-project-0.2.0 (2023-01-06)
-------------------------------
- Merge pull request #1 from jvzantvoort/feature/add_status.
- Updates.
- Merge branch 'develop' into feature/add_status.
- Minor updates.
- Added status column in list.
