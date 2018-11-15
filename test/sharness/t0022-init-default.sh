#!/usr/bin/env bash
#
# Copyright (c) 2014 Christian Couder
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test init command with default config"

. lib/test-lib.sh

cfg_key="Addresses.API"
cfg_val="/ip4/0.0.0.0/tcp/5001"

# test that init succeeds
test_expect_success "udfs init succeeds" '
  export UDFS_PATH="$(pwd)/.udfs" &&
  echo "UDFS_PATH: \"$UDFS_PATH\"" &&
  BITS="2048" &&
  udfs init --bits="$BITS" >actual_init ||
  test_fsh cat actual_init
'

test_expect_success ".udfs/config has been created" '
  test -f "$UDFS_PATH"/config ||
  test_fsh ls -al .udfs
'

test_expect_success "udfs config succeeds" '
  udfs config $cfg_flags "$cfg_key" "$cfg_val"
'

test_expect_success "udfs read config succeeds" '
  UDFS_DEFAULT_CONFIG=$(cat "$UDFS_PATH"/config)
'

test_expect_success "clean up udfs dir" '
  rm -rf "$UDFS_PATH"
'

test_expect_success "udfs init default config succeeds" '
  echo $UDFS_DEFAULT_CONFIG | udfs init - >actual_init ||
  test_fsh cat actual_init
'

test_expect_success "udfs config output looks good" '
  echo "$cfg_val" >expected &&
  udfs config "$cfg_key" >actual &&
  test_cmp expected actual
'

test_done
