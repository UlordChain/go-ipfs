package commands

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	"github.com/ipfs/go-ipfs/core/commands/sms"
	"github.com/ipfs/go-ipfs/core/coreapi/interface"

	"gx/ipfs/QmSXUokcP4TJpFfqozT69AVAYRtzXVMUjzQVkYX41R9Svs/go-ipfs-cmds"
	"gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
)


var CacheCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "Cache IPFS object data.",
		ShortDescription: "Cache the data contained by an IPFS or IPNS object(s) at the given path from the given address.",
	},

	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("address", true, false, "Address of peer to connect to cache from."),
		cmdkit.StringArg("ipfs-path", true, false, "The path to the IPFS object(s) to be cached."),
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption(accountOptionName, "Account of user to check"),
		cmdkit.StringOption(tokenOptionName, "The token value for verify"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		// verify token
		tokenInf := req.Options[tokenOptionName]
		if tokenInf == nil {
			return errors.New("must set option token.")
		}
		token := tokenInf.(string)

		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		cfg, _ := node.Repo.Config()
		var(
			account string
			check string
		)
		if cfg.UOSCheck.Enable {
			acc := req.Options[accountOptionName]
			if acc == nil {
				return errors.New("must set option account.")
			}
			account = acc.(string)

			h := req.Options[checkOptionName]
			if h == nil {
				return errors.New("must set option check.")
			}
			check = h.(string)

			_, err = ValidOnUOS(&cfg.UOSCheck, account, check)
			if err != nil {
				return errors.Wrap(err, "valid failed")
			}
		}

		api, err := cmdenv.GetApi(env)
		if err != nil {
			return err
		}

		if !node.OnlineMode() {
			if err := node.SetupOfflineRouting(); err != nil {
				return err
			}
		}

		addr := req.Arguments[0]

		pis, err := peersWithAddresses([]string{addr})
		if err != nil {
			return err
		}

		pi := pis[0]
		err = api.Swarm().Connect(req.Context, pi)
		if err != nil {
			return fmt.Errorf("connect %s failure: %s", pi.ID.Pretty(), err)
		}

		fpath, err := iface.ParsePath(req.Arguments[1])
		if err != nil {
			return errors.Wrap(err, "pase ipfs-path failed")
		}

		tmCtx, _ := context.WithTimeout(req.Context, 10*time.Second)
		_, err = api.Unixfs().Get(tmCtx, fpath)
		if err != nil {
			return errors.Wrap(err, "get ipfs-path object failed")
		}

		err = sms.Cache(token, fpath.String(), pi.ID.Pretty())
		if err != nil {
			return errors.Wrap(err, "sms cache notify failed")
		}

		ret := fmt.Sprintf("cache %s from %s success.", req.Arguments[1], pi.ID.Pretty())
		return cmds.EmitOnce(res, &stringList{[]string{ret}})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(stringListEncoder),
	},
	Type: stringList{},
}
