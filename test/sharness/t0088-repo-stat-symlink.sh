#!/usr/bin/env bash
#
# Copyright (c) 2017 John Reed
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test 'udfs repo stat' where UDFS_PATH is a symbolic link"

. lib/test-lib.sh

test_expect_success "create symbolic link for UDFS_PATH" '
  mkdir sym_link_target &&
  ln -s sym_link_target .udfs
'

test_init_udfs

# ensure that the RepoSize is reasonable when checked via a symlink.
test_expect_success "'udfs repo stat' RepoSize is correct with sym link" '
  reposize_symlink=$(udfs repo stat | grep RepoSize | awk '\''{ print $2 }'\'') &&
  symlink_size=$(file_size .udfs) &&
  test "${reposize_symlink}" -gt "${symlink_size}"
'

test_done
