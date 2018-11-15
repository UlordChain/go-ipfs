#!/usr/bin/env bash
#
# Copyright (c) 2015 Brian Holder-Chow Lin On
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test user-provided config values"

. lib/test-lib.sh

test_init_udfs

test_expect_success "bootstrap doesn't overwrite user-provided config keys (top-level)" '
  udfs config Foo.Bar baz &&
  udfs bootstrap rm --all &&
  echo "baz" >expected &&
  udfs config Foo.Bar >actual &&
  test_cmp expected actual
'

test_done
