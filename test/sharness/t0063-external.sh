#!/usr/bin/env bash
#
# Copyright (c) 2015 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="test external command functionality"

. lib/test-lib.sh


# set here so daemon launches with it
PATH=`pwd`/bin:$PATH

test_init_udfs

test_expect_success "create fake udfs-update bin" '
  mkdir bin &&
  echo "#!/bin/sh" > bin/udfs-update &&
  echo "pwd" >> bin/udfs-update &&
  echo "test -e \"$UDFS_PATH/repo.lock\" || echo \"repo not locked\" " >> bin/udfs-update &&
  chmod +x bin/udfs-update &&
  mkdir just_for_test
'

test_expect_success "external command runs from current user directory and doesn't lock repo" '
  (cd just_for_test && udfs update) > actual
'

test_expect_success "output looks good" '
  echo `pwd`/just_for_test > exp &&
  echo "repo not locked" >> exp &&
  test_cmp exp actual
'

test_launch_udfs_daemon

test_expect_success "external command runs from current user directory when daemon is running" '
  (cd just_for_test && udfs update) > actual
'

test_expect_success "output looks good" '
  echo `pwd`/just_for_test > exp &&
  test_cmp exp actual
'

test_kill_udfs_daemon

test_done
