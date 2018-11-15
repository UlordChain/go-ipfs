#!/usr/bin/env bash

test_description="Test resolve command"

. lib/test-lib.sh

test_init_udfs

test_expect_success "resolve: prepare files" '
  mkdir -p a/b &&
  echo "a/b/c" >a/b/c &&
  a_hash=$(udfs add -q -r a | tail -n1) &&
  b_hash=$(udfs add -q -r a/b | tail -n1) &&
  c_hash=$(udfs add -q -r a/b/c | tail -n1)
'

test_resolve_setup_name() {
  ref=$1

  test_expect_success "resolve: prepare name" '
    id_hash=$(udfs id -f="<id>") &&
    udfs name publish "$ref" &&
    printf "$ref\n" >expected_nameval &&
    udfs name resolve >actual_nameval &&
    test_cmp expected_nameval actual_nameval
  '
}

test_resolve_setup_name_fail() {
  ref=$1

  test_expect_failure "resolve: prepare name" '
    id_hash=$(udfs id -f="<id>") &&
    udfs name publish "$ref" &&
    printf "$ref" >expected_nameval &&
    udfs name resolve >actual_nameval &&
    test_cmp expected_nameval actual_nameval
  '
}

test_resolve() {
  src=$1
  dst=$2

  test_expect_success "resolve succeeds: $src" '
    udfs resolve -r "$src" >actual
  '

  test_expect_success "resolved correctly: $src -> $dst" '
    printf "$dst\n" >expected &&
    test_cmp expected actual
  '
}

test_resolve_cmd() {

  test_resolve "/udfs/$a_hash" "/udfs/$a_hash"
  test_resolve "/udfs/$a_hash/b" "/udfs/$b_hash"
  test_resolve "/udfs/$a_hash/b/c" "/udfs/$c_hash"
  test_resolve "/udfs/$b_hash/c" "/udfs/$c_hash"

  test_resolve_setup_name "/udfs/$a_hash"
  test_resolve "/ipns/$id_hash" "/udfs/$a_hash"
  test_resolve "/ipns/$id_hash/b" "/udfs/$b_hash"
  test_resolve "/ipns/$id_hash/b/c" "/udfs/$c_hash"

  test_resolve_setup_name "/udfs/$b_hash"
  test_resolve "/ipns/$id_hash" "/udfs/$b_hash"
  test_resolve "/ipns/$id_hash/c" "/udfs/$c_hash"

  test_resolve_setup_name "/udfs/$c_hash"
  test_resolve "/ipns/$id_hash" "/udfs/$c_hash"
}

#todo remove this once the online resolve is fixed
test_resolve_fail() {
  src=$1
  dst=$2

  test_expect_failure "resolve succeeds: $src" '
    udfs resolve "$src" >actual
  '

  test_expect_failure "resolved correctly: $src -> $dst" '
    printf "$dst" >expected &&
    test_cmp expected actual
  '
}

test_resolve_cmd_fail() {
  test_resolve "/udfs/$a_hash" "/udfs/$a_hash"
  test_resolve "/udfs/$a_hash/b" "/udfs/$b_hash"
  test_resolve "/udfs/$a_hash/b/c" "/udfs/$c_hash"
  test_resolve "/udfs/$b_hash/c" "/udfs/$c_hash"

  test_resolve_setup_name_fail "/udfs/$a_hash"
  test_resolve_fail "/ipns/$id_hash" "/udfs/$a_hash"
  test_resolve_fail "/ipns/$id_hash/b" "/udfs/$b_hash"
  test_resolve_fail "/ipns/$id_hash/b/c" "/udfs/$c_hash"

  test_resolve_setup_name_fail "/udfs/$b_hash"
  test_resolve_fail "/ipns/$id_hash" "/udfs/$b_hash"
  test_resolve_fail "/ipns/$id_hash/c" "/udfs/$c_hash"

  test_resolve_setup_name_fail "/udfs/$c_hash"
  test_resolve_fail "/ipns/$id_hash" "/udfs/$c_hash"
}

# should work offline
test_resolve_cmd

# should work online
test_launch_udfs_daemon
test_resolve_cmd_fail
test_kill_udfs_daemon

test_done
