package commands

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/commands/e"
	"github.com/ipfs/go-ipfs/core/corerepo"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"

	"gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	"gx/ipfs/QmSXUokcP4TJpFfqozT69AVAYRtzXVMUjzQVkYX41R9Svs/go-ipfs-cmds"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	inet "gx/ipfs/QmXuRkCR7BNQa9uqfpTiFWsTQLzmTWYg91Ja1w95gnqb6u/go-libp2p-net"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
)

const ProtocolBackup protocol.ID = "/backup/0.0.1"
const numberForBackup int = 1
const timeoutForLookup = 1 * time.Minute

var BackupCmd = &commands.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "Backup objects to remote node storage.",
		ShortDescription: "Stores an IPFS object(s) from a given path locally to remote disk.",
	},

	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("ipfs-path", true, false, "Path to object(s) to be pinned.").EnableStdin(),
	},
	Run: func(req commands.Request, res commands.Response) {
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		if n.Routing == nil {
			res.SetError(ErrNotOnline, cmdkit.ErrNormal)
			return
		}

		// get cid
		c, err := cid.Decode(req.Arguments()[0])
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		output, err := backupFunc(n, c)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		res.SetOutput(output)
	},
	Type: coreiface.BackupOutput{},
	Marshalers: commands.MarshalerMap{
		commands.Text: func(res commands.Response) (io.Reader, error) {
			v, err := unwrapOutput(res.Output())
			if err != nil {
				return nil, err
			}

			out, ok := v.(*coreiface.BackupOutput)
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

func backupFunc(n *core.IpfsNode, c cid.Cid) (*coreiface.BackupOutput, error) {
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

	// 发送cid
	results := make(chan *coreiface.BackupResult, len(peersForBackup))
	var wg sync.WaitGroup
	for p := range peersForBackup {
		wg.Add(1)
		go func(id peer.ID) {
			e := doBackup(n, id, c)
			if e != nil {
				results <- &coreiface.BackupResult{
					ID:  id.Pretty(),
					Msg: e.Error(),
				}
			} else {
				results <- &coreiface.BackupResult{
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

	output := &coreiface.BackupOutput{}
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

func doBackup(n *core.IpfsNode, id peer.ID, c cid.Cid) error {
	s, err := n.PeerHost.NewStream(n.Context(), id, ProtocolBackup)
	if err != nil {
		return err
	}
	defer s.Close()

	// TODO: consider to use protobuf, now just direct send the cid
	_, err = s.Write([]byte(c.String() + "\n"))
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

func SetupBackupHandler(env cmds.Environment) {
	node, _ := cmdenv.GetNode(env)
	api, _ := cmdenv.GetApi(env)

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

		c, err := cid.Decode(bs[:len(bs)-1])
		if err != nil {
			errRet = errors.Wrap(err, "decode cid failed")
			return
		}
		log.Debug("backup-handler cid=", c.String())

		// do pin add
		defer node.Blockstore.PinLock().Unlock()

		_, err = corerepo.Pin(node, api, node.Context(), []string{c.String()}, true)
		if err != nil {
			errRet = errors.Wrapf(err, "backup-handler run pin command for %s failed", c.String())
			return
		}

		log.Debugf("backup-handler run pin add %s success\n", c.String())
	})
}
