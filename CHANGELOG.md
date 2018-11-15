# go-udfs changelog

## 0.4.17 2018-07-27

Udfs 0.4.17 is a quick release to fix a major performance regression in bitswap
(mostly affecting go-udfs -> js-udfs transfers). However, while motivated by
this fix, this release contains a few other goodies that will excite some users.

The headline feature in this release is [urlstore][] support. Urlstore is a
generalization of the filestore backend that can fetch file blocks from remote
URLs on-demand instead of storing them in the local datastore.

Additionally, we've added support for extracting inline blocks from CIDs (blocks
inlined into CIDs using the identity hash function). However, go-udfs won't yet
*create* such CIDs so you're unlikely to see any in the wild.

[urlstore]: https://github.com/udfs/go-udfs/blob/master/docs/experimental-features.md#udfs-urlstore

Features:

* URLStore ([udfs/go-udfs#4896](https://github.com/udfs/go-udfs/pull/4896))
* Add trickle-dag support to the urlstore ([udfs/go-udfs#5245](https://github.com/udfs/go-udfs/pull/5245)).
* Allow specifying how the data field in the `object get` is encoded ([udfs/go-udfs#5139](https://github.com/udfs/go-udfs/pull/5139))
* Add a `-U` flag to `files ls` to disable sorting ([udfs/go-udfs#5219](https://github.com/udfs/go-udfs/pull/5219))
* Add an efficient `--size-only` flag to the `repo stat` ([udfs/go-udfs#5010](https://github.com/udfs/go-udfs/pull/5010))
* Inline blocks in CIDs ([udfs/go-udfs#5117](https://github.com/udfs/go-udfs/pull/5117))

Changes/Fixes:

* Make `udfs files ls -l` correctly report the hash and size of files ([udfs/go-udfs#5045](https://github.com/udfs/go-udfs/pull/5045))
* Fix sorting of `files ls` ([udfs/go-udfs#5219](https://github.com/udfs/go-udfs/pull/5219))
* Improve prefetching in `udfs cat` and related commands ([udfs/go-udfs#5162](https://github.com/udfs/go-udfs/pull/5162))
* Better error message when `udfs cp` fails ([udfs/go-udfs#5218](https://github.com/udfs/go-udfs/pull/5218))
* Don't wait for the peer to close it's end of a bitswap stream before considering the block "sent" ([udfs/go-udfs#5258](https://github.com/udfs/go-udfs/pull/5258))
* Fix resolving links in sharded directories via the gateway ([udfs/go-udfs#5271](https://github.com/udfs/go-udfs/pull/5271))
* Fix building when there's a space in the current directory ([udfs/go-udfs#5261](https://github.com/udfs/go-udfs/pull/5261))

Documentation:

* Improve documentation about the bloomfilter config options ([udfs/go-udfs#4924](https://github.com/udfs/go-udfs/pull/4924))

General refactorings and internal bug fixes:

* Remove the `Offset()` method from the DAGReader ([udfs/go-udfs#5190](https://github.com/udfs/go-udfs/pull/5190))
* Fix TestLargeWriteChunks seek behavior ([udfs/go-udfs#5276](https://github.com/udfs/go-udfs/pull/5276))
* Add a build tag to disable dynamic plugins ([udfs/go-udfs#5274](https://github.com/udfs/go-udfs/pull/5274))
* Use FSNode instead of the Protobuf structure in PBDagReader ([udfs/go-udfs#5189](https://github.com/udfs/go-udfs/pull/5189))
* Remove support for non-directory MFS roots ([udfs/go-udfs#5170](https://github.com/udfs/go-udfs/pull/5170))
* Remove `UnixfsNode` from the balanced builder ([udfs/go-udfs#5118](https://github.com/udfs/go-udfs/pull/5118))
* Fix truncating files (internal) when already at the correct size ([udfs/go-udfs#5253](https://github.com/udfs/go-udfs/pull/5253))
* Fix `dagTruncate` (internal) to preserve the node type ([udfs/go-udfs#5216](https://github.com/udfs/go-udfs/pull/5216))
* Add an internal interface for unixfs directories ([udfs/go-udfs#5160](https://github.com/udfs/go-udfs/pull/5160))
* Refactor the CoreAPI path types and interfaces ([udfs/go-udfs#4672](https://github.com/udfs/go-udfs/pull/4672))
* Refactor `precalcNextBuf` in the dag reader ([udfs/go-udfs#5237](https://github.com/udfs/go-udfs/pull/5237))
* Update a bunch of dependencies that haven't been updated for a while ([udfs/go-udfs#5268](https://github.com/udfs/go-udfs/pull/5268))

## 0.4.16 2018-07-13

Udfs 0.4.16 is a fairly small release in terms of changes to the udfs codebase,
but it contains a huge amount of changes and improvements from the libraries we
depend on, notably libp2p.

This release includes small a repo migration to account for some changes to the
DHT. It should only take a second to run but, depending on your configuration,
you may need to run it manually.

You can run a migration by either:

1. Selecting "Yes" when the daemon prompts you to migrate.
2. Running the daemon with the `--migrate=true` flag.
3. Manually [running](https://github.com/udfs/fs-repo-migrations/blob/master/run.md#running-repo-migrations) the migration.

### Libp2p

This version of udfs contains the changes made in libp2p from v5.0.14 through
v6.0.5. In that time, we have made significant changes to the codebase to allow
for easier integration of future transports and modules along with the usual
performance and reliability improvements. You can find many of these
improvements in the libp2p 6.0 [release blog
post](https://udfs.io/blog/39-go-libp2p-6-0-0/).

The primary motivation for this refactor was adding support for network
transports like QUIC that have built-in support for encryption, authentication,
and stream multiplexing. It will also allow us to plug-in new security
transports (like TLS) without hard-coding them.

For example, our [QUIC
transport](https://github.com/libp2p/go-libp2p-quic-transport) currently works,
and can be plugged into libp2p manually (though note that it is still
experimental, as the upstream spec is still in flux). Further work is needed to
make enabling this inside udfs easy and not require recompilation.

On the user-visible side of things, we've improved our dialing logic and
timeouts. We now abort dials to local subnets after 5 seconds and abort all
dials if the TCP handshake takes longer than 5 seconds. This should
significantly improve performance in some cases as we limit the number of
concurrent dials and slow dials to non-responsive peers have been known to clog
the dialer, blocking dials to reachable peers. Importantly, this should improve
DHT performance as it tends to spend a disproportional amount of time connecting
to peers.

We have also made a few noticeable changes to the DHT: we've significantly
improved the chances of finding a value on the DHT, tightened up some of our
validation logic, and fixed some issues that should reduce traffic to nodes
running in dhtclient mode over time.

Of these, the first one will likely see the most impact. In the past, when
putting a value (e.g., an IPNS entry) into the DHT, we'd try to put the value to
K peers (where K for us is 20). However, we'd often fail to connect to many of
these peers so we'd end up putting the value to significantly fewer than K
peers. We now try to put the value to the K peers we can actually connect to.

Finally, we've fixed JavaScript interoperability in go-multiplex, the one stream
muxer that both go-libp2p and js-libp2p implement. This should significantly
improve go-libp2p and js-libp2p interoperability.

### Multiformats

We are also changing the way that people write 'udfs' multiaddrs. Currently,
udfs multiaddrs look something like
`/ip4/104.131.131.82/tcp/4001/udfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ`.
However, calling them 'udfs' multiaddrs is a bit misleading, as this is actually
the multiaddr of a libp2p peer that happens to run udfs. Other protocols built
on libp2p right now still have to use multiaddrs that say 'udfs', even if they
have nothing to do with udfs. Therefore, we are renaming them to 'p2p'
multiaddrs. Moving forward, these addresses will be written as:
`/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ`.

This release adds support for *parsing* both types of addresses (`.../udfs/...`
and `.../p2p/...`) into the same network format, and the network format is
remaining exactly the same. A future release will have the udfs daemon switch to
*printing* out addresses this way once a large enough portion of the network
has upgraded.

N.B., these addresses are *not* related to UDFS *file* names (`/udfs/Qm...`).
Disambiguating the two was yet another motivation to switch the protocol name to
`/p2p/`.

### UDFS

On the udfs side of things, we've started embedding public keys inside IPNS
records and have enabled the Git plugin by default.

Embedding public keys inside IPNS records allows lookups to be faster as we only
need to fetch the record itself (and not the public key separately). It also
fixes an issue where DHT peers wouldn't store a record for a peer if they didn't
have their public key already. Combined with some of the DHT and dialing fixes,
this should improve the performance of IPNS (once a majority of the network
updates).

Second, our public builds now include the Git plugin (in past builds, you could
add it yourself, but doing so was not easy). With this, udfs can ingest and
operate over Git repositories and commit graphs directly. For more information
on this, see [the go-ipld-git repo](https://github.com/udfs/go-ipld-git).

Finally, we've included many smaller bugfixes, refactorings, improved
documentation, and a good bit more. For the full details, see the changelog
below.

## 0.4.16-rc3 2018-07-09
- Bugfixes
  - Fix dht commands when ipns over pubsub is enabled ([udfs/go-udfs#5200](https://github.com/udfs/go-udfs/pull/5200))
  - Fix content routing when ipns over pubsub is enabled ([udfs/go-udfs#5200](https://github.com/udfs/go-udfs/pull/5200))
  - Correctly handle multi-hop dnslink resolution ([udfs/go-udfs#5202](https://github.com/udfs/go-udfs/pull/5202))

## 0.4.16-rc2 2018-07-05
- Bugfixes
  - Fix usage of file name vs path name in adder ([udfs/go-udfs#5167](https://github.com/udfs/go-udfs/pull/5167))
  - Fix `udfs update` working with migrations ([udfs/go-udfs#5194](https://github.com/udfs/go-udfs/pull/5194))
- Documentation
  - Grammar fix in fuse docs ([udfs/go-udfs#5164](https://github.com/udfs/go-udfs/pull/5164))

## 0.4.16-rc1 2018-06-27
- Features
  - Embed public keys inside ipns records, use for validation ([udfs/go-udfs#5079](https://github.com/udfs/go-udfs/pull/5079))
  - Preload git plugin by default ([udfs/go-udfs#4991](https://github.com/udfs/go-udfs/pull/4991))
- Improvements
  - Only resolve dnslinks once in the gateway ([udfs/go-udfs#4977](https://github.com/udfs/go-udfs/pull/4977))
  - Libp2p transport refactor update ([udfs/go-udfs#4817](https://github.com/udfs/go-udfs/pull/4817))
  - Improve swarm connect/disconnect commands ([udfs/go-udfs#5107](https://github.com/udfs/go-udfs/pull/5107))
- Documentation
  - Fix typo of sudo install command ([udfs/go-udfs#5001](https://github.com/udfs/go-udfs/pull/5001))
  - Fix experimental features Table of Contents ([udfs/go-udfs#4976](https://github.com/udfs/go-udfs/pull/4976))
  - Fix link to systemd init scripts in the README ([udfs/go-udfs#4968](https://github.com/udfs/go-udfs/pull/4968))
  - Add package overview comments to coreapi ([udfs/go-udfs#5108](https://github.com/udfs/go-udfs/pull/5108))
  - Add README to docs folder ([udfs/go-udfs#5095](https://github.com/udfs/go-udfs/pull/5095))
  - Add system requirements to README ([udfs/go-udfs#5137](https://github.com/udfs/go-udfs/pull/5137))
- Bugfixes
  - Fix goroutine leak in pin verify ([udfs/go-udfs#5011](https://github.com/udfs/go-udfs/pull/5011))
  - Fix commit string in version ([udfs/go-udfs#4982](https://github.com/udfs/go-udfs/pull/4982))
  - Fix `key rename` command output error ([udfs/go-udfs#4962](https://github.com/udfs/go-udfs/pull/4962))
  - Report error source when failing to construct private network ([udfs/go-udfs#4952](https://github.com/udfs/go-udfs/pull/4952))
  - Fix build on DragonFlyBSD ([udfs/go-udfs#5031](https://github.com/udfs/go-udfs/pull/5031))
  - Fix goroutine leak in dag put ([udfs/go-udfs#5016](https://github.com/udfs/go-udfs/pull/5016))
  - Fix goroutine leaks in refs.go ([udfs/go-udfs#5018](https://github.com/udfs/go-udfs/pull/5018))
  - Fix panic, Don't handle errors with fallthrough ([udfs/go-udfs#5072](https://github.com/udfs/go-udfs/pull/5072))
  - Fix how filestore is hooked up with caching ([udfs/go-udfs#5122](https://github.com/udfs/go-udfs/pull/5122))
  - Add record validation to offline routing ([udfs/go-udfs#5116](https://github.com/udfs/go-udfs/pull/5116))
  - Fix `udfs update` working with migrations ([udfs/go-udfs#5194](https://github.com/udfs/go-udfs/pull/5194))
- General Changes and Refactorings
  - Remove leftover bits of dead code ([udfs/go-udfs#5022](https://github.com/udfs/go-udfs/pull/5022))
  - Remove fuse platform build constraints ([udfs/go-udfs#5033](https://github.com/udfs/go-udfs/pull/5033))
  - Warning when legacy NoSync setting is set ([udfs/go-udfs#5036](https://github.com/udfs/go-udfs/pull/5036))
  - Clean up and refactor namesys module ([udfs/go-udfs#5007](https://github.com/udfs/go-udfs/pull/5007))
  - When raw-leaves are used for empty files use 'Raw' nodes ([udfs/go-udfs#4693](https://github.com/udfs/go-udfs/pull/4693))
  - Update dist_root in build scripts ([udfs/go-udfs#5093](https://github.com/udfs/go-udfs/pull/5093))
  - Integrate `pb.Data` into `FSNode` to avoid duplicating fields ([udfs/go-udfs#5098](https://github.com/udfs/go-udfs/pull/5098))
  - Reduce log level when we can't republish ([udfs/go-udfs#5091](https://github.com/udfs/go-udfs/pull/5091))
  - Extract ipns record logic to go-ipns ([udfs/go-udfs#5124](https://github.com/udfs/go-udfs/pull/5124))
- Testing
  - Collect test times for sharness ([udfs/go-udfs#4959](https://github.com/udfs/go-udfs/pull/4959))
  - Fix sharness iptb connect timeout ([udfs/go-udfs#4966](https://github.com/udfs/go-udfs/pull/4966))
  - Add more timeouts to the jenkins pipeline ([udfs/go-udfs#4958](https://github.com/udfs/go-udfs/pull/4958))
  - Use go 1.10 on jenkins ([udfs/go-udfs#5009](https://github.com/udfs/go-udfs/pull/5009))
  - Speed up multinode sharness test ([udfs/go-udfs#4967](https://github.com/udfs/go-udfs/pull/4967))
  - Print out iptb logs on iptb test failure (for debugging CI) ([udfs/go-udfs#5069](https://github.com/udfs/go-udfs/pull/5069))
  - Disable the MacOS tests in jenkins ([udfs/go-udfs#5119](https://github.com/udfs/go-udfs/pull/5119))
  - Make republisher test robust against timing issues ([udfs/go-udfs#5125](https://github.com/udfs/go-udfs/pull/5125))
  - Archive sharness trash dirs in jenkins ([udfs/go-udfs#5071](https://github.com/udfs/go-udfs/pull/5071))
  - Fixup DHT sharness tests ([udfs/go-udfs#5114](https://github.com/udfs/go-udfs/pull/5114))
- Dependencies
  - Update go-ipld-git to fix mergetag resolving ([udfs/go-udfs#4988](https://github.com/udfs/go-udfs/pull/4988))
  - Fix duplicate /x/sys imports ([udfs/go-udfs#5068](https://github.com/udfs/go-udfs/pull/5068))
  - Update stream multiplexers ([udfs/go-udfs#5075](https://github.com/udfs/go-udfs/pull/5075))
  - Update dependencies: go-log, sys, go-crypto ([udfs/go-udfs#5100](https://github.com/udfs/go-udfs/pull/5100))
  - Explicitly import go-multiaddr-dns in config/bootstrap_peers ([udfs/go-udfs#5144](https://github.com/udfs/go-udfs/pull/5144))
  - Gx update with dht and dialing improvements ([udfs/go-udfs#5158](https://github.com/udfs/go-udfs/pull/5158))

## 0.4.15 2018-05-09

This release is significantly smaller than the last as much of the work on
improving our datastores, and other libraries libp2p has yet to be merged.
However, it still includes many welcome improvements.

As with 0.4.12 and 0.4.14 (0.4.13 was a patch), this release has a negative
diff-stat. Unfortunately, much of this code isn't actually going away but at
least it's being moved out into separate repositories.

Much of the work that made it into this release is under the hood. We've cleaned
up some code, extracted several packages into their own repositories, and made
some long neglected optimizations (e.g., handling of sharded directories).
Additionally, this release includes a bunch of tests for our CLI commands that
should help us avoid some of the issues we've seen in the past few releases.

More visibly, thanks to @djdv's efforts, this release includes some significant
Windows improvements (with more on the way). Specifically, this release includes
better handling of repo lockfiles (no more `udfs repo fsck`), stdin command-line
support, and, last but not least, UDFS no longer writes random files with scary
garbage in the drive root. To read more about future windows improvements, take
a look at this [blog post](https://blog.udfs.io/36-a-look-at-windows/).

To better support low-power devices, we've added a low-power config profile.
This can be enabled when initializing a repo by running `udfs init` with the
`--profile=lowpower` flag or later by running `udfs config profile apply lowpower`.

Finally, with this release we have begun distributing self-contained source
archives of go-udfs and its dependencies. This should be a welcome improvement
for both packagers and those living in countries with harmonized internet
access.

- Features
  - Add options for record count and timeout for resolving DHT paths ([udfs/go-udfs#4733](https://github.com/udfs/go-udfs/pull/4733))
  - Add low power init profile ([udfs/go-udfs#4154](https://github.com/udfs/go-udfs/pull/4154))
  - Add Opentracing plugin support ([udfs/go-udfs#4506](https://github.com/udfs/go-udfs/pull/4506))
  - Add make target to build source tarballs ([udfs/go-udfs#4920](https://github.com/udfs/go-udfs/pull/4920))

- Improvements
  - Add BlockedFetched/Added/Removed events to Blockservice ([udfs/go-udfs#4649](https://github.com/udfs/go-udfs/pull/4649))
  - Improve performance of HAMT code ([udfs/go-udfs#4889](https://github.com/udfs/go-udfs/pull/4889))
  - Avoid unnecessarily resolving child nodes when listing a sharded directory ([udfs/go-udfs#4884](https://github.com/udfs/go-udfs/pull/4884))
  - Tar writer now supports sharded udfs directories ([udfs/go-udfs#4873](https://github.com/udfs/go-udfs/pull/4873))
  - Infer type from CID when possible in `udfs ls` ([udfs/go-udfs#4890](https://github.com/udfs/go-udfs/pull/4890))
  - Deduplicate keys in GetMany ([udfs/go-udfs#4888](https://github.com/udfs/go-udfs/pull/4888))

- Documentation
  - Fix spelling of retrieval ([udfs/go-udfs#4819](https://github.com/udfs/go-udfs/pull/4819))
  - Update broken links ([udfs/go-udfs#4798](https://github.com/udfs/go-udfs/pull/4798))
  - Remove roadmap.md ([udfs/go-udfs#4834](https://github.com/udfs/go-udfs/pull/4834))
  - Remove link to UDFS paper in contribute.md ([udfs/go-udfs#4812](https://github.com/udfs/go-udfs/pull/4812))
  - Fix broken todo link in readme.md ([udfs/go-udfs#4865](https://github.com/udfs/go-udfs/pull/4865))
  - Document ipns pubsub ([udfs/go-udfs#4903](https://github.com/udfs/go-udfs/pull/4903))
  - Fix missing profile docs ([udfs/go-udfs#4846](https://github.com/udfs/go-udfs/pull/4846))
  - Fix a few typos ([udfs/go-udfs#4835](https://github.com/udfs/go-udfs/pull/4835))
  - Fix typo in fsrepo error message ([udfs/go-udfs#4933](https://github.com/udfs/go-udfs/pull/4933))
  - Remove go-udfs version from issue template ([udfs/go-udfs#4943](https://github.com/udfs/go-udfs/pull/4943))
  - Add docs for --profile=lowpower ([udfs/go-udfs#4970](https://github.com/udfs/go-udfs/pull/4970))
  - Improve Windows build documentation ([udfs/go-udfs#4691](https://github.com/udfs/go-udfs/pull/4691))

- Bugfixes
  - Check CIDs in base case when diffing nodes ([udfs/go-udfs#4767](https://github.com/udfs/go-udfs/pull/4767))
  - Support for CIDv1 with custom mhtype in `udfs block put` ([udfs/go-udfs#4563](https://github.com/udfs/go-udfs/pull/4563))
  - Clean path in DagArchive ([udfs/go-udfs#4743](https://github.com/udfs/go-udfs/pull/4743))
  - Set the prefix for MFS root in `udfs add --hash-only` ([udfs/go-udfs#4755](https://github.com/udfs/go-udfs/pull/4755))
  - Fix get output path ([udfs/go-udfs#4809](https://github.com/udfs/go-udfs/pull/4809))
  - Fix incorrect Read calls ([udfs/go-udfs#4792](https://github.com/udfs/go-udfs/pull/4792))
  - Use prefix in bootstrapWritePeers ([udfs/go-udfs#4832](https://github.com/udfs/go-udfs/pull/4832))
  - Fix mfs Directory.Path not working ([udfs/go-udfs#4844](https://github.com/udfs/go-udfs/pull/4844))
  - Remove header in `udfs stats bw` if not polling ([udfs/go-udfs#4856](https://github.com/udfs/go-udfs/pull/4856))
  - Match Go's GOPATH defaults behaviour in build scripts ([udfs/go-udfs#4678](https://github.com/udfs/go-udfs/pull/4678))
  - Fix default-net profile not reverting bootstrap config ([udfs/go-udfs#4845](https://github.com/udfs/go-udfs/pull/4845))
  - Fix excess goroutines in bitswap caused by insecure CIDs ([udfs/go-udfs#4946](https://github.com/udfs/go-udfs/pull/4946))

- General Changes and Refactorings
  - Refactor trickle DAG builder ([udfs/go-udfs#4730](https://github.com/udfs/go-udfs/pull/4730))
  - Split the coreapi interface into multiple files ([udfs/go-udfs#4802](https://github.com/udfs/go-udfs/pull/4802))
  - Make `udfs init` command use new cmds lib ([udfs/go-udfs#4732](https://github.com/udfs/go-udfs/pull/4732))
  - Extract thirdparty/tar package ([udfs/go-udfs#4857](https://github.com/udfs/go-udfs/pull/4857))
  - Reduce log level when for disconnected peers to info ([udfs/go-udfs#4811](https://github.com/udfs/go-udfs/pull/4811))
  - Only visit nodes in EnumerateChildrenAsync when asked ([udfs/go-udfs#4885](https://github.com/udfs/go-udfs/pull/4885))
  - Refactor coreapi options ([udfs/go-udfs#4807](https://github.com/udfs/go-udfs/pull/4807))
  - Fix error style for most errors ([udfs/go-udfs#4829](https://github.com/udfs/go-udfs/pull/4829))
  - Ensure `--help` always works, even with /dev/null stdin ([udfs/go-udfs#4849](https://github.com/udfs/go-udfs/pull/4849))
  - Deduplicate AddNodeLinkClean into AddNodeLink ([udfs/go-udfs#4940](https://github.com/udfs/go-udfs/pull/4940))
  - Remove some dead code ([udfs/go-udfs#4833](https://github.com/udfs/go-udfs/pull/4833))
  - Remove unused imports ([udfs/go-udfs#4955](https://github.com/udfs/go-udfs/pull/4955))
  - Fix go vet warnings ([udfs/go-udfs#4859](https://github.com/udfs/go-udfs/pull/4859))

- Testing
  - Generate JUnit test reports for sharness tests ([udfs/go-udfs#4530](https://github.com/udfs/go-udfs/pull/4530))
  - Fix t0063-daemon-init.sh by adding test profile to daemon ([udfs/go-udfs#4816](https://github.com/udfs/go-udfs/pull/4816))
  - Remove circular dependencies in merkledag package tests ([udfs/go-udfs#4704](https://github.com/udfs/go-udfs/pull/4704))
  - Check that all the commands fail when passed a bad flag ([udfs/go-udfs#4848](https://github.com/udfs/go-udfs/pull/4848))
  - Allow for some small margin of code coverage dropping on commit ([udfs/go-udfs#4867](https://github.com/udfs/go-udfs/pull/4867))
  - Add confirmation to archive-branches script ([udfs/go-udfs#4797](https://github.com/udfs/go-udfs/pull/4797))

- Dependencies
  - Update lock package ([udfs/go-udfs#4855](https://github.com/udfs/go-udfs/pull/4855))
  - Update to latest go-datastore. Remove thirdparty/datastore2 ([udfs/go-udfs#4742](https://github.com/udfs/go-udfs/pull/4742))
  - Extract fs lock into go-fs-lock ([udfs/go-udfs#4631](https://github.com/udfs/go-udfs/pull/4631))
  - Extract: exchange/interface.go, blocks/blocksutil, exchange/offline ([udfs/go-udfs#4912](https://github.com/udfs/go-udfs/pull/4912))
  - Remove unused lock dep ([udfs/go-udfs#4971](https://github.com/udfs/go-udfs/pull/4971))
  - Update iptb ([udfs/go-udfs#4965](https://github.com/udfs/go-udfs/pull/4965))
  - Update go-udfs-cmds to fix stdin on windows ([udfs/go-udfs#4975](https://github.com/udfs/go-udfs/pull/4975))
  - Update go-ds-flatfs to fix windows corruption issue ([udfs/go-udfs#4872](https://github.com/udfs/go-udfs/pull/4872))

## 0.4.14 2018-03-22

Udfs 0.4.14 is a big release with a large number of improvements and bugfixes.
It is also the first release of 2018, and our first release in over three
months. The release took longer than expected due to our refactoring and
extracting of our commands library. This refactor had two stages.  The first
round of the refactor disentangled the commands code from core udfs code,
allowing us to move it out into a [separate
repository](https://github.com/udfs/go-udfs-cmds).  The code was previously
very entangled with the go-udfs codebase and not usable for other projects.
The second round of the refactor had the goal of fixing several major issues
around streaming outputs, progress bars, and error handling. It also paved the
way for us to more easily provide an API over other transports, such as
websockets and unix domain sockets.  It took a while to flush out all the kinks
on such a massive change.  We're pretty sure we've got most of them, but if you
notice anything weird, please let us know.

Beyond that, we've added a new experimental way to use IPNS. With the new
pubsub IPNS resolver and publisher, you can subscribe to updates of an IPNS
entry, and the owner can publish out changes in real time. With this, IPNS can
become nearly instantaneous. To make use of this, simply start your udfs daemon
with the `--enable-namesys-pubsub` option, and all IPNS resolution and
publishing will use pubsub. Note that resolving an IPNS name via pubsub without
someone publishing it via pubsub will result in a fallback to using the DHT.
Please give this a try and let us know how it goes!

Memory and CPU usage should see a noticeable improvement in this release. We
have spent considerable time fixing excess memory usage throughout the codebase
and down into libp2p. Fixes in peer tracking, bitswap allocation, pinning, and
many other places have brought down both peak and average memory usage. An
upgraded hashing library, base58 encoding library, and improved allocation
patterns all contribute to overall lower CPU usage across the board. See the
full changelist below for more memory and CPU usage improvements.

This release also brings the beginning of the udfs 'Core API'. Once finished,
the Core API will be the primary way to interact with go-udfs using go. Both
embedded nodes and nodes accessed over the http API will have the same
interface. Stay tuned for future updates and documentation.

These are only a sampling of the changes that made it into this release, the
full list (almost 100 PRs!) is below.

Finally, I'd like to thank everyone who contributed to this release, whether
you're just contributing a typo fix or driving new features. We are really
grateful to everyone who has spent their their time pushing udfs forward.

SECURITY NOTE:

This release of udfs disallows the usage of insecure hash functions and
lengths. Udfs does not create these insecure objects for any purpose, but it
did allow manually creating them and fetching them from other peers. If you
currently have objects using insecure hashes in your local udfs repo, please
remove them before updating.

#### Changes from rc2 to rc3
- Fix bug in stdin argument parsing ([udfs/go-udfs#4827](https://github.com/udfs/go-udfs/pull/4827))
- Revert commands back to sending a single response ([udfs/go-udfs#4822](https://github.com/udfs/go-udfs/pull/4822))

#### Changes from rc1 to rc2
- Fix issue in udfs get caused by go1.10 changes ([udfs/go-udfs#4790](https://github.com/udfs/go-udfs/pull/4790))

- Features
  - Pubsub IPNS Publisher and Resolver (experimental) ([udfs/go-udfs#4047](https://github.com/udfs/go-udfs/pull/4047))
  - Implement coreapi Dag interface ([udfs/go-udfs#4471](https://github.com/udfs/go-udfs/pull/4471))
  - Add --offset flag to udfs cat ([udfs/go-udfs#4538](https://github.com/udfs/go-udfs/pull/4538))
  - Command to apply config profile after init ([udfs/go-udfs#4195](https://github.com/udfs/go-udfs/pull/4195))
  - Implement coreapi Name and Key interfaces ([udfs/go-udfs#4477](https://github.com/udfs/go-udfs/pull/4477))
  - Add --length flag to udfs cat ([udfs/go-udfs#4553](https://github.com/udfs/go-udfs/pull/4553))
  - Implement coreapi Object interface ([udfs/go-udfs#4492](https://github.com/udfs/go-udfs/pull/4492))
  - Implement coreapi Block interface ([udfs/go-udfs#4548](https://github.com/udfs/go-udfs/pull/4548))
  - Implement coreapi Pin interface ([udfs/go-udfs#4575](https://github.com/udfs/go-udfs/pull/4575))
  - Add a --with-local flag to udfs files stat ([udfs/go-udfs#4638](https://github.com/udfs/go-udfs/pull/4638))
  - Disallow usage of blocks with insecure hashes ([udfs/go-udfs#4751](https://github.com/udfs/go-udfs/pull/4751))
- Improvements
  - Add uuid to event logs ([udfs/go-udfs#4392](https://github.com/udfs/go-udfs/pull/4392))
  - Add --quiet flag to object put ([udfs/go-udfs#4411](https://github.com/udfs/go-udfs/pull/4411))
  - Pinning memory improvements and fixes ([udfs/go-udfs#4451](https://github.com/udfs/go-udfs/pull/4451))
  - Update WebUI version ([udfs/go-udfs#4449](https://github.com/udfs/go-udfs/pull/4449))
  - Check strong and weak ETag validator ([udfs/go-udfs#3983](https://github.com/udfs/go-udfs/pull/3983))
  - Improve and refactor FD limit handling ([udfs/go-udfs#3801](https://github.com/udfs/go-udfs/pull/3801))
  - Support linking to non-dagpb objects in udfs object patch ([udfs/go-udfs#4460](https://github.com/udfs/go-udfs/pull/4460))
  - Improve allocation patterns of slices in bitswap ([udfs/go-udfs#4458](https://github.com/udfs/go-udfs/pull/4458))
  - Secio handshake now happens synchronously ([libp2p/go-libp2p-secio#25](https://github.com/libp2p/go-libp2p-secio/pull/25))
  - Don't block closing connections on pending writes ([libp2p/go-msgio#7](https://github.com/libp2p/go-msgio/pull/7))
  - Improve memory usage of multiaddr parsing ([multiformats/go-multiaddr#56](https://github.com/multiformats/go-multiaddr/pull/56))
  - Don't lock up 256KiB buffers when adding small files ([udfs/go-udfs#4508](https://github.com/udfs/go-udfs/pull/4508))
  - Clear out memory after reads from the dagreader ([udfs/go-udfs#4525](https://github.com/udfs/go-udfs/pull/4525))
  - Improve error handling in udfs ping ([udfs/go-udfs#4546](https://github.com/udfs/go-udfs/pull/4546))
  - Allow install.sh to be run without being the script dir ([udfs/go-udfs#4547](https://github.com/udfs/go-udfs/pull/4547))
  - Much faster base58 encoding ([libp2p/go-libp2p-peer#24](https://github.com/libp2p/go-libp2p-peer/pull/24))
  - Use faster sha256 and blake2b libs ([multiformats/go-multihash#63](https://github.com/multiformats/go-multihash/pull/63))
  - Greatly improve peerstore memory usage ([libp2p/go-libp2p-peerstore#22](https://github.com/libp2p/go-libp2p-peerstore/pull/22))
  - Improve dht memory usage and peer tracking ([libp2p/go-libp2p-kad-dht#111](https://github.com/libp2p/go-libp2p-kad-dht/pull/111))
  - New libp2p metrics lib with lower overhead ([libp2p/go-libp2p-metrics#8](https://github.com/libp2p/go-libp2p-metrics/pull/8))
  - Fix memory leak that occurred when dialing many peers ([libp2p/go-libp2p-swarm#51](https://github.com/libp2p/go-libp2p-swarm/pull/51))
  - Wire up new dag interfaces to make sessions easier ([udfs/go-udfs#4641](https://github.com/udfs/go-udfs/pull/4641))
- Documentation
  - Correct StorageMax config description ([udfs/go-udfs#4388](https://github.com/udfs/go-udfs/pull/4388))
  - Add how to download UDFS with UDFS doc ([udfs/go-udfs#4390](https://github.com/udfs/go-udfs/pull/4390))
  - Document gx release checklist item ([udfs/go-udfs#4480](https://github.com/udfs/go-udfs/pull/4480))
  - Add some documentation to CoreAPI ([udfs/go-udfs#4493](https://github.com/udfs/go-udfs/pull/4493))
  - Add interop tests to the release checklist ([udfs/go-udfs#4501](https://github.com/udfs/go-udfs/pull/4501))
  - Add badgerds to experimental-features ToC ([udfs/go-udfs#4537](https://github.com/udfs/go-udfs/pull/4537))
  - Fix typos and inconsistencies in commands documentation ([udfs/go-udfs#4552](https://github.com/udfs/go-udfs/pull/4552))
  - Add a document to help troubleshoot data transfers ([udfs/go-udfs#4332](https://github.com/udfs/go-udfs/pull/4332))
  - Add a bunch of documentation on public interfaces ([udfs/go-udfs#4599](https://github.com/udfs/go-udfs/pull/4599))
  - Expand the issue template and remove the severity field ([udfs/go-udfs#4624](https://github.com/udfs/go-udfs/pull/4624))
  - Add godocs for importers module ([udfs/go-udfs#4640](https://github.com/udfs/go-udfs/pull/4640))
  - Document make targets ([udfs/go-udfs#4653](https://github.com/udfs/go-udfs/pull/4653))
  - Add godocs for merkledag module ([udfs/go-udfs#4665](https://github.com/udfs/go-udfs/pull/4665))
  - Add godocs for unixfs module ([udfs/go-udfs#4664](https://github.com/udfs/go-udfs/pull/4664))
  - Add sharding to experimental features list ([udfs/go-udfs#4569](https://github.com/udfs/go-udfs/pull/4569))
  - Add godocs for routing module ([udfs/go-udfs#4676](https://github.com/udfs/go-udfs/pull/4676))
  - Add godocs for path module ([udfs/go-udfs#4689](https://github.com/udfs/go-udfs/pull/4689))
  - Add godocs for pin module ([udfs/go-udfs#4696](https://github.com/udfs/go-udfs/pull/4696))
  - Update link to filestore experimental status ([udfs/go-udfs#4557](https://github.com/udfs/go-udfs/pull/4557))
- Bugfixes
  - Remove trailing slash in udfs get paths, fixes #3729 ([udfs/go-udfs#4365](https://github.com/udfs/go-udfs/pull/4365))
  - fix deadlock in bitswap sessions ([udfs/go-udfs#4407](https://github.com/udfs/go-udfs/pull/4407))
  - Fix two race conditions (and possibly go routine leaks) in commands ([udfs/go-udfs#4406](https://github.com/udfs/go-udfs/pull/4406))
  - Fix output delay in udfs pubsub sub ([udfs/go-udfs#4402](https://github.com/udfs/go-udfs/pull/4402))
  - Use correct context in AddWithContext ([udfs/go-udfs#4433](https://github.com/udfs/go-udfs/pull/4433))
  - Fix various IPNS republisher issues ([udfs/go-udfs#4440](https://github.com/udfs/go-udfs/pull/4440))
  - Fix error handling in commands add and get ([udfs/go-udfs#4454](https://github.com/udfs/go-udfs/pull/4454))
  - Fix hamt (sharding) delete issue ([udfs/go-udfs#4398](https://github.com/udfs/go-udfs/pull/4398))
  - More correctly check for reuseport support ([libp2p/go-reuseport#40](https://github.com/libp2p/go-reuseport/pull/40))
  - Fix goroutine leak in websockets transport ([libp2p/go-ws-transport#21](https://github.com/libp2p/go-ws-transport/pull/21))
  - Update badgerds to fix i386 windows build ([udfs/go-udfs#4464](https://github.com/udfs/go-udfs/pull/4464))
  - Only construct bitswap event loggable if necessary ([udfs/go-udfs#4533](https://github.com/udfs/go-udfs/pull/4533))
  - Ensure that flush on the mfs root flushes its directory ([udfs/go-udfs#4509](https://github.com/udfs/go-udfs/pull/4509))
  - Fix deferred unlock of pin lock in AddR ([udfs/go-udfs#4562](https://github.com/udfs/go-udfs/pull/4562))
  - Fix iOS builds ([udfs/go-udfs#4610](https://github.com/udfs/go-udfs/pull/4610))
  - Calling repo gc now frees up space with badgerds ([udfs/go-udfs#4578](https://github.com/udfs/go-udfs/pull/4578))
  - Fix leak in bitswap sessions shutdown ([udfs/go-udfs#4658](https://github.com/udfs/go-udfs/pull/4658))
  - Fix make on windows ([udfs/go-udfs#4682](https://github.com/udfs/go-udfs/pull/4682))
  - Ignore invalid key files in keystore directory ([udfs/go-udfs#4700](https://github.com/udfs/go-udfs/pull/4700))
- General Changes and Refactorings
  - Extract and refactor commands library ([udfs/go-udfs#3856](https://github.com/udfs/go-udfs/pull/3856))
  - Remove all instances of `Default(false)` ([udfs/go-udfs#4042](https://github.com/udfs/go-udfs/pull/4042))
  - Build for all supported platforms when testing ([udfs/go-udfs#4445](https://github.com/udfs/go-udfs/pull/4445))
  - Refine gateway and namesys logging ([udfs/go-udfs#4428](https://github.com/udfs/go-udfs/pull/4428))
  - Demote bitswap error to an info ([udfs/go-udfs#4472](https://github.com/udfs/go-udfs/pull/4472))
  - Extract posinfo package to github.com/udfs/go-udfs-posinfo ([udfs/go-udfs#4669](https://github.com/udfs/go-udfs/pull/4669))
  - Move signature verification to ipns validator ([udfs/go-udfs#4628](https://github.com/udfs/go-udfs/pull/4628))
  - Extract importers/chunk module as go-udfs-chunker ([udfs/go-udfs#4661](https://github.com/udfs/go-udfs/pull/4661))
  - Extract go-detect-race from Godeps ([udfs/go-udfs#4686](https://github.com/udfs/go-udfs/pull/4686))
  - Extract flags, delay, ds-help ([udfs/go-udfs#4685](https://github.com/udfs/go-udfs/pull/4685))
  - Extract routing package to go-udfs-routing ([udfs/go-udfs#4703](https://github.com/udfs/go-udfs/pull/4703))
  - Extract blocks/blockstore package to go-udfs-blockstore ([udfs/go-udfs#4707](https://github.com/udfs/go-udfs/pull/4707))
  - Add exchange.SessionExchange interface for exchanges that support sessions ([udfs/go-udfs#4709](https://github.com/udfs/go-udfs/pull/4709))
  - Extract thirdparty/pq to go-udfs-pq ([udfs/go-udfs#4711](https://github.com/udfs/go-udfs/pull/4711))
  - Separate "path" from "path/resolver" ([udfs/go-udfs#4713](https://github.com/udfs/go-udfs/pull/4713))
- Testing
  - Increase verbosity of t0088-repo-stat-symlink.sh test ([udfs/go-udfs#4434](https://github.com/udfs/go-udfs/pull/4434))
  - Make repo size test pass deterministically ([udfs/go-udfs#4443](https://github.com/udfs/go-udfs/pull/4443))
  - Always set UDFS_PATH in test-lib.sh ([udfs/go-udfs#4469](https://github.com/udfs/go-udfs/pull/4469))
  - Fix sharness docker ([udfs/go-udfs#4489](https://github.com/udfs/go-udfs/pull/4489))
  - Fix loops in sharness tests to fail the test if the inner command fails ([udfs/go-udfs#4482](https://github.com/udfs/go-udfs/pull/4482))
  - Improve bitswap tests, fix race conditions ([udfs/go-udfs#4499](https://github.com/udfs/go-udfs/pull/4499))
  - Fix circleci cache directory list ([udfs/go-udfs#4564](https://github.com/udfs/go-udfs/pull/4564))
  - Only run the build test on test_go_expensive ([udfs/go-udfs#4645](https://github.com/udfs/go-udfs/pull/4645))
  - Fix go test on Windows ([udfs/go-udfs#4632](https://github.com/udfs/go-udfs/pull/4632))
  - Fix some tests on FreeBSD ([udfs/go-udfs#4662](https://github.com/udfs/go-udfs/pull/4662))

## 0.4.13 2017-11-16

Udfs 0.4.13 is a patch release that fixes two high priority issues that were
discovered in the 0.4.12 release.

Bugfixes:
  - Fix periodic bitswap deadlock ([udfs/go-udfs#4386](https://github.com/udfs/go-udfs/pull/4386))
  - Fix badgerds crash on startup ([udfs/go-udfs#4384](https://github.com/udfs/go-udfs/pull/4384))


## 0.4.12 2017-11-09

Udfs 0.4.12 brings with it many important fixes for the huge spike in network
size we've seen this past month. These changes include the Connection Manager,
faster batching in `udfs add`, libp2p fixes that reduce CPU usage, and a bunch
of new documentation.

The most critical change is the 'Connection Manager': it allows an udfs node to
maintain a limited set of connections to other peers in the network. By default
(and with no config changes required by the user), udfs nodes will now try to
maintain between 600 and 900 open connections. These limits are still likely
higher than needed, and future releases may lower the default recommendation,
but for now we want to make changes gradually. The rationale for this selection
of numbers is as follows:

- The DHT routing table for a large network may rise to around 400 peers
- Bitswap connections tend to be separate from the DHT
- PubSub connections also generally are another distinct set of peers
  (including js-udfs nodes)

Because of this, we selected 600 as a 'LowWater' number, and 900 as a
'HighWater' number to avoid having to clear out connections too frequently.
You can configure different numbers as you see fit via the `Swarm.ConnMgr`
field in your udfs config file. See
[here](https://github.com/udfs/go-udfs/blob/master/docs/config.md#connmgr) for
more details.

Disk utilization during `udfs add` has been optimized for large files by doing
batch writes in parallel. Previously, when adding a large file, users might have
noticed that the add progressed by about 8MB at a time, with brief pauses in between.
This was caused by quickly filling up the batch, then blocking while it was
writing to disk. We now write to disk in the background while continuing to add
the remainder of the file.

Other changes in this release have noticeably reduced memory consumption and CPU
usage. This was done by optimising some frequently called functions in libp2p
that were expensive in terms of both CPU usage and memory allocations. We also
lowered the yamux accept buffer sizes which were raised over a year ago to
combat a separate bug that has since been fixed.

And finally, thank you to everyone who filed bugs, tested out the release candidates,
filed pull requests, and contributed in any other way to this release!

- Features
  - Implement Connection Manager ([udfs/go-udfs#4288](https://github.com/udfs/go-udfs/pull/4288))
  - Support multiple files in dag put ([udfs/go-udfs#4254](https://github.com/udfs/go-udfs/pull/4254))
  - Add 'raw' support to the dag put command ([udfs/go-udfs#4285](https://github.com/udfs/go-udfs/pull/4285))
- Improvements
  - Parallelize dag batch flushing ([udfs/go-udfs#4296](https://github.com/udfs/go-udfs/pull/4296))
  - Update go-peerstream to improve CPU usage ([udfs/go-udfs#4323](https://github.com/udfs/go-udfs/pull/4323))
  - Add full support for CidV1 in Files API and Dag Modifier ([udfs/go-udfs#4026](https://github.com/udfs/go-udfs/pull/4026))
  - Lower yamux accept buffer size ([udfs/go-udfs#4326](https://github.com/udfs/go-udfs/pull/4326))
  - Optimise `udfs pin update` command ([udfs/go-udfs#4348](https://github.com/udfs/go-udfs/pull/4348))
- Documentation
  - Add some docs on plugins ([udfs/go-udfs#4255](https://github.com/udfs/go-udfs/pull/4255))
  - Add more info about private network bootstrap ([udfs/go-udfs#4270](https://github.com/udfs/go-udfs/pull/4270))
  - Add more info about `udfs add` chunker option ([udfs/go-udfs#4306](https://github.com/udfs/go-udfs/pull/4306))
  - Remove cruft in readme and mention discourse forum ([udfs/go-udfs#4345](https://github.com/udfs/go-udfs/pull/4345))
  - Add note about updating before reporting issues ([udfs/go-udfs#4361](https://github.com/udfs/go-udfs/pull/4361))
- Bugfixes
  - Fix FreeBSD build issues ([udfs/go-udfs#4275](https://github.com/udfs/go-udfs/pull/4275))
  - Don't crash when Datastore.StorageMax is not defined ([udfs/go-udfs#4246](https://github.com/udfs/go-udfs/pull/4246))
  - Do not call 'Connect' on NewStream in bitswap ([udfs/go-udfs#4317](https://github.com/udfs/go-udfs/pull/4317))
  - Filter out "" from active peers in bitswap sessions ([udfs/go-udfs#4316](https://github.com/udfs/go-udfs/pull/4316))
  - Fix "seeker can't seek" on specific files ([udfs/go-udfs#4320](https://github.com/udfs/go-udfs/pull/4320))
  - Do not set "gecos" field in Dockerfile ([udfs/go-udfs#4331](https://github.com/udfs/go-udfs/pull/4331))
  - Handle sym links in when calculating repo size ([udfs/go-udfs#4305](https://github.com/udfs/go-udfs/pull/4305))
- General Changes and Refactorings
  - Fix indent in sharness tests ([udfs/go-udfs#4212](https://github.com/udfs/go-udfs/pull/4212))
  - Remove supernode routing ([udfs/go-udfs#4302](https://github.com/udfs/go-udfs/pull/4302))
  - Extract go-udfs-addr ([udfs/go-udfs#4340](https://github.com/udfs/go-udfs/pull/4340))
  - Remove dead code and config files ([udfs/go-udfs#4357](https://github.com/udfs/go-udfs/pull/4357))
  - Update badgerds to 1.0 ([udfs/go-udfs#4327](https://github.com/udfs/go-udfs/pull/4327))
  - Wrap help descriptions under 80 chars ([udfs/go-udfs#4121](https://github.com/udfs/go-udfs/pull/4121))
- Testing
  - Make sharness t0180-p2p less racy ([udfs/go-udfs#4310](https://github.com/udfs/go-udfs/pull/4310))


### 0.4.11 2017-09-14

Udfs 0.4.11 is a larger release that brings many long-awaited features and
performance improvements. These include new datastore options, more efficient
bitswap transfers, greatly improved resource consumption, circuit relay
support, ipld plugins, and more! Take a look at the full changelog below for a
detailed list of every change.

The udfs datastore has, until now, been a combination of leveldb and a custom
git-like storage backend called 'flatfs'. This works well enough for the
average user, but different udfs usecases demand different backend
configurations. To address this, we have changed the configuration file format
for datastores to be a modular way of specifying exactly how you want the
datastore to be structured. You will now be able to configure udfs to use
flatfs, leveldb, badger, an in-memory datastore, and more to suit your needs.
See the new [datastore
documentation](https://github.com/udfs/go-udfs/blob/master/docs/datastores.md)
for more information.

Bitswap received some much needed attention during this release cycle. The
concept of 'Bitswap Sessions' allows bitswap to associate requests for
different blocks to the same underlying session, and from that infer better
ways of requesting that data. In more concrete terms, parts of the udfs
codebase that take advantage of sessions (currently, only `udfs pin add`) will
cause much less extra traffic than before. This is done by making optimistic
guesses about which nodes might be providing given blocks and not sending
wantlist updates to every connected bitswap partner, as well as searching the
DHT for providers less frequently. In future releases we will migrate over more
udfs commands to take advantage of bitswap sessions. As nodes update to this
and future versions, expect to see idle bandwidth usage on the udfs network
go down noticeably.

The never ending effort to reduce resource consumption had a few important
updates this release. First, the bitswap sessions changes discussed above will
help with improving bandwidth usage. Aside from that there are two important
libp2p updates that improved things significantly. The first was a fix to a bug
in the dial limiter code that was causing it to not limit outgoing dials
correctly. This resulted in udfs running out of file descriptors very
frequently (as well as incurring a decent amount of excess outgoing bandwidth),
this has now been fixed. Users who previously received "too many open files"
errors should see this much less often in 0.4.11. The second change was a
memory leak in the DHT that was identified and fixed. Streams being tracked in
a map in the DHT weren't being cleaned up after the peer disconnected leading
to the multiplexer session not being cleaned up properly. This issue has been
resolved, and now memory usage appears to be stable over time. There is still a
lot of work to be done improving memory usage, but we feel this is a solid
victory.

It is often said that NAT traversal is the hardest problem in peer to peer
technology, we tend to agree with this. In an effort to provide a more
ubiquitous p2p mesh, we have implemented a relay mechanism that allows willing
peers to relay traffic for other peers who might not otherwise be able to
communicate with each other.  This feature is still pretty early, and currently
users have to manually connect through a relay. The next step in this endeavour
is automatic relaying, and research for this is currently in progress. We
expect that when it lands, it will improve the perceived performance of udfs by
spending less time attempting connections to hard to reach nodes. A short guide
on using the circuit relay feature can be found
[here](https://github.com/udfs/go-udfs/blob/master/docs/experimental-features.md#circuit-relay).

The last feature we want to highlight (but by no means the last feature in this
release) is our new plugin system. There are many different workflows and
usecases that udfs should be able to support, but not everyone wants to be able
to use every feature. We could simply merge in all these features, but that
causes problems for several reasons: first off, the size of the udfs binary
starts to get very large very quickly. Second, each of these different pieces
needs to be maintained and updated independently, which would cause significant
churn in the codebase. To address this, we have come up with a system that
allows users to install plugins to the vanilla udfs daemon that augment its
capabilities. The first of these plugins are a [git
plugin](https://github.com/udfs/go-udfs/blob/master/plugin/plugins/git/git.go)
that allows udfs to natively address git objects and an [ethereum
plugin](https://github.com/udfs/go-ipld-eth) that lets udfs ingest and operate
on all ethereum blockchain data. Soon to come are plugins for the bitcoin and
zcash data formats. In the future, we will be adding plugins for other things
like datastore backends and specialized libp2p network transports.
You can read more on this topic in [Plugin docs](docs/plugins.md)

In order to simplify its integration with fs-repo-migrations, we've switched
the udfs/go-udfs docker image from a musl base to a glibc base. For most users
this will not be noticeable, but if you've been building your own images based
off this image, you'll have to update your dockerfile. We recommend a
multi-stage dockerfile, where the build stage is based off of a regular Debian or
other glibc-based image, and the assembly stage is based off of the udfs/go-udfs
image, and you copy build artifacts from the build stage to the assembly
stage. Note, if you are using the docker image and see a deprecation message,
please update your usage. We will stop supporting the old method of starting
the dockerfile in the next release.

Finally, I would like to thank all of our contributors, users, supporters, and
friends for helping us along the way. Udfs would not be where it is without
you.


- Features
  - Add `--pin` option to `udfs dag put` ([udfs/go-udfs#4004](https://github.com/udfs/go-udfs/pull/4004))
  - Add `--pin` option to `udfs object put` ([udfs/go-udfs#4095](https://github.com/udfs/go-udfs/pull/4095))
  - Implement `--profile` option on `udfs init` ([udfs/go-udfs#4001](https://github.com/udfs/go-udfs/pull/4001))
  - Add CID Codecs to `udfs block put` ([udfs/go-udfs#4022](https://github.com/udfs/go-udfs/pull/4022))
  - Bitswap sessions ([udfs/go-udfs#3867](https://github.com/udfs/go-udfs/pull/3867))
  - Create plugin API and loader, add ipld-git plugin ([udfs/go-udfs#4033](https://github.com/udfs/go-udfs/pull/4033))
  - Make announced swarm addresses configurable ([udfs/go-udfs#3948](https://github.com/udfs/go-udfs/pull/3948))
  - Reprovider strategies ([udfs/go-udfs#4113](https://github.com/udfs/go-udfs/pull/4113))
  - Circuit Relay integration ([udfs/go-udfs#4091](https://github.com/udfs/go-udfs/pull/4091))
  - More configurable datastore configs ([udfs/go-udfs#3575](https://github.com/udfs/go-udfs/pull/3575))
  - Add experimental support for badger datastore ([udfs/go-udfs#4007](https://github.com/udfs/go-udfs/pull/4007))
- Improvements
  - Add better support for Raw Nodes in MFS and elsewhere ([udfs/go-udfs#3996](https://github.com/udfs/go-udfs/pull/3996))
  - Added file size to response of `udfs add` command ([udfs/go-udfs#4082](https://github.com/udfs/go-udfs/pull/4082))
  - Add /dnsaddr bootstrap nodes ([udfs/go-udfs#4127](https://github.com/udfs/go-udfs/pull/4127))
  - Do not publish public keys extractable from ID ([udfs/go-udfs#4020](https://github.com/udfs/go-udfs/pull/4020))
- Documentation
  - Adding documentation that PubSub Sub can be encoded. ([udfs/go-udfs#3909](https://github.com/udfs/go-udfs/pull/3909))
  - Add Comms items from js-udfs, including blog ([udfs/go-udfs#3936](https://github.com/udfs/go-udfs/pull/3936))
  - Add Developer Certificate of Origin ([udfs/go-udfs#4006](https://github.com/udfs/go-udfs/pull/4006))
  - Add `transports.md` document ([udfs/go-udfs#4034](https://github.com/udfs/go-udfs/pull/4034))
  - Add `experimental-features.md` document ([udfs/go-udfs#4036](https://github.com/udfs/go-udfs/pull/4036))
  - Update release docs ([udfs/go-udfs#4165](https://github.com/udfs/go-udfs/pull/4165))
  - Add documentation for datastore configs ([udfs/go-udfs#4223](https://github.com/udfs/go-udfs/pull/4223))
  - General update and clean-up of docs ([udfs/go-udfs#4222](https://github.com/udfs/go-udfs/pull/4222))
- Bugfixes
  - Fix shutdown check in t0023 ([udfs/go-udfs#3969](https://github.com/udfs/go-udfs/pull/3969))
  - Fix pinning of unixfs sharded directories ([udfs/go-udfs#3975](https://github.com/udfs/go-udfs/pull/3975))
  - Show escaped url in gateway 404 message ([udfs/go-udfs#4005](https://github.com/udfs/go-udfs/pull/4005))
  - Fix early opening of bitswap message sender ([udfs/go-udfs#4069](https://github.com/udfs/go-udfs/pull/4069))
  - Fix determination of 'root' node in dag put ([udfs/go-udfs#4072](https://github.com/udfs/go-udfs/pull/4072))
  - Fix bad multipart message panic in gateway ([udfs/go-udfs#4053](https://github.com/udfs/go-udfs/pull/4053))
  - Add blocks to the blockstore before returning them from blockservice sessions ([udfs/go-udfs#4169](https://github.com/udfs/go-udfs/pull/4169))
  - Various fixes for /udfs fuse code ([udfs/go-udfs#4194](https://github.com/udfs/go-udfs/pull/4194))
  - Fix memory leak in dht stream tracking ([udfs/go-udfs#4251](https://github.com/udfs/go-udfs/pull/4251))
- General Changes and Refactorings
  - Require go 1.8 ([udfs/go-udfs#4044](https://github.com/udfs/go-udfs/pull/4044))
  - Change UDFS to use the new pluggable Block to IPLD decoding framework. ([udfs/go-udfs#4060](https://github.com/udfs/go-udfs/pull/4060))
  - Remove tour command from udfs ([udfs/go-udfs#4123](https://github.com/udfs/go-udfs/pull/4123))
  - Add support for Go 1.9 ([udfs/go-udfs#4156](https://github.com/udfs/go-udfs/pull/4156))
  - Remove some dead code ([udfs/go-udfs#4204](https://github.com/udfs/go-udfs/pull/4204))
  - Switch docker image from musl to glibc ([udfs/go-udfs#4219](https://github.com/udfs/go-udfs/pull/4219))

### 0.4.10 - 2017-06-27

Udfs 0.4.10 is a patch release that contains several exciting new features,
bugfixes and general improvements. Including new commands, easier corruption
recovery, and a generally cleaner codebase.

The `udfs pin` command has two new subcommands, `verify` and `update`. `udfs
pin verify` is used to scan the repo for pinned object graphs and check their
integrity. Any issues are reported back with helpful error text to make error
recovery simpler.  This subcommand was added to help recover from datastore
corruptions, particularly if using the experimental filestore and accidentally
deleting tracked files.
`udfs pin update` was added to make the task of keeping a large, frequently
changing object graph pinned. Previously users had to call `udfs pin rm` on the
old pin, and `udfs pin add` on the new one. The 'new' `udfs pin add` call would
be very expensive as it would need to verify the entirety of the graph again.
The `udfs pin update` command takes shortcuts, portions of the graph that were
covered under the old pin are assumed to be fine, and the command skips
checking them.

Next up, we have finally implemented an `udfs shutdown` command so users can
shut down their udfs daemons via the API. This is especially useful on
platforms that make it difficult to control processes (Android, for example),
and is also useful when needing to shut down a node remotely and you do not
have access to the machine itself.

`udfs add` has gained a new flag; the `--hash` flag allows you to select which
hash function to use and we have given it the ability to select `blake2b-256`.
This pushes us one step closer to shifting over to using blake2b as the
default. Blake2b is significantly faster than sha2-256, and also is conjectured
to provide superior security.

We have also finally implemented a very early (and experimental) `udfs p2p`.
This command and its subcommands will allow you to open up arbitrary streams to
other udfs peers through libp2p. The interfaces are a little bit clunky right
now, but shouldn't get in the way of anyone wanting to try building a fully
peer to peer application on top of udfs and libp2p. For more info on this
command, to ask questions, or to provide feedback, head over to the [feedback
issue](https://github.com/udfs/go-udfs/issues/3994) for the command.

A few other subcommands and flags were added around the API, as well as many
other requested improvements. See below for the full list of changes.


- Features
  - Add support for specifying the hash function in `udfs add` ([udfs/go-udfs#3919](https://github.com/udfs/go-udfs/pull/3919))
  - Implement `udfs key {rm, rename}` ([udfs/go-udfs#3892](https://github.com/udfs/go-udfs/pull/3892))
  - Implement `udfs shutdown` command ([udfs/go-udfs#3884](https://github.com/udfs/go-udfs/pull/3884))
  - Implement `udfs pin update` ([udfs/go-udfs#3846](https://github.com/udfs/go-udfs/pull/3846))
  - Implement `udfs pin verify` ([udfs/go-udfs#3843](https://github.com/udfs/go-udfs/pull/3843))
  - Implemented experimental p2p commands ([udfs/go-udfs#3943](https://github.com/udfs/go-udfs/pull/3943))
- Improvements
  - Add MaxStorage field to output of "repo stat" ([udfs/go-udfs#3915](https://github.com/udfs/go-udfs/pull/3915))
  - Add Suborigin header to gateway responses ([udfs/go-udfs#3914](https://github.com/udfs/go-udfs/pull/3914))
  - Add "--file-order" option to "filestore ls" and "verify" ([udfs/go-udfs#3938](https://github.com/udfs/go-udfs/pull/3938))
  - Allow selecting ipns keys by Peer ID ([udfs/go-udfs#3882](https://github.com/udfs/go-udfs/pull/3882))
  - Don't redirect to trailing slash in gateway for `go get` ([udfs/go-udfs#3963](https://github.com/udfs/go-udfs/pull/3963))
  - Add 'udfs dht findprovs --num-providers' to allow choosing number of providers to find ([udfs/go-udfs#3966](https://github.com/udfs/go-udfs/pull/3966))
  - Make sure all keystore keys get republished ([udfs/go-udfs#3951](https://github.com/udfs/go-udfs/pull/3951))
- Documentation
  - Adding documentation on PubSub encodings ([udfs/go-udfs#3909](https://github.com/udfs/go-udfs/pull/3909))
  - Change 'neccessary' to 'necessary' ([udfs/go-udfs#3941](https://github.com/udfs/go-udfs/pull/3941))
  - README.md: add Nix to the linux package managers ([udfs/go-udfs#3939](https://github.com/udfs/go-udfs/pull/3939))
  - More verbose errors in filestore ([udfs/go-udfs#3964](https://github.com/udfs/go-udfs/pull/3964))
- Bugfixes
  - Fix typo in message when file size check fails ([udfs/go-udfs#3895](https://github.com/udfs/go-udfs/pull/3895))
  - Clean up bitswap ledgers when disconnecting ([udfs/go-udfs#3437](https://github.com/udfs/go-udfs/pull/3437))
  - Make odds of 'process added after close' panic less likely ([udfs/go-udfs#3940](https://github.com/udfs/go-udfs/pull/3940))
- General Changes and Refactorings
  - Remove 'udfs diag net' from codebase ([udfs/go-udfs#3916](https://github.com/udfs/go-udfs/pull/3916))
  - Update to dht code with provide announce option ([udfs/go-udfs#3928](https://github.com/udfs/go-udfs/pull/3928))
  - Apply the megacheck code vetting tool ([udfs/go-udfs#3949](https://github.com/udfs/go-udfs/pull/3949))
  - Expose port 8081 in docker container for /ws listener ([udfs/go-udfs#3954](https://github.com/udfs/go-udfs/pull/3954))

### 0.4.9 - 2017-04-30

Udfs 0.4.9 is a maintenance release that contains several useful bugfixes and
improvements. Notably, `udfs add` has gained the ability to select which CID
version will be output. The common udfs hash that looks like this:
`QmRjNgF2mRLDT8AzCPsQbw1EYF2hDTFgfUmJokJPhCApYP` is a multihash. Multihashes
allow us to specify the hashing algorithm that was used to verify the data, but
it doesn't give us any indication of what format that data might be. To address
that issue, we are adding another couple of bytes to the prefix that will allow us
to indicate the format of the data referenced by the hash. This new format is
called a Content ID, or CID for short. The previous bare multihashes will still
be fully supported throughout the entire application as CID version 0. The new
format with the type information will be CID version 1. To give an example,
the content referenced by the hash above is "Hello Udfs!". That same content,
in the same format (dag-protobuf) using CIDv1 is
`zb2rhkgXZVkT2xvDiuUsJENPSbWJy7fdYnsboLBzzEjjZMRoG`.

CIDv1 hashes are supported in udfs versions back to 0.4.5. Nodes running 0.4.4
and older will not be able to load content via CIDv1 and we recommend that they
update to a newer version.

There are many other use cases for CIDs. Plugins can be written to
allow udfs to natively address content from any other merkletree based system,
such as git, bitcoin, zcash and ethereum -- a few systems we've already started work on.

Aside from the CID flag, there were many other changes as noted below:

- Features
  - Add support for using CidV1 in 'udfs add' ([udfs/go-udfs#3743](https://github.com/udfs/go-udfs/pull/3743))
- Improvements
  - Use CID as an ETag strong validator ([udfs/go-udfs#3869](https://github.com/udfs/go-udfs/pull/3869))
  - Update go-multihash with keccak and bitcoin hashes ([udfs/go-udfs#3833](https://github.com/udfs/go-udfs/pull/3833))
  - Update go-is-domain to contain new gTLD ([udfs/go-udfs#3873](https://github.com/udfs/go-udfs/pull/3873))
  - Periodically flush cached directories during udfs add ([udfs/go-udfs#3888](https://github.com/udfs/go-udfs/pull/3888))
  - improved gateway directory listing for sharded nodes ([udfs/go-udfs#3897](https://github.com/udfs/go-udfs/pull/3897))
- Documentation
  - Change issue template to use Severity instead of Priority ([udfs/go-udfs#3834](https://github.com/udfs/go-udfs/pull/3834))
  - Fix link to commit hook script in contribute.md ([udfs/go-udfs#3863](https://github.com/udfs/go-udfs/pull/3863))
  - Fix install_unsupported for openbsd, add docs ([udfs/go-udfs#3880](https://github.com/udfs/go-udfs/pull/3880))
- Bugfixes
  - Fix wanlist typo in prometheus metric name ([udfs/go-udfs#3841](https://github.com/udfs/go-udfs/pull/3841))
  - Fix `make install` not using ldflags for git hash ([udfs/go-udfs#3838](https://github.com/udfs/go-udfs/pull/3838))
  - Fix `make install` not installing dependencies ([udfs/go-udfs#3848](https://github.com/udfs/go-udfs/pull/3848))
  - Fix erroneous Cache-Control: immutable on dir listings ([udfs/go-udfs#3870](https://github.com/udfs/go-udfs/pull/3870))
  - Fix bitswap accounting of 'BytesSent' in ledger ([udfs/go-udfs#3876](https://github.com/udfs/go-udfs/pull/3876))
  - Fix gateway handling of sharded directories ([udfs/go-udfs#3889](https://github.com/udfs/go-udfs/pull/3889))
  - Fix sharding memory growth, and fix resolver for unixfs paths ([udfs/go-udfs#3890](https://github.com/udfs/go-udfs/pull/3890))
- General Changes and Refactorings
  - Use ctx var consistently in daemon.go ([udfs/go-udfs#3864](https://github.com/udfs/go-udfs/pull/3864))
  - Handle 404 correctly in dist_get tool ([udfs/go-udfs#3879](https://github.com/udfs/go-udfs/pull/3879))
- Testing
  - Fix go fuse tests ([udfs/go-udfs#3840](https://github.com/udfs/go-udfs/pull/3840))

### 0.4.8 - 2017-03-29

Udfs 0.4.8 brings with it several improvements, bugfixes, documentation
improvements, and the long awaited directory sharding code.

Currently, when too many items are added into a unixfs directory, the object
gets too large and you may experience issues. To pervent this problem, and
generally make working really large directories more efficient, we have
implemented a HAMT structure for unixfs. To enable this feature, run:
```
udfs config --json Experimental.ShardingEnabled true
```

And restart your daemon if it was running.

Note: With this setting enabled, the hashes of any newly added directories will
be different than they previously were, as the new code will use the sharded
HAMT structure for all directories. Also, nodes running udfs 0.4.7 and earlier
will not be able to access directories created with this option.

That said, please do give it a try, let us know how it goes, and then take a
look at all the other cool things added in 0.4.8 below.

- Features
	- Implement unixfs directory sharding ([udfs/go-udfs#3042](https://github.com/udfs/go-udfs/pull/3042))
	- Add DisableNatPortMap option ([udfs/go-udfs#3798](https://github.com/udfs/go-udfs/pull/3798))
	- Basic Filestore utilty commands ([udfs/go-udfs#3653](https://github.com/udfs/go-udfs/pull/3653))
- Improvements
	- More Robust GC ([udfs/go-udfs#3712](https://github.com/udfs/go-udfs/pull/3712))
	- Automatically fix permissions for docker volumes ([udfs/go-udfs#3744](https://github.com/udfs/go-udfs/pull/3744))
	- Core API refinements and efficiency improvements ([udfs/go-udfs#3493](https://github.com/udfs/go-udfs/pull/3493))
	- Improve IsPinned() lookups for indirect pins ([udfs/go-udfs#3809](https://github.com/udfs/go-udfs/pull/3809))
- Documentation
	- Improve 'name' and 'key' helptexts ([udfs/go-udfs#3806](https://github.com/udfs/go-udfs/pull/3806))
	- Update link to paper in dev.md ([udfs/go-udfs#3812](https://github.com/udfs/go-udfs/pull/3812))
	- Add test to enforce helptext on commands ([udfs/go-udfs#2648](https://github.com/udfs/go-udfs/pull/2648))
- Bugfixes
	- Remove bloom filter check on Put call in blockstore ([udfs/go-udfs#3782](https://github.com/udfs/go-udfs/pull/3782))
	- Re-add the GOPATH checking functionality ([udfs/go-udfs#3787](https://github.com/udfs/go-udfs/pull/3787))
	- Use fsrepo.IsInitialized to test for initialization ([udfs/go-udfs#3805](https://github.com/udfs/go-udfs/pull/3805))
	- Return 404 Not Found for failed path resolutions ([udfs/go-udfs#3777](https://github.com/udfs/go-udfs/pull/3777))
	- Fix 'dist\_get' failing without failing ([udfs/go-udfs#3818](https://github.com/udfs/go-udfs/pull/3818))
	- Update iptb with fix for t0130 hanging issue ([udfs/go-udfs#3823](https://github.com/udfs/go-udfs/pull/3823))
	- fix hidden file detection on windows ([udfs/go-udfs#3829](https://github.com/udfs/go-udfs/pull/3829))
- General Changes and Refactorings
	- Fix multiple govet warnings ([udfs/go-udfs#3824](https://github.com/udfs/go-udfs/pull/3824))
	- Make Golint happy in the blocks submodule ([udfs/go-udfs#3827](https://github.com/udfs/go-udfs/pull/3827))
- Testing
	- Enable codeclimate for automated linting and vetting ([udfs/go-udfs#3821](https://github.com/udfs/go-udfs/pull/3821))
	- Fix EOF test failure with Multipart.Read ([udfs/go-udfs#3804](https://github.com/udfs/go-udfs/pull/3804))

### 0.4.7 - 2017-03-15

Udfs 0.4.7 contains several exciting new features!
First off, The long awaited filestore feature has been merged, allowing users
the option to not have udfs store chunked copies of added files in the
blockstore, pushing to burden of ensuring those files are not changed to the
user. The filestore feature is currently still experimental, and must be
enabled in your config with:
```
udfs config --json Experimental.FilestoreEnabled true
```
before it can be used. Please see [this issue](https://github.com/udfs/go-udfs/issues/3397#issuecomment-284337564) for more details.

Next up, We have merged initial support for udfs 'Private Networks'. This
feature allows users to run udfs in a mode that will only connect to other
peers in the private network. This feature, like the filestore is being
released experimentally, but if you're interested please try it out.
Instructions for setting it up can be found
[here](https://github.com/udfs/go-udfs/issues/3397#issuecomment-284341649).

This release also enables support for the 'mplex' stream muxer by default. This
stream multiplexing protocol was available previously via the
`--enable-mplex-experiment` daemon flag, but has now graduated to being 'less
experimental' and no longer requires the flag to use it.

Aside from those, we have a good number of bugfixes, perf improvements and new
tests. Heres a list of highlights:

- Features
	- Implement basic filestore 'no-copy' functionality ([udfs/go-udfs#3629](https://github.com/udfs/go-udfs/pull/3629))
	- Add support for private udfs networks ([udfs/go-udfs#3697](https://github.com/udfs/go-udfs/pull/3697))
	- Enable 'mplex' stream muxer by default ([udfs/go-udfs#3725](https://github.com/udfs/go-udfs/pull/3725))
	- Add `--quieter` option to `udfs add` ([udfs/go-udfs#3770](https://github.com/udfs/go-udfs/pull/3770))
	- Report progress during `pin add` via `--progress` ([udfs/go-udfs#3671](https://github.com/udfs/go-udfs/pull/3671))
- Improvements
	- Allow `udfs get` to handle content added with raw leaves option ([udfs/go-udfs#3757](https://github.com/udfs/go-udfs/pull/3757))
	- Fix accuracy of progress bar on `udfs get` ([udfs/go-udfs#3758](https://github.com/udfs/go-udfs/pull/3758))
	- Limit number of objects in batches to prevent too many fds issue ([udfs/go-udfs#3756](https://github.com/udfs/go-udfs/pull/3756))
	- Add more info to bitswap stat ([udfs/go-udfs#3635](https://github.com/udfs/go-udfs/pull/3635))
	- Add multiple performance metrics ([udfs/go-udfs#3615](https://github.com/udfs/go-udfs/pull/3615))
	- Make `dist_get` fall back to other downloaders if one fails ([udfs/go-udfs#3692](https://github.com/udfs/go-udfs/pull/3692))
- Documentation
	- Add Arch Linux install instructions to readme ([udfs/go-udfs#3742](https://github.com/udfs/go-udfs/pull/3742))
	- Improve release checklist document ([udfs/go-udfs#3717](https://github.com/udfs/go-udfs/pull/3717))
- Bugfixes
	- Fix drive root parsing on windows ([udfs/go-udfs#3328](https://github.com/udfs/go-udfs/pull/3328))
	- Fix panic in udfs get when passing no parameters to API ([udfs/go-udfs#3768](https://github.com/udfs/go-udfs/pull/3768))
	- Fix breakage of `udfs pin add` api output ([udfs/go-udfs#3760](https://github.com/udfs/go-udfs/pull/3760))
	- Fix issue in DHT queries that was causing poor record replication ([udfs/go-udfs#3748](https://github.com/udfs/go-udfs/pull/3748))
	- Fix `udfs mount` crashing if no name was published before ([udfs/go-udfs#3728](https://github.com/udfs/go-udfs/pull/3728))
	- Add `self` key to the `udfs key list` listing ([udfs/go-udfs#3734](https://github.com/udfs/go-udfs/pull/3734))
	- Fix panic when shutting down `udfs daemon` pre gateway setup ([udfs/go-udfs#3723](https://github.com/udfs/go-udfs/pull/3723))
- General Changes and Refactorings
	- Refactor `EnumerateChildren` to avoid need for bestEffort parameter ([udfs/go-udfs#3700](https://github.com/udfs/go-udfs/pull/3700))
	- Update fuse dependency, fixing several issues ([udfs/go-udfs#3727](https://github.com/udfs/go-udfs/pull/3727))
	- Add `install_unsupported` makefile target for 'exotic' systems ([udfs/go-udfs#3719](https://github.com/udfs/go-udfs/pull/3719))
	- Deprecate implicit daemon argument in Dockerfile ([udfs/go-udfs#3685](https://github.com/udfs/go-udfs/pull/3685))
- Testing
	- Add test to ensure helptext is under 80 columns wide ([udfs/go-udfs#3774](https://github.com/udfs/go-udfs/pull/3774))
	- Add unit tests for auto migration code ([udfs/go-udfs#3618](https://github.com/udfs/go-udfs/pull/3618))
	- Fix iptb stop issue in sharness tests  ([udfs/go-udfs#3714](https://github.com/udfs/go-udfs/pull/3714))


### 0.4.6 - 2017-02-21

Udfs 0.4.6 contains several bugfixes related to migrations and also contains a
few other improvements to other parts of the codebase. Notably:

- The default config will now contain some ipv6 addresses for bootstrap nodes.
- `udfs pin add` should be faster and consume less memory.
- Pinning thousands of files no longer causes superlinear usage of storage space.

- Improvements
	- Make pinset sharding deterministic ([udfs/go-udfs#3640](https://github.com/udfs/go-udfs/pull/3640))
	- Update to go-multihash with blake2 ([udfs/go-udfs#3649](https://github.com/udfs/go-udfs/pull/3649))
	- Pass cids instead of nodes around in EnumerateChildrenAsync ([udfs/go-udfs#3598](https://github.com/udfs/go-udfs/pull/3598))
	- Add /ip6 bootstrap nodes ([udfs/go-udfs#3523](https://github.com/udfs/go-udfs/pull/3523))
	- Add sub-object support to `dag get` command ([udfs/go-udfs#3687](https://github.com/udfs/go-udfs/pull/3687))
	- Add half-closed streams support to multiplex experiment ([udfs/go-udfs#3695](https://github.com/udfs/go-udfs/pull/3695))
- Documentation
	- Add the snap installation instructions ([udfs/go-udfs#3663](https://github.com/udfs/go-udfs/pull/3663))
	- Add closed PRs, Issues throughput ([udfs/go-udfs#3602](https://github.com/udfs/go-udfs/pull/3602))
- Bugfixes
	- Fix auto-migration on docker nodes ([udfs/go-udfs#3698](https://github.com/udfs/go-udfs/pull/3698))
	- Update flatfs to v1.1.2, fixing directory fd issue ([udfs/go-udfs#3711](https://github.com/udfs/go-udfs/pull/3711))
- General Changes and Refactorings
	- Remove `FindProviders` from routing mocks ([udfs/go-udfs#3617](https://github.com/udfs/go-udfs/pull/3617))
	- Use Marshalers instead of PostRun to process `block rm` output ([udfs/go-udfs#3708](https://github.com/udfs/go-udfs/pull/3708))
- Testing
	- Makefile rework and sharness test coverage ([udfs/go-udfs#3504](https://github.com/udfs/go-udfs/pull/3504))
	- Print out all daemon stderr files when iptb stop fails ([udfs/go-udfs#3701](https://github.com/udfs/go-udfs/pull/3701))
	- Add tests for recursively pinning a dag ([udfs/go-udfs#3691](https://github.com/udfs/go-udfs/pull/3691))
	- Fix lack of commit hash during build ([udfs/go-udfs#3705](https://github.com/udfs/go-udfs/pull/3705))

### 0.4.5 - 2017-02-11

#### Changes from rc3 to rc4
- Update to fixed webui. ([udfs/go-udfs#3669](https://github.com/udfs/go-udfs/pull/3669))

#### Changes from rc2 to rc3
- Fix handling of null arrays in cbor ipld objects.  ([udfs/go-udfs#3666](https://github.com/udfs/go-udfs/pull/3666))
- Add env var to enable yamux debug logging.  ([udfs/go-udfs#3668](https://github.com/udfs/go-udfs/pull/3668))
- Fix libc check during auto-migrations.  ([udfs/go-udfs#3665](https://github.com/udfs/go-udfs/pull/3665))

#### Changes from rc1 to rc2
- Fixed json output of ipld objects in `udfs dag get` ([udfs/go-udfs#3655](https://github.com/udfs/go-udfs/pull/3655))

#### Changes since 0.4.4

- Notable changes
	- IPLD and CIDs
	  - Rework go-udfs to use Content IDs  ([udfs/go-udfs#3187](https://github.com/udfs/go-udfs/pull/3187))  ([udfs/go-udfs#3290](https://github.com/udfs/go-udfs/pull/3290))
	  - Turn merkledag.Node into an interface ([udfs/go-udfs#3301](https://github.com/udfs/go-udfs/pull/3301))
	  - Implement cbor ipld nodes  ([udfs/go-udfs#3325](https://github.com/udfs/go-udfs/pull/3325))
	  - Allow cid format selection in block put command  ([udfs/go-udfs#3324](https://github.com/udfs/go-udfs/pull/3324))  ([udfs/go-udfs#3483](https://github.com/udfs/go-udfs/pull/3483))
	  - Bitswap protocol extension to handle cids  ([udfs/go-udfs#3297](https://github.com/udfs/go-udfs/pull/3297))
	  - Add dag get to read-only api  ([udfs/go-udfs#3499](https://github.com/udfs/go-udfs/pull/3499))
	- Raw Nodes
	  - Implement 'Raw Node' node type for addressing raw data  ([udfs/go-udfs#3307](https://github.com/udfs/go-udfs/pull/3307))
	  - Optimize DagService GetLinks for Raw Nodes.  ([udfs/go-udfs#3351](https://github.com/udfs/go-udfs/pull/3351))
	- Experimental PubSub
	  - Added a very basic pubsub implementation  ([udfs/go-udfs#3202](https://github.com/udfs/go-udfs/pull/3202))
	- Core API
	  - gateway: use core api for serving GET/HEAD/POST  ([udfs/go-udfs#3244](https://github.com/udfs/go-udfs/pull/3244))

- Improvements
	- Disable auto-gc check in 'udfs cat'  ([udfs/go-udfs#3100](https://github.com/udfs/go-udfs/pull/3100))
	- Add `bitswap ledger` command  ([udfs/go-udfs#2852](https://github.com/udfs/go-udfs/pull/2852))
	- Add `udfs block rm` command.  ([udfs/go-udfs#2962](https://github.com/udfs/go-udfs/pull/2962))
	- Add config option to disable bandwidth metrics   ([udfs/go-udfs#3381](https://github.com/udfs/go-udfs/pull/3381))
	- Add experimental dht 'client mode' flag  ([udfs/go-udfs#3269](https://github.com/udfs/go-udfs/pull/3269))
	- Add config option to set reprovider interval  ([udfs/go-udfs#3101](https://github.com/udfs/go-udfs/pull/3101))
	- Add `udfs dht provide` command  ([udfs/go-udfs#3106](https://github.com/udfs/go-udfs/pull/3106))
	- Add stream info to `udfs swarm peers -v`  ([udfs/go-udfs#3352](https://github.com/udfs/go-udfs/pull/3352))
	- Add option to enable go-multiplex experiment  ([udfs/go-udfs#3447](https://github.com/udfs/go-udfs/pull/3447))
	- Basic Keystore implementation  ([udfs/go-udfs#3472](https://github.com/udfs/go-udfs/pull/3472))
	- Make `udfs add --local` not send providers messages  ([udfs/go-udfs#3102](https://github.com/udfs/go-udfs/pull/3102))
	- Fix bug in `udfs tar add` that buffered input in memory  ([udfs/go-udfs#3334](https://github.com/udfs/go-udfs/pull/3334))
	- Make blockstore retry operations on temporary errors  ([udfs/go-udfs#3091](https://github.com/udfs/go-udfs/pull/3091))
	- Don't hold the PinLock in adder when not pinning.  ([udfs/go-udfs#3222](https://github.com/udfs/go-udfs/pull/3222))
	- Validate repo/api file and improve error message  ([udfs/go-udfs#3219](https://github.com/udfs/go-udfs/pull/3219))
	- no longer hard code gomaxprocs  ([udfs/go-udfs#3357](https://github.com/udfs/go-udfs/pull/3357))
	- Updated Bash complete script  ([udfs/go-udfs#3377](https://github.com/udfs/go-udfs/pull/3377))
	- Remove expensive debug statement in blockstore AllKeysChan  ([udfs/go-udfs#3384](https://github.com/udfs/go-udfs/pull/3384))
	- Remove GC timeout, fix GC tests  ([udfs/go-udfs#3494](https://github.com/udfs/go-udfs/pull/3494))
	- Fix `udfs pin add` resource consumption  ([udfs/go-udfs#3495](https://github.com/udfs/go-udfs/pull/3495))  ([udfs/go-udfs#3571](https://github.com/udfs/go-udfs/pull/3571))
	- Add IPNS entry to DHT cache after publish  ([udfs/go-udfs#3501](https://github.com/udfs/go-udfs/pull/3501))
	- Add in `--routing=none` daemon option  ([udfs/go-udfs#3605](https://github.com/udfs/go-udfs/pull/3605))

- Bitswap
	- Don't re-provide blocks we've provided very recently  ([udfs/go-udfs#3105](https://github.com/udfs/go-udfs/pull/3105))
	- Add a deadline to sendmsg calls ([udfs/go-udfs#3445](https://github.com/udfs/go-udfs/pull/3445))
	- cleanup bitswap and handle message send failure slightly better  ([udfs/go-udfs#3408](https://github.com/udfs/go-udfs/pull/3408))
	- Increase wantlist resend delay to one minute  ([udfs/go-udfs#3448](https://github.com/udfs/go-udfs/pull/3448))
	- Fix issue where wantlist fullness wasn't included in messages  ([udfs/go-udfs#3461](https://github.com/udfs/go-udfs/pull/3461))
	- Only pass keys down newBlocks chan in bitswap   ([udfs/go-udfs#3271](https://github.com/udfs/go-udfs/pull/3271))

- Bugfixes
	- gateway: fix --writable flag  ([udfs/go-udfs#3206](https://github.com/udfs/go-udfs/pull/3206))
	- Fix relative seek in unixfs not expanding file properly   ([udfs/go-udfs#3095](https://github.com/udfs/go-udfs/pull/3095))
	- Update multicodec service names for udfs services  ([udfs/go-udfs#3132](https://github.com/udfs/go-udfs/pull/3132))
	- dht: add missing protocol ID to newStream call  ([udfs/go-udfs#3203](https://github.com/udfs/go-udfs/pull/3203))
	- Return immediately on namesys error  ([udfs/go-udfs#3345](https://github.com/udfs/go-udfs/pull/3345))
	- Improve osxfuse handling  ([udfs/go-udfs#3098](https://github.com/udfs/go-udfs/pull/3098))  ([udfs/go-udfs#3413](https://github.com/udfs/go-udfs/pull/3413))
	- commands: fix opt.Description panic when desc was empty  ([udfs/go-udfs#3521](https://github.com/udfs/go-udfs/pull/3521))
	- Fixes #3133: Properly handle release candidates in version comparison  ([udfs/go-udfs#3136](https://github.com/udfs/go-udfs/pull/3136))
	- Don't drop error in readStreamedJson.  ([udfs/go-udfs#3276](https://github.com/udfs/go-udfs/pull/3276))
	- Error out on invalid `--routing` option  ([udfs/go-udfs#3482](https://github.com/udfs/go-udfs/pull/3482))
	- Respect contexts when returning diagnostics responses  ([udfs/go-udfs#3353](https://github.com/udfs/go-udfs/pull/3353))
	- Fix json marshalling of pbnode  ([udfs/go-udfs#3507](https://github.com/udfs/go-udfs/pull/3507))

- General changes and refactorings
	- Disable Suborigins the spec changed and our impl conflicts  ([udfs/go-udfs#3519](https://github.com/udfs/go-udfs/pull/3519))
	- Avoid sending provide messages for pinsets  ([udfs/go-udfs#3103](https://github.com/udfs/go-udfs/pull/3103))
	- Refactor cli handling to expose argument parsing functionality  ([udfs/go-udfs#3308](https://github.com/udfs/go-udfs/pull/3308))
	- Create a FilestoreNode object to carry PosInfo  ([udfs/go-udfs#3314](https://github.com/udfs/go-udfs/pull/3314))
	- Print 'n/a' instead of zero latency in `udfs swarm peers`  ([udfs/go-udfs#3491](https://github.com/udfs/go-udfs/pull/3491))
	- Add DAGService.GetLinks() method to optimize traversals.  ([udfs/go-udfs#3255](https://github.com/udfs/go-udfs/pull/3255))
	- Make path resolver no longer require whole UdfsNode for construction  ([udfs/go-udfs#3321](https://github.com/udfs/go-udfs/pull/3321))
	- Distinguish between Offline and Local Modes of daemon operation.  ([udfs/go-udfs#3259](https://github.com/udfs/go-udfs/pull/3259))
	- Separate out the GC Locking from the Blockstore interface.  ([udfs/go-udfs#3348](https://github.com/udfs/go-udfs/pull/3348))
	- Avoid unnecessary allocs in datastore key handling  ([udfs/go-udfs#3407](https://github.com/udfs/go-udfs/pull/3407))
	- Use NextSync method for datastore queries ([udfs/go-udfs#3386](https://github.com/udfs/go-udfs/pull/3386))
	- Switch unixfs.Metadata.MimeType to optional ([udfs/go-udfs#3458](https://github.com/udfs/go-udfs/pull/3458))
	- Fix path parsing in `udfs name publish`   ([udfs/go-udfs#3592](https://github.com/udfs/go-udfs/pull/3592))
	- Fix inconsistent `udfs stats bw` formatting  ([udfs/go-udfs#3554](https://github.com/udfs/go-udfs/pull/3554))
	- Set the libp2p agent version based on version string  ([udfs/go-udfs#3569](https://github.com/udfs/go-udfs/pull/3569))

- Cross Platform Changes
	- Fix 'dist_get' script on BSDs.  ([udfs/go-udfs#3264](https://github.com/udfs/go-udfs/pull/3264))
	- ulimit: Tune resource limits on BSDs  ([udfs/go-udfs#3374](https://github.com/udfs/go-udfs/pull/3374))

- Metrics
	- Introduce go-metrics-interface  ([udfs/go-udfs#3189](https://github.com/udfs/go-udfs/pull/3189))
	- Fix metrics injection  ([udfs/go-udfs#3315](https://github.com/udfs/go-udfs/pull/3315))

- Misc
	- Bump Go requirement to 1.7  ([udfs/go-udfs#3111](https://github.com/udfs/go-udfs/pull/3111))
	- Merge 0.4.3 release candidate changes back into master  ([udfs/go-udfs#3248](https://github.com/udfs/go-udfs/pull/3248))
	- Add security@udfs.io GPG key to assets  ([udfs/go-udfs#2997](https://github.com/udfs/go-udfs/pull/2997))
	- Improve makefiles  ([udfs/go-udfs#2999](https://github.com/udfs/go-udfs/pull/2999))  ([udfs/go-udfs#3265](https://github.com/udfs/go-udfs/pull/3265))
	- Refactor install.sh script  ([udfs/go-udfs#3194](https://github.com/udfs/go-udfs/pull/3194))
	- Add test check for go code formatting  ([udfs/go-udfs#3421](https://github.com/udfs/go-udfs/pull/3421))
	- bin: dist_get script: prevents get_go_vars() returns same values twice  ([udfs/go-udfs#3079](https://github.com/udfs/go-udfs/pull/3079))

- Dependencies
	- Update libp2p to have fixed spdystream dep  ([udfs/go-udfs#3210](https://github.com/udfs/go-udfs/pull/3210))
	- Update libp2p and dht packages  ([udfs/go-udfs#3263](https://github.com/udfs/go-udfs/pull/3263))
	- Update to libp2p 4.0.1 and propogate other changes  ([udfs/go-udfs#3284](https://github.com/udfs/go-udfs/pull/3284))
	- Update to libp2p 4.0.4  ([udfs/go-udfs#3361](https://github.com/udfs/go-udfs/pull/3361))
	- Update go-libp2p across codebase  ([udfs/go-udfs#3406](https://github.com/udfs/go-udfs/pull/3406))
	- Update to go-libp2p 4.1.0  ([udfs/go-udfs#3373](https://github.com/udfs/go-udfs/pull/3373))
	- Update deps for libp2p 3.4.0  ([udfs/go-udfs#3110](https://github.com/udfs/go-udfs/pull/3110))
	- Update go-libp2p-swarm with deadlock fixes  ([udfs/go-udfs#3339](https://github.com/udfs/go-udfs/pull/3339))
	- Update to new cid and ipld node packages  ([udfs/go-udfs#3326](https://github.com/udfs/go-udfs/pull/3326))
	- Update to newer ipld node interface with Copy and better Tree  ([udfs/go-udfs#3391](https://github.com/udfs/go-udfs/pull/3391))
	- Update experimental go-multiplex to 0.2.6  ([udfs/go-udfs#3475](https://github.com/udfs/go-udfs/pull/3475))
	- Rework routing interfaces to make separation easier  ([udfs/go-udfs#3107](https://github.com/udfs/go-udfs/pull/3107))
	- Update to dht code with fixed GetClosestPeers  ([udfs/go-udfs#3346](https://github.com/udfs/go-udfs/pull/3346))
	- Move go-is-domain to gx  ([udfs/go-udfs#3077](https://github.com/udfs/go-udfs/pull/3077))
	- Extract thirdparty/loggables and thirdparty/peerset  ([udfs/go-udfs#3204](https://github.com/udfs/go-udfs/pull/3204))
	- Completely remove go-key dep  ([udfs/go-udfs#3439](https://github.com/udfs/go-udfs/pull/3439))
	- Remove randbo dep, its no longer needed  ([udfs/go-udfs#3118](https://github.com/udfs/go-udfs/pull/3118))
	- Update libp2p for identify configuration updates  ([udfs/go-udfs#3539](https://github.com/udfs/go-udfs/pull/3539))
	- Use newer flatfs sharding scheme  ([udfs/go-udfs#3608](https://github.com/udfs/go-udfs/pull/3608))

- Testing
	- fix test_fsh arg quoting in udfs-test-lib  ([udfs/go-udfs#3085](https://github.com/udfs/go-udfs/pull/3085))
	- 100% coverage for blocks/blocksutil  ([udfs/go-udfs#3090](https://github.com/udfs/go-udfs/pull/3090))
	- 100% coverage on blocks/set  ([udfs/go-udfs#3084](https://github.com/udfs/go-udfs/pull/3084))
	- 81% coverage on blockstore  ([udfs/go-udfs#3074](https://github.com/udfs/go-udfs/pull/3074))
	- 80% coverage of unixfs/mod  ([udfs/go-udfs#3096](https://github.com/udfs/go-udfs/pull/3096))
	- 82% coverage on blocks  ([udfs/go-udfs#3086](https://github.com/udfs/go-udfs/pull/3086))
	- 87% coverage on unixfs   ([udfs/go-udfs#3492](https://github.com/udfs/go-udfs/pull/3492)) 
	- Improve coverage on routing/offline  ([udfs/go-udfs#3516](https://github.com/udfs/go-udfs/pull/3516))
	- Add test for flags package   ([udfs/go-udfs#3449](https://github.com/udfs/go-udfs/pull/3449))
	- improve test coverage on merkledag package  ([udfs/go-udfs#3113](https://github.com/udfs/go-udfs/pull/3113))
	- 80% coverage of unixfs/io ([udfs/go-udfs#3097](https://github.com/udfs/go-udfs/pull/3097))
	- Accept more than one digit in repo version tests  ([udfs/go-udfs#3130](https://github.com/udfs/go-udfs/pull/3130))
	- Fix typo in hash in t0050  ([udfs/go-udfs#3170](https://github.com/udfs/go-udfs/pull/3170))
	- fix bug in pinsets and add a stress test for the scenario  ([udfs/go-udfs#3273](https://github.com/udfs/go-udfs/pull/3273))  ([udfs/go-udfs#3302](https://github.com/udfs/go-udfs/pull/3302))
	- Report coverage to codecov  ([udfs/go-udfs#3473](https://github.com/udfs/go-udfs/pull/3473))
	- Add test for 'udfs config replace'  ([udfs/go-udfs#3073](https://github.com/udfs/go-udfs/pull/3073))
	- Fix netcat on macOS not closing socket when the stdin sends EOF  ([udfs/go-udfs#3515](https://github.com/udfs/go-udfs/pull/3515))

- Documentation
	- Update dns help with a correct domain name  ([udfs/go-udfs#3087](https://github.com/udfs/go-udfs/pull/3087))
	- Add period to `udfs pin rm`  ([udfs/go-udfs#3088](https://github.com/udfs/go-udfs/pull/3088))
	- Make all Taglines use imperative mood  ([udfs/go-udfs#3041](https://github.com/udfs/go-udfs/pull/3041))
	- Document listing commands better  ([udfs/go-udfs#3083](https://github.com/udfs/go-udfs/pull/3083))
	- Add notes to readme on building for uncommon systems  ([udfs/go-udfs#3051](https://github.com/udfs/go-udfs/pull/3051))
	- Add branch naming conventions doc  ([udfs/go-udfs#3035](https://github.com/udfs/go-udfs/pull/3035))
	- Replace <default> keyword with <<default>>  ([udfs/go-udfs#3129](https://github.com/udfs/go-udfs/pull/3129))
	- Fix Add() docs regarding pinning  ([udfs/go-udfs#3513](https://github.com/udfs/go-udfs/pull/3513))
	- Add sudo to install commands.  ([udfs/go-udfs#3201](https://github.com/udfs/go-udfs/pull/3201))
	- Add docs for `"commands".Command.Run`  ([udfs/go-udfs#3382](https://github.com/udfs/go-udfs/pull/3382))
	- Put config keys in proper case  ([udfs/go-udfs#3365](https://github.com/udfs/go-udfs/pull/3365))
	- Fix link in `udfs stats bw` help message  ([udfs/go-udfs#3620](https://github.com/udfs/go-udfs/pull/3620))


### 0.4.4 - 2016-10-11

This release contains an important hotfix for a bug we discovered in how pinning works.
If you had a large number of pins, new pins would overwrite existing pins.
Apart from the hotfix, this release is equal to the previous release 0.4.3.

- Fix bug in pinsets fanout, and add stress test. (@whyrusleeping, [udfs/go-udfs#3273](https://github.com/udfs/go-udfs/pull/3273))

We published a [detailed account of the bug and fix in a blog post](https://udfs.io/blog/21-go-udfs-0-4-4-released/).

### 0.4.3 - 2016-09-20

There have been no changes since the last release candidate 0.4.3-rc4. \o/

### 0.4.3-rc4 - 2016-09-09

This release candidate fixes issues in Bitswap and the `udfs add` command, and improves testing.
We plan for this to be the last release candidate before the release of go-udfs v0.4.3.

With this release candidate, we're also moving go-udfs to Go 1.7, which we expect will yield improvements in runtime performance, memory usage, build time and size of the release binaries.

- Require Go 1.7. (@whyrusleeping, @Kubuxu, @lgierth, [udfs/go-udfs#3163](https://github.com/udfs/go-udfs/pull/3163))
  - For this purpose, switch Docker image from Alpine 3.4 to Alpine Edge.
- Fix cancellation of Bitswap `wantlist` entries. (@whyrusleeping, [udfs/go-udfs#3182](https://github.com/udfs/go-udfs/pull/3182))
- Fix clearing of `active` state of Bitswap provider queries. (@whyrusleeping, [udfs/go-udfs#3169](https://github.com/udfs/go-udfs/pull/3169))
- Fix a panic in the DHT code. (@Kubuxu, [udfs/go-udfs#3200](https://github.com/udfs/go-udfs/pull/3200))
- Improve handling of `Identity` field in `udfs config` command. (@Kubuxu, @whyrusleeping, [udfs/go-udfs#3141](https://github.com/udfs/go-udfs/pull/3141))
- Fix explicit adding of symlinked files and directories. (@kevina, [udfs/go-udfs#3135](https://github.com/udfs/go-udfs/pull/3135))
- Fix bash auto-completion of `udfs daemon --unrestricted-api` option. (@lgierth, [udfs/go-udfs#3159](https://github.com/udfs/go-udfs/pull/3159))
- Introduce a new timeout tool for tests to avoid licensing issues. (@Kubuxu, [udfs/go-udfs#3152](https://github.com/udfs/go-udfs/pull/3152))
- Improve output for migrations of fs-repo. (@lgierth, [udfs/go-udfs#3158](https://github.com/udfs/go-udfs/pull/3158))
- Fix info notice of commands taking input from stdin. (@Kubuxu, [udfs/go-udfs#3134](https://github.com/udfs/go-udfs/pull/3134))
- Bring back a few tests for stdin handling of `udfs cat` and `udfs add`. (@Kubuxu, [udfs/go-udfs#3144](https://github.com/udfs/go-udfs/pull/3144))
- Improve sharness tests for `udfs repo verify` command. (@whyrusleeping, [udfs/go-udfs#3148](https://github.com/udfs/go-udfs/pull/3148))
- Improve sharness tests for CORS headers on the gateway. (@Kubuxu, [udfs/go-udfs#3142](https://github.com/udfs/go-udfs/pull/3142))
- Improve tests for pinning within `udfs files`. (@kevina, [udfs/go-udfs#3151](https://github.com/udfs/go-udfs/pull/3151))
- Improve tests for the automatic raising of file descriptor limits. (@whyrusleeping, [udfs/go-udfs#3149](https://github.com/udfs/go-udfs/pull/3149))

### 0.4.3-rc3 - 2016-08-11

This release candidate fixes a panic that occurs when input from stdin was
expected, but none was given: [udfs/go-udfs#3050](https://github.com/udfs/go-udfs/pull/3050)

### 0.4.3-rc2 - 2016-08-04

This release includes bugfixes and fixes for regressions that were introduced
between 0.4.2 and 0.4.3-rc1.

- Regressions
  - Fix daemon panic when there is no multipart input provided over the HTTP API.
  (@whyrusleeping, [udfs/go-udfs#2989](https://github.com/udfs/go-udfs/pull/2989))
  - Fix `udfs refs --edges` not printing edges.
  (@Kubuxu, [udfs/go-udfs#3007](https://github.com/udfs/go-udfs/pull/3007))
  - Fix progress option for `udfs add` defaulting to true on the HTTP API.
  (@whyrusleeping, [udfs/go-udfs#3025](https://github.com/udfs/go-udfs/pull/3025))
  - Fix erroneous printing of stdin reading message.
  (@whyrusleeping, [udfs/go-udfs#3033](https://github.com/udfs/go-udfs/pull/3033))
  - Fix panic caused by passing `--mount` and `--offline` flags to `udfs daemon`.
  (@Kubuxu, [udfs/go-udfs#3022](https://github.com/udfs/go-udfs/pull/3022))
  - Fix symlink path resolution on windows.
  (@Kubuxu, [udfs/go-udfs#3023](https://github.com/udfs/go-udfs/pull/3023))
  - Add in code to prevent issue 3032 from crashing the daemon.
  (@whyrusleeping, [udfs/go-udfs#3037](https://github.com/udfs/go-udfs/pull/3037))


### 0.4.3-rc1 - 2016-07-23

This is a maintenance release which comes with a couple of nice enhancements, and improves the performance of Storage, Bitswap, as well as Content and Peer Routing. It also introduces a handful of new commands and options, and fixes a good bunch of bugs.

This is the first Release Candidate. Unless there are vulnerabilities or regressions discovered, the final 0.4.3 release will happen about one week from now.

- Security Vulnerability

  - The `master` branch if go-udfs suffered from a vulnerability for about 3 weeks. It allowed an attacker to use an iframe to request malicious HTML and JS from the API of a local go-udfs node. The attacker could then gain unrestricted access to the node's API, and e.g. extract the private key. We fixed this issue by reintroducing restrictions on which particular objects can be loaded through the API (@lgierth, [udfs/go-udfs#2949](https://github.com/udfs/go-udfs/pull/2949)), and by completely excluding the private key from the API (@Kubuxu, [udfs/go-udfs#2957](https://github.com/udfs/go-udfs/pull/2957)). We will also work on more hardening of the API in the next release.
  - **The previous release 0.4.2 is not vulnerable. That means if you're using official binaries from [dist.udfs.io](https://dist.udfs.io) you're not affected.** If you're running go-udfs built from the `master` branch between June 17th ([udfs/go-udfs@1afebc21](https://github.com/udfs/go-udfs/commit/1afebc21f324982141ca8a29710da0d6f83ca804)) and July 7th ([udfs/go-udfs@39bef0d5](https://github.com/udfs/go-udfs/commit/39bef0d5b01f70abf679fca2c4d078a2d55620e2)), please update to v0.4.3-rc1 immediately.
  - We are grateful to the group of independent researchers who made us aware of this vulnerability. We wanna use this opportunity to reiterate that we're very happy about any additional review of pull requests and releases. You can contact us any time at security@udfs.io (GPG [4B9665FB 92636D17 7C7A86D3 50AAE8A9 59B13AF3](https://pgp.mit.edu/pks/lookup?op=get&search=0x50AAE8A959B13AF3)).

- Notable changes

  - Improve Bitswap performance. (@whyrusleeping, [udfs/go-udfs#2727](https://github.com/udfs/go-udfs/pull/2727), [udfs/go-udfs#2798](https://github.com/udfs/go-udfs/pull/2798))
  - Improve Content Routing and Peer Routing performance. (@whyrusleeping, [udfs/go-udfs#2817](https://github.com/udfs/go-udfs/pull/2817), [udfs/go-udfs#2841](https://github.com/udfs/go-udfs/pull/2841))
  - Improve datastore, blockstore, and dagstore performance. (@kevina, @Kubuxu, @whyrusleeping [udfs/go-datastore#43](https://github.com/udfs/go-datastore/pull/43), [udfs/go-udfs#2885](https://github.com/udfs/go-udfs/pull/2885), [udfs/go-udfs#2961](https://github.com/udfs/go-udfs/pull/2961), [udfs/go-udfs#2953](https://github.com/udfs/go-udfs/pull/2953), [udfs/go-udfs#2960](https://github.com/udfs/go-udfs/pull/2960))
  - Content Providers are now stored on disk to gain savings on process memory. (@whyrusleeping, [udfs/go-udfs#2804](https://github.com/udfs/go-udfs/pull/2804), [udfs/go-udfs#2860](https://github.com/udfs/go-udfs/pull/2860))
  - Migrations of the fs-repo (usually stored at `~/.udfs`) now run automatically. If there's a TTY available, you'll get prompted when running `udfs daemon`, and in addition you can use the `--migrate=true` or `--migrate=false` options to avoid the prompt. (@whyrusleeping, @lgierth, [udfs/go-udfs#2939](https://github.com/udfs/go-udfs/pull/2939))
  - The internal naming of blocks in the blockstore has changed, which requires a migration of the fs-repo, from version 3 to 4. (@whyrusleeping, [udfs/go-udfs#2903](https://github.com/udfs/go-udfs/pull/2903))
  - We now automatically raise the file descriptor limit to 1024 if neccessary. (@whyrusleeping, [udfs/go-udfs#2884](https://github.com/udfs/go-udfs/pull/2884), [udfs/go-udfs#2891](https://github.com/udfs/go-udfs/pull/2891))
  - After a long struggle with deadlocks and hanging connections, we've decided to disable the uTP transport by default for now. (@whyrusleeping, [udfs/go-udfs#2840](https://github.com/udfs/go-udfs/pull/2840), [udfs/go-libp2p-transport@88244000](https://github.com/udfs/go-libp2p-transport/commit/88244000f0ce8851ffcfbac746ebc0794b71d2a4))
  - There is now documentation for the configuration options in `docs/config.md`. (@whyrusleeping, [udfs/go-udfs#2974](https://github.com/udfs/go-udfs/pull/2974))
  - All commands now sanely handle the combination of stdin and optional flags in certain edge cases. (@lgierth, [udfs/go-udfs#2952](https://github.com/udfs/go-udfs/pull/2952))

- New Features

  - Add `--offline` option to `udfs daemon` command, which disables all swarm networking. (@Kubuxu, [udfs/go-udfs#2696](https://github.com/udfs/go-udfs/pull/2696), [udfs/go-udfs#2867](https://github.com/udfs/go-udfs/pull/2867))
  - Add `Datastore.HashOnRead` option for verifying block hashes on read access. (@Kubuxu, [udfs/go-udfs#2904](https://github.com/udfs/go-udfs/pull/2904))
  - Add `Datastore.BloomFilterSize` option for tuning the blockstore's new lookup bloom filter. (@Kubuxu, [udfs/go-udfs#2973](https://github.com/udfs/go-udfs/pull/2973))

- Bugfixes

  - Fix publishing of local IPNS entries, and more. (@whyrusleeping, [udfs/go-udfs#2943](https://github.com/udfs/go-udfs/pull/2943))
  - Fix progress bars in `udfs add` and `udfs get`. (@whyrusleeping, [udfs/go-udfs#2893](https://github.com/udfs/go-udfs/pull/2893), [udfs/go-udfs#2948](https://github.com/udfs/go-udfs/pull/2948))
  - Make sure files added through `udfs files` are pinned and don't get GC'd. (@kevina, [udfs/go-udfs#2872](https://github.com/udfs/go-udfs/pull/2872))
  - Fix copying into directory using `udfs files cp`. (@whyrusleeping, [udfs/go-udfs#2977](https://github.com/udfs/go-udfs/pull/2977))
  - Fix `udfs version --commit` with Docker containers. (@lgierth, [udfs/go-udfs#2734](https://github.com/udfs/go-udfs/pull/2734))
  - Run `udfs diag` commands in the daemon instead of the CLI. (@Kubuxu, [udfs/go-udfs#2761](https://github.com/udfs/go-udfs/pull/2761))
  - Fix protobuf encoding on the API and in commands. (@stebalien, [udfs/go-udfs#2516](https://github.com/udfs/go-udfs/pull/2516))
  - Fix goroutine leak in `/udfs/ping` protocol handler. (@whyrusleeping, [udfs/go-libp2p#58](https://github.com/udfs/go-libp2p/pull/58))
  - Fix `--flags` option on `udfs commands`. (@Kubuxu, [udfs/go-udfs#2773](https://github.com/udfs/go-udfs/pull/2773))
  - Fix the error channels in `namesys`. (@whyrusleeping, [udfs/go-udfs#2788](https://github.com/udfs/go-udfs/pull/2788))
  - Fix consumptions of observed swarm addresses. (@whyrusleeping, [udfs/go-libp2p#63](https://github.com/udfs/go-libp2p/pull/63), [udfs/go-udfs#2771](https://github.com/udfs/go-udfs/issues/2771))
  - Fix a rare DHT panic. (@whyrusleeping, [udfs/go-udfs#2856](https://github.com/udfs/go-udfs/pull/2856))
  - Fix go-udfs/js-udfs interoperability issues in SPDY. (@whyrusleeping, [whyrusleeping/go-smux-spdystream@fae17783](https://github.com/whyrusleeping/go-smux-spdystream/commit/fae1778302a9e029bb308cf71cf33f857f2d89e8))
  - Fix a logging race condition during shutdown. (@Kubuxu, [udfs/go-log#3](https://github.com/udfs/go-log/pull/3))
  - Prevent DHT connection hangs. (@whyrusleeping, [udfs/go-udfs#2826](https://github.com/udfs/go-udfs/pull/2826), [udfs/go-udfs#2863](https://github.com/udfs/go-udfs/pull/2863))
  - Fix NDJSON output of `udfs refs local`. (@Kubuxu, [udfs/go-udfs#2812](https://github.com/udfs/go-udfs/pull/2812))
  - Fix race condition in NAT detection. (@whyrusleeping, [udfs/go-libp2p#69](https://github.com/udfs/go-libp2p/pull/69))
  - Fix error messages. (@whyrusleeping, @Kubuxu, [udfs/go-udfs#2905](https://github.com/udfs/go-udfs/pull/2905), [udfs/go-udfs#2928](https://github.com/udfs/go-udfs/pull/2928))

- Enhancements

  - Increase maximum object size on `udfs put` from 1 MiB to 2 MiB. The maximum object size on the wire including all framing is 4 MiB. (@kpcyrd, [udfs/go-udfs#2980](https://github.com/udfs/go-udfs/pull/2980))
  - Add CORS headers to the Gateway's default config. (@Kubuxu, [udfs/go-udfs#2778](https://github.com/udfs/go-udfs/pull/2778))
  - Clear the dial backoff for a peer when using `udfs swarm connect`. (@whyrusleeping, [udfs/go-udfs#2941](https://github.com/udfs/go-udfs/pull/2941))
  - Allow passing options to daemon in Docker container. (@lgierth, [udfs/go-udfs#2955](https://github.com/udfs/go-udfs/pull/2955))
  - Add `-v/--verbose` to `ìpfs swarm peers` command. (@csasarak, [udfs/go-udfs#2713](https://github.com/udfs/go-udfs/pull/2713))
  - Add `--format`, `--hash`, and `--size` options to `udfs files stat` command. (@Kubuxu, [udfs/go-udfs#2706](https://github.com/udfs/go-udfs/pull/2706))
  - Add `--all` option to `udfs version` command. (@Kubuxu, [udfs/go-udfs#2790](https://github.com/udfs/go-udfs/pull/2790))
  - Add `udfs repo version` command. (@pfista, [udfs/go-udfs#2598](https://github.com/udfs/go-udfs/pull/2598))
  - Add `udfs repo verify` command. (@whyrusleeping, [udfs/go-udfs#2924](https://github.com/udfs/go-udfs/pull/2924), [udfs/go-udfs#2951](https://github.com/udfs/go-udfs/pull/2951))
  - Add `udfs stats repo` and `udfs stats bitswap` command aliases. (@pfista, [udfs/go-udfs#2810](https://github.com/udfs/go-udfs/pull/2810))
  - Add success indication to responses of `udfs ping` command. (@Kubuxu, [udfs/go-udfs#2813](https://github.com/udfs/go-udfs/pull/2813))
  - Save changes made via `udfs swarm filter` to the config file. (@yuvallanger, [udfs/go-udfs#2880](https://github.com/udfs/go-udfs/pull/2880))
  - Expand `udfs_p2p_peers` metric to include libp2p transport. (@lgierth, [udfs/go-udfs#2728](https://github.com/udfs/go-udfs/pull/2728))
  - Rework `udfs files add` internals to avoid caching and prevent memory leaks. (@whyrusleeping, [udfs/go-udfs#2795](https://github.com/udfs/go-udfs/pull/2795))
  - Support `GOPATH` with multiple path components. (@karalabe, @lgierth, @djdv, [udfs/go-udfs#2808](https://github.com/udfs/go-udfs/pull/2808), [udfs/go-udfs#2862](https://github.com/udfs/go-udfs/pull/2862), [udfs/go-udfs#2975](https://github.com/udfs/go-udfs/pull/2975))

- General Codebase

  - Take steps towards the `filestore` datastore. (@kevina, [udfs/go-udfs#2792](https://github.com/udfs/go-udfs/pull/2792), [udfs/go-udfs#2634](https://github.com/udfs/go-udfs/pull/2634))
  - Update recommended Golang version to 1.6.2 (@Kubuxu, [udfs/go-udfs#2724](https://github.com/udfs/go-udfs/pull/2724))
  - Update to Gx 0.8.0 and Gx-Go 1.2.1, which is faster and less noisy. (@whyrusleeping, [udfs/go-udfs#2979](https://github.com/udfs/go-udfs/pull/2979))
  - Use `go4.org/lock` instead of `camlistore/lock` for locking. (@whyrusleeping, [udfs/go-udfs#2887](https://github.com/udfs/go-udfs/pull/2887))
  - Manage `go.uuid`, `hamming`, `backoff`, `proquint`, `pb`, `go-context`, `cors`, `go-datastore` packages with Gx. (@Kubuxu, [udfs/go-udfs#2733](https://github.com/udfs/go-udfs/pull/2733), [udfs/go-udfs#2736](https://github.com/udfs/go-udfs/pull/2736), [udfs/go-udfs#2757](https://github.com/udfs/go-udfs/pull/2757), [udfs/go-udfs#2825](https://github.com/udfs/go-udfs/pull/2825), [udfs/go-udfs#2838](https://github.com/udfs/go-udfs/pull/2838))
  - Clean up the gateway's surface. (@lgierth, [udfs/go-udfs#2874](https://github.com/udfs/go-udfs/pull/2874))
  - Simplify the API gateway's access restrictions. (@lgierth, [udfs/go-udfs#2949](https://github.com/udfs/go-udfs/pull/2949), [udfs/go-udfs#2956](https://github.com/udfs/go-udfs/pull/2956))
  - Update docker image to Alpine Linux 3.4 and remove Go version constraint. (@lgierth, [udfs/go-udfs#2901](https://github.com/udfs/go-udfs/pull/2901), [udfs/go-udfs#2929](https://github.com/udfs/go-udfs/pull/2929))
  - Clarify `Dockerfile` and `Dockerfile.fast`. (@lgierth, [udfs/go-udfs#2796](https://github.com/udfs/go-udfs/pull/2796))
  - Simplify resolution of Git commit refs in Dockerfiles. (@lgierth, [udfs/go-udfs#2754](https://github.com/udfs/go-udfs/pull/2754))
  - Consolidate `--verbose` description across commands. (@Kubuxu, [udfs/go-udfs#2746](https://github.com/udfs/go-udfs/pull/2746))
  - Allow setting position of default values in command option descriptions. (@Kubuxu, [udfs/go-udfs#2744](https://github.com/udfs/go-udfs/pull/2744))
  - Set explicit default values for boolean command options. (@RichardLitt, [udfs/go-udfs#2657](https://github.com/udfs/go-udfs/pull/2657))
  - Autogenerate command synopsises. (@Kubuxu, [udfs/go-udfs#2785](https://github.com/udfs/go-udfs/pull/2785))
  - Fix and improve lots of documentation. (@RichardLitt, [udfs/go-udfs#2741](https://github.com/udfs/go-udfs/pull/2741), [udfs/go-udfs#2781](https://github.com/udfs/go-udfs/pull/2781))
  - Improve command descriptions to fit a width of 78 characters. (@RichardLitt, [udfs/go-udfs#2779](https://github.com/udfs/go-udfs/pull/2779), [udfs/go-udfs#2780](https://github.com/udfs/go-udfs/pull/2780), [udfs/go-udfs#2782](https://github.com/udfs/go-udfs/pull/2782))
  - Fix filename conflict in the debugging guide. (@Kubuxu, [udfs/go-udfs#2752](https://github.com/udfs/go-udfs/pull/2752))
  - Decapitalize log messages, according to Golang style guides. (@RichardLitt, [udfs/go-udfs#2853](https://github.com/udfs/go-udfs/pull/2853))
  - Add Github Issues HowTo guide. (@RichardLitt, @chriscool, [udfs/go-udfs#2889](https://github.com/udfs/go-udfs/pull/2889), [udfs/go-udfs#2895](https://github.com/udfs/go-udfs/pull/2895))
  - Add Github Issue template. (@chriscool, [udfs/go-udfs#2786](https://github.com/udfs/go-udfs/pull/2786))
  - Apply standard-readme to the README file. (@RichardLitt, [udfs/go-udfs#2883](https://github.com/udfs/go-udfs/pull/2883))
  - Fix issues pointed out by `govet`. (@Kubuxu, [udfs/go-udfs#2854](https://github.com/udfs/go-udfs/pull/2854))
  - Clarify `udfs get` error message. (@whyrusleeping, [udfs/go-udfs#2886](https://github.com/udfs/go-udfs/pull/2886))
  - Remove dead code. (@whyrusleeping, [udfs/go-udfs#2819](https://github.com/udfs/go-udfs/pull/2819))
  - Add changelog for v0.4.3. (@lgierth, [udfs/go-udfs#2984](https://github.com/udfs/go-udfs/pull/2984))

- Tests & CI

  - Fix flaky `udfs mount` sharness test by using the `iptb` tool. (@noffle, [udfs/go-udfs#2707](https://github.com/udfs/go-udfs/pull/2707))
  - Fix flaky IP port selection in tests. (@Kubuxu, [udfs/go-udfs#2855](https://github.com/udfs/go-udfs/pull/2855))
  - Fix CLI tests on OSX by resolving /tmp symlink. (@Kubuxu, [udfs/go-udfs#2926](https://github.com/udfs/go-udfs/pull/2926))
  - Fix flaky GC test by running the daemon in offline mode. (@Kubuxu, [udfs/go-udfs#2908](https://github.com/udfs/go-udfs/pull/2908))
  - Add tests for `udfs add` with hidden files. (@Kubuxu, [udfs/go-udfs#2756](https://github.com/udfs/go-udfs/pull/2756))
  - Add test to make sure the body of HEAD responses is empty. (@Kubuxu, [udfs/go-udfs#2775](https://github.com/udfs/go-udfs/pull/2775))
  - Add test to catch misdials. (@Kubuxu, [udfs/go-udfs#2831](https://github.com/udfs/go-udfs/pull/2831))
  - Mark flaky tests for `udfs dht query` as known failure. (@noffle, [udfs/go-udfs#2720](https://github.com/udfs/go-udfs/pull/2720))
  - Remove failing blockstore-without-context test. (@Kubuxu, [udfs/go-udfs#2857](https://github.com/udfs/go-udfs/pull/2857))
  - Fix `--version` tests for versions with a suffix like `-dev` or `-rc1`. (@lgierth, [udfs/go-udfs#2937](https://github.com/udfs/go-udfs/pull/2937))
  - Make sharness tests work in cases where go-udfs is symlinked into GOPATH. (@lgierth, [udfs/go-udfs#2937](https://github.com/udfs/go-udfs/pull/2937))
  - Add variable delays to blockstore mocks. (@rikonor, [udfs/go-udfs#2871](https://github.com/udfs/go-udfs/pull/2871))
  - Disable Travis CI email notifications. (@Kubuxu, [udfs/go-udfs#2896](https://github.com/udfs/go-udfs/pull/2896))


### 0.4.2 - 2016-05-17

This is a patch release which fixes performance and networking bugs in go-libp2p,
You should see improvements in CPU and RAM usage, as well as speed of object lookups.
There are also a few other nice improvements.

* Notable Fixes
  * Set a deadline for dialing attempts. This prevents a node from accumulating
    failed connections. (@whyrusleeping)
  * Avoid unnecessary string/byte conversions in go-multihash. (@whyrusleeping)
  * Fix a deadlock around the yamux stream muxer. (@whyrusleeping)
  * Fix a bug that left channels open, causing hangs. (@whyrusleeping)
  * Fix a bug around yamux which caused connection hangs. (@whyrusleeping)
  * Fix a crash caused by nil multiaddrs. (@whyrusleeping)

* Enhancements
  * Add NetBSD support. (@erde74)
  * Set Cache-Control: immutable on /udfs responses. (@kpcyrd)
  * Have `udfs init` optionally accept a default configuration from stdin. (@sivachandran)
  * Add `udfs log ls` command for listing logging subsystems. (@hsanjuan)
  * Allow bitswap to read multiple messages per stream. (@whyrusleeping)
  * Remove `make toolkit_upgrade` step. (@chriscool)

* Documentation
  * Add a debug-guidelines document. (@richardlitt)
  * Update the contribute document. (@richardlitt)
  * Fix documentation of many `udfs` commands. (@richardlitt)
  * Fall back to ShortDesc if LongDesc is missing. (@Kubuxu)

* Removals
  * Remove -f option from `udfs init` command. (@whyrusleeping)

* Bugfixes
  * Fix `udfs object patch` argument handling and validation. (@jbenet)
  * Fix `udfs config edit` command by running it client-side. (@Kubuxu)
  * Set default value for `udfs refs` arguments. (@richardlitt)
  * Fix parsing of incorrect command and argument permutations. (@thomas-gardner)
  * Update Dockerfile to latest go1.5.4-r0. (@chriscool)
  * Allow passing UDFS_LOGGING to Docker image. (@lgierth)
  * Fix dot path parsing on Windows. (@djdv)
  * Fix formatting of `udfs log ls` output. (@richardlitt)

* General Codebase
  * Refactor Makefile. (@kevina)
  * Wire context into bitswap requests more deeply. (@whyrusleeping)
  * Use gx for iptb. (@chriscool)
  * Update gx and gx-go. (@chriscool)
  * Make blocks.Block an interface. (@kevina)
  * Silence check for Docker existance. (@chriscool)
  * Add dist_get script for fetching tools from dist.udfs.io. (@whyrusleeping)
  * Add proper defaults to all `udfs` commands. (@richardlitt)
  * Remove dead `count` option from `udfs pin ls`. (@richardlitt)
  * Initialize pin mode strings only once. (@chriscool)
  * Add changelog for v0.4.2. (@lgierth)
  * Specify a dist.udfs.io hash for tool downloads instead of trusting DNS. (@lgierth)

* CI
  * Fix t0170-dht sharness test. (@chriscool)
  * Increase timeout in t0060-daemon sharness test. (@Kubuxu)
  * Have CircleCI use `make deps` instead of `gx` directly. (@whyrusleeping)


### 0.4.1 - 2016-04-25

This is a patch release that fixes a few bugs, and adds a few small (but not
insignificant) features. The primary reason for this release is the listener
hang bugfix that was shipped in the 0.4.0 release.

* Features
  * implemented udfs object diff (@whyrusleeping)
  * allow promises (used in get, refs) to fail (@whyrusleeping)

* Tool changes
  * Adds 'toolkit_upgrade' to the makefile help target (@achin)

* General Codebase
  * Use extracted go-libp2p-crypto, -secio, -peer packages (@lgierth)
  * Update go-libp2p (@lgierth)
  * Fix package manifest fields (@lgierth)
  * remove incfusever dead-code (@whyrusleeping)
  * remove a ton of unused godeps (@whyrusleeping)
  * metrics: add prometheus back (@lgierth)
  * clean up dead code and config fields (@whyrusleeping)
  * Add log events when blocks are added/removed from the blockstore (@michealmure)
  * repo: don't create logs directory, not used any longer (@lgierth)

* Bugfixes
  * fixed udfs name resolve --local multihash error (@pfista)
  * udfs patch commands won't return null links field anymore (@whyrusleeping)
  * Make non recursive resolve print the result (@Kubuxu)
  * Output dirs on udfs add -rn (@noffle)
  * update libp2p dep to fix hanging listeners problem (@whyrusleeping)
  * Fix Swarm.AddrFilters config setting with regard to `/ip6` addresses (@lgierth)
  * fix dht command key escaping (@whyrusleeping)

* Testing
  * Adds tests to make sure 'object patch' writes. (@noffle)
  * small sharness test for promise failure checking (@whyrusleeping)
  * sharness/Makefile: clean all BINS when cleaning (@chriscool)

* Documentation
  * Fix disconnect argument description (@richardlitt)
  * Added a note about swarm disconnect (@richardlitt)
  * Also fixed syntax for comment (@richardlitt)
  * Alphabetized swarm subcmds (@richardlitt)
  * Added note to udfs stats bw interval option (@richardlitt)
  * Small syntax changes to repo stat man (@richardlitt)
  * update log command help text (@pfista)
  * Added a long description to add (@richardlitt)
  * Edited object patch set-data doc (@richardlitt)
  * add roadmap.md (@Jeromy)
  * Adds files api cmd to helptext (@noffle)


### 0.4.0 - 2016-04-05

This is a major release with plenty of new features and bugfixes.
It also includes breaking changes which make it incompatible with v0.3.x
on the networking layer.

* Major Changes
  * Multistream
    * The addition of multistream is a breaking change on the networking layer,
      but gives UDFS implementations the ability to mix and match different
      stream multiplexers, e.g. yamux, spdystream, or muxado.
      This adds a ton of flexibility on one of the lower layers of the protocol,
      and will help us avoid further breaking protocol changes in the future.
  * Files API
    * The new `files` command and API allow a program to interact with UDFS
      using familiar filesystem operations, namely: creating directories,
      reading, writing, and deleting files, listing out different directories,
      and so on. This feature enables any other application that uses a
      filesystem-like backend for storage, to use UDFS as its storage driver
      without having change the application logic at all.
  * Gx
    * go-udfs now uses [gx](https://github.com/whyrusleeping/gx) to manage its
      dependencies. This means that under the hood, go-udfs's dependencies are
      backed by UDFS itself! It also means that go-udfs is no longer installed
      using `go get`. Use `make install` instead.
* New Features
  * Web UI
    * Update to new version which is compatible with 0.4.0. (@dignifiedquire)
  * Networking
    * Implement uTP transport. (@whyrusleeping)
    * Allow multiple addresses per configured bootstrap node. (@whyrusleeping)
  * IPNS
    * Improve IPNS resolution performance. (@whyrusleeping)
    * Have dnslink prefer `TXT _dnslink.example.com`, allows usage of CNAME records. (@Kubuxu)
    * Prevent `udfs name publish` when `/ipns` is mounted. (@noffle)
  * Repo
    * Improve performance of `udfs add`. (@whyrusleeping)
    * Add `Datastore.NoSync` config option for flatfs. (@rht)
    * Implement mark-and-sweep GC. (@whyrusleeping)
    * Allow for GC during `udfs add`. (@whyrusleeping)
    * Add `udfs repo stat` command. (@tmg, @diasdavid)
  * General
    * Add support for HTTP OPTIONS requests. (@lidel)
    * Add `udfs diag cmds` to view active API requests (@whyrusleeping)
    * Add an `UDFS_LOW_MEM` environment variable which relaxes Bitswap's memory usage. (@whyrusleeping)
    * The Docker image now lives at `udfs/go-udfs` and has been completely reworked. (@lgierth)
* Security fixes
  * The gateway path prefix added in v0.3.10 was vulnerable to cross-site
    scripting attacks. This release introduces a configurable list of allowed
    path prefixes. It's called `Gateway.PathPrefixes` and takes a list of
    strings, e.g. `["/blog", "/foo/bar"]`. The v0.3.x line will not receive any
    further updates, so please update to v0.4.0 as soon as possible. (@lgierth)
* Incompatible Changes
  * Install using `make install` instead of `go get` (@whyrusleeping)
  * Rewrite pinning to store pins in UDFS objects. (@tv42)
  * Bump fs-repo version to 3. (@whyrusleeping)
  * Use multistream muxer (@whyrusleeping)
  * The default for `--type` in `udfs pin ls` is now `all`. (@chriscool)
* Bug Fixes
  * Remove msgio double wrap. (@jbenet)
  * Buffer msgio. (@whyrusleeping)
  * Perform various fixes to the FUSE code. (@tv42)
  * Compute `udfs add` size in background to not stall add operation. (@whyrusleeping)
  * Add option to have `udfs add` include top-level hidden files. (@noffle)
  * Fix CORS checks on the API. (@rht)
  * Fix `udfs update` error message. (@tomgg)
  * Resolve paths in `udfs pin rm` without network lookup. (@noffle)
  * Detect FUSE unmounts and track mount state. (@noffle)
  * Fix go1.6rc2 panic caused by CloseNotify being called from wrong goroutine. (@rwcarlsen)
  * Bump DHT kvalue from 10 to 20. (@whyrusleeping)
  * Put public key and IPNS entry to DHT in parallel. (@whyrusleeping)
  * Fix panic in CLI argument parsing. (@whyrusleeping)
  * Fix range error by using larger-than-zero-length buffer. (@noffle)
  * Fix yamux hanging issue by increasing AcceptBacklog. (@whyrusleeping)
  * Fix double Transport-Encoding header bug. (@whyrusleeping)
  * Fix uTP panic and file descriptor leak. (@whyrusleeping)
* Tool Changes
  * Add `--pin` option to `udfs add`, which defaults to `true` and allows `--pin=false`. (@eminence)
  * Add arguments to `udfs pin ls`. (@chriscool)
  * Add `dns` and `resolve` commands to read-only API. (@Kubuxu)
  * Add option to display headers for `udfs object links`. (@palkeo)
* General Codebase Changes
  * Check Golang version in Makefile. (@chriscool)
  * Improve Makefile. (@tomgg)
  * Remove dead Jenkins CI code. (@lgierth)
  * Add locking interface to blockstore. (@whyrusleeping)
  * Add Merkledag FetchGraph and EnumerateChildren. (@whyrusleeping)
  * Rename Lock/RLock to GCLock/PinLock. (@jbenet)
  * Implement pluggable datastore types. (@tv42)
  * Record datastore metrics for non-default datastores. (@tv42)
  * Allow multistream to have zero-rtt stream opening. (@whyrusleeping)
  * Refactor `ipnsfs` into a more generic and well tested `mfs`. (@whyrusleeping)
  * Grab more peers if bucket doesn't contain enough. (@whyrusleeping)
  * Use CloseNotify in gateway. (@whyrusleeping)
  * Flatten multipart file transfers. (@whyrusleeping)
  * Send updated DHT record fixes to peers who sent outdated records. (@whyrusleeping)
  * Replace go-psutil with go-sysinfo. (@whyrusleeping)
  * Use ServeContent for index.html. (@AtnNn)
  * Refactor `object patch` API to not store data in URL. (@whyrusleeping)
  * Use mfs for `udfs add`. (@whyrusleeping)
  * Add `Server` header to API responses. (@Kubuxu)
  * Wire context directly into HTTP requests. (@rht)
  * Wire context directly into GetDAG operations within GC. (@rht)
  * Vendor libp2p using gx. (@whyrusleeping)
  * Use gx vendored packages instead of Godeps. (@whyrusleeping)
  * Simplify merkledag package interface to ease IPLD inclusion. (@mildred)
  * Add default option value support to commands lib. (@whyrusleeping)
  * Refactor merkledag fetching methods. (@whyrusleeping)
  * Use net/url to escape paths within Web UI. (@noffle)
  * Deprecated key.Pretty(). (@MichealMure)
* Documentation
  * Fix and update help text for **every** `udfs` command. (@RichardLitt)
  * Change sample API origin settings from wildcard (`*`) to `example.com`. (@Kubuxu)
  * Improve documentation of installation process in README. (@whyrusleeping)
  * Improve windows.md. (@chriscool)
  * Clarify instructions for installing from source. (@noffle)
  * Make version checking more robust. (@jedahan)
  * Assert the source code is located within GOPATH. (@whyrusleeping)
  * Remove mentions of `/dns` from `udfs dns` command docs. (@lgierth)
* Testing
  * Refactor iptb tests. (@chriscool)
  * Improve t0240 sharness test. (@chriscool)
  * Make bitswap tests less flaky. (@whyrusleeping)
  * Use TCP port zero for udfs daemon in sharness tests. (@whyrusleeping)
  * Improve sharness tests on AppVeyor. (@chriscool)
  * Add a pause to fix timing on t0065. (@whyrusleeping)
  * Add support for arbitrary TCP ports to t0060-daemon.sh. (@noffle)
  * Make t0060 sharness test use TCP port zero. (@whyrusleeping)
  * Randomized udfs stress testing via randor (@dignifiedquire)
  * Stress test pinning and migrations (@whyrusleeping)

### 0.3.11 - 2016-01-12

This is the final udfs version before the transition to v0.4.0.
It introduces a few stability improvements, bugfixes, and increased
test coverage.

* Features
  * Add 'get' and 'patch' to the allowed gateway commands (@whyrusleeping)
  * Updated webui version (@dignifiedquire)

* BugFixes
  * Fix path parsing for add command (@djdv)
  * namesys: Make paths with multiple segments work. Fixes #2059 (@Kubuxu)
  * Fix up panic catching in http handler funcs (@whyrusleeping)
  * Add correct access control headers to the default api config (@dignifiedquire)
  * Fix closenotify by not sending empty file set (@whyrusleeping)

* Tool Changes
  * Have install.sh use the full path to udfs binary if detected (@jedahan)
  * Install daemon system-wide if on El Capitan (@jedahan)
  * makefile: add -ldflags to install and nofuse tasks (@lgierth)

* General Codebase
  * Clean up http client code (@whyrusleeping)
  * Move api version check to header (@rht)

* Documentation
  * Improved release checklist (@jbenet)
  * Added quotes around command in long description (@RichardLitt)
  * Added a shutdown note to daemon description (@RichardLitt)

* Testing
  * t0080: improve last tests (@chriscool)
  * t0080: improve 'udfs refs --unique' test (@chriscool)
  * Fix t.Fatal usage in goroutines (@chriscool)
  * Add docker testing support to sharness (@chriscool)
  * sharness: add t0300-docker-image.sh (@chriscool)
  * Included more namesys tests. (@Kubuxu)
  * Add sharness test to verify requests look good (@whyrusleeping)
  * Re-enable ipns sharness test now that iptb is fixed (@whyrusleeping)
  * Force use of ipv4 in test (@whyrusleeping)
  * Travis-CI: use go 1.5.2 (@jbenet)

### 0.3.10 - 2015-12-07

This patch update introduces the 'udfs update' command which will be used for
future udfs updates along with a few other bugfixes and documentation
improvements.


* Features
  * support for 'udfs update' to call external binary (@whyrusleeping)
  * cache ipns entries to speed things up a little (@whyrusleeping)
  * add option to version command to print repo version (@whyrusleeping)
  * Add in some more notifications to help profile queries (@whyrusleeping)
  * gateway: add path prefix for directory listings (@lgierth)
  * gateway: add CurrentCommit to /version (@lgierth)

* BugFixes
  * set data and links nil if not present (@whyrusleeping)
  * fix log hanging issue, and implement close-notify for commands (@whyrusleeping)
  * fix dial backoff (@whyrusleeping)
  * proper ndjson implementation (@whyrusleeping)
  * seccat: fix secio context (@lgierth)
  * Add newline to end of the output for a few commands. (@nham)
  * Add fixed period repo GC + test (@rht)

* Tool Changes
  * Allow `udfs cat` on ipns path (@rht)

* General Codebase
  * rewrite of backoff mechanism (@whyrusleeping)
  * refactor net code to use transports, in rough accordance with libp2p (@whyrusleeping)
  * disable building fuse stuff on windows (@whyrusleeping)
  * repo: remove Log config (@lgierth)
  * commands: fix description of --api (@lgierth)

* Documentation
  * --help: Add a note on using UDFS_PATH to the footer of the helptext.  (@sahib)
  * Moved email juan to udfs/contribute (@richardlitt)
  * Added commit sign off section (@richardlitt)
  * Added a security section (@richardlitt)
  * Moved TODO doc to issue #1929 (@richardlitt)

* Testing
  * gateway: add tests for /version (@lgierth)
  * Add gc auto test (@rht)
  * t0020: cleanup dir with bad perms (@chriscool)

Note: this commit introduces fixed-period repo gc, which will trigger gc
after a fixed period of time. This feature is introduced now, disabled by
default, and can be enabled with `udfs daemon --enable-gc`. If all goes well,
in the future, it will be enabled by default.

### 0.3.9 - 2015-10-30

This patch update includes a good number of bugfixes, notably, it fixes
builds on windows, and puts newlines between streaming json objects for a
proper ndjson format.

* Features
  * Writable gateway enabled again (@cryptix)

* Bugfixes
  * fix windows builds (@whyrusleeping)
  * content type on command responses default to text (@whyrusleeping)
  * add check to makefile to ensure windows builds don't fail silently (@whyrusleeping)
  * put newlines between streaming json output objects (@whyrusleeping)
  * fix streaming output to flush per write (@whyrusleeping)
  * purposely fail builds pre go1.5 (@whyrusleeping)
  * fix udfs id <self> (@whyrusleeping)
  * fix a few race conditions in mocknet (@whyrusleeping)
  * fix makefile failing when not in a git repo (@whyrusleeping)
  * fix cli flag orders (long, short) (@rht)
  * fix races in http cors (@miolini)
  * small webui update (some bugfixes) (@jbenet)

* Tool Changes
  * make swarm connect return an error when it fails (@whyrusleeping)
  * Add short flag for `udfs ls --headers` (v for verbose) (@rht)

* General Codebase
  * bitswap: clean log printf and humanize dup data count (@cryptix)
  * config: update pluto's peerID (@lgierth)
  * config: update bootstrap list hostname (@lgierth)

* Documentation
  * Pared down contribute to link to new go guidelines (@richardlitt)

* Testing
  * t0010: add tests for 'udfs commands --flags' (@chriscool)
  * ipns_test: fix namesys.NewNameSystem() call (@chriscool)
  * t0060: fail if no nc (@chriscool)

### 0.3.8 - 2015-10-09

This patch update includes changes to make ipns more consistent and reliable,
symlink support in unixfs, mild performance improvements, new tooling features,
a plethora of bugfixes, and greatly improved tests.

NOTICE: Version 0.3.8 also requires golang version 1.5.1 or higher.

* Bugfixes
  * refactor ipns to be more consistent and reliable (@whyrusleeping)
  * fix 'udfs refs' json output (@whyrusleeping)
  * fix setting null config maps (@rht)
  * fix output of dht commands (@whyrusleeping)
  * fix NAT spam dialing (@whyrusleeping)
  * fix random panics on 32 bit systems (@whyrusleeping)
  * limit total number of network fd's (@whyrusleeping)
  * fix http api content type (@WeMeetAgain)
  * fix writing of api file for port zero daemons (@whyrusleeping)
  * windows connection refused fixes (@mjanczyk)
  * use go1.5's built in trailers, no more failures (@whyrusleeping)
  * fix random bitswap hangs (@whyrusleeping)
  * rate limit fd usage (@whyrusleeping)
  * fix panic in bitswap ratelimiting (@whyrusleeping)

* Tool Changes
  * --empty-repo option for init (@prusnak)
  * implement symlinks (@whyrusleeping)
  * improve cmds lib files processing (@rht)
  * properly return errors through commands (@whyrusleeping)
  * bitswap unwant command (@whyrusleeping)
  * tar add/cat commands (@whyrusleeping)
  * fix gzip compression in get (@klauspost)
  * bitswap stat logs wasted bytes (@whyrusleeping)
  * resolve command now uses core.Resolve (@rht)
  * add `--local` flag to 'name resolve' (@whyrusleeping)
  * add `udfs diag sys` command for debugging help (@whyrusleeping)

* General Codebase
  * improvements to dag editor (@whyrusleeping)
  * swarm IPv6 in default config (Baptiste Jonglez)
  * improve dir listing css (@rht)
  * removed elliptic.P224 usage (@prusnak)
  * improve bitswap providing speed (@jbenet)
  * print panics that occur in cmds lib (@whyrusleeping)
  * udfs api check test fixes (@rht)
  * update peerstream and datastore (@whyrusleeping)
  * cleaned up tar-reader code (@jbenet)
  * write context into coreunix.Cat (@rht)
  * move assets to separate repo (@rht)
  * fix proc/ctx wiring in bitswap (@jbenet)
  * rabin fingerprinting chunker (@whyrusleeping)
  * better notification on daemon ready (@rht)
  * coreunix cat cleanup (@rht)
  * extract logging into go-log (@whyrusleeping)
  * blockservice.New no longer errors (@whyrusleeping)
  * refactor udfs get (@rht)
  * readonly api on gateway (@rht)
  * cleanup context usage all over (@rht)
  * add xml decoding to 'object put' (@ForrestWeston)
  * replace nodebuilder with NewNode method (@whyrusleeping)
  * add metrics to http handlers (@lgierth)
  * rm blockservice workers (@whyrusleeping)
  * decompose maybeGzWriter (@rht)
  * makefile sets git commit sha on build (@CaioAlonso)

* Documentation
  * add contribute file (@RichardLitt)
  * add go devel guide to contribute.md (@whyrusleeping)

* Testing
  * fix mock notifs test (@whyrusleeping)
  * test utf8 with object cmd (@chriscool)
  * make mocknet conn close idempotent (@jbenet)
  * fix fuse tests (@pnelson)
  * improve sharness test quoting (@chriscool)
  * sharness tests for chunker and add-cat (@rht)
  * generalize peerid check in sharness (@chriscool)
  * test_cmp argument cleanup (@chriscool)

### 0.3.7 - 2015-08-02

This patch update fixes a problem we introduced in 0.3.6 and did not
catch: the webui failed to work with out-of-the-box CORS configs.
This has been fixed and now should work correctly. @jbenet

### 0.3.6 - 2015-07-30

This patch improves the resource consumption of go-udfs,
introduces a few new options on the CLI, and also
fixes (yet again) windows builds.

* Resource consumption:
  * fixed goprocess memory leak @rht
  * implement batching on datastore @whyrusleeping
  * Fix bitswap memory leak @whyrusleeping
  * let bitswap ignore temporary write errors @whyrusleeping
  * remove logging to disk in favor of api endpoint @whyrusleeping
  * --only-hash option for add to skip writing to disk @whyrusleeping

* Tool changes
  * improved `udfs daemon` output with all addresses @jbenet
  * improved `udfs id -f` output, added `<addrs>` and  `\n \t` support @jbenet
  * `udfs swarm addrs local` now shows the local node's addrs @jbenet
  * improved config json parsing @rht
  * improved Dockerfile to use alpine linux @Luzifer @lgierth
  * improved bash completion @MichaelMure
  * Improved 404 for gateway @cryptix
  * add unixfs ls to list correct filesizes @wking
  * ignore hidden files by default @gatesvp
  * global --timeout flag @whyrusleeping
  * fix random API failures by closing resp bodies @whyrusleeping
  * udfs swarm filters @whyrusleeping
  * api returns errors in http trailers @whyrusleeping @jbenet
  * `udfs patch` learned to create intermediate nodes @whyrusleeping
  * `udfs object stat` now shows Hash @whyrusleeping
  * `udfs cat` now clears progressbar on exit @rht
  * `udfs add -w -r <dir>` now wraps directories @jbenet
  * `udfs add -w <file1> <file2>` now wraps with one dir @jbenet
  * API + Gateway now support arbitrary HTTP Headers from config @jbenet
  * API now supports CORS properly from config @jbenet
  * **Deprecated:** `API_ORIGIN` env var (use config, see `udfs daemon --help`) @jbenet

* General Codebase
  * `nofuse` tag for windows @Luzifer
  * improved `udfs add` code @gatesvp
  * started requiring license trailers @chriscool @jbenet
  * removed CtxCloser for goprocess @rht
  * remove deadcode @lgierth @whyrusleeping
  * reduced number of logging libs to 2 (soon to be 1) @rht
  * dial address filtering @whyrusleeping
  * prometheus metrics @lgierth
  * new index page for gateway @krl @cryptix
  * move ping to separate protocol @whyrusleeping
  * add events to bitswap for a dashboard @whyrusleeping
  * add latency and bandwidth options to mocknet @heems
  * levenshtein distance cmd autosuggest @sbruce
  * refactor/cleanup of cmds http handler @whyrusleeping
  * cmds http stream reports errors in trailers @whyrusleeping

* Bugfixes
  * fixed path resolution and validation @rht
  * fixed `udfs get -C` output and progress bar @rht
  * Fixed install pkg dist bug @jbenet @Luzifer
  * Fix `udfs get` silent failure   @whyrusleeping
  * `udfs get` tarx no longer times out @jbenet
  * `udfs refs -r -u` is now correct @gatesvp
  * Fix `udfs add -w -r <dir>` wrapping bugs @jbenet
  * Fixed FUSE unmount failures @jbenet
  * Fixed `udfs log tail` command (api + cli) @whyrusleeping

* Testing
  * sharness updates @chriscool
  * ability to disable secio for testing @jbenet
  * fixed many random test failures, more reliable CI @whyrusleeping
  * Fixed racey notifier failures @whyrusleeping
  * `udfs refs -r -u` test cases @jbenet
  * Fix failing pinning test @jbenet
  * Better CORS + Referer tests @jbenet
  * Added reversible gc test @rht
  * Fixed bugs in FUSE IPNS tests @whyrusleeping
  * Fixed bugs in FUSE UDFS tests @jbenet
  * Added `random-files` tool for easier sharness tests @jbenet

* Documentation
  * Add link to init system examples @slang800
  * Add CORS documentation to daemon init @carver  (Note: this will change soon)

### 0.3.5 - 2015-06-11

This patch improves overall stability and performance

* added 'object patch' and 'object new' commands @whyrusleeping
* improved symmetric NAT avoidance @jbenet
* move util.Key to blocks.Key @whyrusleeping
* fix memory leak in provider store @whyrusleeping
* updated webui to 0.2.0 @krl
* improved bitswap performance @whyrusleeping
* update fuse lib @cryptix
* fix path resolution @wking
* implement test_seq() in sharness @chriscool
* improve parsing of stdin for commands @chriscool
* fix 'udfs refs' failing silently @whyrusleeping
* fix serial dialing bug @jbenet
* improved testing @chriscool @rht @jbenet
* fixed domain resolving @luzifer
* fix parsing of unwanted stdin @lgierth
* added CORS handlers to gateway @NodeGuy
* added `udfs daemon --unrestricted-api` option @krl
* general cleanup of dependencies

### 0.3.4 - 2015-05-10

* fix ipns append bug @whyrusleeping
* fix out of memory panic @whyrusleeping
* add in expvar metrics @tv42
* bitswap improvements @whyrusleeping
* fix write-cache in blockstore @tv42
* vendoring cleanup @cryptix
* added `launchctl` plist for OSX @grncdr
* improved Dockerfile, changed root and mount paths @ehd
* improved `pin ls` output to show types @vitorbaptista

### 0.3.3 - 2015-04-28

This patch update fixes various issues, in particular:
- windows support (0.3.0 had broken it)
- commandline parses spaces correctly.

* much improved commandline parsing by @AtnNn
* improved dockerfile by @luzifer
* add cmd cleanup by @wking
* fix flatfs windows support by @tv42 and @gatesvp
* test case improvements by @chriscool
* ipns resolution timeout bug fix by @whyrusleeping
* new cluster tests with iptb by @whyrusleeping
* fix log callstack printing bug by @whyrusleeping
* document bash completion by @dylanPowers

### 0.3.2 - 2015-04-22

This patch update implements multicast dns as well as fxing a few test issues.

* implement mdns peer discovery @whyrusleeping
* fix mounting issues in sharness tests @chriscool

### 0.3.1 - 2015-04-21

This patch update fixes a few bugs:

* harden shutdown logic by @torarnv
* daemon locking fixes by @travisperson
* don't re-add entire dirs by @whyrusleeping
* tests now wait for graceful shutdown by @jbenet
* default key size is now 2048 by @jbenet

### 0.3.0 - 2015-04-20

We've just released version 0.3.0, which contains many
performance improvements, bugfixes, and new features.
Perhaps the most noticeable change is moving block storage
from leveldb to flat files in the filesystem.

What to expect:

* _much faster_ performance

* Repo format 2
  * moved default location from ~/.go-udfs -> ~/.udfs
  * renamed lock filename daemon.lock -> repo.lock
  * now using a flat-file datastore for local blocks

* Fixed lots of bugs
  * proper udfs-path in various commands
  * fixed two pinning bugs (recursive pins)
  * increased yamux streams window (for speed)
  * increased bitswap workers (+ env var)
  * fixed memory leaks
  * udfs add error returns
  * daemon exit bugfix
  * set proper UID and GID on fuse mounts

* Gateway
  * Added support for HEAD requests

* configuration
  * env var to turn off SO_REUSEPORT: UDFS_REUSEPORT=false
  * env var to increase bitswap workers: UDFS_BITSWAP_TASK_WORKERS=n

* other
  * bash completion is now available
  * udfs stats bw -- bandwidth meetering

And many more things.

### 0.2.3 - 2015-03-01

* Alpha Release

### 2015-01-31:

* bootstrap addresses now have .../udfs/... in format
  config file Bootstrap field changed accordingly. users
  can upgrade cleanly with:

      udfs bootstrap >boostrap_peers
      udfs bootstrap rm --all
      <install new udfs>
      <manually add .../udfs/... to addrs in bootstrap_peers>
      udfs bootstrap add <bootstrap_peers
