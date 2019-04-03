package commands

import (
	"fmt"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	"github.com/ipfs/go-ipfs/core/commands/sms"
	"github.com/ipfs/go-ipfs/core/coreapi/interface"
	"github.com/pkg/errors"
	"gx/ipfs/QmT3rzed1ppXefourpmoZ7tyVQfsGPQZ1pHDngLmCvXxd3/go-path"
	"os"
	"strings"

	"github.com/ipfs/go-ipfs/core/commands/e"
	"github.com/ipfs/go-ipfs/core/corerepo"
	"gx/ipfs/QmSXUokcP4TJpFfqozT69AVAYRtzXVMUjzQVkYX41R9Svs/go-ipfs-cmds"
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
		//cmdkit.BoolOption("clear", "", "Clear the cache from repo.").WithDefault(false),
		cmdkit.StringOption(tokenOptionName, "The token value for verify"),
	},
	Type: PinOutput{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		if len(req.Arguments) > 1 {
			return errors.New("Do not allow multiple files to be remove once now")
		}

		// verify token
		tokenInf := req.Options[tokenOptionName]
		if tokenInf == nil {
			return errors.New("must set option token.")
		}
		token := tokenInf.(string)
		p := path.Path(req.Arguments[0])

		err := sms.Delete(token, p.String())
		if err != nil {
			return err
		}

		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		api, err := cmdenv.GetApi(env)
		if err != nil {
			return err
		}

		// set recursive flag
		recursive, _ := req.Options["recursive"].(bool)

		args := req.Arguments
		cids := make([]cid.Cid, len(args))
		for i, a := range args {
			pth, err := iface.ParsePath(a)
			if err != nil {
				return err
			}

			c, err := api.ResolvePath(req.Context, pth)
			if err != nil {
				return err
			}

			cids[i] = c.Cid()
		}

		err = corerepo.Remove(node, req.Context, cids, recursive, false)
		if err != nil {
			return err
		}

		_, err = corerepo.Unpin(node, api, req.Context, req.Arguments, recursive)
		if err != nil && !strings.Contains(err.Error(), "not pinned"){
			return err
		}

		res.Emit(&PinOutput{cidsToStrings(cids)})
		return nil
	},
	PostRun: cmds.PostRunMap{
		cmds.CLI: func(res cmds.Response, re cmds.ResponseEmitter) error {
			v, err := res.Next()
			if err != nil {
				return err
			}

			added, ok := v.(*PinOutput)
			if !ok {
				return e.New(e.TypeErr(added, v))
			}

			for _, k := range added.Pins {
				fmt.Fprintf(os.Stdout, "removed %s\n", k)
			}
			return nil
		},
	},
}
