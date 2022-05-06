#!/usr/bin/env bash

run() {
  local tag_name="${1:?ERROR => must provide tag_name}"
  local id_file="${2:?ERROR => must provide id file}"

  declare -A valid_tag_names_table=(
    [server]=true
    [bastion]=true
    [terrainfra]=true
  )

  local is_valid="${valid_tag_names_table[$tag_name]-}"
  if [ -z "$is_valid" ]; then
    _error "unsupported tag '$tag_name'"
  fi

  local ip="$(
    aws ec2 describe-addresses \
      --filters "Name=tag:Name,Values=$tag_name" \
      --query 'Addresses[0].PublicIp' \
      --output text
  )"

  # local ip="$(
  #   aws ec2 describe-instances \
  #     --filters "Name=tag:Name,Values=*$tag_name" \
  #     --query 'Reservations[*].Instances[*].PublicIpAddress' \
  #     --output text
  # )"

  grep '/^[ ]*$/' <<<"$ip" && _error 'IP must not be blank'

  echo "... SSH to '$ip'..."
  set +o errexit
  echo -ne '\033[22;0t'
  ssh -i "$id_file" -o StrictHostKeyChecking=no remote@"$ip"
  echo -ne '\033[23;0t'
  set -o errexit
}

_error() {
  echo "... ERROR: $@" >&2
  exit 1
}

set -o errexit
set -o pipefail
set -o nounset
run "$@"
