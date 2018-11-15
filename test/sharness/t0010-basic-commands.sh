#!/usr/bin/env bash
#
# Copyright (c) 2014 Christian Couder
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test installation and some basic commands"

. lib/test-lib.sh

test_expect_success "current dir is writable" '
  echo "It works!" >test.txt
'

test_expect_success "udfs version succeeds" '
  udfs version >version.txt
'

test_expect_success "udfs --version success" '
  udfs --version
'

test_expect_success "udfs version output looks good" '
  egrep "^udfs version [0-9]+\.[0-9]+\.[0-9]" version.txt >/dev/null ||
  test_fsh cat version.txt
'

test_expect_success "udfs versions matches udfs --version" '
  udfs version > version.txt &&
  udfs --version > version2.txt &&
  diff version2.txt version.txt ||
  test_fsh udfs --version

'

test_expect_success "udfs version --all has all required fields" '
  udfs version --all > version_all.txt &&
  grep "go-udfs version" version_all.txt &&
  grep "Repo version" version_all.txt &&
  grep "System version" version_all.txt &&
  grep "Golang version" version_all.txt
'

test_expect_success "udfs help succeeds" '
  udfs help >help.txt
'

test_expect_success "udfs help output looks good" '
  egrep -i "^Usage" help.txt >/dev/null &&
  egrep "udfs <command>" help.txt >/dev/null ||
  test_fsh cat help.txt
'

test_expect_success "'udfs commands' succeeds" '
  udfs commands >commands.txt
'

test_expect_success "'udfs commands' output looks good" '
  grep "udfs add" commands.txt &&
  grep "udfs daemon" commands.txt &&
  grep "udfs update" commands.txt
'

test_expect_success "All commands accept --help" '
  echo 0 > fail
  while read -r cmd
  do
    $cmd --help >/dev/null ||
      { echo "$cmd doesnt accept --help"; echo 1 > fail; }
    echo stuff | $cmd --help >/dev/null ||
      { echo "$cmd doesnt accept --help when using stdin"; echo 1 > fail; }
  done <commands.txt

  if [ $(cat fail) = 1 ]; then
    return 1
  fi
'

test_expect_failure "All udfs root commands are mentioned in base helptext" '
  echo 0 > fail
  cut -d" " -f 2 commands.txt | grep -v udfs | sort -u | \
  while read cmd
  do
    grep "  $cmd" help.txt > /dev/null ||
      { echo "missing $cmd from helptext"; echo 1 > fail; }
  done

  if [ $(cat fail) = 1 ]; then
    return 1
  fi
'

test_expect_failure "All udfs commands docs are 80 columns or less" '
  echo 0 > fail
  while read cmd
  do
    LENGTH="$($cmd --help | awk "{ print length }" | sort -nr | head -1)"
    [ $LENGTH -gt 80 ] &&
      { echo "$cmd help text is longer than 79 chars ($LENGTH)"; echo 1 > fail; }
  done <commands.txt

  if [ $(cat fail) = 1 ]; then
    return 1
  fi
'

test_expect_success "All udfs commands fail when passed a bad flag" '
  echo 0 > fail
  while read -r cmd
  do
    test_must_fail $cmd --badflag >/dev/null ||
      { echo "$cmd exit with code 0 when passed --badflag"; echo 1 > fail; }
  done <commands.txt

  if [ $(cat fail) = 1 ]; then
    return 1
  fi
'

test_expect_success "'udfs commands --flags' succeeds" '
  udfs commands --flags >commands.txt
'

test_expect_success "'udfs commands --flags' output looks good" '
  grep "udfs pin add --recursive / udfs pin add -r" commands.txt &&
  grep "udfs id --format / udfs id -f" commands.txt &&
  grep "udfs repo gc --quiet / udfs repo gc -q" commands.txt
'



test_done
