# CHAST CLI

This is the cli for CHAST.
All the command are located in `cmd`.
Subcommands are prefixed with the parent command, e.g., `run` => `runRefactoring` = `run refactoring`.

Each command has a `Short` and if necessary a `Long` description.
With `Args`, certain arguments validators can be configured.
`Run` defines the actual command that is being executed. Use this to call the corresponding function in the CHAST API.

The `init` function binds the command to the CLI. Call `AddCommand` on the corresponding parent.

Additional functionality for handing flags, arguments, help functions, etc. can be configured.


## Generating new commands

Generating command are done through `cobra` (https://github.com/spf13/cobra):

1. Download cobra. It is installed into the GO PATH.
```bash
go install github.com/spf13/cobra-cli@latest
```

2. Generate the command
```bash
cobra-cli add <command> # for main commands

cobra-cli add <command> --parent <parent> # for subcommands
```
For detailed instructions see: https://github.com/spf13/cobra-cli/blob/main/README.md

## Helpers
For displaying all `help` sections. Run `./cli_help_check.sh`.


## Viper

Viper (https://github.com/spf13/viper) is used for configuration management.
It is currently unused as no config support is added.
