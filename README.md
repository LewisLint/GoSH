# GoSH

GoSH is a small interactive command shell written in Go. It provides a prompt, a few built-in conveniences, and runs external programs from your `PATH`.

## Build and Run

Requires the Go version declared in [`go.mod`](go.mod).

```sh
go build -o gosh .
./gosh
```

You can also start it without building a binary first:

```sh
go run .
```

## Usage

- `cd [directory]` changes the shell's working directory. Running `cd` with no argument, or with `~`, changes to your home directory.
- `exit` quits GoSH.
- Other commands are run as external executables, for example `ls`, `find`, or `git`.
- Single and double quotes group text into an argument, and backslash escapes the next character.
- `*`, `?`, and character-class patterns in arguments are expanded to matching paths.
- Use the up and down arrow keys to navigate command history while running in a terminal.

## Aliases

GoSH loads aliases from `~/.goshrc` when it starts. Put one alias on each line in the form `name=command arguments`:

```text
ll=ls -la
gs=git status
```

An alias replaces the first word of a command. Alias definitions are read from the configuration file; `alias` is not an interactive built-in command.

## Command History

Non-empty commands are saved to `~/.gosh_history` and loaded the next time GoSH starts. In an interactive terminal, the up and down arrows recall saved commands.

## Limitations

GoSH is not a POSIX-compatible shell. It parses one command at a time and does not implement shell operators or expansions, including:

- Command chaining such as `command1 && command2`, `||`, or `;`
- Pipelines and input/output redirection, such as `command | other-command`
- Command substitution such as `$(command)`
- Environment-variable expansion such as `$HOME`
- General tilde expansion: `~` is recognized as the home directory only as the argument to the built-in `cd`

To combine commands, use a shell such as `bash` or `sh` explicitly, for example `bash -c 'command1 && command2'`.