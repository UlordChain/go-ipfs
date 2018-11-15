#!/usr/bin/env bash

test_description="Test datastore config"

. lib/test-lib.sh

test_init_udfs

test_launch_udfs_daemon
test_kill_udfs_daemon

SPEC_NOSYNC=$(cat ../t0024-files/spec-nosync)

SPEC_NEWSHARDFUN=$(cat ../t0024-files/spec-newshardfun)

test_expect_success "change runtime value in spec config" '
  udfs config --json Datastore.Spec "$SPEC_NOSYNC"
'

test_launch_udfs_daemon
test_kill_udfs_daemon

test_expect_success "change on-disk value in spec config" '
  udfs config --json Datastore.Spec "$SPEC_NEWSHARDFUN"
'

test_expect_success "can not launch daemon after on-disk value change" '
  test_must_fail udfs daemon
'

test_done
