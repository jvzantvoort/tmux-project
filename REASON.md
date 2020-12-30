# Reason

The reason I'm writing this thing.

At the moment I'm working in a project based on tickets. For each
ticket I re-checkout what ever repositories I need. Seems silly
until you work on 4 projects at one and start to lose sight of
things. In my case I wrote a small wrapper in my bash profile that
allows me to resume working on a project by executing:

  resume <projectname>

This solution consists of a few distinct targets:

* ``HOME/.bash/tmux.d/<project>.rc``, the tmux configuration used
  for this.
* ``HOME/.bash/tmux.d/<project>.env``, the bash configuration
  sourced when resuming.
* ``PROJECSTDIR`` the location where projects are checked out.

For the longest time I had only one type of project to work on and
the original client/organization specific solution I wrote in Python
covered this neatly. However others recently came. Different ticket
name format, different archive, etc.. And instead of re-writing my
python thing I instead opted for a golang based approach. Why?
Because I'm shit at golang, it's the Christmas holiday and I have
nothing better to do.

