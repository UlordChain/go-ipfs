# go-udfs

![banner](https://udfs.io/udfs/QmVk7srrwahXLNmcDYvyUEJptyoxpndnRa57YJ11L4jV26/udfs.go.png)

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://ipn.io)
[![](https://img.shields.io/badge/project-UDFS-blue.svg?style=flat-square)](http://udfs.io/)
[![](https://img.shields.io/badge/freenode-%23udfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23udfs)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
[![GoDoc](https://godoc.org/github.com/udfs/go-udfs?status.svg)](https://godoc.org/github.com/udfs/go-udfs)
[![Build Status](https://travis-ci.org/udfs/go-udfs.svg?branch=master)](https://travis-ci.org/udfs/go-udfs)

[![Throughput Graph](https://graphs.waffle.io/udfs/go-udfs/throughput.svg)](https://waffle.io/udfs/go-udfs/metrics/throughput)

> UDFS implementation in Go

UDFS is a global, versioned, peer-to-peer filesystem. It combines good ideas from
Git, BitTorrent, Kademlia, SFS, and the Web. It is like a single bittorrent swarm,
exchanging git objects. UDFS provides an interface as simple as the HTTP web, but
with permanence built in. You can also mount the world at /udfs.

For more info see: https://github.com/udfs/udfs.

Please put all issues regarding UDFS _design_ in the
[udfs repo issues](https://github.com/udfs/udfs/issues).
Please put all issues regarding the Go UDFS _implementation_ in [this repo](https://github.com/udfs/go-udfs/issues).

## Table of Contents

- [Security Issues](#security-issues)
- [Install](#install)
  - [System Requirements](#system-requirements)
  - [Install prebuilt packages](#install-prebuilt-packages)
  - [From Linux package managers](#from-linux-package-managers)
  - [Build from Source](#build-from-source)
    - [Install Go](#install-go)
    - [Download and Compile UDFS](#download-and-compile-udfs)
    - [Troubleshooting](#troubleshooting)
  - [Development Dependencies](#development-dependencies)
  - [Updating](#updating)
- [Usage](#usage)
- [Getting Started](#getting-started)
  - [Some things to try](#some-things-to-try)
  - [Docker usage](#docker-usage)
  - [Troubleshooting](#troubleshooting-1)
- [Contributing](#contributing)
  - [Want to hack on UDFS?](#want-to-hack-on-udfs)
  - [Want to read our code?](#want-to-read-our-code)
- [License](#license)

## Security Issues

The UDFS protocol and its implementations are still in heavy development. This means that there may be problems in our protocols, or there may be mistakes in our implementations. And -- though UDFS is not production-ready yet -- many people are already running nodes in their machines. So we take security vulnerabilities very seriously. If you discover a security issue, please bring it to our attention right away!

If you find a vulnerability that may affect live deployments -- for example, by exposing a remote execution exploit -- please send your report privately to security@udfs.io. Please DO NOT file a public issue. The GPG key for security@udfs.io is [4B9665FB 92636D17 7C7A86D3 50AAE8A9 59B13AF3](https://pgp.mit.edu/pks/lookup?op=get&search=0x50AAE8A959B13AF3).

If the issue is a protocol weakness that cannot be immediately exploited or something not yet deployed, just discuss it openly.

## Install

The canonical download instructions for UDFS are over at: http://udfs.io/docs/install/. It is **highly suggested** you follow those instructions if you are not interested in working on UDFS development.

### System Requirements

UDFS can run on most Linux, macOS, and Windows systems. We recommend running it on a machine with at least 2 GB of RAM (it’ll do fine with only one CPU core), but it should run fine with as little as 1 GB of RAM. On systems with less memory, it may not be completely stable.

### Install prebuilt packages

We host prebuilt binaries over at our [distributions page](https://udfs.io/ipns/dist.udfs.io#go-udfs).

From there:
- Click the blue "Download go-udfs" on the right side of the page.
- Open/extract the archive.
- Move `udfs` to your path (`install.sh` can do it for you).

### From Linux package managers

- [Arch Linux](#arch-linux)
- [Nix](#nix)
- [Snap](#snap)

#### Arch Linux

In Arch Linux go-udfs is available as
[go-udfs](https://www.archlinux.org/packages/community/x86_64/go-udfs/) package.

	$ sudo pacman -S go-udfs

Development version of go-udfs is also on AUR under
[go-udfs-git](https://aur.archlinux.org/packages/go-udfs-git/).
You can install it using your favourite AUR Helper or manually from AUR.

### Nix

For Linux and MacOSX you can use the purely functional package manager [Nix](https://nixos.org/nix/):

```
$ nix-env -i udfs
```
You can also install the Package by using it's attribute name, which is also `udfs`.

#### Snap

With snap, in any of the [supported Linux distributions](https://snapcraft.io/docs/core/install):

    $ sudo snap install udfs

### Build from Source

#### Install Go

The build process for udfs requires Go 1.10 or higher. If you don't have it: [Download Go 1.10+](https://golang.org/dl/).


You'll need to add Go's bin directories to your `$PATH` environment variable e.g., by adding these lines to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile`:

```
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
```

(If you run into trouble, see the [Go install instructions](https://golang.org/doc/install)).

#### Download and Compile UDFS

```
$ go get -u -d github.com/udfs/go-udfs

$ cd $GOPATH/src/github.com/udfs/go-udfs
$ make install
```

If you are building on FreeBSD instead of `make install` use `gmake install`.

#### Building on less common systems

If your operating system isn't officially supported, but you still want to try
building udfs anyways (it should work fine in most cases), you can do the
following instead of `make install`:

```
$ make install_unsupported
```

Note: This process may break if [gx](https://github.com/whyrusleeping/gx)
(used for dependency management) or any of its dependencies break as `go get`
will always select the latest code for every dependency, often resulting in
mismatched APIs.

#### Troubleshooting

* Separate [instructions are available for building on Windows](docs/windows.md).
* Also, [instructions for OpenBSD](docs/openbsd.md).
* `git` is required in order for `go get` to fetch all dependencies.
* Package managers often contain out-of-date `golang` packages.
  Ensure that `go version` reports at least 1.10. See above for how to install go.
* If you are interested in development, please install the development
dependencies as well.
* *WARNING: Older versions of OSX FUSE (for Mac OS X) can cause kernel panics when mounting!*
  We strongly recommend you use the [latest version of OSX FUSE](http://osxfuse.github.io/).
  (See https://github.com/udfs/go-udfs/issues/177)
* For more details on setting up FUSE (so that you can mount the filesystem), see the docs folder.
* Shell command completion is available in `misc/completion/udfs-completion.bash`. Read [docs/command-completion.md](docs/command-completion.md) to learn how to install it.
* See the [init examples](https://github.com/udfs/website/tree/master/static/docs/examples/init) for how to connect UDFS to systemd or whatever init system your distro uses.

### Development Dependencies

If you make changes to the protocol buffers, you will need to install the [protoc compiler](https://github.com/google/protobuf).

### Updating

#### Updating using udfs-update
UDFS has an updating tool that can be accessed through `udfs update`. The tool is
not installed alongside UDFS in order to keep that logic independent of the main
codebase. To install `udfs update`, [download it here](https://udfs.io/ipns/dist.udfs.io/#udfs-update).

#### Downloading UDFS builds using UDFS
List the available versions of go-udfs:
```
$ udfs cat /ipns/dist.udfs.io/go-udfs/versions
```

Then, to view available builds for a version from the previous command ($VERSION):
```
$ udfs ls /ipns/dist.udfs.io/go-udfs/$VERSION
```

To download a given build of a version:
```
$ udfs get /ipns/dist.udfs.io/go-udfs/$VERSION/go-udfs_$VERSION_darwin-386.tar.gz # darwin 32-bit build
$ udfs get /ipns/dist.udfs.io/go-udfs/$VERSION/go-udfs_$VERSION_darwin-amd64.tar.gz # darwin 64-bit build
$ udfs get /ipns/dist.udfs.io/go-udfs/$VERSION/go-udfs_$VERSION_freebsd-amd64.tar.gz # freebsd 64-bit build
$ udfs get /ipns/dist.udfs.io/go-udfs/$VERSION/go-udfs_$VERSION_linux-386.tar.gz # linux 32-bit build
$ udfs get /ipns/dist.udfs.io/go-udfs/$VERSION/go-udfs_$VERSION_linux-amd64.tar.gz # linux 64-bit build
$ udfs get /ipns/dist.udfs.io/go-udfs/$VERSION/go-udfs_$VERSION_linux-arm.tar.gz # linux arm build
$ udfs get /ipns/dist.udfs.io/go-udfs/$VERSION/go-udfs_$VERSION_windows-amd64.zip # windows 64-bit build
```

## Usage

```
  udfs - Global p2p merkle-dag filesystem.

  udfs [<flags>] <command> [<arg>] ...

SUBCOMMANDS
  BASIC COMMANDS
    init          Initialize udfs local configuration
    add <path>    Add a file to udfs
    cat <ref>     Show udfs object data
    get <ref>     Download udfs objects
    ls <ref>      List links from an object
    refs <ref>    List hashes of links from an object

  DATA STRUCTURE COMMANDS
    block         Interact with raw blocks in the datastore
    object        Interact with raw dag nodes
    files         Interact with objects as if they were a unix filesystem

  ADVANCED COMMANDS
    daemon        Start a long-running daemon process
    mount         Mount an udfs read-only mountpoint
    resolve       Resolve any type of name
    name          Publish or resolve IPNS names
    dns           Resolve DNS links
    pin           Pin objects to local storage
    repo          Manipulate an UDFS repository

  NETWORK COMMANDS
    id            Show info about udfs peers
    bootstrap     Add or remove bootstrap peers
    swarm         Manage connections to the p2p network
    dht           Query the DHT for values or peers
    ping          Measure the latency of a connection
    diag          Print diagnostics

  TOOL COMMANDS
    config        Manage configuration
    version       Show udfs version information
    update        Download and apply go-udfs updates
    commands      List all available commands

  Use 'udfs <command> --help' to learn more about each command.

  udfs uses a repository in the local file system. By default, the repo is located
  at ~/.udfs. To change the repo location, set the $UDFS_PATH environment variable:

    export UDFS_PATH=/path/to/udfsrepo
```

## Getting Started

See also: http://udfs.io/docs/getting-started/

To start using UDFS, you must first initialize UDFS's config files on your
system, this is done with `udfs init`. See `udfs init --help` for information on
the optional arguments it takes. After initialization is complete, you can use
`udfs mount`, `udfs add` and any of the other commands to explore!

### Some things to try

Basic proof of 'udfs working' locally:

	echo "hello world" > hello
	udfs add hello
	# This should output a hash string that looks something like:
	# QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o
	udfs cat <that hash>


### Docker usage

An UDFS docker image is hosted at [hub.docker.com/r/udfs/go-udfs](https://hub.docker.com/r/udfs/go-udfs/).
To make files visible inside the container you need to mount a host directory
with the `-v` option to docker. Choose a directory that you want to use to
import/export files from UDFS. You should also choose a directory to store
UDFS files that will persist when you restart the container.

    export udfs_staging=</absolute/path/to/somewhere/>
    export udfs_data=</absolute/path/to/somewhere_else/>

Start a container running udfs and expose ports 4001, 5001 and 8080:

    docker run -d --name udfs_host -v $udfs_staging:/export -v $udfs_data:/data/udfs -p 4001:4001 -p 127.0.0.1:8080:8080 -p 127.0.0.1:5001:5001 udfs/go-udfs:latest

Watch the udfs log:

    docker logs -f udfs_host

Wait for udfs to start. udfs is running when you see:

    Gateway (readonly) server
    listening on /ip4/0.0.0.0/tcp/8080

You can now stop watching the log.

Run udfs commands:

    docker exec udfs_host udfs <args...>

For example: connect to peers

    docker exec udfs_host udfs swarm peers

Add files:

    cp -r <something> $udfs_staging
    docker exec udfs_host udfs add -r /export/<something>

Stop the running container:

    docker stop udfs_host

### Troubleshooting

If you have previously installed UDFS before and you are running into
problems getting a newer version to work, try deleting (or backing up somewhere
else) your UDFS config directory (~/.udfs by default) and rerunning `udfs init`.
This will reinitialize the config file to its defaults and clear out the local
datastore of any bad entries.

Please direct general questions and help requests to our
[forum](https://discuss.udfs.io) or our IRC channel (freenode #udfs).

If you believe you've found a bug, check the [issues list](https://github.com/udfs/go-udfs/issues)
and, if you don't see your problem there, either come talk to us on IRC (freenode #udfs) or
file an issue of your own!

## Contributing

We ❤️ all [our contributors](docs/AUTHORS); this project wouldn’t be what it is without you! If you want to help out, please see [Contribute.md](contribute.md).

This repository falls under the UDFS [Code of Conduct](https://github.com/udfs/community/blob/master/code-of-conduct.md).

### Want to hack on UDFS?

[![](https://cdn.rawgit.com/jbenet/contribute-udfs-gif/master/img/contribute.gif)](https://github.com/udfs/community/blob/master/contributing.md)

### Want to read our code?

Some places to get you started. (WIP)

Main file: [cmd/udfs/main.go](https://github.com/udfs/go-udfs/blob/master/cmd/udfs/main.go) <br>
CLI Commands: [core/commands/](https://github.com/udfs/go-udfs/tree/master/core/commands) <br>
Bitswap (the data trading engine): [exchange/bitswap/](https://github.com/udfs/go-udfs/tree/master/exchange/bitswap)

DHT: https://github.com/libp2p/go-libp2p-kad-dht <br>
PubSub: https://github.com/libp2p/go-floodsub <br>
libp2p: https://github.com/libp2p/go-libp2p

## License

MIT
