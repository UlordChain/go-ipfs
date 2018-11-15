#!/usr/bin/env bash
#
# Copyright (c) 2014 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test udfs repo operations"

. lib/test-lib.sh

test_init_udfs
test_launch_udfs_daemon --offline

test_expect_success "'udfs repo gc' succeeds" '
  udfs repo gc >gc_out_actual
'

test_expect_success "'udfs add afile' succeeds" '
  echo "some text" >afile &&
  HASH=`udfs add -q afile`
'

test_expect_success "added file was pinned" '
  udfs pin ls --type=recursive >actual &&
  grep "$HASH" actual
'

test_expect_success "'udfs repo gc' succeeds" '
  udfs repo gc >gc_out_actual
'

test_expect_success "'udfs repo gc' looks good (patch root)" '
  PATCH_ROOT=QmQXirSbubiySKnqaFyfs5YzziXRB5JEVQVjU6xsd7innr &&
  grep "removed $PATCH_ROOT" gc_out_actual
'

test_expect_success "'udfs repo gc' doesnt remove file" '
  udfs cat "$HASH" >out &&
  test_cmp out afile
'

test_expect_success "'udfs pin rm' succeeds" '
  udfs pin rm -r "$HASH" >actual1
'

test_expect_success "'udfs pin rm' output looks good" '
  echo "unpinned $HASH" >expected1 &&
  test_cmp expected1 actual1
'

test_expect_failure "udfs repo gc fully reverse udfs add" '
  udfs repo gc &&
  random 100000 41 >gcfile &&
  disk_usage "$UDFS_PATH/blocks" >expected &&
  hash=`udfs add -q gcfile` &&
  udfs pin rm -r $hash &&
  udfs repo gc &&
  disk_usage "$UDFS_PATH/blocks" >actual &&
  test_cmp expected actual
'

test_expect_success "file no longer pinned" '
  udfs pin ls --type=recursive --quiet >actual2 &&
  test_expect_code 1 grep $HASH actual2
'

test_expect_success "recursively pin afile(default action)" '
  HASH=`udfs add -q afile` &&
  udfs pin add "$HASH"
'

test_expect_success "recursively pin rm afile (default action)" '
  udfs pin rm "$HASH"
'

test_expect_success "recursively pin afile" '
  udfs pin add -r "$HASH"
'

test_expect_success "pinning directly should fail now" '
  echo "Error: pin: $HASH already pinned recursively" >expected3 &&
  test_must_fail udfs pin add -r=false "$HASH" 2>actual3 &&
  test_cmp expected3 actual3
'

test_expect_success "'udfs pin rm -r=false <hash>' should fail" '
  echo "Error: $HASH is pinned recursively" >expected4 &&
  test_must_fail udfs pin rm -r=false "$HASH" 2>actual4 &&
  test_cmp expected4 actual4
'

test_expect_success "remove recursive pin, add direct" '
  echo "unpinned $HASH" >expected5 &&
  udfs pin rm -r "$HASH" >actual5 &&
  test_cmp expected5 actual5 &&
  udfs pin add -r=false "$HASH"
'

test_expect_success "remove direct pin" '
  echo "unpinned $HASH" >expected6 &&
  udfs pin rm "$HASH" >actual6 &&
  test_cmp expected6 actual6
'

test_expect_success "'udfs repo gc' removes file" '
  udfs repo gc >actual7 &&
  grep "removed $HASH" actual7 &&
  grep "removed $PATCH_ROOT" actual7
'

test_expect_success "'udfs refs local' no longer shows file" '
  EMPTY_DIR=QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn &&
  udfs refs local >actual8 &&
  grep "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y" actual8 &&
  grep "$EMPTY_DIR" actual8 &&
  grep "$HASH_WELCOME_DOCS" actual8 &&
  test_must_fail grep "$HASH" actual8 &&
  test_must_fail grep "$PATCH_ROOT" actual8
'

test_expect_success "adding multiblock random file succeeds" '
  random 1000000 >multiblock &&
  MBLOCKHASH=`udfs add -q multiblock`
'

test_expect_success "'udfs pin ls --type=indirect' is correct" '
  udfs refs "$MBLOCKHASH" >refsout &&
  udfs refs -r "$HASH_WELCOME_DOCS" >>refsout &&
  sed -i"~" "s/\(.*\)/\1 indirect/g" refsout &&
  udfs pin ls --type=indirect >indirectpins &&
  test_sort_cmp refsout indirectpins
'

test_expect_success "pin something directly" '
  echo "udfs is so awesome" >awesome &&
  DIRECTPIN=`udfs add -q awesome` &&
  echo "unpinned $DIRECTPIN" >expected9 &&
  udfs pin rm -r "$DIRECTPIN" >actual9 &&
  test_cmp expected9 actual9  &&

  echo "pinned $DIRECTPIN directly" >expected10 &&
  udfs pin add -r=false "$DIRECTPIN" >actual10 &&
  test_cmp expected10 actual10
'

test_expect_success "'udfs pin ls --type=direct' is correct" '
  echo "$DIRECTPIN direct" >directpinexpected &&
  udfs pin ls --type=direct >directpinout &&
  test_sort_cmp directpinexpected directpinout
'

test_expect_success "'udfs pin ls --type=recursive' is correct" '
  echo "$MBLOCKHASH" >rp_expected &&
  echo "$HASH_WELCOME_DOCS" >>rp_expected &&
  echo "$EMPTY_DIR" >>rp_expected &&
  sed -i"~" "s/\(.*\)/\1 recursive/g" rp_expected &&
  udfs pin ls --type=recursive >rp_actual &&
  test_sort_cmp rp_expected rp_actual
'

test_expect_success "'udfs pin ls --type=all --quiet' is correct" '
  cat directpinout >allpins &&
  cat rp_actual >>allpins &&
  cat indirectpins >>allpins &&
  cut -f1 -d " " allpins | sort | uniq >> allpins_uniq_hashes &&
  udfs pin ls --type=all --quiet >actual_allpins &&
  test_sort_cmp allpins_uniq_hashes actual_allpins
'

test_expect_success "'udfs refs --unique' is correct" '
  mkdir -p uniques &&
  echo "content1" > uniques/file1 &&
  echo "content1" > uniques/file2 &&
  udfs add -r -q uniques > add_output &&
  ROOT=$(tail -n1 add_output) &&
  udfs refs --unique $ROOT >expected &&
  udfs add -q uniques/file1 >unique_hash &&
  test_cmp expected unique_hash || test_fsh cat add_output
'

test_expect_success "'udfs refs --unique --recursive' is correct" '
  mkdir -p a/b/c &&
  echo "c1" > a/f1 &&
  echo "c1" > a/b/f1 &&
  echo "c1" > a/b/c/f1 &&
  echo "c2" > a/b/c/f2 &&
  udfs add -r -q a >add_output &&
  ROOT=$(tail -n1 add_output) &&
  udfs refs --unique --recursive $ROOT >refs_output &&
  wc -l refs_output | sed "s/^ *//g" >line_count &&
  echo "4 refs_output" >expected &&
  test_cmp expected line_count || test_fsh cat add_output || test_fsh cat refs_output
'

test_expect_success "'udfs refs --recursive (bigger)'" '
  mkdir -p b/c/d/e &&
  echo "content1" >b/f &&
  echo "content1" >b/c/f1 &&
  echo "content1" >b/c/d/f2 &&
  echo "content2" >b/c/f2 &&
  echo "content2" >b/c/d/f1 &&
  echo "content2" >b/c/d/e/f &&
  cp -r b b2 && mv b2 b/b2 &&
  cp -r b b3 && mv b3 b/b3 &&
  cp -r b b4 && mv b4 b/b4 &&
  udfs add -r -q b >add_output &&
  hash=$(tail -n1 add_output) &&
  udfs refs -r "$hash" >refs_output &&
  wc -l refs_output | sed "s/^ *//g" >actual &&
  echo "79 refs_output" >expected &&
  test_cmp expected actual || test_fsh cat add_output || test_fsh cat refs_output
'

test_expect_success "'udfs refs --unique --recursive (bigger)'" '
  udfs refs -r "$hash" >refs_output &&
  sort refs_output | uniq >expected &&
  udfs refs -r -u "$hash" >actual &&
  test_sort_cmp expected actual || test_fsh cat refs_output
'

get_field_num() {
  field=$1
  file=$2
  num=$(grep "$field" "$file" | awk '{ print $2 }')
  echo $num
}

test_expect_success "'udfs repo stat' succeeds" '
  udfs repo stat > repo-stats
'
test_expect_success "repo stats came out correct" '
  grep "RepoPath" repo-stats &&
  grep "RepoSize" repo-stats &&
  grep "NumObjects" repo-stats &&
  grep "Version" repo-stats &&
  grep "StorageMax" repo-stats
'

test_expect_success "'udfs repo stat' after adding a file" '
  udfs add repo-stats &&
  udfs repo stat > repo-stats-2
'

test_expect_success "repo stats are updated correctly" '
  test $(get_field_num "RepoSize" repo-stats-2) -ge $(get_field_num "RepoSize" repo-stats)
'

test_expect_success "'udfs repo stat --size-only' succeeds" '
  udfs repo stat --size-only > repo-stats-size-only
'

test_expect_success "repo stats came out correct for --size-only" '
  grep "RepoSize" repo-stats-size-only &&
  grep "StorageMax" repo-stats-size-only &&
  grep -v "RepoPath" repo-stats-size-only &&
  grep -v "NumObjects" repo-stats-size-only &&
  grep -v "Version" repo-stats-size-only
'

test_expect_success "'udfs repo version' succeeds" '
  udfs repo version > repo-version
'

test_expect_success "repo version came out correct" '
  egrep "^udfs repo version fs-repo@[0-9]+" repo-version >/dev/null
'

test_expect_success "'udfs repo version -q' succeeds" '
  udfs repo version -q > repo-version-q
'
test_expect_success "repo version came out correct" '
  egrep "^fs-repo@[0-9]+" repo-version-q >/dev/null
'

test_kill_udfs_daemon

test_expect_success "remove Datastore.StorageMax from config" '
  udfs config Datastore.StorageMax ""
'
test_expect_success "'udfs repo stat' still succeeds" '
  udfs repo stat > repo-stats
'

test_done
