#!/usr/bin/env bash

test_description="Test udfs cli cmd suggest"

. lib/test-lib.sh

test_suggest() {


  test_expect_success "test command fails" '
    test_must_fail udfs kog 2>actual
  '

  test_expect_success "test one command is suggested" '
    grep "Did you mean this?" actual &&
    grep "log" actual ||
    test_fsh cat actual
  '

  test_expect_success "test command fails" '
    test_must_fail udfs lis 2>actual
  '

  test_expect_success "test multiple commands are suggested" '
    grep "Did you mean any of these?" actual &&
    grep "ls" actual &&
    grep "id" actual ||
    test_fsh cat actual
  '

}

test_init_udfs

test_suggest

test_launch_udfs_daemon

test_suggest

test_kill_udfs_daemon

test_done
