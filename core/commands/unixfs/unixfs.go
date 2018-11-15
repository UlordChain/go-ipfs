package unixfs

import (
	cmds "github.com/udfs/go-udfs/commands"
	e "github.com/udfs/go-udfs/core/commands/e"

	"gx/udfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-udfs-cmdkit"
)

var UnixFSCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Interact with UDFS objects representing Unix filesystems.",
		ShortDescription: `
'udfs file' provides a familiar interface to file systems represented
by UDFS objects, which hides udfs implementation details like layout
objects (e.g. fanout and chunking).
`,
		LongDescription: `
'udfs file' provides a familiar interface to file systems represented
by UDFS objects, which hides udfs implementation details like layout
objects (e.g. fanout and chunking).
`,
	},

	Subcommands: map[string]*cmds.Command{
		"ls": LsCmd,
	},
}

// copy+pasted from ../commands.go
func unwrapOutput(i interface{}) (interface{}, error) {
	var (
		ch <-chan interface{}
		ok bool
	)

	if ch, ok = i.(<-chan interface{}); !ok {
		return nil, e.TypeErr(ch, i)
	}

	return <-ch, nil
}
