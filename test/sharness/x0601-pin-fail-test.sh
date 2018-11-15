#!/usr/bin/env bash
#
# Copyright (c) 2016 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test very large number of pins"

. lib/test-lib.sh

test_init_udfs

test_launch_udfs_daemon

test_expect_success "pre-test setup" '
  printf "" > pins &&
  udfs pin ls --type=recursive -q > rec_pins_before
'


for i in `seq 9000`
do
  test_expect_success "udfs add (and pin) a file" '
    echo $i | udfs add -q >> pins
  '
done

test_expect_success "get pinset afterwards" '
  udfs pin ls --type=recursive -q | sort > rec_pins_after &&
  cat pins rec_pins_before | sort | uniq > exp_pins_after &&
  test_cmp rec_pins_after exp_pins_after
'

test_kill_udfs_daemon

test_done

