#!/usr/bin/env bash
#
# Copyright (c) 2016 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test migrations auto update prompt"

. lib/test-lib.sh

test_init_udfs

test_expect_success "setup mock migrations" '
  mkdir bin &&
  echo "#!/bin/bash" > bin/fs-repo-migrations &&
  echo "echo 5" >> bin/fs-repo-migrations &&
  chmod +x bin/fs-repo-migrations &&
  export PATH="$(pwd)/bin":$PATH
'

test_expect_success "manually reset repo version to 3" '
  echo "3" > "$UDFS_PATH"/version
'

test_expect_success "udfs daemon --migrate=false fails" '
  test_expect_code 1 udfs daemon --migrate=false > false_out
'

test_expect_success "output looks good" '
  grep "Please get fs-repo-migrations from https://dist.udfs.io" false_out
'

test_expect_success "udfs daemon --migrate=true runs migration" '
  test_expect_code 1 udfs daemon --migrate=true > true_out
'

test_expect_failure "output looks good" '
  grep "Running: " true_out > /dev/null &&
  grep "Success: fs-repo has been migrated to version 5." true_out > /dev/null
'

test_expect_success "'udfs daemon' prompts to auto migrate" '
  test_expect_code 1 udfs daemon > daemon_out 2> daemon_err
'

test_expect_success "output looks good" '
  grep "Found outdated fs-repo" daemon_out > /dev/null &&
  grep "Run migrations now?" daemon_out > /dev/null &&
  grep "Please get fs-repo-migrations from https://dist.udfs.io" daemon_out > /dev/null
'

test_done
