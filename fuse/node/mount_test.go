// +build !nofuse

package node

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"context"

	core "github.com/udfs/go-udfs/core"
	ipns "github.com/udfs/go-udfs/fuse/ipns"
	mount "github.com/udfs/go-udfs/fuse/mount"
	namesys "github.com/udfs/go-udfs/namesys"

	offroute "gx/udfs/QmbFRJeEmEU16y3BmKKaD4a9fm5oHsEAMHe2vSB1UnfLMi/go-udfs-routing/offline"
	ci "gx/udfs/QmcW4FGAt24fdK1jBgWQn3yP4R9ZLyWQqjozv9QK7epRhL/go-testutil/ci"
)

func maybeSkipFuseTests(t *testing.T) {
	if ci.NoFuse() {
		t.Skip("Skipping FUSE tests")
	}
}

func mkdir(t *testing.T, path string) {
	err := os.Mkdir(path, os.ModeDir|os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

// Test externally unmounting, then trying to unmount in code
func TestExternalUnmount(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	// TODO: needed?
	maybeSkipFuseTests(t)

	node, err := core.NewNode(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	err = node.LoadPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	node.Routing = offroute.NewOfflineRouter(node.Repo.Datastore(), node.RecordValidator)
	node.Namesys = namesys.NewNameSystem(node.Routing, node.Repo.Datastore(), 0)

	err = ipns.InitializeKeyspace(node, node.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	// get the test dir paths (/tmp/fusetestXXXX)
	dir, err := ioutil.TempDir("", "fusetest")
	if err != nil {
		t.Fatal(err)
	}

	udfsDir := dir + "/udfs"
	ipnsDir := dir + "/ipns"
	mkdir(t, udfsDir)
	mkdir(t, ipnsDir)

	err = Mount(node, udfsDir, ipnsDir)
	if err != nil {
		t.Fatal(err)
	}

	// Run shell command to externally unmount the directory
	cmd := "fusermount"
	args := []string{"-u", ipnsDir}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		t.Fatal(err)
	}

	// TODO(noffle): it takes a moment for the goroutine that's running fs.Serve to be notified and do its cleanup.
	time.Sleep(time.Millisecond * 100)

	// Attempt to unmount IPNS; check that it was already unmounted.
	err = node.Mounts.Ipns.Unmount()
	if err != mount.ErrNotMounted {
		t.Fatal("Unmount should have failed")
	}

	// Attempt to unmount UDFS; it should unmount successfully.
	err = node.Mounts.Udfs.Unmount()
	if err != nil {
		t.Fatal(err)
	}
}
