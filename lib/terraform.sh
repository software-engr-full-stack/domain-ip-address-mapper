#!/usr/bin/env bash

run() {
  local action="${1:?ERROR => must provide action}"
  shift

  declare -A valid_actions_table=(
    [apply]=true
    [destroy]=true
    [plan]=true
  )

  local is_valid="${valid_actions_table[$action]-}"
  if [ -z "$is_valid" ]; then
    _error "unsupported action '$action'"
  fi

  declare -A args_table=(
    [apply]="-auto-approve"
    [destroy]="-auto-approve"
    [plan]=''
  )

  local args="${args_table[$action]-}"

  local arg_files=''
  local af=''
  for af in "$@"; do
    arg_files="$arg_files -var-file=$af"
  done

  echo "... executing: terraform '$action' '$args' '$arg_files'"

  terraform "$action" $args $arg_files
}

_error() {
  echo "... ERROR: $@" >&2
  exit 1
}

set -o errexit
set -o pipefail
set -o nounset
run "$@"
