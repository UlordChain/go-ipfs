#!/usr/bin/env bash
#
# Copyright (c) 2014 Christian Couder
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test init command"

. lib/test-lib.sh

# test that udfs fails to init if UDFS_PATH isnt writeable
test_expect_success "create dir and change perms succeeds" '
  export UDFS_PATH="$(pwd)/.badudfs" &&
  mkdir "$UDFS_PATH" &&
  chmod 000 "$UDFS_PATH"
'

test_expect_success "udfs init fails" '
  test_must_fail udfs init 2> init_fail_out
'

# Under Windows/Cygwin the error message is different,
# so we use the STD_ERR_MSG prereq.
if test_have_prereq STD_ERR_MSG; then
  init_err_msg="Error: error opening repository at $UDFS_PATH: permission denied"
else
  init_err_msg="Error: mkdir $UDFS_PATH: The system cannot find the path specified."
fi

test_expect_success "udfs init output looks good" '
  echo "$init_err_msg" >init_fail_exp &&
  test_cmp init_fail_exp init_fail_out
'

test_expect_success "cleanup dir with bad perms" '
  chmod 775 "$UDFS_PATH" &&
  rmdir "$UDFS_PATH"
'

# test no repo error message
# this applies to `udfs add sth`, `udfs refs <hash>`
test_expect_success "udfs cat fails" '
  export UDFS_PATH="$(pwd)/.udfs" &&
  test_must_fail udfs cat Qmaa4Rw81a3a1VEx4LxB7HADUAXvZFhCoRdBzsMZyZmqHD 2> cat_fail_out
'

test_expect_success "udfs cat no repo message looks good" '
  echo "Error: no UDFS repo found in $UDFS_PATH." > cat_fail_exp &&
  echo "please run: '"'"'udfs init'"'"'" >> cat_fail_exp &&
  test_path_cmp cat_fail_exp cat_fail_out
'

# test that init succeeds
test_expect_success "udfs init succeeds" '
  export UDFS_PATH="$(pwd)/.udfs" &&
  echo "UDFS_PATH: \"$UDFS_PATH\"" &&
  BITS="2048" &&
  udfs init --bits="$BITS" >actual_init ||
  test_fsh cat actual_init
'

test_expect_success ".udfs/ has been created" '
  test -d ".udfs" &&
  test -f ".udfs/config" &&
  test -d ".udfs/datastore" &&
  test -d ".udfs/blocks" &&
  test ! -f ._check_writeable ||
  test_fsh ls -al .udfs
'

test_expect_success "udfs config succeeds" '
  echo /udfs >expected_config &&
  udfs config Mounts.UDFS >actual_config &&
  test_cmp expected_config actual_config
'

test_expect_success "udfs peer id looks good" '
  PEERID=$(udfs config Identity.PeerID) &&
  test_check_peerid "$PEERID"
'

test_expect_success "udfs init output looks good" '
  STARTFILE="udfs cat /udfs/$HASH_WELCOME_DOCS/readme" &&
  echo "initializing UDFS node at $UDFS_PATH" >expected &&
  echo "generating $BITS-bit RSA keypair...done" >>expected &&
  echo "peer identity: $PEERID" >>expected &&
  echo "to get started, enter:" >>expected &&
  printf "\\n\\t$STARTFILE\\n\\n" >>expected &&
  test_cmp expected actual_init
'

test_expect_success "Welcome readme exists" '
  udfs cat /udfs/$HASH_WELCOME_DOCS/readme
'

test_expect_success "clean up udfs dir" '
  rm -rf "$UDFS_PATH"
'

test_expect_success "'udfs init --empty-repo' succeeds" '
  BITS="1024" &&
  udfs init --bits="$BITS" --empty-repo >actual_init
'

test_expect_success "udfs peer id looks good" '
  PEERID=$(udfs config Identity.PeerID) &&
  test_check_peerid "$PEERID"
'

test_expect_success "'udfs init --empty-repo' output looks good" '
  echo "initializing UDFS node at $UDFS_PATH" >expected &&
  echo "generating $BITS-bit RSA keypair...done" >>expected &&
  echo "peer identity: $PEERID" >>expected &&
  test_cmp expected actual_init
'

test_expect_success "Welcome readme doesn't exists" '
  test_must_fail udfs cat /udfs/$HASH_WELCOME_DOCS/readme
'

test_expect_success "udfs id agent string contains correct version" '
  udfs id -f "<aver>" | grep $(udfs version -n)
'

test_expect_success "clean up udfs dir" '
  rm -rf "$UDFS_PATH"
'

# test init profiles
test_expect_success "'udfs init --profile' with invalid profile fails" '
  BITS="1024" &&
  test_must_fail udfs init --bits="$BITS" --profile=nonexistent_profile 2> invalid_profile_out
  EXPECT="Error: invalid configuration profile: nonexistent_profile" &&
  grep "$EXPECT" invalid_profile_out
'

test_expect_success "'udfs init --profile' succeeds" '
  BITS="1024" &&
  udfs init --bits="$BITS" --profile=server
'

test_expect_success "'udfs config Swarm.AddrFilters' looks good" '
  udfs config Swarm.AddrFilters > actual_config &&
  test $(cat actual_config | wc -l) = 17
'

test_expect_success "clean up udfs dir" '
  rm -rf "$UDFS_PATH"
'

test_expect_success "'udfs init --profile=test' succeeds" '
  BITS="1024" &&
  udfs init --bits="$BITS" --profile=test
'

test_expect_success "'udfs config Bootstrap' looks good" '
  udfs config Bootstrap > actual_config &&
  test $(cat actual_config) = "[]"
'

test_expect_success "'udfs config Addresses.API' looks good" '
  udfs config Addresses.API > actual_config &&
  test $(cat actual_config) = "/ip4/127.0.0.1/tcp/0"
'

test_expect_success "clean up udfs dir" '
  rm -rf "$UDFS_PATH"
'

test_expect_success "'udfs init --profile=lowpower' succeeds" '
  BITS="1024" &&
  udfs init --bits="$BITS" --profile=lowpower
'

test_expect_success "'udfs config Discovery.Routing' looks good" '
  udfs config Routing.Type > actual_config &&
  test $(cat actual_config) = "dhtclient"
'

test_expect_success "clean up udfs dir" '
  rm -rf "$UDFS_PATH"
'

test_init_udfs

test_launch_udfs_daemon

test_expect_success "udfs init should not run while daemon is running" '
  test_must_fail udfs init 2> daemon_running_err &&
  EXPECT="Error: udfs daemon is running. please stop it to run this command" &&
  grep "$EXPECT" daemon_running_err
'

test_kill_udfs_daemon

test_done
