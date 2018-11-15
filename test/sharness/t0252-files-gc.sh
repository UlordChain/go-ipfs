#!/usr/bin/env bash
#
# Copyright (c) 2016 Kevin Atkinson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="test how the unix files api interacts with the gc"

. lib/test-lib.sh

test_init_udfs

test_expect_success "object not removed after gc" '
  echo "hello world" > hello.txt &&
  cat hello.txt | udfs files write --create /hello.txt &&
  udfs repo gc &&
  udfs cat QmVib14uvPnCP73XaCDpwugRuwfTsVbGyWbatHAmLSdZUS
'

test_expect_success "/hello.txt still accessible after gc" '
  udfs files read /hello.txt > hello-actual &&
  test_cmp hello.txt hello-actual
'

ADIR_HASH=QmbCgoMYVuZq8m1vK31JQx9DorwQdLMF1M3sJ7kygLLqnW
FILE1_HASH=QmX4eaSJz39mNhdu5ACUwTDpyA6y24HmrQNnAape6u3buS

test_expect_success "gc okay after adding incomplete node -- prep" '
  udfs files mkdir /adir &&
  echo "file1" |  udfs files write --create /adir/file1 &&
  echo "file2" |  udfs files write --create /adir/file2 &&
  udfs pin add --recursive=false $ADIR_HASH &&
  udfs files rm -r /adir &&
  udfs repo gc && # will remove /adir/file1 and /adir/file2 but not /adir
  test_must_fail udfs cat $FILE1_HASH &&
  udfs files cp /udfs/$ADIR_HASH /adir &&
  udfs pin rm $ADIR_HASH
'

test_expect_success "gc okay after adding incomplete node" '
  udfs refs $ADIR_HASH &&
  udfs repo gc &&
  udfs refs $ADIR_HASH
'

test_expect_success "add directory with direct pin" '
  mkdir mydir/ &&
  echo "hello world!" > mydir/hello.txt &&
  FILE_UNPINNED=$(udfs add --pin=false -q -r mydir/hello.txt) &&
  DIR_PINNED=$(udfs add --pin=false -q -r mydir | tail -n1) &&
  udfs add --pin=false -r mydir &&
  udfs pin add --recursive=false $DIR_PINNED &&
  udfs cat $FILE_UNPINNED
'

test_expect_success "run gc and make sure directory contents are removed" '
  udfs repo gc &&
  test_must_fail udfs cat $FILE_UNPINNED
'

test_expect_success "add incomplete directory and make sure gc is okay" '
  udfs files cp /udfs/$DIR_PINNED /mydir &&
  udfs repo gc &&
  test_must_fail udfs cat $FILE_UNPINNED
'

test_expect_success "add back directory contents and run gc" '
  udfs add --pin=false mydir/hello.txt &&
  udfs repo gc
'

test_expect_success "make sure directory contents are not removed" '
  udfs cat $FILE_UNPINNED
'

test_done
