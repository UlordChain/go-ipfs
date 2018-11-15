#!/usr/bin/env bash
#
# Copyright (c) 2016 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test active request commands"

. lib/test-lib.sh

test_init_udfs
test_launch_udfs_daemon

test_expect_success "command works" '
  udfs diag cmds > cmd_out
'

test_expect_success "invoc shows up in output" '
  grep "diag/cmds" cmd_out > /dev/null
'

test_expect_success "start longer running command" '
  udfs log tail &
  LOGPID=$!
  go-sleep 100ms
'

test_expect_success "long running command shows up" '
  udfs diag cmds > cmd_out2
'

test_expect_success "output looks good" '
  grep "log/tail" cmd_out2 | grep "true" > /dev/null
'

test_expect_success "kill log cmd" '
  kill $LOGPID
  go-sleep 0.5s
  kill $LOGPID

  wait $LOGPID || true
'

test_expect_success "long running command inactive" '
  udfs diag cmds > cmd_out3
'

test_expect_success "command shows up as inactive" '
  grep "log/tail" cmd_out3 | grep "false"
'

test_kill_udfs_daemon
test_done
