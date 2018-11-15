#!/usr/bin/env bash
#
# Copyright (c) 2017 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test keystore commands"

. lib/test-lib.sh

test_init_udfs

test_key_cmd() {
  test_expect_success "create a new rsa key" '
    rsahash=$(udfs key gen foobarsa --type=rsa --size=2048)
  '

  test_expect_success "create a new ed25519 key" '
    edhash=$(udfs key gen bazed --type=ed25519)
  '

  test_expect_success "both keys show up in list output" '
    echo bazed > list_exp &&
    echo foobarsa >> list_exp &&
    echo self >> list_exp
    udfs key list | sort > list_out &&
    test_cmp list_exp list_out
  '

  test_expect_success "key hashes show up in long list output" '
    udfs key list -l | grep $edhash > /dev/null &&
    udfs key list -l | grep $rsahash > /dev/null
  '

  test_expect_success "key list -l contains self key with peerID" '
    PeerID="$(udfs config Identity.PeerID)"
    udfs key list -l | grep "$PeerID\s\+self"
  '

  test_expect_success "key rm remove a key" '
    udfs key rm foobarsa
    echo bazed > list_exp &&
    echo self >> list_exp
    udfs key list | sort > list_out &&
    test_cmp list_exp list_out
  '

  test_expect_success "key rm can't remove self" '
    test_must_fail udfs key rm self 2>&1 | tee key_rm_out &&
    grep -q "Error: cannot remove key with name" key_rm_out
  '

  test_expect_success "key rename rename a key" '
    udfs key rename bazed fooed
    echo fooed > list_exp &&
    echo self >> list_exp
    udfs key list | sort > list_out &&
    test_cmp list_exp list_out
  '

  test_expect_success "key rename rename key output succeeds" '
    key_content=$(udfs key gen key1 --type=rsa --size=2048) &&
    udfs key rename key1 key2 >rs &&
    echo "Key $key_content renamed to key2" >expect &&
    test_cmp rs expect
  '

  test_expect_success "key rename can't rename self" '
    test_must_fail udfs key rename self bar 2>&1 | tee key_rename_out &&
    grep -q "Error: cannot rename key with name" key_rename_out
  '

  test_expect_success "key rename can't overwrite self, even with force" '
    test_must_fail udfs key rename -f fooed self 2>&1 | tee key_rename_out &&
    grep -q "Error: cannot overwrite key with name" key_rename_out
  '
}

test_key_cmd

test_done
