#!/usr/bin/env bash
#
# Copyright (c) 2016 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="Test udfs pinning operations"

. lib/test-lib.sh


test_pins() {
  EXTRA_ARGS=$1

  test_expect_success "create some hashes" '
    HASH_A=$(echo "A" | udfs add -q --pin=false) &&
    HASH_B=$(echo "B" | udfs add -q --pin=false) &&
    HASH_C=$(echo "C" | udfs add -q --pin=false) &&
    HASH_D=$(echo "D" | udfs add -q --pin=false) &&
    HASH_E=$(echo "E" | udfs add -q --pin=false) &&
    HASH_F=$(echo "F" | udfs add -q --pin=false) &&
    HASH_G=$(echo "G" | udfs add -q --pin=false)
  '

  test_expect_success "put all those hashes in a file" '
    echo $HASH_A > hashes &&
    echo $HASH_B >> hashes &&
    echo $HASH_C >> hashes &&
    echo $HASH_D >> hashes &&
    echo $HASH_E >> hashes &&
    echo $HASH_F >> hashes &&
    echo $HASH_G >> hashes
  '

  test_expect_success "'udfs pin add $EXTRA_ARGS' via stdin" '
    cat hashes | udfs pin add $EXTRA_ARGS
  '

  test_expect_success "see if verify works" '
    udfs pin verify
  '

  test_expect_success "see if verify --verbose works" '
    udfs pin verify --verbose > verify_out &&
    test $(cat verify_out | wc -l) > 8
  '

  test_expect_success "unpin those hashes" '
    cat hashes | udfs pin rm
  '
}

RANDOM_HASH=Qme8uX5n9hn15pw9p6WcVKoziyyC9LXv4LEgvsmKMULjnV

test_pins_error_reporting() {
  EXTRA_ARGS=$1

  test_expect_success "'udfs pin add $EXTRA_ARGS' on non-existent hash should fail" '
    test_must_fail udfs pin add $EXTRA_ARGS $RANDOM_HASH 2> err &&
    grep -q "not found" err
  '
}

test_pin_dag_init() {
  EXTRA_ARGS=$1

  test_expect_success "'udfs add $EXTRA_ARGS --pin=false' 1MB file" '
    random 1048576 56 > afile &&
    HASH=`udfs add $EXTRA_ARGS --pin=false -q afile`
  '
}

test_pin_dag() {
  test_pin_dag_init $1

  test_expect_success "'udfs pin add --progress' file" '
    udfs pin add --recursive=true $HASH
  '

  test_expect_success "'udfs pin rm' file" '
    udfs pin rm $HASH
  '

  test_expect_success "remove part of the dag" '
    PART=`udfs refs $HASH | head -1` &&
    udfs block rm $PART
  '

  test_expect_success "pin file, should fail" '
    test_must_fail udfs pin add --recursive=true $HASH 2> err &&
    cat err &&
    grep -q "not found" err
  '
}

test_pin_progress() {
  test_pin_dag_init

  test_expect_success "'udfs pin add --progress' file" '
    udfs pin add --progress $HASH 2> err
  '

  test_expect_success "pin progress reported correctly" '
    cat err
    grep -q " 5 nodes" err
  '
}

test_init_udfs

test_pins
test_pins --progress

test_pins_error_reporting
test_pins_error_reporting --progress

test_pin_dag
test_pin_dag --raw-leaves

test_pin_progress

test_launch_udfs_daemon --offline

test_pins
test_pins --progress

test_pins_error_reporting
test_pins_error_reporting --progress

test_pin_dag
test_pin_dag --raw-leaves

test_pin_progress

test_kill_udfs_daemon

test_done
