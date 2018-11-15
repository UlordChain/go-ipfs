#!/usr/bin/env bash
#
# Copyright (c) 2016 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test object diff command"

. lib/test-lib.sh

test_init_udfs

test_expect_success "create some objects for testing diffs" '
  mkdir foo &&
  echo "stuff" > foo/bar &&
  mkdir foo/baz &&
  A=$(udfs add -r -q foo | tail -n1) &&
  echo "more things" > foo/cat &&
  B=$(udfs add -r -q foo | tail -n1) &&
  echo "nested" > foo/baz/dog &&
  C=$(udfs add -r -q foo | tail -n1)
  echo "changed" > foo/bar &&
  D=$(udfs add -r -q foo | tail -n1) &&
  echo "" > single_file &&
  SINGLE_FILE=$(udfs add -r -q single_file | tail -n1) &&
  mkdir empty_dir
  EMPTY_DIR=$(udfs add -r -q empty_dir | tail -n1)
'

test_expect_success "diff against self is empty" '
  udfs object diff $A $A > diff_out
'

test_expect_success "identity diff output looks good" '
  printf "" > diff_exp &&
  test_cmp diff_exp diff_out
'

test_expect_success "diff against self (single file) is empty" '
  udfs object diff $SINGLE_FILE $SINGLE_FILE > diff_out
  printf "" > diff_exp &&
  test_cmp diff_exp diff_out
'

test_expect_success "diff against self (empty dir) is empty" '
  udfs object diff $EMPTY_DIR $EMPTY_DIR > diff_out
  printf "" > diff_exp &&
  test_cmp diff_exp diff_out
'

test_expect_success "diff added link works" '
  udfs object diff $A $B > diff_out
'

test_expect_success "diff added link looks right" '
  echo + QmUSvcqzhdfYM1KLDbM76eLPdS9ANFtkJvFuPYeZt73d7A \"cat\" > diff_exp &&
  test_cmp diff_exp diff_out
'

test_expect_success "verbose diff added link works" '
  udfs object diff -v $A $B > diff_out
'

test_expect_success "verbose diff added link looks right" '
  echo Added new link \"cat\" pointing to QmUSvcqzhdfYM1KLDbM76eLPdS9ANFtkJvFuPYeZt73d7A. > diff_exp &&
  test_cmp diff_exp diff_out
'

test_expect_success "diff removed link works" '
  udfs object diff -v $B $A > diff_out
'

test_expect_success "diff removed link looks right" '
  echo Removed link \"cat\" \(was QmUSvcqzhdfYM1KLDbM76eLPdS9ANFtkJvFuPYeZt73d7A\). > diff_exp &&
  test_cmp diff_exp diff_out
'

test_expect_success "diff nested add works" '
  udfs object diff -v $B $C > diff_out
'

test_expect_success "diff looks right" '
  echo Added new link \"baz/dog\" pointing to QmdNJQUTZuDpsUcec7YDuCfRfvw1w4J13DCm7YcU4VMZdS. > diff_exp &&
  test_cmp diff_exp diff_out
'

test_expect_success "diff changed link works" '
  udfs object diff -v $C $D > diff_out
'

test_expect_success "diff looks right" '
  echo Changed \"bar\" from QmNgd5cz2jNftnAHBhcRUGdtiaMzb5Rhjqd4etondHHST8 to QmRfFVsjSXkhFxrfWnLpMae2M4GBVsry6VAuYYcji5MiZb. > diff_exp &&
  test_cmp diff_exp diff_out
'

test_done
