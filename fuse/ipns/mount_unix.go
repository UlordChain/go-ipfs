// +build linux darwin freebsd netbsd openbsd
// +build !nofuse

package ipns

import (
	core "github.com/udfs/go-udfs/core"
	mount "github.com/udfs/go-udfs/fuse/mount"
)

// Mount mounts ipns at a given location, and returns a mount.Mount instance.
func Mount(udfs *core.UdfsNode, ipnsmp, udfsmp string) (mount.Mount, error) {
	cfg, err := udfs.Repo.Config()
	if err != nil {
		return nil, err
	}

	allow_other := cfg.Mounts.FuseAllowOther

	fsys, err := NewFileSystem(udfs, udfs.PrivateKey, udfsmp, ipnsmp)
	if err != nil {
		return nil, err
	}

	return mount.NewMount(udfs.Process(), fsys, ipnsmp, allow_other)
}
