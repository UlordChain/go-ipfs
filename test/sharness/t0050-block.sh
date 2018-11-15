#!/usr/bin/env bash
#
# Copyright (c) 2014 Christian Couder
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test block command"

. lib/test-lib.sh

test_init_udfs

HASH="QmRKqGMAM6EZngbpjSqrvYzq5Qd8b1bSWymjSUY9zQSNDk"

#
# "block put tests"
#

test_expect_success "'udfs block put' succeeds" '
  echo "Hello Mars!" >expected_in &&
  udfs block put <expected_in >actual_out
'

test_expect_success "'udfs block put' output looks good" '
  echo "$HASH" >expected_out &&
  test_cmp expected_out actual_out
'

#
# "block get" tests
#

test_expect_success "'udfs block get' succeeds" '
  udfs block get $HASH >actual_in
'

test_expect_success "'udfs block get' output looks good" '
  test_cmp expected_in actual_in
'

#
# "block stat" tests
#

test_expect_success "'udfs block stat' succeeds" '
  udfs block stat $HASH >actual_stat
'

test_expect_success "'udfs block stat' output looks good" '
  echo "Key: $HASH" >expected_stat &&
  echo "Size: 12" >>expected_stat &&
  test_cmp expected_stat actual_stat
'

#
# "block rm" tests
#

test_expect_success "'udfs block rm' succeeds" '
  udfs block rm $HASH >actual_rm
'

test_expect_success "'udfs block rm' output looks good" '
  echo "removed $HASH" > expected_rm &&
  test_cmp expected_rm actual_rm
'

test_expect_success "'udfs block rm' block actually removed" '
  test_must_fail udfs block stat $HASH
'

DIRHASH=QmdWmVmM6W2abTgkEfpbtA1CJyTWS2rhuUB9uP1xV8Uwtf
FILE1HASH=Qmae3RedM7SNkWGsdzYzsr6svmsFdsva4WoTvYYsWhUSVz
FILE2HASH=QmUtkGLvPf63NwVzLPKPUYgwhn8ZYPWF6vKWN3fZ2amfJF
FILE3HASH=Qmesmmf1EEG1orJb6XdK6DabxexsseJnCfw8pqWgonbkoj

test_expect_success "add and pin directory" '
  mkdir adir &&
  echo "file1" > adir/file1 &&
  echo "file2" > adir/file2 &&
  echo "file3" > adir/file3 &&
  udfs add -r adir
  udfs pin add -r $DIRHASH
'

test_expect_success "can't remove pinned block" '
  test_must_fail udfs block rm $DIRHASH 2> block_rm_err
'

test_expect_success "can't remove pinned block: output looks good" '
  grep -q "$DIRHASH: pinned: recursive" block_rm_err
'

test_expect_success "can't remove indirectly pinned block" '
  test_must_fail udfs block rm $FILE1HASH 2> block_rm_err
'

test_expect_success "can't remove indirectly pinned block: output looks good" '
  grep -q "$FILE1HASH: pinned via $DIRHASH" block_rm_err
'

test_expect_success "remove pin" '
  udfs pin rm -r $DIRHASH
'

test_expect_success "multi-block 'udfs block rm' succeeds" '
  udfs block rm $FILE1HASH $FILE2HASH $FILE3HASH > actual_rm
'

test_expect_success "multi-block 'udfs block rm' output looks good" '
  grep -F -q "removed $FILE1HASH" actual_rm &&
  grep -F -q "removed $FILE2HASH" actual_rm &&
  grep -F -q "removed $FILE3HASH" actual_rm
'

test_expect_success "'add some blocks' succeeds" '
  echo "Hello Mars!" | udfs block put &&
  echo "Hello Venus!" | udfs block put
'

test_expect_success "add and pin directory" '
  udfs add -r adir
  udfs pin add -r $DIRHASH
'

HASH=QmRKqGMAM6EZngbpjSqrvYzq5Qd8b1bSWymjSUY9zQSNDk
HASH2=QmdnpnsaEj69isdw5sNzp3h3HkaDz7xKq7BmvFFBzNr5e7
RANDOMHASH=QmRKqGMAM6EbngbZjSqrvYzq5Qd8b1bSWymjSUY9zQSNDq

test_expect_success "multi-block 'udfs block rm' mixed" '
  test_must_fail udfs block rm $FILE1HASH $DIRHASH $HASH $FILE3HASH $RANDOMHASH $HASH2 2> block_rm_err
'

test_expect_success "pinned block not removed" '
  udfs block stat $FILE1HASH &&
  udfs block stat $FILE3HASH
'

test_expect_success "non-pinned blocks removed" '
  test_must_fail udfs block stat $HASH &&
  test_must_fail udfs block stat $HASH2
'

test_expect_success "error reported on removing non-existent block" '
  grep -q "cannot remove $RANDOMHASH" block_rm_err
'

test_expect_success "'add some blocks' succeeds" '
  echo "Hello Mars!" | udfs block put &&
  echo "Hello Venus!" | udfs block put
'

test_expect_success "multi-block 'udfs block rm -f' with non existent blocks succeed" '
  udfs block rm -f $HASH $RANDOMHASH $HASH2
'

test_expect_success "existent blocks removed" '
  test_must_fail udfs block stat $HASH &&
  test_must_fail udfs block stat $HASH2
'

test_expect_success "'add some blocks' succeeds" '
  echo "Hello Mars!" | udfs block put &&
  echo "Hello Venus!" | udfs block put
'

test_expect_success "multi-block 'udfs block rm -q' produces no output" '
  udfs block rm -q $HASH $HASH2 > block_rm_out &&
  test ! -s block_rm_out
'

test_expect_success "can set cid format on block put" '
  HASH=$(udfs block put --format=protobuf ../t0051-object-data/testPut.pb)
'

test_expect_success "created an object correctly!" '
  udfs object get $HASH > obj_out &&
  echo "{\"Links\":[],\"Data\":\"test json for sharness test\"}" > obj_exp &&
  test_cmp obj_out obj_exp
'

test_expect_success "block get output looks right" '
  udfs block get $HASH > pb_block_out &&
  test_cmp pb_block_out ../t0051-object-data/testPut.pb
'

test_expect_success "can set multihash type and length on block put" '
  HASH=$(echo "foooo" | udfs block put --format=raw --mhtype=sha3 --mhlen=20)
'

test_expect_success "output looks good" '
  test "z83bYcqyBkbx5fuNAcvbdv4pr5RYQiEpK" = "$HASH"
'

test_expect_success "can read block with different hash" '
  udfs block get $HASH > blk_get_out &&
  echo "foooo" > blk_get_exp &&
  test_cmp blk_get_exp blk_get_out
'
#
# Misc tests
#

test_expect_success "'udfs block stat' with nothing from stdin doesnt crash" '
  test_expect_code 1 udfs block stat < /dev/null 2> stat_out
'

test_expect_success "no panic in output" '
  test_expect_code 1 grep "panic" stat_out
'

test_expect_success "can set multihash type and length on block put without format" '
  HASH=$(echo "foooo" | udfs block put --mhtype=sha3 --mhlen=20)
'

test_expect_success "output looks good" '
  test "z8bwYCvQPhyDY7VUTsUdGdE8Evm1ktSPV" = "$HASH"
'

test_expect_success "put with sha3 and cidv0 fails" '
  echo "foooo" | test_must_fail udfs block put --mhtype=sha3 --mhlen=20 --format=v0
'

test_done
