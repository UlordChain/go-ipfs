#!/usr/bin/env bash
#
# Copyright (c) 2017 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test shutdown command"

. lib/test-lib.sh

test_init_udfs

test_launch_udfs_daemon

test_expect_success "shutdown succeeds" '
  udfs shutdown
'

test_expect_success "daemon no longer running" '
  for i in $(test_seq 1 100)
  do
    go-sleep 100ms
    ! kill -0 $UDFS_PID 2>/dev/null && return
  done
'

test_launch_udfs_daemon --offline

test_expect_success "shutdown succeeds" '
  udfs shutdown
'

test_expect_success "daemon no longer running" '
  for i in $(test_seq 1 100)
  do
    go-sleep 100ms
    ! kill -0 $UDFS_PID 2>/dev/null && return
  done
'
test_done
