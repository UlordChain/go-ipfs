# Plugins

Since 0.4.11 go-udfs has an experimental plugin system that allows augmenting
the daemons functionality without recompiling.

When an UDFS node is created, it will load plugins from the `$UDFS_PATH/plugins`
directory (by default `~/.udfs/plugins`).

### Plugin types

#### IPLD
IPLD plugins add support for additional formats to `udfs dag` and other IPLD
related commands.

### Supported plugins

| Name | Type |
|------|------|
|  git | IPLD |

#### Installation

##### Linux

1. Build included plugins:
```bash
go-udfs$ make build_plugins
go-udfs$ ls plugin/plugins/*.so
```

3. Copy desired plugins to `$UDFS_PATH/plugins`
```bash
go-udfs$ mkdir -p ~/.udfs/plugins/
go-udfs$ cp plugin/plugins/git.so ~/.udfs/plugins/
go-udfs$ chmod +x ~/.udfs/plugins/git.so # ensure plugin is executable
```

4. Restart daemon if it is running

##### Other

Go currently only supports plugins on Linux, for other platforms you will need
to compile them into UDFS binary.

1. Uncomment plugin entries in `plugin/loader/preload_list`
2. Build udfs
```bash
go-udfs$ make build
```
