#!/usr/bin/env bash
#
# Copyright (c) 2014 Christian Couder
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test add and cat commands"

. lib/test-lib.sh

test_add_cat_file() {
  test_expect_success "udfs add --help works" '
    udfs add --help 2> add_help_err > /dev/null
  '

  test_expect_success "stdin reading message doesnt show up" '
    test_expect_code 1 grep "udfs: Reading from" add_help_err &&
    test_expect_code 1 grep "send Ctrl-d to stop." add_help_err
  '

  test_expect_success "udfs add succeeds" '
    echo "Hello Worlds!" >mountdir/hello.txt &&
    udfs add mountdir/hello.txt >actual
  '

  test_expect_success "udfs add output looks good" '
    HASH="QmVr26fY1tKyspEJBniVhqxQeEjhF78XerGiqWAwraVLQH" &&
    echo "added $HASH hello.txt" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs add --only-hash succeeds" '
    udfs add --only-hash mountdir/hello.txt > oh_actual
  '

  test_expect_success "udfs add --only-hash output looks good" '
    test_cmp expected oh_actual
  '

  test_expect_success "udfs cat succeeds" '
    udfs cat "$HASH" >actual
  '

  test_expect_success "udfs cat output looks good" '
    echo "Hello Worlds!" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat with offset succeeds" '
    udfs cat --offset 10 "$HASH" >actual
  '

  test_expect_success "udfs cat from offset output looks good" '
    echo "ds!" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat multiple hashes with offset succeeds" '
    udfs cat --offset 10 "$HASH" "$HASH" >actual
  '

  test_expect_success "udfs cat from offset output looks good" '
    echo "ds!" >expected &&
    echo "Hello Worlds!" >>expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat multiple hashes with offset succeeds" '
    udfs cat --offset 16 "$HASH" "$HASH" >actual
  '

  test_expect_success "udfs cat from offset output looks good" '
    echo "llo Worlds!" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat from negitive offset should fail" '
    test_expect_code 1 udfs cat --offset -102 "$HASH" > actual
  '

  test_expect_success "udfs cat with length succeeds" '
    udfs cat --length 8 "$HASH" >actual
  '

  test_expect_success "udfs cat with length output looks good" '
    printf "Hello Wo" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat multiple hashes with offset and length succeeds" '
    udfs cat --offset 5 --length 15 "$HASH" "$HASH" "$HASH" >actual
  '

  test_expect_success "udfs cat multiple hashes with offset and length looks good" '
    printf " Worlds!\nHello " >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat with exact length succeeds" '
    udfs cat --length $(udfs cat "$HASH" | wc -c) "$HASH" >actual
  '

  test_expect_success "udfs cat with exact length looks good" '
    echo "Hello Worlds!" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat with 0 length succeeds" '
    udfs cat --length 0 "$HASH" >actual
  '

  test_expect_success "udfs cat with 0 length looks good" '
    : >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat with oversized length succeeds" '
    udfs cat --length 100 "$HASH" >actual
  '

  test_expect_success "udfs cat with oversized length looks good" '
    echo "Hello Worlds!" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat with negitive length should fail" '
    test_expect_code 1 udfs cat --length -102 "$HASH" > actual
  '

  test_expect_success "udfs cat /udfs/file succeeds" '
    udfs cat /udfs/$HASH >actual
  '

  test_expect_success "output looks good" '
    echo "Hello Worlds!" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs add -t succeeds" '
    udfs add -t mountdir/hello.txt >actual
  '

  test_expect_success "udfs add -t output looks good" '
    HASH="QmUkUQgxXeggyaD5Ckv8ZqfW8wHBX6cYyeiyqvVZYzq5Bi" &&
    echo "added $HASH hello.txt" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs add --chunker size-32 succeeds" '
    udfs add --chunker rabin mountdir/hello.txt >actual
  '

  test_expect_success "udfs add --chunker size-32 output looks good" '
    HASH="QmVr26fY1tKyspEJBniVhqxQeEjhF78XerGiqWAwraVLQH" &&
    echo "added $HASH hello.txt" >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs add on hidden file succeeds" '
    echo "Hello Worlds!" >mountdir/.hello.txt &&
    udfs add mountdir/.hello.txt >actual
  '

  test_expect_success "udfs add on hidden file output looks good" '
    HASH="QmVr26fY1tKyspEJBniVhqxQeEjhF78XerGiqWAwraVLQH" &&
    echo "added $HASH .hello.txt" >expected &&
    test_cmp expected actual
  '

  test_expect_success "add zero length file" '
    touch zero-length-file &&
    ZEROHASH=$(udfs add -q zero-length-file) &&
    echo $ZEROHASH
  '

  test_expect_success "zero length file has correct hash" '
    test "$ZEROHASH" = QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH
  '

  test_expect_success "cat zero length file" '
    udfs cat $ZEROHASH > zero-length-file_out
  '

  test_expect_success "make sure it looks good" '
    test_cmp zero-length-file zero-length-file_out
  '
}

test_add_cat_5MB() {
  ADD_FLAGS="$1"
  EXP_HASH="$2"

  test_expect_success "generate 5MB file using go-random" '
    random 5242880 41 >mountdir/bigfile
  '

  test_expect_success "sha1 of the file looks ok" '
    echo "11145620fb92eb5a49c9986b5c6844efda37e471660e" >sha1_expected &&
    multihash -a=sha1 -e=hex mountdir/bigfile >sha1_actual &&
    test_cmp sha1_expected sha1_actual
  '

  test_expect_success "'udfs add $ADD_FLAGS bigfile' succeeds" '
    udfs add $ADD_FLAGS mountdir/bigfile >actual ||
    test_fsh cat daemon_err
  '

  test_expect_success "'udfs add bigfile' output looks good" '
    echo "added $EXP_HASH bigfile" >expected &&
    test_cmp expected actual
  '
  test_expect_success "'udfs cat' succeeds" '
    udfs cat "$EXP_HASH" >actual
  '

  test_expect_success "'udfs cat' output looks good" '
    test_cmp mountdir/bigfile actual
  '

  test_expect_success FUSE "cat udfs/bigfile succeeds" '
    cat "udfs/$EXP_HASH" >actual
  '

  test_expect_success FUSE "cat udfs/bigfile looks good" '
    test_cmp mountdir/bigfile actual
  '
}

test_add_cat_raw() {
  test_expect_success "add a small file with raw-leaves" '
    echo "foobar" > afile &&
    HASH=$(udfs add -q --raw-leaves afile)
  '

  test_expect_success "cat that small file" '
    udfs cat $HASH > afile_out
  '

  test_expect_success "make sure it looks good" '
    test_cmp afile afile_out
  '

  test_expect_success "add zero length file with raw-leaves" '
    touch zero-length-file &&
    ZEROHASH=$(udfs add -q --raw-leaves zero-length-file) &&
    echo $ZEROHASH
  '

  test_expect_success "zero length file has correct hash" '
    test "$ZEROHASH" = zb2rhmy65F3REf8SZp7De11gxtECBGgUKaLdiDj7MCGCHxbDW
  '

  test_expect_success "cat zero length file" '
    udfs cat $ZEROHASH > zero-length-file_out
  '

  test_expect_success "make sure it looks good" '
    test_cmp zero-length-file zero-length-file_out
  '
}

test_add_cat_expensive() {
  ADD_FLAGS="$1"
  HASH="$2"

  test_expect_success EXPENSIVE "generate 100MB file using go-random" '
    random 104857600 42 >mountdir/bigfile
  '

  test_expect_success EXPENSIVE "sha1 of the file looks ok" '
    echo "1114885b197b01e0f7ff584458dc236cb9477d2e736d" >sha1_expected &&
    multihash -a=sha1 -e=hex mountdir/bigfile >sha1_actual &&
    test_cmp sha1_expected sha1_actual
  '

  test_expect_success EXPENSIVE "udfs add $ADD_FLAGS bigfile succeeds" '
    udfs add $ADD_FLAGS mountdir/bigfile >actual
  '

  test_expect_success EXPENSIVE "udfs add bigfile output looks good" '
    echo "added $HASH bigfile" >expected &&
    test_cmp expected actual
  '

  test_expect_success EXPENSIVE "udfs cat succeeds" '
    udfs cat "$HASH" | multihash -a=sha1 -e=hex >sha1_actual
  '

  test_expect_success EXPENSIVE "udfs cat output looks good" '
    udfs cat "$HASH" >actual &&
    test_cmp mountdir/bigfile actual
  '

  test_expect_success EXPENSIVE "udfs cat output hashed looks good" '
    echo "1114885b197b01e0f7ff584458dc236cb9477d2e736d" >sha1_expected &&
    test_cmp sha1_expected sha1_actual
  '

  test_expect_success FUSE,EXPENSIVE "cat udfs/bigfile succeeds" '
    cat "udfs/$HASH" | multihash -a=sha1 -e=hex >sha1_actual
  '

  test_expect_success FUSE,EXPENSIVE "cat udfs/bigfile looks good" '
    test_cmp sha1_expected sha1_actual
  '
}

test_add_named_pipe() {
  err_prefix=$1
  test_expect_success "useful error message when adding a named pipe" '
    mkfifo named-pipe &&
    test_expect_code 1 udfs add named-pipe 2>actual &&
    STAT=$(generic_stat named-pipe) &&
    rm named-pipe &&
    grep "Error: Unrecognized file type for named-pipe: $STAT" actual &&
    grep USAGE actual &&
    grep "udfs add" actual
  '

  test_expect_success "useful error message when recursively adding a named pipe" '
    mkdir -p named-pipe-dir &&
    mkfifo named-pipe-dir/named-pipe &&
    STAT=$(generic_stat named-pipe-dir/named-pipe) &&
    test_expect_code 1 udfs add -r named-pipe-dir 2>actual &&
    printf "Error:$err_prefix Unrecognized file type for named-pipe-dir/named-pipe: $STAT\n" >expected &&
    rm named-pipe-dir/named-pipe &&
    rmdir named-pipe-dir &&
    test_cmp expected actual
  '
}

test_add_pwd_is_symlink() {
  test_expect_success "udfs add -r adds directory content when ./ is symlink" '
    mkdir hellodir &&
    echo "World" > hellodir/world &&
    ln -s hellodir hellolink &&
    ( cd hellolink &&
    udfs add -r . > ../actual ) &&
    grep "added Qma9CyFdG5ffrZCcYSin2uAETygB25cswVwEYYzwfQuhTe" actual &&
    rm -r hellodir
  '
}

test_launch_udfs_daemon_and_mount

test_expect_success "'udfs add --help' succeeds" '
  udfs add --help >actual
'

test_expect_success "'udfs add --help' output looks good" '
  egrep "udfs add.*<path>" actual >/dev/null ||
  test_fsh cat actual
'

test_expect_success "'udfs cat --help' succeeds" '
  udfs cat --help >actual
'

test_expect_success "'udfs cat --help' output looks good" '
  egrep "udfs cat.*<udfs-path>" actual >/dev/null ||
  test_fsh cat actual
'

test_add_cat_file

test_expect_success "udfs cat succeeds with stdin opened (issue #1141)" '
  cat mountdir/hello.txt | while read line; do udfs cat "$HASH" >actual || exit; done
'

test_expect_success "udfs cat output looks good" '
  cat mountdir/hello.txt >expected &&
  test_cmp expected actual
'

test_expect_success "udfs cat accept hash from built input" '
  echo "$HASH" | udfs cat >actual
'

test_expect_success "udfs cat output looks good" '
  test_cmp expected actual
'

test_expect_success FUSE "cat udfs/stuff succeeds" '
  cat "udfs/$HASH" >actual
'

test_expect_success FUSE "cat udfs/stuff looks good" '
  test_cmp expected actual
'

test_expect_success "'udfs add -q' succeeds" '
  echo "Hello Venus!" >mountdir/venus.txt &&
  udfs add -q mountdir/venus.txt >actual
'

test_expect_success "'udfs add -q' output looks good" '
  HASH="QmU5kp3BH3B8tnWUU2Pikdb2maksBNkb92FHRr56hyghh4" &&
  echo "$HASH" >expected &&
  test_cmp expected actual
'

test_expect_success "'udfs add -q' with stdin input succeeds" '
  echo "Hello Jupiter!" | udfs add -q >actual
'

test_expect_success "'udfs add -q' output looks good" '
  HASH="QmUnvPcBctVTAcJpigv6KMqDvmDewksPWrNVoy1E1WP5fh" &&
  echo "$HASH" >expected &&
  test_cmp expected actual
'

test_expect_success "'udfs cat' succeeds" '
  udfs cat "$HASH" >actual
'

test_expect_success "udfs cat output looks good" '
  echo "Hello Jupiter!" >expected &&
  test_cmp expected actual
'

test_expect_success "'udfs add' with stdin input succeeds" '
  printf "Hello Neptune!\nHello Pluton!" | udfs add >actual
'

test_expect_success "'udfs add' output looks good" '
  HASH="QmZDhWpi8NvKrekaYYhxKCdNVGWsFFe1CREnAjP1QbPaB3" &&
  echo "added $HASH $HASH" >expected &&
  test_cmp expected actual
'

test_expect_success "'udfs cat' with built input succeeds" '
  echo "$HASH" | udfs cat >actual
'

test_expect_success "udfs cat with built input output looks good" '
  printf "Hello Neptune!\nHello Pluton!" >expected &&
  test_cmp expected actual
'

add_directory() {
  EXTRA_ARGS=$1

  test_expect_success "'udfs add -r $EXTRA_ARGS' succeeds" '
    mkdir mountdir/planets &&
    echo "Hello Mars!" >mountdir/planets/mars.txt &&
    echo "Hello Venus!" >mountdir/planets/venus.txt &&
    udfs add -r $EXTRA_ARGS mountdir/planets >actual
  '

  test_expect_success "'udfs add -r $EXTRA_ARGS' output looks good" '
    echo "added $MARS planets/mars.txt" >expected &&
    echo "added $VENUS planets/venus.txt" >>expected &&
    echo "added $PLANETS planets" >>expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat accept many hashes from built input" '
    { echo "$MARS"; echo "$VENUS"; } | udfs cat >actual
  '

  test_expect_success "udfs cat output looks good" '
    cat mountdir/planets/mars.txt mountdir/planets/venus.txt >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat accept many hashes as args" '
    udfs cat "$MARS" "$VENUS" >actual
  '

  test_expect_success "udfs cat output looks good" '
    test_cmp expected actual
  '

  test_expect_success "udfs cat with both arg and stdin" '
    echo "$MARS" | udfs cat "$VENUS" >actual
  '

  test_expect_success "udfs cat output looks good" '
    cat mountdir/planets/venus.txt >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs cat with two args and stdin" '
    echo "$MARS" | udfs cat "$VENUS" "$VENUS" >actual
  '

  test_expect_success "udfs cat output looks good" '
    cat mountdir/planets/venus.txt mountdir/planets/venus.txt >expected &&
    test_cmp expected actual
  '

  test_expect_success "udfs add --quieter succeeds" '
    udfs add -r -Q $EXTRA_ARGS mountdir/planets >actual
  '

  test_expect_success "udfs add --quieter returns only one correct hash" '
    echo "$PLANETS" > expected &&
    test_cmp expected actual
  '

  test_expect_success "cleanup" '
    rm -r mountdir/planets
  '
}

PLANETS="QmWSgS32xQEcXMeqd3YPJLrNBLSdsfYCep2U7CFkyrjXwY"
MARS="QmPrrHqJzto9m7SyiRzarwkqPcCSsKR2EB1AyqJfe8L8tN"
VENUS="QmU5kp3BH3B8tnWUU2Pikdb2maksBNkb92FHRr56hyghh4"
add_directory

PLANETS="QmfWfQfKCY5Ukv9peBbxM5vqWM9BzmqUSXvdCgjT2wsiBT"
MARS="zb2rhZdTkQNawVajsTNiYc9cTPHqgLdJVvBRkZok9RjkgQYRU"
VENUS="zb2rhn6TGvnUaMAg4VV4y9HVx5W42HihcH4jsyrDv8mkepFqq"
add_directory '--raw-leaves'

PLANETS="zdj7Wnbun6P41Z5ddTkNvZaDTmQ8ZLdiKFcJrL9sV87rPScMP"
MARS="zb2rhZdTkQNawVajsTNiYc9cTPHqgLdJVvBRkZok9RjkgQYRU"
VENUS="zb2rhn6TGvnUaMAg4VV4y9HVx5W42HihcH4jsyrDv8mkepFqq"
add_directory '--cid-version=1'

PLANETS="zdj7WiC51v78BjBcmZR7uuBvmDWxSn5EDr5MiyTwE18e8qvb7"
MARS="zdj7WWx6fGNrNGkdpkuTAxCjKbQ1pPtarqA6VQhedhLTZu34J"
VENUS="zdj7WbB1BUF8WejmVpQCmMLd1RbPnxJtvAj1Lep6eTmXRFbrz"
add_directory '--cid-version=1 --raw-leaves=false'

PLANETS="zDMZof1kqxDAx9myQbXsyWwyWP9qRPsXsWH7XuTz6isT7Rh1S6nM"
MARS="zCT5htkdz1ZBHYVQXFQn51ngPXLVqaHSWoae87V1d6e9qWpSAjXw"
VENUS="zCT5htke5JcdoMM4WhmUKXWf2QC3TnQToqGZHH1WsZERv6kPhFPg"
add_directory '--hash=blake2b-256'

test_expect_success "'udfs add -rn' succeeds" '
  mkdir -p mountdir/moons/jupiter &&
  mkdir -p mountdir/moons/saturn &&
  echo "Hello Europa!" >mountdir/moons/jupiter/europa.txt &&
  echo "Hello Titan!" >mountdir/moons/saturn/titan.txt &&
  echo "hey youre no moon!" >mountdir/moons/mercury.txt &&
  udfs add -rn mountdir/moons >actual
'

test_expect_success "'udfs add -rn' output looks good" '
  MOONS="QmVKvomp91nMih5j6hYBA8KjbiaYvEetU2Q7KvtZkLe9nQ" &&
  EUROPA="Qmbjg7zWdqdMaK2BucPncJQDxiALExph5k3NkQv5RHpccu" &&
  JUPITER="QmS5mZddhFPLWFX3w6FzAy9QxyYkaxvUpsWCtZ3r7jub9J" &&
  SATURN="QmaMagZT4rTE7Nonw8KGSK4oe1bh533yhZrCo1HihSG8FK" &&
  TITAN="QmZzppb9WHn552rmRqpPfgU5FEiHH6gDwi3MrB9cTdPwdb" &&
  MERCURY="QmUJjVtnN8YEeYcS8VmUeWffTWhnMQAkk5DzZdKnPhqUdK" &&
  echo "added $EUROPA moons/jupiter/europa.txt" >expected &&
  echo "added $MERCURY moons/mercury.txt" >>expected &&
  echo "added $TITAN moons/saturn/titan.txt" >>expected &&
  echo "added $JUPITER moons/jupiter" >>expected &&
  echo "added $SATURN moons/saturn" >>expected &&
  echo "added $MOONS moons" >>expected &&
  test_cmp expected actual
'

test_expect_success "go-random is installed" '
  type random
'

test_add_cat_5MB "" "QmSr7FqYkxYWGoSfy8ZiaMWQ5vosb18DQGCzjwEQnVHkTb"

test_add_cat_5MB --raw-leaves "QmbdLHCmdi48eM8T7D67oXjA1S2Puo8eMfngdHhdPukFd6"

# note: the specified hash implies that internal nodes are stored
# using CidV1 and leaves are stored using raw blocks
test_add_cat_5MB --cid-version=1 "zdj7WiiaedqVBXjX4SNqx3jfuZideDqdLYnDzCDJ66JDMK9o2"

# note: the specified hash implies that internal nodes are stored
# using CidV1 and leaves are stored using CidV1 but using the legacy
# format (i.e. not raw)
test_add_cat_5MB '--cid-version=1 --raw-leaves=false' "zdj7WfgEsj897BBZj2mcfsRLhaPZcCixPV2G7DkWgF1Wdr64P"

# note: --hash=blake2b-256 implies --cid-version=1 which implies --raw-leaves=true
# the specified hash represents the leaf nodes stored as raw leaves and
# encoded with the blake2b-256 hash funtion
test_add_cat_5MB '--hash=blake2b-256' "zDMZof1kuxn7ebvKyvmkYLPvocSvFYxxAWT1yQBN1wWiXXr7w5mY"

# the specified hash represents the leaf nodes stored as protoful nodes and
# encoded with the blake2b-256 hash funtion
test_add_cat_5MB '--hash=blake2b-256 --raw-leaves=false' "zDMZof1krz3SFTyhboRyWZyUP2qNgVdn9wjtaX211aHJ8WgeyT9v"

test_add_cat_expensive "" "QmU9SWAPPmNEKZB8umYMmjYvN7VyHqABNvdA6GUi4MMEz3"

# note: the specified hash implies that internal nodes are stored
# using CidV1 and leaves are stored using raw blocks
test_add_cat_expensive "--cid-version=1" "zdj7WcatQrtuE4WMkS4XsfsMixuQN2po4irkYhqxeJyW1wgCq"

# note: --hash=blake2b-256 implies --cid-version=1 which implies --raw-leaves=true
# the specified hash represents the leaf nodes stored as raw leaves and
# encoded with the blake2b-256 hash funtion
test_add_cat_expensive '--hash=blake2b-256' "zDMZof1kwndounDzQCANUHjiE3zt1mPEgx7RE3JTHoZrRRa79xcv"

test_add_named_pipe " Post http://$API_ADDR/api/v0/add?chunker=size-262144&encoding=json&hash=sha2-256&pin=true&progress=true&recursive=true&stream-channels=true:"

test_add_pwd_is_symlink

test_add_cat_raw

test_expect_success "udfs add --cid-version=9 fails" '
  echo "context" > afile.txt &&
  test_must_fail udfs add --cid-version=9 afile.txt 2>&1 | tee add_out &&
  grep -q "unknown CID version" add_out
'

test_kill_udfs_daemon

# should work offline

test_add_cat_file

test_add_cat_raw

test_expect_success "udfs add --only-hash succeeds" '
  echo "unknown content for only-hash" | udfs add --only-hash -q > oh_hash
'

#TODO: this doesn't work when online hence separated out from test_add_cat_file
test_expect_success "udfs cat file fails" '
  test_must_fail udfs cat $(cat oh_hash)
'

test_add_named_pipe ""

test_add_pwd_is_symlink

# Test daemon in offline mode
test_launch_udfs_daemon --offline

test_add_cat_file

test_kill_udfs_daemon

test_done
