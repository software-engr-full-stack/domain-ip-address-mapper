_digitalocean_live_dev_server_mkfile_dir := $(dir $(abspath $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))
_digitalocean_live_dev_server_mkfile_dir := $(_digitalocean_live_dev_server_mkfile_dir:/=)

_digitalocean_live_dev_server_app_dir := ${_digitalocean_live_dev_server_mkfile_dir}/../../../../../..

_digitalocean_live_dev_server_secrets_dir := ${_digitalocean_live_dev_server_app_dir}/secrets

include ${_digitalocean_live_dev_server_app_dir}/lib/Makefile

action :=
-digitalocean-live-dev-server-action-validation:
	@:$(call check_defined, action)

digitalocean-live-dev-server-run: -digitalocean-live-dev-server-action-validation
	cd "${_digitalocean_live_dev_server_mkfile_dir}" && \
	${_digitalocean_live_dev_server_app_dir}/lib/terraform.sh \
		"${action}" \
		"${_digitalocean_live_dev_server_secrets_dir}/digitalocean.json" \
		"${_digitalocean_live_dev_server_secrets_dir}/main-user-ssh-fingerprints.json" \
		"${_digitalocean_live_dev_server_mkfile_dir}/../../config.json"

digitalocean-live-dev-server-ssh:
	${_digitalocean_live_dev_server_mkfile_dir}/ssh.py
