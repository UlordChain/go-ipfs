#!/usr/bin/env bash
#
# Copyright (c) 2016 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="test the unix files api flushing"

. lib/test-lib.sh

test_init_udfs

verify_path_exists() {
  # simply running ls on a file should be a good 'check'
  udfs files ls $1
}

verify_dir_contents() {
  dir=$1
  shift
  rm -f expected
  touch expected
  for e in $@
  do
    echo $e >> expected
  done

  test_expect_success "can list dir" '
    udfs files ls $dir > output
  '

  test_expect_success "dir entries look good" '
    test_sort_cmp output expected
  '
}

test_launch_udfs_daemon

test_expect_success "can copy a file in" '
  HASH=$(echo "foo" | udfs add -q) &&
  udfs files cp /udfs/$HASH /file
'

test_kill_udfs_daemon
test_launch_udfs_daemon

test_expect_success "file is still there" '
  verify_path_exists /file
'

test_kill_udfs_daemon

test_done
