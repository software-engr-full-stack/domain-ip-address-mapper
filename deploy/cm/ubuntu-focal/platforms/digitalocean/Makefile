_cm_ubuntu_focal_digitalocean_mkfile_dir := $(dir $(abspath $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))
_cm_ubuntu_focal_digitalocean_mkfile_dir := $(_cm_ubuntu_focal_digitalocean_mkfile_dir:/=)

_cm_ubuntu_focal_digitalocean_app_dir := ${_cm_ubuntu_focal_digitalocean_mkfile_dir}/../../../../..
_cm_ubuntu_focal_digitalocean_cm_dir := ${_cm_ubuntu_focal_digitalocean_mkfile_dir}/../..

_cm_ubuntu_focal_digitalocean_inventory_file := /tmp/ansible/ubuntu-focal_digitalocean.inventory

include ${_cm_ubuntu_focal_digitalocean_app_dir}/lib/Makefile

# DEBUG: put in main Makefile later
remote_ips := $$(${_cm_ubuntu_focal_digitalocean_app_dir}/deploy/platforms/digitalocean/lib/remote-ips.py)

-cm-ubuntu-focal-digitalocean-remote-ip-validation:
	@:$(call check_defined, remote_ips)

-cm-ubuntu-focal-digitalocean-build-inventory: -cm-ubuntu-focal-digitalocean-remote-ip-validation
	cd "${_cm_ubuntu_focal_digitalocean_mkfile_dir}" && \
	${_cm_ubuntu_focal_digitalocean_app_dir}/lib/create-inventory-file.py \
		--remote-ips "${remote_ips}" \
		--inventory-file "${_cm_ubuntu_focal_digitalocean_inventory_file}"

_cm_ubuntu_focal_digitalocean_secrets_dir := ${_cm_ubuntu_focal_digitalocean_app_dir}/secrets

# do_reset_and_migrate=yes: run on one server first, then
#   do_reset_and_migrate=no on the rest of the servers
cm-ubuntu-focal-digitalocean-run: -cm-ubuntu-focal-digitalocean-build-inventory
	cd "${_cm_ubuntu_focal_digitalocean_mkfile_dir}" && \
	ansible-playbook \
		--inventory-file "${_cm_ubuntu_focal_digitalocean_inventory_file}" \
		--extra-vars secrets_dir="${_cm_ubuntu_focal_digitalocean_secrets_dir}" \
		--extra-vars do_reset_and_migrate='yes' \
    --extra-vars "@${_cm_ubuntu_focal_digitalocean_mkfile_dir}/../../vars.yml" \
    --extra-vars "cm_dir=${_cm_ubuntu_focal_digitalocean_cm_dir}" \
    --extra-vars "@${_cm_ubuntu_focal_digitalocean_app_dir}/db/postgres.yml" \
    --extra-vars '{"tls_container":${tls}}' \
    --extra-vars "fetched_tls_files_dir=${fetched_tls_files_dir}" \
    --extra-vars '@${domain_config_file}' \
		${_cm_ubuntu_focal_digitalocean_mkfile_dir}/playbook.yml
