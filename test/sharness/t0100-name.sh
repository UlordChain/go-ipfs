#!/usr/bin/env bash
#
# Copyright (c) 2014 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test udfs repo operations"

. lib/test-lib.sh

test_init_udfs

# test publishing a hash

test_expect_success "'udfs name publish' succeeds" '
  PEERID=`udfs id --format="<id>"` &&
  test_check_peerid "${PEERID}" &&
  udfs name publish "/udfs/$HASH_WELCOME_DOCS" >publish_out
'

test_expect_success "publish output looks good" '
  echo "Published to ${PEERID}: /udfs/$HASH_WELCOME_DOCS" >expected1 &&
  test_cmp expected1 publish_out
'

test_expect_success "'udfs name resolve' succeeds" '
  udfs name resolve "$PEERID" >output
'

test_expect_success "resolve output looks good" '
  printf "/udfs/%s\n" "$HASH_WELCOME_DOCS" >expected2 &&
  test_cmp expected2 output
'

# now test with a path

test_expect_success "'udfs name publish' succeeds" '
  PEERID=`udfs id --format="<id>"` &&
  test_check_peerid "${PEERID}" &&
  udfs name publish "/udfs/$HASH_WELCOME_DOCS/help" >publish_out
'

test_expect_success "publish a path looks good" '
  echo "Published to ${PEERID}: /udfs/$HASH_WELCOME_DOCS/help" >expected3 &&
  test_cmp expected3 publish_out
'

test_expect_success "'udfs name resolve' succeeds" '
  udfs name resolve "$PEERID" >output
'

test_expect_success "resolve output looks good" '
  printf "/udfs/%s/help\n" "$HASH_WELCOME_DOCS" >expected4 &&
  test_cmp expected4 output
'

test_expect_success "udfs cat on published content succeeds" '
  udfs cat "/udfs/$HASH_WELCOME_DOCS/help" >expected &&
  udfs cat "/ipns/$PEERID" >actual &&
  test_cmp expected actual
'

# publish with an explicit node ID

test_expect_failure "'udfs name publish <local-id> <hash>' succeeds" '
  PEERID=`udfs id --format="<id>"` &&
  test_check_peerid "${PEERID}" &&
  echo udfs name publish "${PEERID}" "/udfs/$HASH_WELCOME_DOCS" &&
  udfs name publish "${PEERID}" "/udfs/$HASH_WELCOME_DOCS" >actual_node_id_publish
'

test_expect_failure "publish with our explicit node ID looks good" '
  echo "Published to ${PEERID}: /udfs/$HASH_WELCOME_DOCS" >expected_node_id_publish &&
  test_cmp expected_node_id_publish actual_node_id_publish
'

# publish with an explicit node ID as key name

test_expect_success "generate and verify a new key" '
  NEWID=`udfs key gen --type=rsa --size=2048 keyname` &&
  test_check_peerid "${NEWID}"
'

test_expect_success "'udfs name publish --key=<peer-id> <hash>' succeeds" '
  udfs name publish --key=${NEWID} "/udfs/$HASH_WELCOME_DOCS" >actual_node_id_publish
'

test_expect_success "publish an explicit node ID as key name looks good" '
  echo "Published to ${NEWID}: /udfs/$HASH_WELCOME_DOCS" >expected_node_id_publish &&
  test_cmp expected_node_id_publish actual_node_id_publish
'


# test publishing nothing

test_expect_success "'udfs name publish' fails" '
  printf '' | test_expect_code 1 udfs name publish >publish_out 2>&1
'

test_expect_success "publish output has the correct error" '
  grep "argument \"udfs-path\" is required" publish_out
'

test_expect_success "'udfs name publish --help' succeeds" '
  udfs name publish --help
'

test_launch_udfs_daemon

test_expect_success "empty request to name publish doesn't panic and returns error" '
  curl "http://$API_ADDR/api/v0/name/publish" > curl_out || true &&
    grep "argument \"udfs-path\" is required" curl_out
'

test_kill_udfs_daemon


test_done
