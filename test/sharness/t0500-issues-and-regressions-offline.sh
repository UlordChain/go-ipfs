#!/usr/bin/env bash

test_description="Tests for various fixed issues and regressions."

. lib/test-lib.sh

# Tests go here

test_expect_success "udfs init with occupied input works - #2748" '
  export UDFS_PATH="udfs_path"
  echo "" | go-timeout 10 udfs init &&
  rm -rf udfs_path
'
test_init_udfs

test_expect_success "udfs cat --help succeeds when input remains open" '
  yes | go-timeout 1 udfs cat --help
'

test_expect_success "udfs pin ls --help succeeds when input remains open" '
  yes | go-timeout 1 udfs pin ls --help
'

test_expect_success "udfs add on 1MB from stdin woks" '
  random 1048576 42 | udfs add -q > 1MB.hash
'

test_expect_success "'udfs refs -r -e \$(cat 1MB.hash)' succeeds" '
  udfs refs -r -e $(cat 1MB.hash) > refs-e.out
'

test_expect_success "output of 'udfs refs -e' links to separate blocks" '
  grep "$(cat 1MB.hash) ->" refs-e.out
'

test_expect_success "output of 'udfs refs -e' contains all first level links" '
  grep "$(cat 1MB.hash) ->" refs-e.out | sed -e '\''s/.* -> //'\'' | sort > refs-s.out &&
  udfs refs "$(cat 1MB.hash)" | sort > refs-one.out &&
  test_cmp refs-s.out refs-one.out
'

test_done
