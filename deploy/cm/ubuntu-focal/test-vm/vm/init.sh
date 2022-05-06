#!/usr/bin/bash

run() {
  export DEBIAN_FRONTEND=noninteractive
  handle_packages
  handle_services

  local ru="${REMOTE_USER-}"
  if [ -z "$ru" ]; then
    ru="${TMPL_REMOTE_USER}"
  fi
  local remote_user="${ru:?ERROR => must provide REMOTE_USER or TMPL_REMOTE_USER}"

  setup_remote_user "$remote_user"
  convenience "$remote_user"
}

handle_packages() {
  apt-get update || error 'apt failed to update'
  apt-get -y upgrade || error 'apt failed to upgrade'

  declare -a packages=(
    'openssh-server'
    'sudo'
    'python3.8'
    'bash-completion'
    'vim-tiny'
    'tree'
    'wget'
    'ca-certificates'
    - git
    # So can unarchive *.xz files
    - xz-utils

    # Tools, convenience
    - nmap
    - tcpdump
    - aptitude
  )

  local pkg=
  for pkg in "${packages[@]}"; do
    if ! dpkg -s "$pkg"; then
      apt-get install --yes --no-install-recommends "$pkg" || error "failed to install '$pkg' package"
    fi
  done
}

handle_services() {
  declare -a services=(
    'ssh'
    # 'sudo'
  )
  local svc=
  for svc in "${services[@]}"; do
    if type systemctl; then
      systemctl enable "$svc"
    elif type update-rc.d; then
      update-rc.d "$svc" defaults
      update-rc.d "$svc" enable
      service "$svc" start
    else
      error "systemctl and update-rc.d unavailable, cannot start '$svc' service"
    fi
  done
}

setup_remote_user() {
  local remote_user="${1:?ERROR => must provide remote user name}"

  local id='1455'
  groupadd --gid "$id" "$remote_user"

  useradd --create-home \
    --uid "$id" --gid "$id" \
    --shell /bin/bash \
    "$remote_user" || error "failed to create user '$remote_user'"

  echo "$remote_user ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/remote

  usermod -a -G adm "$remote_user"
  usermod -a -G sudo "$remote_user"

  local home_dir="$(get_home_dir "$remote_user")"
  [ -n "$home_dir" ] || error "home dir for '$remote_user' not found"

  local ssh_dir="$home_dir/.ssh"
  local auth_keys_path="$ssh_dir/authorized_keys"
  mkdir -p "$ssh_dir"
  mv '/tmp/ssh-pub-key' "$auth_keys_path"
  chmod 600 "$auth_keys_path"
  chown "$remote_user:$remote_user" -R "$ssh_dir"
}

convenience() {
  local remote_user="${1:?ERROR => must provide remote user name}"

  local bashrc='/etc/bash.bashrc'
  local home_dir="$(get_home_dir "$remote_user")"
  [ -n "$home_dir" ] || error "home dir for '$remote_user' not found"

  local user='root'
  local root_home="$(get_home_dir "$user")"
  [ -n "$root_home" ] || error "home dir for '$user' not found"
  (
    echo
    echo "# ... $(date)"
    echo "alias ls='/bin/ls -l --color=auto --group-directories-first'"
    echo "alias grep='grep --color=auto'"
  ) | tee -a '/etc/bash.bashrc' "$home_dir/.bashrc" "$root_home/.bashrc"
}

get_home_dir() {
  local user="${1:?ERROR => must pass user name}"

  echo "$(getent passwd "$user" | cut -f 6 -d ':')"
}

_logger() {
  local msg="${1:?ERROR => must pass log message}"
  local msg_type="${2-INFO}"
  logger "... $msg_type: $msg"
}

error() {
  local msg="$${1-}"
  [ -n "$msg" ] || read msg
  if [ -z "$msg" ]; then
    _logger 'must pass or pipe error message' 'ERROR'
    exit 1
  fi

  _logger "$msg" 'ERROR'

  exit 1
}

set -o errexit
set -o pipefail
set -o nounset
run "$@"
