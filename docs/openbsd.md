# Building on OpenBSD

## Prepare your system

Make sure you have `git`, `go` and `gmake` installed on your system.

```
$ doas pkg_add -v git go gmake
```

## Prepare go environment

Make sure your gopath is set:

```
$ export GOPATH=~/go
$ echo "$GOPATH"
$ export PATH="$PATH:$GOPATH/bin"
```

## Build

The `install_unsupported` target works nicely for openbsd. This will install
`gx`, `gx-go` and run `go install -tags nofuse ./cmd/udfs`.

```
$ go get -v -u -d github.com/udfs/go-udfs

$ cd $GOPATH/src/github.com/udfs/go-udfs
$ gmake install_unsupported
```

if everything went well, your udfs binary should be ready at `$GOPATH/bin/udfs`.

```
$ udfs version
```
