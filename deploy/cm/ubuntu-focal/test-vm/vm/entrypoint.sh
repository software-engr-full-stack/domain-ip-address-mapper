#!/usr/bin/env bash

run() {
  /usr/sbin/sshd -o ListenAddress=0.0.0.0

  tail -f /dev/null
}

set -o errexit
set -o pipefail
set -o nounset
run "$@"
