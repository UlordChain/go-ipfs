#!/usr/bin/env bash

test_description="Test mount command in conjunction with publishing"

# imports
. lib/test-lib.sh

# if in travis CI, dont test mount (no fuse)
if ! test_have_prereq FUSE; then
  skip_all='skipping mount tests, fuse not available'

  test_done
fi

test_init_udfs

# start iptb + wait for peering
NUM_NODES=3
test_expect_success 'init iptb' '
  iptb init -n $NUM_NODES -f --bootstrap=none --port=0 &&
  startup_cluster $NUM_NODES
'

# pre-mount publish
HASH=$(echo 'hello warld' | udfsi 0 add -q)
test_expect_success "can publish before mounting /ipns" '
  udfsi 0 name publish '$HASH'
'

# mount
UDFS_MOUNT_DIR="$PWD/udfs"
IPNS_MOUNT_DIR="$PWD/ipns"
test_expect_success FUSE "'udfs mount' succeeds" '
  udfsi 0 mount -f "'"$UDFS_MOUNT_DIR"'" -n "'"$IPNS_MOUNT_DIR"'" >actual
'
test_expect_success FUSE "'udfs mount' output looks good" '
  echo "UDFS mounted at: $PWD/udfs" >expected &&
  echo "IPNS mounted at: $PWD/ipns" >>expected &&
  test_cmp expected actual
'

test_expect_success "cannot publish after mounting /ipns" '
  echo "Error: cannot manually publish while IPNS is mounted" >expected &&
  test_must_fail udfsi 0 name publish '$HASH' 2>actual &&
  test_cmp expected actual
'

test_expect_success "unmount /ipns out-of-band" '
  fusermount -u "'"$IPNS_MOUNT_DIR"'"
'

test_expect_success "can publish after unmounting /ipns" '
  udfsi 0 name publish '$HASH'
'

# clean-up udfs
test_expect_success "unmount /udfs" '
  fusermount -u "'"$UDFS_MOUNT_DIR"'"
'
iptb stop

test_done
