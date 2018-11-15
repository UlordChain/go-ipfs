#!/usr/bin/env bash
#
# Copyright (c) 2017 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test out the urlstore functionality"

. lib/test-lib.sh

test_init_udfs

test_expect_success "create some random files" '
  random 2222     7 > file1 &&
  random 500000   7 > file2 &&
  random 50000000 7 > file3
'

test_expect_success "add files using trickle dag format without raw leaves" '
  HASH1a=$(udfs add -q --trickle --raw-leaves=false file1) &&
  HASH2a=$(udfs add -q --trickle --raw-leaves=false file2) &&
  HASH3a=$(udfs add -q --trickle --raw-leaves=false file3)
'
test_launch_udfs_daemon --offline

test_expect_success "make sure files can be retrived via the gateway" '
  curl http://127.0.0.1:$GWAY_PORT/udfs/$HASH1a -o file1.actual &&
  test_cmp file1 file1.actual &&
  curl http://127.0.0.1:$GWAY_PORT/udfs/$HASH2a -o file2.actual &&
  test_cmp file2 file2.actual &&
  curl http://127.0.0.1:$GWAY_PORT/udfs/$HASH3a -o file3.actual &&
  test_cmp file3 file3.actual 
'

test_expect_success "add files without enabling url store" '
  test_must_fail udfs urlstore add http://127.0.0.1:$GWAY_PORT/udfs/$HASH1a &&
  test_must_fail udfs urlstore add http://127.0.0.1:$GWAY_PORT/udfs/$HASH2a
'

test_kill_udfs_daemon

test_expect_success "enable urlstore" '
  udfs config --json Experimental.UrlstoreEnabled true
'

test_launch_udfs_daemon --offline

test_expect_success "add files using gateway address via url store" '
  HASH1=$(udfs urlstore add http://127.0.0.1:$GWAY_PORT/udfs/$HASH1a) &&
  HASH2=$(udfs urlstore add http://127.0.0.1:$GWAY_PORT/udfs/$HASH2a) &&
  udfs pin add $HASH1 $HASH2
'

test_expect_success "make sure hashes are different" '
  test $HASH1a != $HASH1 &&
  test $HASH2a != $HASH2
'

test_expect_success "get files via urlstore" '
  udfs get $HASH1 -o file1.actual &&
  test_cmp file1 file1.actual &&
  udfs get $HASH2 -o file2.actual &&
  test_cmp file2 file2.actual
'

cat <<EOF | sort > ls_expect
zb2rhX1q5oFFzEkPNsTe1Y8osUdFqSQGjUWRZsqC9fbY6WVSk  262144 http://127.0.0.1:$GWAY_PORT/udfs/QmUow2T4P69nEsqTQDZCt8yg9CPS8GFmpuDAr5YtsPhTdM 0
zb2rhYbKFn1UWGHXaAitcdVTkDGTykX8RFpGWzRFuLpoe9VE4  237856 http://127.0.0.1:$GWAY_PORT/udfs/QmUow2T4P69nEsqTQDZCt8yg9CPS8GFmpuDAr5YtsPhTdM 262144
zb2rhjddJ5DNzBrFu8G6CP1ApY25BukwCeskXHzN1H18CiVVZ    2222 http://127.0.0.1:$GWAY_PORT/udfs/QmcHm3BL2cXuQ6rJdKQgPrmT9suqGkfy2KzH3MkXPEBXU6 0
EOF

test_expect_success "udfs filestore ls works with urls" '
  udfs filestore ls | sort > ls_actual &&
  test_cmp ls_expect ls_actual
'

cat <<EOF | sort > verify_expect
ok      zb2rhX1q5oFFzEkPNsTe1Y8osUdFqSQGjUWRZsqC9fbY6WVSk  262144 http://127.0.0.1:$GWAY_PORT/udfs/QmUow2T4P69nEsqTQDZCt8yg9CPS8GFmpuDAr5YtsPhTdM 0
ok      zb2rhYbKFn1UWGHXaAitcdVTkDGTykX8RFpGWzRFuLpoe9VE4  237856 http://127.0.0.1:$GWAY_PORT/udfs/QmUow2T4P69nEsqTQDZCt8yg9CPS8GFmpuDAr5YtsPhTdM 262144
ok      zb2rhjddJ5DNzBrFu8G6CP1ApY25BukwCeskXHzN1H18CiVVZ    2222 http://127.0.0.1:$GWAY_PORT/udfs/QmcHm3BL2cXuQ6rJdKQgPrmT9suqGkfy2KzH3MkXPEBXU6 0
EOF

test_expect_success "udfs filestore verify works with urls" '
  udfs filestore verify | sort > verify_actual &&
  test_cmp verify_expect verify_actual
'

test_expect_success "remove original hashes from local gateway" '
  udfs pin rm $HASH1a $HASH2a &&
  udfs repo gc > /dev/null
'

test_expect_success "gatway no longer has files" '
  test_must_fail curl -f http://127.0.0.1:$GWAY_PORT/udfs/$HASH1a -o file1.actual
  test_must_fail curl -f http://127.0.0.1:$GWAY_PORT/udfs/$HASH2a -o file2.actual
'

cat <<EOF | sort > verify_expect_2
error   zb2rhX1q5oFFzEkPNsTe1Y8osUdFqSQGjUWRZsqC9fbY6WVSk  262144 http://127.0.0.1:$GWAY_PORT/udfs/QmUow2T4P69nEsqTQDZCt8yg9CPS8GFmpuDAr5YtsPhTdM 0
error   zb2rhYbKFn1UWGHXaAitcdVTkDGTykX8RFpGWzRFuLpoe9VE4  237856 http://127.0.0.1:$GWAY_PORT/udfs/QmUow2T4P69nEsqTQDZCt8yg9CPS8GFmpuDAr5YtsPhTdM 262144
error   zb2rhjddJ5DNzBrFu8G6CP1ApY25BukwCeskXHzN1H18CiVVZ    2222 http://127.0.0.1:$GWAY_PORT/udfs/QmcHm3BL2cXuQ6rJdKQgPrmT9suqGkfy2KzH3MkXPEBXU6 0
EOF

test_expect_success "udfs filestore verify is correct" '
  udfs filestore verify | sort > verify_actual_2 &&
  test_cmp verify_expect_2 verify_actual_2
'

test_expect_success "files can not be retrieved via the urlstore" '
  test_must_fail udfs cat $HASH1 > /dev/null &&
  test_must_fail udfs cat $HASH2 > /dev/null
'

test_expect_success "add large file using gateway address via url store" '
  HASH3=$(udfs urlstore add http://127.0.0.1:$GWAY_PORT/udfs/$HASH3a)
'

test_expect_success "make sure hashes are different" '
  test $HASH3a != $HASH3
'

test_expect_success "get large file via urlstore" '
  udfs get $HASH3 -o file3.actual &&
  test_cmp file3 file3.actual
'

test_expect_success "check that the trickle option works" '
  HASHat=$(udfs add -q --cid-version=1 --raw-leaves=true -n --trickle file3) &&
  HASHut=$(udfs urlstore add --trickle http://127.0.0.1:$GWAY_PORT/udfs/$HASH3a) &&
  test $HASHat = $HASHut
'

test_kill_udfs_daemon

test_expect_success "files can not be retrieved via the urlstore" '
  test_must_fail udfs cat $HASH1 > /dev/null &&
  test_must_fail udfs cat $HASH2 > /dev/null &&
  test_must_fail udfs cat $HASH3 > /dev/null
'

test_expect_success "check that the hashes were correct" '
  HASH1e=$(udfs add -q -n --cid-version=1 --raw-leaves=true file1) &&
  HASH2e=$(udfs add -q -n --cid-version=1 --raw-leaves=true file2) &&
  HASH3e=$(udfs add -q -n --cid-version=1 --raw-leaves=true file3) &&
  test $HASH1e = $HASH1 &&
  test $HASH2e = $HASH2 &&
  test $HASH3e = $HASH3
'

test_done
