// for *inux system

// +build !windows,amd64

package commands

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"
	"github.com/ipfs/go-ipfs/core/coreapi/interface/options"
	"github.com/ipfs/go-ipfs/core/corerepo"
	"gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"github.com/pkg/errors"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmPtj12fdwuAqj9sBSTNUxBNu8kCGNp8b3o8yUzMm5GHpq/pb"
	"gx/ipfs/QmSXUokcP4TJpFfqozT69AVAYRtzXVMUjzQVkYX41R9Svs/go-ipfs-cmds"
	"gx/ipfs/QmZMWMvWMVKCbHetJ4RgndbuEF1io2UpUxwQwtNjtYPzSC/go-ipfs-files"
	"gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
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
		cmdkit.StringOption(stdinPathName, "Assign a name if the file source is stdin."),
		cmdkit.BoolOption(hiddenOptionName, "H", "Include files that are hidden. Only takes effect on recursive add."),
		cmdkit.StringOption(chunkerOptionName, "s", "Chunking algorithm, size-[bytes] or rabin-[min]-[avg]-[max]").WithDefault("size-262144"),
		cmdkit.BoolOption(pinOptionName, "Pin this object when pushing.").WithDefault(false),
		cmdkit.BoolOption(rawLeavesOptionName, "Use raw blocks for leaf nodes. (experimental)"),
		cmdkit.BoolOption(noCopyOptionName, "Add the file using filestore. Implies raw-leaves. (experimental)"),
		cmdkit.BoolOption(fstoreCacheOptionName, "Check the filestore for pre-existing blocks. (experimental)"),
		cmdkit.IntOption(cidVersionOptionName, "CID version. Defaults to 0 unless an option that depends on CIDv1 is passed. (experimental)"),
		cmdkit.StringOption(hashOptionName, "Hash function to use. Implies CIDv1 if not sha2-256. (experimental)").WithDefault("sha2-256"),
		cmdkit.BoolOption(inlineOptionName, "Inline small blocks into CIDs. (experimental)"),
		cmdkit.IntOption(inlineLimitOptionName, "Maximum block size to inline. (experimental)").WithDefault(32),
		cmdkit.StringOption(accountOptionName, "Account of user to check"),
		cmdkit.StringOption(checkOptionName, "The hash value for check"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		cfg, _ := node.Repo.Config()
		var(
			account string
			check string
			size uint64
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

			size, err = ValidOnUOS(&cfg.UOSCheck, account, check)
			if err != nil {
				return errors.Wrap(err, "valid failed")
			}
		}

		// Must be online!
		if !node.OnlineMode() {
			return cmdkit.Errorf(cmdkit.ErrClient, ErrNotOnline.Error())
		}

		api, err := cmdenv.GetApi(env)
		if err != nil {
			return err
		}

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
		inline, _ := req.Options[inlineOptionName].(bool)
		inlineLimit, _ := req.Options[inlineLimitOptionName].(int)
		pathName, _ := req.Options[stdinPathName].(string)
		local, _ := req.Options["local"].(bool)

		hashFunCode, ok := mh.Names[strings.ToLower(hashFunStr)]
		if !ok {
			return fmt.Errorf("unrecognized hash function: %s", strings.ToLower(hashFunStr))
		}

		events := make(chan interface{}, adderOutChanSize)

		opts := []options.UnixfsAddOption{
			options.Unixfs.Hash(hashFunCode),

			options.Unixfs.Inline(inline),
			options.Unixfs.InlineLimit(inlineLimit),

			options.Unixfs.Chunker(chunker),

			options.Unixfs.Pin(true),
			options.Unixfs.HashOnly(hash),
			options.Unixfs.Local(local),
			options.Unixfs.FsCache(fscache),
			options.Unixfs.Nocopy(nocopy),

			options.Unixfs.Wrap(wrap),
			options.Unixfs.Hidden(hidden),
			options.Unixfs.StdinName(pathName),

			options.Unixfs.Progress(progress),
			options.Unixfs.Silent(silent),
			//options.Unixfs.Events(events),
		}

		if cidVerSet {
			opts = append(opts, options.Unixfs.CidVersion(cidVer))
		}

		if rbset {
			opts = append(opts, options.Unixfs.RawLeaves(rawblks))
		}

		if trickle {
			opts = append(opts, options.Unixfs.Layout(options.TrickleLayout))
		}

		needUnpin := !dopin
		errCh := make(chan error)
		go func() {
			var err error
			defer func() { errCh <- err }()
			defer close(events)


			wg := sync.WaitGroup{}
			tempEvents := make(chan interface{}, adderOutChanSize)
			opts = append(opts, options.Unixfs.Events(tempEvents))

			var backupErr error
			wg.Add(1)
			go func(){
				defer wg.Done()

				for e := range tempEvents{
					if addevent, ok := e.(*coreiface.AddEvent); !ok {
						events <- e
					}else if backupErr == nil {

						wg.Add(1)
						go func(ar *coreiface.AddEvent){
							defer wg.Done()

							if needUnpin {
								PushRecorder.Write(ar.Hash, "1")
								defer func() {
									_, e := corerepo.Unpin(node, api, req.Context, []string{ar.Hash}, true)
									if e != nil {
										if err != nil {
											log.Warning(err.Error() + "(unpin failed" + e.Error() + ")")
										} else {
											log.Warning("unpin failed: "+ e.Error())
										}
									} else {
										PushRecorder.Write(ar.Hash, "-1")
									}
								}()
							}

							c, err := cid.Parse(ar.Hash)
							if err != nil {
								backupErr = errors.Wrapf(err, "parse hash %s to cid failed", ar.Hash)
								return
							}

							if cfg.UOSCheck.Enable {
								// check size

								s, _ := strconv.ParseUint(ar.Size, 10, 64)
								if size*1024 < s {
									// remove the content
									err = corerepo.Remove(node, req.Context, []cid.Cid{c}, true, false)
									if err != nil {
										err = errors.Wrap(err, "unpin the content failed")
										return
									}

									err = errors.New("the content size not matched on uos")
									return
								}

								// check hash
								if c.String() != check {

									// remove the content
									err = corerepo.Remove(node, req.Context, []cid.Cid{c}, true, false)
									if err != nil {
										err = errors.Wrap(err, "unpin the content failed")
										return
									}

									err = errors.New("the content hash not matched on uos")
									return
								}
							}

							log.Debug(("add success, ready to push"))

							// do backup
							backupOutput, err := backupFunc(node, c)

							if err != nil {
								backupErr = errors.Wrap(err, "backup failed")
								return
							}

							ar.Extend = backupOutput
							events <- ar

						}(addevent)
					}
				}
			}()

			_, err = api.Unixfs().Add(req.Context, req.Files, opts...)
			if err != nil {
				return
			}
			close(tempEvents)

			wg.Wait()
			err = backupErr
		}()

		err = res.Emit(events)
		if err != nil {
			return err
		}

		return <-errCh
	},
	PostRun: cmds.PostRunMap{
		cmds.CLI: func(res cmds.Response, re cmds.ResponseEmitter) error {
			sizeChan := make(chan int64, 1)
			outChan := make(chan interface{})
			req := res.Request()

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
						output := out.(*coreiface.AddEvent)
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

							if output.Extend != nil {
								for _, s := range output.Extend.Success {
									fmt.Fprintf(os.Stdout, "backup success to %s\n", s.ID)
								}
								for _, f := range output.Extend.Failed {
									fmt.Fprintf(os.Stdout, "backup failed to %s : %s\n", f.ID, f.Msg)
								}
							}
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

			if e := res.Error(); e != nil {
				close(outChan)
				return e
			}

			wait := make(chan struct{})
			go progressBar(wait)

			defer func() { <-wait }()
			defer close(outChan)

			for {
				v, err := res.Next()
				if err != nil {
					if err == io.EOF {
						return nil
					}

					return err
				}

				select {
				case outChan <- v:
				case <-req.Context.Done():
					return req.Context.Err()
				}
			}
		},
	},
	Type: coreiface.AddEvent{},
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
}

//释放锁
func (t *pushRecord) unlock(fd uintptr) error {
	return syscall.Flock(int(fd), syscall.LOCK_UN)
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
