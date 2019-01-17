package commands

import (
	"bytes"
	"fmt"
	"github.com/ipfs/go-ipfs/core/coreapi/interface"
	"io"

	cmds "github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core/commands/e"
	"github.com/ipfs/go-ipfs/core/corerepo"

	"gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	"gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
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

		api, err := req.InvocContext().GetApi()
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


		args := req.Arguments()
		cids := make([]cid.Cid, len(args))
		for i, a := range args {
			pth, err := iface.ParsePath(a)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}

			c, err := api.ResolvePath(req.Context(), pth)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}

			cids[i] = c.Cid()
		}

		removed, err := corerepo.Unpin(n, api, req.Context(), req.Arguments(), recursive)
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
