# Development

## Project layout

* `cli` - subcommands and flags
* `config` - configuration helpers
* `ntypes` - notification types
* `runstat` - run statistics
* `triggers` - notification triggers

## Configuration

Noti has multiple layers of configuration.

* Default settings
* Config file
* CLI flags

Each subsequent layer of configuration overrides the previous layer.

## Templates

You can pass the following template strings to certain notification options.

```
{{.Cmd}}
{{.Args}}
{{.ExitStatus}}
{{.Err}}
{{.Duration}}
{{.State}}
{{.ExpandedAlias}}
```

For example, `noti banner -subtitle 'Exit status: {{.ExitStatus}}`.
