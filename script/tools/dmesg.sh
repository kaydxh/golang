#!/usr/bin/env bash
#
# exit by command return non-zero exit code
set -o errexit
# Indicate an error when it encounters an undefined variable
set -o nounset
# Fail on any error.
set -o pipefail
# set -o xtrace
#
# bash dmesg.sh "abc" 10

grep_pattern=${1:-''}
grep_n=${2:-0}
tail_n=${3:-$((${grep_n} + 1))}

dmesg -T | grep -C ${grep_n} "${grep_pattern}" | tail -n ${tail_n}
