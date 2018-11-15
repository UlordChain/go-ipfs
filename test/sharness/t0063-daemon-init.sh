#!/usr/bin/env bash
#
# Copyright (c) 2014 Juan Batiz-Benet
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test daemon --init command"

. lib/test-lib.sh

# We don't want the normal test_init_udfs but we need to make sure the
# UDFS_PATH is set correctly.
export UDFS_PATH="$(pwd)/.udfs"

# safety check since we will be removing the directory
if [ -e "$UDFS_PATH" ]; then
  echo "$UDFS_PATH exists"
  exit 1
fi

test_udfs_daemon_init() {
  # Doing it manually since we want to launch the daemon with an
  # empty or non-existent repo; the normal
  # test_launch_udfs_daemon does not work since it assumes the
  # repo was created a particular way with regard to the API
  # server.

  test_expect_success "'udfs daemon --init' succeeds" '
    udfs daemon --init --init-profile=test >actual_daemon 2>daemon_err &
    UDFS_PID=$!
    sleep 2 &&
    if ! kill -0 $UDFS_PID; then cat daemon_err; return 1; fi
  '

  test_expect_success "'udfs daemon' can be killed" '
    test_kill_repeat_10_sec $UDFS_PID
  '
}

test_expect_success "remove \$UDFS_PATH dir" '
  rm -rf "$UDFS_PATH"
'
test_udfs_daemon_init

test_expect_success "create empty \$UDFS_PATH dir" '
  rm -rf "$UDFS_PATH" &&
  mkdir "$UDFS_PATH"
'

test_udfs_daemon_init

test_done
