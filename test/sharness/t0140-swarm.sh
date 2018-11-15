#!/usr/bin/env bash
#
# Copyright (c) 2014 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test udfs swarm command"

. lib/test-lib.sh

test_init_udfs

test_launch_udfs_daemon

test_expect_success 'disconnected: peers is empty' '
  udfs swarm peers >actual &&
  test_must_be_empty actual
'

test_expect_success 'disconnected: addrs local has localhost' '
  udfs swarm addrs local >actual &&
  grep "/ip4/127.0.0.1" actual
'

test_expect_success 'disconnected: addrs local matches udfs id' '
  udfs id -f="<addrs>\\n" | sort >expected &&
  udfs swarm addrs local --id | sort >actual &&
  test_cmp expected actual
'

test_expect_success "udfs id self works" '
  myid=$(udfs id -f="<id>") &&
  udfs id --timeout=1s $myid > output
'

test_expect_success "output looks good" '
  grep $myid output &&
  grep PublicKey output
'

addr="/ip4/127.0.0.1/tcp/9898/udfs/QmUWKoHbjsqsSMesRC2Zoscs8edyFz6F77auBB1YBBhgpX"

test_expect_success "cant trigger a dial backoff with swarm connect" '
  test_expect_code 1 udfs swarm connect $addr 2> connect_out
  test_expect_code 1 udfs swarm connect $addr 2>> connect_out
  test_expect_code 1 udfs swarm connect $addr 2>> connect_out
  test_expect_code 1 grep "backoff" connect_out
'

test_kill_udfs_daemon

announceCfg='["/ip4/127.0.0.1/tcp/4001", "/ip4/1.2.3.4/tcp/1234"]'
test_expect_success "test_config_set succeeds" "
  udfs config --json Addresses.Announce '$announceCfg'
"

test_launch_udfs_daemon

test_expect_success 'Addresses.Announce affects addresses' '
  udfs swarm addrs local >actual &&
  grep "/ip4/1.2.3.4/tcp/1234" actual &&
  udfs id -f"<addrs>" | xargs -n1 echo >actual &&
  grep "/ip4/1.2.3.4/tcp/1234" actual
'

test_kill_udfs_daemon

noAnnounceCfg='["/ip4/1.2.3.4/tcp/1234"]'
test_expect_success "test_config_set succeeds" "
  udfs config --json Addresses.NoAnnounce '$noAnnounceCfg'
"

test_launch_udfs_daemon

test_expect_success "Addresses.NoAnnounce affects addresses" '
  udfs swarm addrs local >actual &&
  grep -v "/ip4/1.2.3.4/tcp/1234" actual &&
  udfs id -f"<addrs>" | xargs -n1 echo >actual &&
  grep -v "/ip4/1.2.3.4/tcp/1234" actual
'

test_kill_udfs_daemon

noAnnounceCfg='["/ip4/1.2.3.4/ipcidr/16"]'
test_expect_success "test_config_set succeeds" "
  udfs config --json Addresses.NoAnnounce '$noAnnounceCfg'
"

test_launch_udfs_daemon

test_expect_success "Addresses.NoAnnounce with /ipcidr affects addresses" '
  udfs swarm addrs local >actual &&
  grep -v "/ip4/1.2.3.4/tcp/1234" actual &&
  udfs id -f"<addrs>" | xargs -n1 echo >actual &&
  grep -v "/ip4/1.2.3.4/tcp/1234" actual
'

test_kill_udfs_daemon

test_expect_success "set up tcp testbed" '
  iptb init -n 2 -p 0 -f --bootstrap=none
'

startup_cluster 2

test_expect_success "disconnect work without specifying a transport address" '
  [ $(udfsi 0 swarm peers | wc -l) -eq 1 ] &&
  udfsi 0 swarm disconnect "/udfs/$(iptb get id 1)" &&
  [ $(udfsi 0 swarm peers | wc -l) -eq 0 ]
'

test_expect_success "connect work without specifying a transport address" '
  [ $(udfsi 0 swarm peers | wc -l) -eq 0 ] &&
  udfsi 0 swarm connect "/udfs/$(iptb get id 1)" &&
  [ $(udfsi 0 swarm peers | wc -l) -eq 1 ]
'

test_expect_success "/p2p addresses work" '
  [ $(udfsi 0 swarm peers | wc -l) -eq 1 ] &&
  udfsi 0 swarm disconnect "/p2p/$(iptb get id 1)" &&
  [ $(udfsi 0 swarm peers | wc -l) -eq 0 ] &&
  udfsi 0 swarm connect "/p2p/$(iptb get id 1)" &&
  [ $(udfsi 0 swarm peers | wc -l) -eq 1 ]
'

test_expect_success "stopping cluster" '
  iptb stop
'

test_done
