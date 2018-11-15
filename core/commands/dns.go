package commands

import (
	"io"
	"strings"

	cmds "github.com/udfs/go-udfs/commands"
	e "github.com/udfs/go-udfs/core/commands/e"
	namesys "github.com/udfs/go-udfs/namesys"
	nsopts "github.com/udfs/go-udfs/namesys/opts"

	"gx/udfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-udfs-cmdkit"
)

var DNSCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Resolve DNS links.",
		ShortDescription: `
Multihashes are hard to remember, but domain names are usually easy to
remember.  To create memorable aliases for multihashes, DNS TXT
records can point to other DNS links, UDFS objects, IPNS keys, etc.
This command resolves those links to the referenced object.
`,
		LongDescription: `
Multihashes are hard to remember, but domain names are usually easy to
remember.  To create memorable aliases for multihashes, DNS TXT
records can point to other DNS links, UDFS objects, IPNS keys, etc.
This command resolves those links to the referenced object.

For example, with this DNS TXT record:

	> dig +short TXT _dnslink.udfs.io
	dnslink=/udfs/QmRzTuh2Lpuz7Gr39stNr6mTFdqAghsZec1JoUnfySUzcy

The resolver will give:

	> udfs dns udfs.io
	/udfs/QmRzTuh2Lpuz7Gr39stNr6mTFdqAghsZec1JoUnfySUzcy

The resolver can recursively resolve:

	> dig +short TXT recursive.udfs.io
	dnslink=/ipns/udfs.io
	> udfs dns -r recursive.udfs.io
	/udfs/QmRzTuh2Lpuz7Gr39stNr6mTFdqAghsZec1JoUnfySUzcy
`,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("domain-name", true, false, "The domain-name name to resolve.").EnableStdin(),
	},
	Options: []cmdkit.Option{
		cmdkit.BoolOption("recursive", "r", "Resolve until the result is not a DNS link."),
	},
	Run: func(req cmds.Request, res cmds.Response) {

		recursive, _, _ := req.Option("recursive").Bool()
		name := req.Arguments()[0]
		resolver := namesys.NewDNSResolver()

		var ropts []nsopts.ResolveOpt
		if !recursive {
			ropts = append(ropts, nsopts.Depth(1))
		}

		output, err := resolver.Resolve(req.Context(), name, ropts...)
		if err == namesys.ErrResolveFailed {
			res.SetError(err, cmdkit.ErrNotFound)
			return
		}
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}
		res.SetOutput(&ResolvedPath{output})
	},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v, err := unwrapOutput(res.Output())
			if err != nil {
				return nil, err
			}

			output, ok := v.(*ResolvedPath)
			if !ok {
				return nil, e.TypeErr(output, v)
			}
			return strings.NewReader(output.Path.String() + "\n"), nil
		},
	},
	Type: ResolvedPath{},
}
