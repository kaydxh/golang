#!/usr/bin/env bash

# Fail on any error.
set -euo pipefail
# set -o xtrace

pushd "$(dirname $0)" > /dev/null || exit
SCRIPT_PATH=$(pwd -P)
popd > /dev/null || exit

for i in $(find . -name '*.go')
do
  if  ! grep -q Copyright "$i" ; then
     # exclud third_party and tutorial
     if [[ "$i" != *"third_party"* && "$i" != *"tutorial"* ]]; then
       #echo "$i"
       cat "${SCRIPT_PATH}"/copyright.txt "$i" >"$i".new && mv "$i".new "$i"
     fi
  fi
done
