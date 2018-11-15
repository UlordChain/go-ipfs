# Documentation and Guides

If you’re experiencing an issue with UDFS, **please follow [our issue guide](github-issue-guide.md) when filing an issue!**

Otherwise, check out the following guides to using and developing UDFS:


## General Usage

- [Transferring a File Over UDFS](file-transfer.md)
- [Configuration reference](config.md)
    - [Datastore configuration](datastores.md)
    - [Experimental features](experimental-features.md)
- [Installing command completion](command-completion.md)
- [Mounting UDFS with FUSE](fuse.md)
- [Installing plugins](plugins.md)


## API Support

- [How to Implement an API Client](implement-api-bindings.md)
- [Connecting with Websockets](transports.md) — if you want `js-udfs` nodes in web browsers to connect to your `go-udfs` node, you will need to turn on websocket support in your `go-udfs` node.


## Developing `go-udfs`

- Building on…
    - [Windows](windows.md)
    - [OpenBSD](openbsd.md)
- [Performance Debugging Guidelines](debug-guide.md)
- [Release Checklist](releases.md)


## Other

- [Thanks to all our contributors ❤️](AUTHORS) (We use the `generate-authors.sh` script to regenerate this list.)
- Our [Developer Certificate of Origin (DCO)](developer-certificate-of-origin) — when you sign your commits with `Signed-off-by: <your name>`, you are agreeing to this document.
- [How to file a GitHub Issue](github-issue-guide.md)
