# Development

## Project layout

* `noti/{banner,speech,slack,etc...}` - general purpose notification library
* `noti/cmd/noti` - noti command line utility
* `noti/docs` - documentation for the library and command
* `noti/tests` - tests for the library and command

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
{{.Stdout}}
{{.Stderr}}
{{.ExitCode}}
{{.ExecErr}}
{{.Duration}}
```

For example, `noti banner -subtitle 'Exit code: {{.ExitCode}}`.
