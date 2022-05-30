#!/usr/bin/env bash

# exit by command return non-zero exit code
set -o errexit
# Indicate an error when it encounters an undefined variable
set -o nounset
# Fail on any error.
set -o pipefail
# set -o xtrace


function get_version_from_git() {
  GIT="git"
  GIT_COMMIT=$("${GIT}" rev-parse "HEAD^{commit}")
  GIT_TAG=$(git describe --long --tags --dirty --tags --always)
  GIT_BUILD_TIME=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')


  GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)

  GIT_TREE_STATE=${GIT_TREE_STATE-}
  if git_status=$("${GIT}" status --porcelain) && [[ -z ${git_status} ]]; then
     GIT_TREE_STATE="clean"
  else
     GIT_TREE_STATE="dirty"
  fi
}

function gitinfos() {
  get_version_from_git

  local -a gitinfos
  function add_gitinfo() {
    local key=${1}
    local val=${2}

    # update the list github.com/kaydxh/golang/pkg/app.
     gitinfos+=(
      "${key}=${val}"
    )
  }

  add_gitinfo "buildDate" "${GIT_BUILD_TIME}"
  add_gitinfo "gitVersion" "${GIT_TAG}"
  add_gitinfo "gitCommit" "${GIT_COMMIT}"
  add_gitinfo "gitTreeState" "${GIT_TREE_STATE}"

 # "$*" => get arg1 arg2 arg3 as a single argument "a1 a2 a3"
 # "$@" => gets arg1, arg2 and arg3 as a separate arguments "a1" "a2" "a3"
 # if no quotes, $* is the same to $@, as a separate arguments "a1" "a2" "a3"
 for gitinfo in "${gitinfos[@]}"
 do
   echo ${gitinfo}
 done
}

# Allows to call a function based on arguments passed to the script
$*
