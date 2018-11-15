#!/usr/bin/env bash

test_description="Test non-standard datastores"

. lib/test-lib.sh

test_expect_success "'udfs init --profile=badgerds' succeeds" '
  BITS="1024" &&
  udfs init --bits="$BITS" --profile=badgerds
'

test_expect_success "'udfs pin ls' works" '
  udfs pin ls | wc -l | grep 9
'

test_done
