#!/usr/bin/env bash
#
# Copyright (c) 2015 Matt Bell
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test get command"

. lib/test-lib.sh

test_init_udfs

test_udfs_get_flag() {
  ext="$1"; shift
  tar_flag="$1"; shift
  flag="$@"

  test_expect_success "udfs get $flag succeeds" '
    udfs get "$HASH" '"$flag"' >actual
  '

  test_expect_success "udfs get $flag output looks good" '
    printf "%s\n" "Saving archive to $HASH$ext" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs get $flag archive output is valid" '
    tar "$tar_flag" "$HASH$ext" &&
    test_cmp "$HASH" data &&
    rm "$HASH$ext" &&
    rm "$HASH"
  '
}

# we use a function so that we can run it both offline + online
test_get_cmd() {

  test_expect_success "'udfs get --help' succeeds" '
    udfs get --help >actual
  '

  test_expect_success "'udfs get --help' output looks good" '
    egrep "udfs get.*<udfs-path>" actual >/dev/null ||
    test_fsh cat actual
  '

  test_expect_success "udfs get succeeds" '
    echo "Hello Worlds!" >data &&
    HASH=`udfs add -q data` &&
    udfs get "$HASH" >actual
  '

  test_expect_success "udfs get output looks good" '
    printf "%s\n" "Saving file(s) to $HASH" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs get file output looks good" '
    test_cmp "$HASH" data
  '

  test_expect_success "udfs get DOES NOT error when trying to overwrite a file" '
    udfs get "$HASH" >actual &&
    rm "$HASH"
  '

  test_expect_success "udfs get works with raw leaves" '
  HASH2=$(udfs add --raw-leaves -q data) &&
    udfs get "$HASH2" >actual2
  '

  test_expect_success "udfs get output looks good" '
    printf "%s\n" "Saving file(s) to $HASH2" >expected2 &&
    test_cmp expected2 actual2
  '

  test_expect_success "udfs get file output looks good" '
    test_cmp "$HASH2" data
  '

  test_udfs_get_flag ".tar" "-xf" -a

  test_udfs_get_flag ".tar.gz" "-zxf" -a -C

  test_udfs_get_flag ".tar.gz" "-zxf" -a -C -l 9

  test_expect_success "udfs get succeeds (directory)" '
    mkdir -p dir &&
    touch dir/a &&
    mkdir -p dir/b &&
    echo "Hello, Worlds!" >dir/b/c &&
    HASH2=`udfs add -r -q dir | tail -n 1` &&
    udfs get "$HASH2" >actual
  '

  test_expect_success "udfs get output looks good (directory)" '
    printf "%s\n" "Saving file(s) to $HASH2" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs get output is valid (directory)" '
    test_cmp dir/a "$HASH2"/a &&
    test_cmp dir/b/c "$HASH2"/b/c &&
    rm -r "$HASH2"
  '

  # Test issue #4720: problems when path contains a trailing slash.
  test_expect_success "udfs get with slash (directory)" '
    udfs get "$HASH2/" &&
    test_cmp dir/a "$HASH2"/a &&
    test_cmp dir/b/c "$HASH2"/b/c &&
    rm -r "$HASH2"
  '

  test_expect_success "udfs get -a -C succeeds (directory)" '
    udfs get "$HASH2" -a -C >actual
  '

  test_expect_success "udfs get -a -C output looks good (directory)" '
    printf "%s\n" "Saving archive to $HASH2.tar.gz" >expected &&
    test_cmp expected actual
  '

  test_expect_success "gzipped tar archive output is valid (directory)" '
    tar -zxf "$HASH2".tar.gz &&
    test_cmp dir/a "$HASH2"/a &&
    test_cmp dir/b/c "$HASH2"/b/c &&
    rm -r "$HASH2"
  '

  test_expect_success "udfs get ../.. should fail" '
    echo "Error: invalid 'udfs ref' path" >expected &&
    test_must_fail udfs get ../.. 2>actual &&
    test_cmp expected actual
  '

  test_expect_success "create small file" '
    echo "foo" > small &&
    udfs add -q small > hash_small
  '

  test_expect_success "get small file" '
    udfs get -o out_small $(cat hash_small) &&
    test_cmp small out_small
  '

  test_expect_success "create medium file" '
    head -c 16000 > medium &&
    udfs add -q medium > hash_medium
  '

  test_expect_success "get medium file" '
    udfs get -o out_medium $(cat hash_medium) &&
    test_cmp medium out_medium
  '
}

test_get_fail() {
  test_expect_success "create an object that has unresolveable links" '
    cat <<-\EOF >bad_object &&
{ "Links": [ { "Name": "foo", "Hash": "QmZzaC6ydNXiR65W8VjGA73ET9MZ6VFAqUT1ngYMXcpihn", "Size": 1897 }, { "Name": "bar", "Hash": "Qmd4mG6pDFDmDTn6p3hX1srP8qTbkyXKj5yjpEsiHDX3u8", "Size": 56 }, { "Name": "baz", "Hash": "QmUTjwRnG28dSrFFVTYgbr6LiDLsBmRr2SaUSTGheK2YqG", "Size": 24266 } ], "Data": "\b\u0001" }
EOF
    cat bad_object | udfs object put > put_out
  '

  test_expect_success "output looks good" '
    echo "added QmaGidyrnX8FMbWJoxp8HVwZ1uRKwCyxBJzABnR1S2FVUr" > put_exp &&
    test_cmp put_exp put_out
  '

  test_expect_success "udfs get fails" '
    test_expect_code 1 udfs get QmaGidyrnX8FMbWJoxp8HVwZ1uRKwCyxBJzABnR1S2FVUr
  '
}

# should work offline
test_get_cmd

# only really works offline, will try and search network when online
test_get_fail

# should work online
test_launch_udfs_daemon
test_get_cmd

test_expect_success "empty request to get doesn't panic and returns error" '
  curl "http://$API_ADDR/api/v0/get" > curl_out || true &&
    grep "argument \"udfs-path\" is required" curl_out
'
test_kill_udfs_daemon

test_done
