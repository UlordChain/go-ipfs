// +build linux darwin freebsd netbsd openbsd
// +build !nofuse

package readonly

import (
	core "github.com/udfs/go-udfs/core"
	mount "github.com/udfs/go-udfs/fuse/mount"
)

// Mount mounts UDFS at a given location, and returns a mount.Mount instance.
func Mount(udfs *core.UdfsNode, mountpoint string) (mount.Mount, error) {
	cfg, err := udfs.Repo.Config()
	if err != nil {
		return nil, err
	}
	allow_other := cfg.Mounts.FuseAllowOther
	fsys := NewFileSystem(udfs)
	return mount.NewMount(udfs.Process(), fsys, mountpoint, allow_other)
}
