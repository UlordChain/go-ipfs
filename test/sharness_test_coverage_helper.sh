#!/bin/sh

USAGE="$0 [-h] [-v]"

usage() {
    echo "$USAGE"
    echo "	Print sharness test coverage"
    echo "	Options:"
    echo "		-h|--help: print this usage message and exit"
    echo "		-v|--verbose: print logs of what happens"
    exit 0
}

log() {
    test -z "$VERBOSE" || echo "->" "$@"
}

die() {
    printf >&2 "fatal: %s\n" "$@"
    exit 1
}

# get user options
while [ "$#" -gt "0" ]; do
    # get options
    arg="$1"
    shift

    case "$arg" in
	-h|--help)
	    usage ;;
	-v|--verbose)
	    VERBOSE=1 ;;
	-*)
	    die "unrecognised option: '$arg'\n$USAGE" ;;
	*)
	    die "too many arguments\n$USAGE" ;;
    esac
done

log "Create temporary directory"
DATE=$(date +"%Y-%m-%dT%H:%M:%SZ")
TMPDIR=$(mktemp -d "/tmp/coverage_helper.$DATE.XXXXXX") ||
die "could not 'mktemp -d /tmp/coverage_helper.$DATE.XXXXXX'"

log "Grep the sharness tests for udfs commands"
CMD_RAW="$TMPDIR/udfs_cmd_raw.txt"
git grep -n -E '\Wudfs\W' -- sharness/t*-*.sh >"$CMD_RAW" ||
die "Could not grep udfs in the sharness tests"

grep_out() {
    pattern="$1"
    src="$TMPDIR/udfs_cmd_${2}.txt"
    dst="$TMPDIR/udfs_cmd_${3}.txt"
    desc="$4"

    log "Remove $desc"
    egrep -v "$pattern" "$src" >"$dst" || die "Could not remove $desc"
}

grep_out 'test_expect_.*udfs' raw expect "test_expect_{success,failure} lines"
grep_out '^[^:]+:[^:]+:\s*#' expect comment "comments"
grep_out 'test_description=' comment desc "test_description lines"
grep_out '^[^:]+:[^:]+:\s*\w+="[^"]*"\s*(\&\&)?\s*$' desc def "variable definition lines"
grep_out '^[^:]+:[^:]+:\s*e?grep\W[^|]*\Wudfs' def grep "grep lines"
grep_out '^[^:]+:[^:]+:\s*cat\W[^|]*\Wudfs' grep cat "cat lines"
grep_out '^[^:]+:[^:]+:\s*rmdir\W[^|]*\Wudfs' cat rmdir "rmdir lines"
grep_out '^[^:]+:[^:]+:\s*echo\W[^|]*\Wudfs' cat echo "echo lines"

grep_in() {
    pattern="$1"
    src="$TMPDIR/udfs_cmd_${2}.txt"
    dst="$TMPDIR/udfs_cmd_${3}.txt"
    desc="$4"

    log "Keep $desc"
    egrep "$pattern" "$src" >"$dst"
}

grep_in '\Wudfs\W.*/udfs/' echo slash_in1 "udfs.*/udfs/"
grep_in '/udfs/.*\Wudfs\W' echo slash_in2 "/udfs/.*udfs"

grep_out '/udfs/' echo slash "/udfs/"

grep_in '\Wudfs\W.*\.udfs' slash dot_in1 "udfs.*\.udfs"
grep_in '\.udfs.*\Wudfs\W' slash dot_in2 "\.udfs.*udfs"

grep_out '\.udfs' slash dot ".udfs"

log "Print result"
CMD_RES="$TMPDIR/udfs_cmd_result.txt"
for f in dot slash_in1 slash_in2 dot_in1 dot_in2
do
    fname="$TMPDIR/udfs_cmd_${f}.txt"
    cat "$fname" || die "Could not cat '$fname'"
done | sort | uniq >"$CMD_RES" || die "Could not write '$CMD_RES'"

log "Get all the udfs commands from 'udfs commands'"
CMD_CMDS="$TMPDIR/commands.txt"
udfs commands --flags >"$CMD_CMDS" || die "'udfs commands' failed"

# Portable function to reverse lines in a file
reverse() {
    if type tac >/dev/null
    then
	tac "$@"
    else
	tail -r "$@"
    fi
}

log "Match the test line commands with the commands they use"
GLOBAL_REV="$TMPDIR/global_results_reversed.txt"

process_command() {
    udfs="$1"
    cmd="$2"
    sub1="$3"
    sub2="$4"
    sub3="$5"

    if test -n "$cmd"
    then
	CMD_OUT="$TMPDIR/res_${udfs}_${cmd}"
	PATTERN="$udfs(\W.*)*\W$cmd"
	NAME="$udfs $cmd"

	if test -n "$sub1"
	then
	    CMD_OUT="${CMD_OUT}_${sub1}"
	    PATTERN="$PATTERN(\W.*)*\W$sub1"
	    NAME="$NAME $sub1"

	    if test -n "$sub2"
	    then
		CMD_OUT="${CMD_OUT}_${sub2}"
		PATTERN="$PATTERN(\W.*)*\W$sub2"
		NAME="$NAME $sub2"

		if test -n "$sub3"
		then
		    CMD_OUT="${CMD_OUT}_${sub3}"
		    PATTERN="$PATTERN(\W.*)*\W$sub3"
		    NAME="$NAME $sub3"
		fi
	    fi
	fi

	egrep "$PATTERN" "$CMD_RES" >"$CMD_OUT.txt"
	reverse "$CMD_OUT.txt" | sed -e 's/^sharness\///' | cut -d- -f1 | uniq -c >>"$GLOBAL_REV"
    fi
}

reverse "$CMD_CMDS" | while read -r line
do
    LONG_CMD=$(echo "$line" | cut -d/ -f1)
    SHORT_CMD=$(expr "$line" : "[^/]*/*\(.*\)")

    log "Processing $LONG_CMD"
    process_command $LONG_CMD
    LONG_NAME="$NAME"

    log "Processing $SHORT_CMD"
    process_command $SHORT_CMD
    SHORT_NAME="$NAME"

    test -n "$SHORT_CMD" && echo "$SHORT_NAME" >>"$GLOBAL_REV"
    test "$LONG_CMD" != "udfs" && echo "$LONG_NAME" >>"$GLOBAL_REV"
    echo >>"$GLOBAL_REV"
done

# The following will allow us to check that
# we are properly excuding enough stuff using:
# diff -u udfs_cmd_result.txt cmd_found.txt
log "Get all the line commands that matched"
CMD_FOUND="$TMPDIR/cmd_found.txt"
cat $TMPDIR/res_*.txt | sort -n | uniq >"$CMD_FOUND"

log "Print results"
reverse "$GLOBAL_REV"

# Remove temp directory...
