package commands

import (
	"bytes"
	"fmt"
	"gx/ipfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-ipfs-cmdkit"
	"io"

	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"

	cmds "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"
	e "github.com/ipfs/go-ipfs/core/commands/e"
	"github.com/ipfs/go-ipfs/core/corerepo"
	"github.com/ipfs/go-ipfs/path"
	"github.com/ipfs/go-ipfs/path/resolver"
	uio "github.com/ipfs/go-ipfs/unixfs/io"
)

var LocalrmCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Remove objects from pin and repo.",
		ShortDescription: `
'ipfs localrm' is a plumbing command that will remove the objects that are pinned and cached.
`,
	},
	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("ipfs-path", true, true, "Path to object(s) to be removed.").EnableStdin(),
	},
	Options: []cmdkit.Option{
		cmdkit.BoolOption("recursive", "r", "Recursively unpin the object linked to by the specified object(s).").WithDefault(true),
		cmdkit.BoolOption("clear", "", "Clear the cache from repo.").WithDefault(false),
	},
	Type: PinOutput{},
	Run: func(req cmds.Request, res cmds.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		// set recursive flag
		recursive, _, err := req.Option("recursive").Bool()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		r := &resolver.Resolver{
			DAG:         n.DAG,
			ResolveOnce: uio.ResolveUnixfsOnce,
		}

		args := req.Arguments()
		cids := make([]*cid.Cid, len(args))
		for i, a := range args {
			p, err := path.ParsePath(a)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}

			k, err := core.ResolveToCid(req.Context(), n.Namesys, r, p)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}

			cids[i] = k
		}

		removed, err := corerepo.Unpin(n, req.Context(), req.Arguments(), recursive)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		err = corerepo.Remove(n, req.Context(), removed, recursive, false)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		res.SetOutput(&PinOutput{cidsToStrings(removed)})
	},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v, err := unwrapOutput(res.Output())
			if err != nil {
				return nil, err
			}

			added, ok := v.(*PinOutput)
			if !ok {
				return nil, e.TypeErr(added, v)
			}

			buf := new(bytes.Buffer)
			for _, k := range added.Pins {
				fmt.Fprintf(buf, "unpinned %s\n", k)
			}
			return buf, nil
		},
	},
}
