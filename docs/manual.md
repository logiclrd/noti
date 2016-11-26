# Manual

## Usage

Basically, put `noti` at the beginning or end of your commands.

```
noti ls
ls; noti
```

If you start a long-running process, then remember about `noti`. That's okay too.
Just press `ctrl-z` to stop the process and restart it like this.

```
$ sleep 10
^Z
zsh: suspended  sleep 10
$ fg; noti
```

## Services

There are different notification services that `noti` supports. Each type has
different flags.

```
noti banner ls
noti slack ls
```

You can use `-h` to see specific help for a particular notification type.

## Configuration

`noti` will look for a file named `.noti.yaml` in your current directory. If it
doesn't find it, then it'll keep going up the directory structure until it does.
Finally, if one doesn't exist anywhere, `noti` will use baked-in defaults.

### File format

The configuration file is standard YAML. This sets `banner` as my default
notification and `exit` as the trigger. You can also see that I've customized
what information the Banner notification displays.

```
---
DefaultNotifications:
    - banner
DefaultTriggers:
    - exit
Banner:
    Title: "Command: {{.Cmd}}"
    Subtitle: "Arguments: {{.Args}}"
    InformativeText: "Duration: {{.Duration}}"
    ContentImage: ""
    SoundName: Ping
```

### Custom info with templates

String fields can use Go template notation to pull certain information about the
utility that was run. Here is some information you can customize your
notification with.

```
Cmd           string
Args          []string
ExitStatus    int
Err           error
Duration      time.Duration
State         string
ExpandedAlias []string
```

## Triggers

Triggers are conditions under which a notification should be sent. Here are some
examples. You can also set these in the config file.

### Exit

Send a notification when a process exits.

```
noti -trigger 'exit' tar -cjf music.tar.bz2 Music/
```

### Match

Send a notification when the utility prints out a certain message.

```
noti -trigger 'match=ready for connections' mysqld_safe
```

### Timeout

Send a notification and kill the process after a certain time.

```
noti -trigger 'timeout=5m' cmd_that_hangs_sometimes
```
