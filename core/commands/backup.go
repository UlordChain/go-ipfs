package commands

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	inet "gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"

	"gx/ipfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-ipfs-cmdkit"

	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"

	"strings"

	cmds "github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/commands/e"
	"github.com/ipfs/go-ipfs/core/corerepo"
	"github.com/ipfs/go-ipfs/core/coreunix"
)

const ProtocolBackup protocol.ID = "/backup/0.0.1"
const numberForBackup int = 1
const timeoutForLookup = 1 * time.Minute

var BackupCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "Backup objects to remote node storage.",
		ShortDescription: "Stores an IPFS object(s) from a given path locally to remote disk.",
	},

	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("ipfs-path", true, false, "Path to object(s) to be pinned.").EnableStdin(),
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption(accountOptionName, "Account of user to check"),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		acc := req.Options()[accountOptionName]
		if acc == nil {
			res.SetError(errors.New("must set option account."), cmdkit.ErrNormal)
			return
		}
		account := acc.(string)
		check := req.StringArguments()[0]

		err := ValidOnUOS(account, check)
		if err != nil {
			res.SetError(errors.Wrap(err, "valid failed"), cmdkit.ErrNormal)
			return
		}

		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		if n.Routing == nil {
			res.SetError(errNotOnline, cmdkit.ErrNormal)
			return
		}

		// get cid
		c, err := cid.Decode(req.Arguments()[0])
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		output, err := backupFunc(n, c, account)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		res.SetOutput(output)
	},
	Type: coreunix.BackupOutput{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			v, err := unwrapOutput(res.Output())
			if err != nil {
				return nil, err
			}

			out, ok := v.(*coreunix.BackupOutput)
			if !ok {
				return nil, e.TypeErr(out, v)
			}

			buf := new(bytes.Buffer)
			for _, s := range out.Success {
				fmt.Fprintf(buf, "backup success to %s\n", s.ID)
			}
			for _, f := range out.Failed {
				fmt.Fprintf(buf, "backup failed to %s : %s\n", f.ID, f.Msg)
			}

			return buf, nil
		},
	},
}

func backupFunc(n *core.IpfsNode, c *cid.Cid, account string) (*coreunix.BackupOutput, error) {
	// get peers for backup
	toctx, cancel := context.WithTimeout(n.Context(), timeoutForLookup)
	defer cancel()
	closestPeers, err := n.DHT.GetClosestMasterPeers(toctx, c.KeyString())
	if err != nil {
		return nil, errors.Wrap(err, "got closest master peers timeout")
	}

	peers := make(map[peer.ID]struct{}, 0)
	for p := range closestPeers {
		peers[p] = struct{}{}

		if len(peers) >= numberForBackup {
			cancel()
			break
		}
	}

	if len(peers) < numberForBackup {
		return nil, errors.Errorf("Failed to find the minimum number of closest peers required: %d/%d", len(peers),
			numberForBackup)
	}

	log.Debug("found the peers to backup:", peers)
	peersForBackup := peers

	// do backup
	results := make(chan *coreunix.BackupResult, len(peersForBackup))
	var wg sync.WaitGroup
	for p := range peersForBackup {
		wg.Add(1)
		go func(id peer.ID) {
			e := doBackup(n, id, c, account)
			if e != nil {
				results <- &coreunix.BackupResult{
					ID:  id.Pretty(),
					Msg: e.Error(),
				}
			} else {
				results <- &coreunix.BackupResult{
					ID: id.Pretty(),
				}
			}
			wg.Done()
		}(p)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	output := &coreunix.BackupOutput{}
	for r := range results {
		if r.Msg != "" {
			output.Failed = append(output.Failed, r)
		} else {
			output.Success = append(output.Success, r)
		}
	}

	if len(output.Failed) > 0 {
		return nil, errors.New("backup failed")
	}

	return output, nil
}

func doBackup(n *core.IpfsNode, id peer.ID, c *cid.Cid, account string) error {
	s, err := n.PeerHost.NewStream(n.Context(), id, ProtocolBackup)
	if err != nil {
		return err
	}
	defer s.Close()

	// TODO: consider to use protobuf, now just direct send the cid and account
	_, err = s.Write([]byte(account + "," + c.String() + "\n"))
	if err != nil {
		return err
	}

	// read result
	buf := bufio.NewReader(s)
	bs, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	if bs == "\n" {
		log.Debugf("backup %s successd to node %s\n", c.String(), id.Pretty())
		return nil
	}

	return errors.New(bs)
}

func SetupBackupHandler(node *core.IpfsNode) {
	node.PeerHost.SetStreamHandler(ProtocolBackup, func(s inet.Stream) {
		var errRet error
		defer func() {
			var e error
			if errRet != nil {
				log.Error("backup-hander failed: ", errRet.Error())
				_, e = s.Write([]byte(errRet.Error() + "\n"))
			} else {
				_, e = s.Write([]byte("\n"))
			}

			if e != nil {
				log.Error("backup-handler send result failed: ", e.Error())
			} else {
				log.Debug("backup-handler send result success")
			}

			s.Close()
		}()

		select {
		case <-node.Context().Done():
			return
		default:
		}

		log.Debug("backup-handler receive request from", s.Conn().RemoteMultiaddr().String(), "/", s.Conn().RemotePeer().Pretty())

		// TODO: consider use protobuf, now just direct get cid
		buf := bufio.NewReader(s)
		bs, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			errRet = errors.Wrap(err, "backup-handler read bytes failed")
			return
		}

		n := strings.Index(bs, ",")
		if n == -1 {
			errRet = errors.Wrap(err, "backup-handler read ',' failed")
			return
		}

		account := bs[:n]

		c, err := cid.Decode(bs[n+1 : len(bs)-1])
		if err != nil {
			errRet = errors.Wrap(err, "decode cid failed")
			return
		}
		log.Debug("backup-handler cid=", c.String())

		err = ValidOnUOS(account, c.String())
		if err != nil {
			errRet = errors.Wrapf(err, "valid ipfs-hash for user=%s error", account)
			return
		}

		// do pin add
		defer node.Blockstore.PinLock().Unlock()

		_, err = corerepo.Pin(node, node.Context(), []string{c.String()}, true)
		if err != nil {
			errRet = errors.Wrapf(err, "backup-handler run pin command for %s failed", c.String())
			return
		}

		log.Debugf("backup-handler run pin add %s success\n", c.String())
	})
}
