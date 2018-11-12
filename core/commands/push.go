package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/udfs/go-udfs/blockservice"
	"github.com/udfs/go-udfs/core"
	"github.com/udfs/go-udfs/core/coreunix"
	"github.com/udfs/go-udfs/filestore"
	dag "github.com/udfs/go-udfs/merkledag"
	dagtest "github.com/udfs/go-udfs/merkledag/test"
	"github.com/udfs/go-udfs/mfs"
	ft "github.com/udfs/go-udfs/unixfs"
	"github.com/pkg/errors"

	"gx/ipfs/QmNueRyPRQiV7PUEpnP4GgGLuK1rKQLaRW7sfPvUetYig1/go-ipfs-cmds"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmPtj12fdwuAqj9sBSTNUxBNu8kCGNp8b3o8yUzMm5GHpq/pb"
	"gx/ipfs/QmS6mo1dPpHdYsVkm27BRZDLxpKBCiJKUH8fHX15XFfMez/go-ipfs-exchange-offline"
	bstore "gx/ipfs/QmadMhXJLHMFjpRmh85XjpmVDkEtQpNYEZNRpWRvYVLrvb/go-ipfs-blockstore"
	"gx/ipfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-ipfs-cmdkit"
	"gx/ipfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-ipfs-cmdkit/files"

	"path/filepath"

	"encoding/csv"

	"strconv"

	"sync"

	"os/exec"

	"syscall"

	"github.com/udfs/go-udfs/core/corerepo"
)

var PushCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Push a file or directory to ipfs.",
		ShortDescription: `
Pushs contents of <path> to ipfs. Use -r to add directories (recursively).
`,
		LongDescription: `
Push do the same thing like command add first (but with default not pin). Then do the same thing like command backup.
`,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.FileArg("path", true, true, "The path to a file to be added to ipfs.").EnableRecursive().EnableStdin(),
	},
	Options: []cmdkit.Option{
		cmds.OptionRecursivePath, // a builtin option that allows recursive paths (-r, --recursive)
		cmdkit.BoolOption(quietOptionName, "q", "Write minimal output."),
		cmdkit.BoolOption(quieterOptionName, "Q", "Write only final hash."),
		cmdkit.BoolOption(silentOptionName, "Write no output."),
		cmdkit.BoolOption(progressOptionName, "p", "Stream progress data."),
		cmdkit.BoolOption(trickleOptionName, "t", "Use trickle-dag format for dag generation."),
		cmdkit.BoolOption(onlyHashOptionName, "n", "Only chunk and hash - do not write to disk."),
		cmdkit.BoolOption(wrapOptionName, "w", "Wrap files with a directory object."),
		cmdkit.BoolOption(hiddenOptionName, "H", "Include files that are hidden. Only takes effect on recursive add."),
		cmdkit.StringOption(chunkerOptionName, "s", "Chunking algorithm, size-[bytes] or rabin-[min]-[avg]-[max]").WithDefault("size-262144"),
		cmdkit.BoolOption(pinOptionName, "Pin this object when pushing.").WithDefault(false),
		cmdkit.BoolOption(rawLeavesOptionName, "Use raw blocks for leaf nodes. (experimental)"),
		cmdkit.BoolOption(noCopyOptionName, "Add the file using filestore. Implies raw-leaves. (experimental)"),
		cmdkit.BoolOption(fstoreCacheOptionName, "Check the filestore for pre-existing blocks. (experimental)"),
		cmdkit.IntOption(cidVersionOptionName, "CID version. Defaults to 0 unless an option that depends on CIDv1 is passed. (experimental)"),
		cmdkit.StringOption(hashOptionName, "Hash function to use. Implies CIDv1 if not sha2-256. (experimental)").WithDefault("sha2-256"),
	},
	PreRun: func(req *cmds.Request, env cmds.Environment) error {
		quiet, _ := req.Options[quietOptionName].(bool)
		quieter, _ := req.Options[quieterOptionName].(bool)
		quiet = quiet || quieter

		silent, _ := req.Options[silentOptionName].(bool)

		if quiet || silent {
			return nil
		}

		// ipfs cli progress bar defaults to true unless quiet or silent is used
		_, found := req.Options[progressOptionName].(bool)
		if !found {
			req.Options[progressOptionName] = true
		}

		return nil
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		n, err := GetNode(env)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		if n.Routing == nil {
			res.SetError(errNotOnline, cmdkit.ErrNormal)
			return
		}

		cfg, err := n.Repo.Config()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}
		// check if repo will exceed storage limit if added
		// TODO: this doesn't handle the case if the hashed file is already in blocks (deduplicated)
		// TODO: conditional GC is disabled due to it is somehow not possible to pass the size to the daemon
		//if err := corerepo.ConditionalGC(req.Context(), n, uint64(size)); err != nil {
		//	res.SetError(err, cmdkit.ErrNormal)
		//	return
		//}

		progress, _ := req.Options[progressOptionName].(bool)
		trickle, _ := req.Options[trickleOptionName].(bool)
		wrap, _ := req.Options[wrapOptionName].(bool)
		hash, _ := req.Options[onlyHashOptionName].(bool)
		hidden, _ := req.Options[hiddenOptionName].(bool)
		silent, _ := req.Options[silentOptionName].(bool)
		chunker, _ := req.Options[chunkerOptionName].(string)
		dopin, _ := req.Options[pinOptionName].(bool)
		rawblks, rbset := req.Options[rawLeavesOptionName].(bool)
		nocopy, _ := req.Options[noCopyOptionName].(bool)
		fscache, _ := req.Options[fstoreCacheOptionName].(bool)
		cidVer, cidVerSet := req.Options[cidVersionOptionName].(int)
		hashFunStr, _ := req.Options[hashOptionName].(string)

		// The arguments are subject to the following constraints.
		//
		// nocopy -> filestoreEnabled
		// nocopy -> rawblocks
		// (hash != sha2-256) -> cidv1

		// NOTE: 'rawblocks -> cidv1' is missing. Legacy reasons.

		// nocopy -> filestoreEnabled
		if nocopy && !cfg.Experimental.FilestoreEnabled {
			res.SetError(filestore.ErrFilestoreNotEnabled, cmdkit.ErrClient)
			return
		}

		// nocopy -> rawblocks
		if nocopy && !rawblks {
			// fixed?
			if rbset {
				res.SetError(
					fmt.Errorf("nocopy option requires '--raw-leaves' to be enabled as well"),
					cmdkit.ErrNormal,
				)
				return
			}
			// No, satisfy mandatory constraint.
			rawblks = true
		}

		// (hash != "sha2-256") -> CIDv1
		if hashFunStr != "sha2-256" && cidVer == 0 {
			if cidVerSet {
				res.SetError(
					errors.New("CIDv0 only supports sha2-256"),
					cmdkit.ErrClient,
				)
				return
			}
			cidVer = 1
		}

		// cidV1 -> raw blocks (by default)
		if cidVer > 0 && !rbset {
			rawblks = true
		}

		prefix, err := dag.PrefixForCidVersion(cidVer)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		hashFunCode, ok := mh.Names[strings.ToLower(hashFunStr)]
		if !ok {
			res.SetError(fmt.Errorf("unrecognized hash function: %s", strings.ToLower(hashFunStr)), cmdkit.ErrNormal)
			return
		}

		prefix.MhType = hashFunCode
		prefix.MhLength = -1

		if hash {
			nilnode, err := core.NewNode(n.Context(), &core.BuildCfg{
				//TODO: need this to be true or all files
				// hashed will be stored in memory!
				NilRepo: true,
			})
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}
			n = nilnode
		}

		addblockstore := n.Blockstore
		if !(fscache || nocopy) {
			addblockstore = bstore.NewGCBlockstore(n.BaseBlocks, n.GCLocker)
		}

		exch := n.Exchange
		local, _ := req.Options["local"].(bool)
		if local {
			exch = offline.Exchange(addblockstore)
		}

		bserv := blockservice.New(addblockstore, exch) // hash security 001
		dserv := dag.NewDAGService(bserv)

		outChan := make(chan interface{}, adderOutChanSize)

		fileAdder, err := coreunix.NewAdder(req.Context, n.Pinning, n.Blockstore, dserv)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		fileAdder.Out = outChan
		fileAdder.Chunker = chunker
		fileAdder.Progress = progress
		fileAdder.Hidden = hidden
		fileAdder.Trickle = trickle
		fileAdder.Wrap = wrap
		fileAdder.Pin = dopin
		fileAdder.Silent = silent
		fileAdder.RawLeaves = rawblks
		fileAdder.NoCopy = nocopy
		fileAdder.Prefix = &prefix

		needUnpin := !fileAdder.Pin
		fileAdder.Pin = true

		if hash {
			md := dagtest.Mock()
			emptyDirNode := ft.EmptyDirNode()
			// Use the same prefix for the "empty" MFS root as for the file adder.
			emptyDirNode.Prefix = *fileAdder.Prefix
			mr, err := mfs.NewRoot(req.Context, md, emptyDirNode, nil)
			if err != nil {
				res.SetError(err, cmdkit.ErrNormal)
				return
			}

			fileAdder.SetMfsRoot(mr)
		}

		addAllAndPin := func(f files.File) error {
			// Iterate over each top-level file and add individually. Otherwise the
			// single files.File f is treated as a directory, affecting hidden file
			// semantics.
			for {
				file, err := f.NextFile()
				if err == io.EOF {
					// Finished the list of files.
					break
				} else if err != nil {
					return err
				}
				if err := fileAdder.AddFile(file); err != nil {
					return err
				}
			}

			// copy intermediary nodes from editor to our actual dagservice
			_, err := fileAdder.Finalize()
			if err != nil {
				return err
			}

			if hash {
				return nil
			}

			return fileAdder.PinRoot()
		}

		errCh := make(chan error)
		go func() {
			var err error
			defer func() { errCh <- err }()
			defer close(outChan)
			err = addAllAndPin(req.Files)
			if err != nil {
				return
			}

			log.Debug(("add success, ready to push"))

			// got root hash
			root, err := fileAdder.RootNode()
			if err != nil {
				log.Warning("cant got the root node")
				return
			}
			c := root.Cid()

			if needUnpin {
				PushRecorder.Write(c.String(), "1")
				defer func() {
					_, e := corerepo.Unpin(n, req.Context, []string{c.String()}, true)
					if e != nil {
						if err != nil {
							err = errors.New(err.Error() + "(unpin failed" + e.Error() + ")")
						} else {
							err = errors.Wrap(e, "unpin failed")
						}
					} else {
						PushRecorder.Write(c.String(), "-1")
					}
				}()
			}

			// do backup
			backupOutput, err := backupFunc(n, c)

			if err != nil {
				err = errors.Wrap(err, "backup failed:")
				return
			}

			outChan <- coreunix.AddedObject{
				Extend: backupOutput,
			}
		}()

		defer res.Close()

		err = res.Emit(outChan)
		if err != nil {
			log.Error(err)
			return
		}
		err = <-errCh
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
		}
	},
	PostRun: cmds.PostRunMap{
		cmds.CLI: func(req *cmds.Request, re cmds.ResponseEmitter) cmds.ResponseEmitter {
			reNext, res := cmds.NewChanResponsePair(req)
			outChan := make(chan interface{})

			sizeChan := make(chan int64, 1)

			sizeFile, ok := req.Files.(files.SizeFile)
			if ok {
				// Could be slow.
				go func() {
					size, err := sizeFile.Size()
					if err != nil {
						log.Warningf("error getting files size: %s", err)
						// see comment above
						return
					}

					sizeChan <- size
				}()
			} else {
				// we don't need to error, the progress bar just
				// won't know how big the files are
				log.Warning("cannot determine size of input file")
			}

			progressBar := func(wait chan struct{}) {
				defer close(wait)

				quiet, _ := req.Options[quietOptionName].(bool)
				quieter, _ := req.Options[quieterOptionName].(bool)
				quiet = quiet || quieter

				progress, _ := req.Options[progressOptionName].(bool)

				var bar *pb.ProgressBar
				if progress {
					bar = pb.New64(0).SetUnits(pb.U_BYTES)
					bar.ManualUpdate = true
					bar.ShowTimeLeft = false
					bar.ShowPercent = false
					bar.Output = os.Stderr
					bar.Start()
				}

				lastFile := ""
				lastHash := ""
				var totalProgress, prevFiles, lastBytes int64

			LOOP:
				for {
					select {
					case out, ok := <-outChan:
						if !ok {
							if quieter {
								fmt.Fprintln(os.Stdout, lastHash)
							}

							break LOOP
						}
						output := out.(*coreunix.AddedObject)
						if len(output.Hash) > 0 {
							lastHash = output.Hash
							if quieter {
								continue
							}

							if progress {
								// clear progress bar line before we print "added x" output
								fmt.Fprintf(os.Stderr, "\033[2K\r")
							}
							if quiet {
								fmt.Fprintf(os.Stdout, "%s\n", output.Hash)
							} else {
								fmt.Fprintf(os.Stdout, "added %s %s\n", output.Hash, output.Name)
							}

						} else if output.Extend != nil {
							for _, s := range output.Extend.Success {
								fmt.Fprintf(os.Stdout, "backup success to %s\n", s.ID)
							}
							for _, f := range output.Extend.Failed {
								fmt.Fprintf(os.Stdout, "backup failed to %s : %s\n", f.ID, f.Msg)
							}
							continue
						} else {
							if !progress {
								continue
							}

							if len(lastFile) == 0 {
								lastFile = output.Name
							}
							if output.Name != lastFile || output.Bytes < lastBytes {
								prevFiles += lastBytes
								lastFile = output.Name
							}
							lastBytes = output.Bytes
							delta := prevFiles + lastBytes - totalProgress
							totalProgress = bar.Add64(delta)
						}

						if progress {
							bar.Update()
						}
					case size := <-sizeChan:
						if progress {
							bar.Total = size
							bar.ShowPercent = true
							bar.ShowBar = true
							bar.ShowTimeLeft = true
						}
					case <-req.Context.Done():
						// don't set or print error here, that happens in the goroutine below
						return
					}
				}
			}

			go func() {
				// defer order important! First close outChan, then wait for output to finish, then close re
				defer re.Close()

				if e := res.Error(); e != nil {
					defer close(outChan)
					re.SetError(e.Message, e.Code)
					return
				}

				wait := make(chan struct{})
				go progressBar(wait)

				defer func() { <-wait }()
				defer close(outChan)

				for {
					v, err := res.Next()
					if !cmds.HandleError(err, res, re) {
						break
					}

					select {
					case outChan <- v:
					case <-req.Context.Done():
						re.SetError(req.Context.Err(), cmdkit.ErrNormal)
						return
					}
				}
			}()

			return reNext
		},
	},
	Type: coreunix.AddedObject{},
}

type pushRecord struct {
	filename string
	f        *os.File
	w        *csv.Writer
	once     sync.Once
}

const pushRecordFileName = "push.record"

var PushRecorder = pushRecord{}

func (t *pushRecord) Init(repoDir string) {
	t.filename = filepath.Join(repoDir, pushRecordFileName)
}

func (t *pushRecord) Write(k, v string) {
	t.once.Do(func() {
		if t.filename == "" {
			log.Warning("please call init file of pushRecord before write")
			return
		}

		var err error
		t.f, err = os.OpenFile(t.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Warningf("open file %s failed: %v\n", t.filename, err)
			return
		}

		err = t.lock(t.f.Fd())
		if err != nil {
			log.Warning("locak file %s failed: %v\n", t.filename, err)
			t.f.Close()
			return
		}

		t.w = csv.NewWriter(t.f)
	})

	if t.w == nil {
		log.Warning("record file not opened")
		return
	}

	err := t.w.Write([]string{k, v})
	if err != nil {
		log.Warningf("record for %s=%s failed: %v\n ", k, v, err)
		return
	}
	t.w.Flush()

	log.Debug("write record:", k, v)
}

//加锁
func (t *pushRecord) lock(fd uintptr) error {
	return syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB)
	// return nil
}

//释放锁
func (t *pushRecord) unlock(fd uintptr) error {
	return syscall.Flock(int(fd), syscall.LOCK_UN)
	// return nil
}

func (t *pushRecord) Clear(ctx context.Context) {
	// file already opened
	if t.f != nil {
		return
	}

	f, err := os.Open(t.filename)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Warningf("read push record file %s failed: %v\n", t.filename, err)
		}
		return
	}
	err = t.lock(f.Fd())
	if err == syscall.EWOULDBLOCK {
		f.Close()
		return
	}

	rmFile := true
	defer func() {
		t.unlock(f.Fd())
		f.Close()
		if rmFile {
			os.Remove(t.filename)
		}
	}()

	// parse the bs
	r := csv.NewReader(f)
	r.FieldsPerRecord = 2

	hashes := make(map[string]int, 0)
	rcd, err := r.Read()
	for err == nil && rcd != nil {
		n, err := strconv.Atoi(rcd[1])
		if err != nil {
			log.Warningf("push record file %s record %s convert failed: %v\n", t.filename, rcd, err)
			rmFile = false
			return
		}

		if _, found := hashes[rcd[0]]; found {
			hashes[rcd[0]] += n
		} else {
			hashes[rcd[0]] = n
		}

		rcd, err = r.Read()
	}

	if err != nil && err != io.EOF {
		log.Warningf("read line from push record file %s failed: %v\n", t.filename, err)
		rmFile = false
		return
	}

	// handle
	pathes := make([]string, 0, len(hashes))
	for k, v := range hashes {
		if v == 0 {
			continue
		}

		if v < 0 {
			log.Warningf("push record file %s record %s value = %d\n", t.filename, k, v)
			rmFile = false
			continue
		}

		pathes = append(pathes, k)
	}

	log.Debug("hashes need to unpined: ", pathes)
	if len(pathes) > 0 {
		args := []string{"pin", "rm"}
		args = append(args, pathes...)
		bs, err := exec.CommandContext(ctx, "ipfs", args...).CombinedOutput()
		if err != nil && !strings.Contains(err.Error(), "exit status 1") {
			log.Warning("do unpin failed:", err, string(bs))
			rmFile = false
			return
		}
	}
}
