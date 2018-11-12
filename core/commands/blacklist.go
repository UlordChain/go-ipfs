package commands

import (
	"encoding/csv"

	"context"

	"io"

	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	format "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"

	"gx/ipfs/QmNueRyPRQiV7PUEpnP4GgGLuK1rKQLaRW7sfPvUetYig1/go-ipfs-cmds"
	"gx/ipfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-ipfs-cmdkit"

	"time"

	"fmt"

	"sort"

	"io/ioutil"

	"os"

	"path/filepath"

	"github.com/udfs/go-udfs/core"
	"github.com/udfs/go-udfs/core/corerepo"
	"github.com/udfs/go-udfs/core/coreunix"
	"github.com/udfs/go-udfs/merkledag"
	"github.com/udfs/go-udfs/path"
	"github.com/pkg/errors"
)

const (
	blacklistDir     = "/ipns/QmbETUnWes7zdwZkkMGgPRtpZAYpFPxrUrCYy7fWi7JjFY/blacklistdir"
	timeFormatLayout = "2006-01-02 15:04:05 -0700 MST"
)

var (
	lastHandleTimeFile = filepath.Join(os.TempDir(), "blacklist_last_handle_time")
	lastHandleTime     = time.Unix(0, 0)
	period             = 24 * time.Hour
)

func init() {
	t, err := ioutil.ReadFile(lastHandleTimeFile)
	if err != nil && !os.IsNotExist(err) {
		log.Error("read last blacklist name file failed:", err.Error())
	} else if len(t) > 0 {
		last, err := time.Parse(timeFormatLayout, string(t))
		if err != nil {
			log.Error("parse blacklist last handle time failed:", err.Error())
		} else {
			lastHandleTime = last
			fmt.Println("blacklist last handle time:", string(t))
		}
	}
}

var BlacklistCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "Run blacklist service.",
		ShortDescription: "run a blacklist refresh operation right now.",
	},

	Arguments: []cmdkit.Argument{},
	Options:   []cmdkit.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		node, err := GetNode(env)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		if !node.OnlineMode() {
			if err := node.SetupOfflineRouting(); err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
		}

		err = refreshBlacklist(req.Context, node, 1)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

	},
}

func getBlacklistFiles(ctx context.Context, n *core.IpfsNode) ([]*format.Link, error) {
	node, err := core.Resolve(ctx, n.Namesys, n.Resolver, blacklistDir)
	if err != nil {
		return nil, errors.Errorf("read blacklist directory failed: %v\n", err.Error())
	}

	if len(node.Links()) == 0 {
		return nil, nil
	}

	links := node.Links()
	sort.Stable(merkledag.LinkSlice(links))

	deadline := int64(0)
	if lastHandleTime.Unix() > 0 {
		deadline = lastHandleTime.Unix() - int64(period.Seconds())
	}
	deadlineStr := fmt.Sprintf("%10d", deadline)

	for i, link := range links {
		if link.Name <= deadlineStr {
			continue
		}

		return links[i:], nil
	}
	return nil, nil
}

func refreshBlacklist(ctx context.Context, n *core.IpfsNode, minFailed int) error {
	links, err := getBlacklistFiles(ctx, n)
	if err != nil {
		return errors.Wrap(err, "get blacklist files failed")
	}

	if len(links) == 0 {
		log.Debug("no new blacklist file need to handle, last blacklist handle time :", lastHandleTime)
		return nil
	}

	newBlacklistHandleTime := lastHandleTime
	defer func() {
		if newBlacklistHandleTime == lastHandleTime {
			return
		}
		lastHandleTime = newBlacklistHandleTime

		newTimeStr := newBlacklistHandleTime.Format(timeFormatLayout)

		// save new blacklist name to file
		e := ioutil.WriteFile(lastHandleTimeFile, []byte(newTimeStr), 0644)
		if e != nil {
			log.Warningf("save new blacklist name %s to file faied: %v\n", newTimeStr, e.Error())
		}

	}()

	for _, link := range links {
		log.Debug("handle blacklist file ", link.Name)
		err = handleBlacklistFile(ctx, n, minFailed, link.Cid)
		if err != nil {
			return errors.Wrapf(err, "refresh blacklist file %s failed", link.Name)
		}

		newBlacklistHandleTime = time.Now()
	}

	return nil
}

func handleBlacklistFile(ctx context.Context, n *core.IpfsNode, minFailed int, c *cid.Cid) error {

	dagReader, err := coreunix.Cat(ctx, n, path.FromCid(c).String())
	if err != nil {
		return errors.Errorf("read blacklist failed: %v\n", err.Error())
	}

	csvReader := csv.NewReader(dagReader)
	failedCount := 0
	for {
		record, err := csvReader.Read()
		if record == nil || err == io.EOF {
			return nil
		}

		if err != nil {
			return errors.Errorf("read record from blacklist failed: %v\n", err.Error())
		}
		log.Debug("blacklist record:", record)

		err = handleBlacklistRecord(ctx, n, record)
		if err != nil {
			failedCount++
			if minFailed > 0 && failedCount >= minFailed {
				return err
			}

			log.Error(err)
		}
	}

}

func RunBlacklistRefreshService(ctx context.Context, n *core.IpfsNode) error {
	conf, err := n.Repo.Config()
	if err != nil {
		return errors.Wrap(err, "got config failed")
	}

	dur := 10 * time.Second
	if conf.Blacklist.Interval != "" {
		d, err := time.ParseDuration(conf.Blacklist.Interval)
		if err != nil {
			return errors.Wrap(err, "parse config.Blacklist.Interval failed")
		}
		dur = d
	}
	if conf.Blacklist.Period != "" {
		p, err := time.ParseDuration(conf.Blacklist.Period)
		if err != nil {
			return errors.Wrap(err, "parse config.Blacklist.Period failed")
		}

		period = p
	}

	log.Debug("current blacklist file handle interval = ", dur)
	log.Debug("current blacklist file invalid period = ", period)

	tm := time.NewTimer(dur)

	go func() {
		defer tm.Stop()

		for {
			select {
			case <-tm.C:
				refreshBlacklist(ctx, n, -1)

				tm.Reset(dur)

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil

}

func handleBlacklistRecord(ctx context.Context, n *core.IpfsNode, record []string) error {
	c, err := pathToCid(record[0])
	if err != nil {
		return errors.Errorf("got blacklist record cid failed from %v: %v\n", record, err.Error())
	}

	has, err := n.Blockstore.Has(c)
	if err != nil {
		return errors.Errorf("check blacklist record cid %s exist failed: %v\n", c.String(), err.Error())
	}

	if !has {
		return nil
	}

	_, pined, err := n.Pinning.IsPinned(c)
	if err != nil {
		return errors.Errorf("check blacklist record cid %s pined failed: %v\n", c.String(), err.Error())
	}

	if pined {
		_, err := corerepo.Unpin(n, ctx, record[:1], true)
		if err != nil {
			return errors.Errorf("unpin blacklist record cid %s failed: %v\n", c.String(), err.Error())
		}
		fmt.Println("unpin ", record[0])
	}

	err = corerepo.Remove(n, ctx, []*cid.Cid{c}, true, false)
	if err != nil {
		return errors.Errorf("blacklist record cid %s remove from repo failed: %v\n", c.String(), err.Error())
	}

	fmt.Println("remove ", c.String())

	return nil
}

func pathToCid(pstr string) (*cid.Cid, error) {
	p, err := path.ParsePath(pstr)
	if err != nil {
		return nil, err
	}

	c, _, err := path.SplitAbsPath(p)
	if err != nil {
		return nil, err
	}

	return c, nil
}
