// +build !windows,nofuse

package commands

import (
	cmds "github.com/udfs/go-udfs/commands"

	"gx/udfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-udfs-cmdkit"
)

var MountCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Mounts udfs to the filesystem (disabled).",
		ShortDescription: `
This version of udfs is compiled without fuse support, which is required
for mounting. If you'd like to be able to mount, please use a version of
udfs compiled with fuse.

For the latest instructions, please check the project's repository:
  http://github.com/udfs/go-udfs
`,
	},
}
