// +build !windows,!nofuse

package commands

import (
	"fmt"
	"io"
	"strings"

	cmds "github.com/udfs/go-udfs/commands"
	e "github.com/udfs/go-udfs/core/commands/e"
	nodeMount "github.com/udfs/go-udfs/fuse/node"
	config "github.com/udfs/go-udfs/repo/config"

	"gx/udfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-udfs-cmdkit"
)

var MountCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Mounts UDFS to the filesystem (read-only).",
		ShortDescription: `
Mount UDFS at a read-only mountpoint on the OS (default: /udfs and /ipns).
All UDFS objects will be accessible under that directory. Note that the
root will not be listable, as it is virtual. Access known paths directly.

You may have to create /udfs and /ipns before using 'udfs mount':

> sudo mkdir /udfs /ipns
> sudo chown $(whoami) /udfs /ipns
> udfs daemon &
> udfs mount
`,
		LongDescription: `
Mount UDFS at a read-only mountpoint on the OS. The default, /udfs and /ipns,
are set in the configuration file, but can be overriden by the options.
All UDFS objects will be accessible under this directory. Note that the
root will not be listable, as it is virtual. Access known paths directly.

You may have to create /udfs and /ipns before using 'udfs mount':

> sudo mkdir /udfs /ipns
> sudo chown $(whoami) /udfs /ipns
> udfs daemon &
> udfs mount

Example:

# setup
> mkdir foo
> echo "baz" > foo/bar
> udfs add -r foo
added QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR foo/bar
added QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC foo
> udfs ls QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC
QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR 12 bar
> udfs cat QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR
baz

# mount
> udfs daemon &
> udfs mount
UDFS mounted at: /udfs
IPNS mounted at: /ipns
> cd /udfs/QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC
> ls
bar
> cat bar
baz
> cat /udfs/QmSh5e7S6fdcu75LAbXNZAFY2nGyZUJXyLCJDvn2zRkWyC/bar
baz
> cat /udfs/QmWLdkp93sNxGRjnFHPaYg8tCQ35NBY3XPn6KiETd3Z4WR
baz
`,
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption("udfs-path", "f", "The path where UDFS should be mounted."),
		cmdkit.StringOption("ipns-path", "n", "The path where IPNS should be mounted."),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		cfg, err := req.InvocContext().GetConfig()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		node, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		// error if we aren't running node in online mode
		if node.LocalMode() {
			res.SetError(errNotOnline, cmdkit.ErrClient)
			return
		}

		fsdir, found, err := req.Option("f").String()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}
		if !found {
			fsdir = cfg.Mounts.UDFS // use default value
		}

		// get default mount points
		nsdir, found, err := req.Option("n").String()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}
		if !found {
			nsdir = cfg.Mounts.IPNS // NB: be sure to not redeclare!
		}

		err = nodeMount.Mount(node, fsdir, nsdir)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		var output config.Mounts
		output.UDFS = fsdir
		output.IPNS = nsdir
		res.SetOutput(&output)
	},
	Type: config.Mounts{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v, err := unwrapOutput(res.Output())
			if err != nil {
				return nil, err
			}

			mnts, ok := v.(*config.Mounts)
			if !ok {
				return nil, e.TypeErr(mnts, v)
			}

			s := fmt.Sprintf("UDFS mounted at: %s\n", mnts.UDFS)
			s += fmt.Sprintf("IPNS mounted at: %s\n", mnts.IPNS)
			return strings.NewReader(s), nil
		},
	},
}
