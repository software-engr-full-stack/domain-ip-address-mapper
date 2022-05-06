# https://ftp.gnu.org/old-gnu/Manuals/make-3.80/html_node/make_17.html
_cm_ubfoc_tvm_cm_mkfile_dir := $(dir $(abspath $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))
_cm_ubfoc_tvm_cm_mkfile_dir := $(_cm_ubfoc_tvm_cm_mkfile_dir:/=)

_cm_ubfoc_tvm_cm_app_dir := $(abspath ${_cm_ubfoc_tvm_cm_mkfile_dir}/../../../../..)

_cm_ubfoc_tvm_cm_secrets_dir := ${_cm_ubfoc_tvm_cm_app_dir}/secrets
_cm_ubfoc_tvm_cm_secrets_file := ${_cm_ubfoc_tvm_cm_secrets_dir}/secrets.yml

include ${_cm_ubfoc_tvm_cm_app_dir}/lib/Makefile
include ${_cm_ubfoc_tvm_cm_mkfile_dir}/../lib/Makefile

# Must be after include because it's using some common variables found in the other Makefiles
_cm_ubfoc_tvm_cm_remote_ip = "$$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "${_cm_ubfoc_tvm_container_name}")"

-cm-ubuntu-focal-test-vm-cm-validate-remote-ip:
	@:$(call check_defined, _cm_ubfoc_tvm_cm_remote_ip)

# Saving terminal title to stack: https://stackoverflow.com/a/20048874
test-vm-ssh: -cm-ubuntu-focal-test-vm-cm-validate-remote-ip
	@/bin/echo -ne '\033[22;0t' && \
		ssh -o StrictHostKeyChecking=no "remote@${_cm_ubfoc_tvm_cm_remote_ip}"; \
		/bin/echo -ne '\033[23;0t'

test-vm-http: build cm-ubuntu-focal-test-vm-generate-config-files
	SECRETS_FILE="${_cm_ubfoc_tvm_cm_secrets_file}" \
		APP_ENV_DB_CONFIG_FILE=${_cm_ubfoc_tvm_db_config_within_containers} \
		${_cm_ubfoc_tvm_cm_app_dir}/bin/serv/http

_cm_ubfoc_tvm_inventory_file := /tmp/ansible/ubuntu-focal_test-vm.inventory
-cm-ubuntu-focal-test-vm-cm-build-inventory: -cm-ubuntu-focal-test-vm-cm-validate-remote-ip
	cd "${_cm_ubfoc_tvm_cm_mkfile_dir}" && \
	${_cm_ubfoc_tvm_cm_app_dir}/lib/create-inventory-file.py \
		--remote-ips "[\"${_cm_ubfoc_tvm_cm_remote_ip}\"]" \
		--inventory-file "${_cm_ubfoc_tvm_inventory_file}"

_cm_ubfoc_tvm_cm_main_cm_dir := $(abspath ${_cm_ubfoc_tvm_cm_mkfile_dir}/../..)

# TODO: refactor somehow
_cm_ubfoc_tvm_cm_makefile_path_in_remote := /deploy/app/deploy/cm/ubuntu-focal/test-vm/cm/Makefile

cm-ubuntu-focal-test-vm-cm: -cm-ubuntu-focal-test-vm-cm-build-inventory cm-ubuntu-focal-test-vm-generate-config-files
	cd "${_cm_ubfoc_tvm_cm_mkfile_dir}" && \
	ansible-playbook \
		--inventory-file "${_cm_ubfoc_tvm_inventory_file}" \
		--extra-vars 'sysv_service_cmd="make --file=${_cm_ubfoc_tvm_cm_makefile_path_in_remote} test-vm-http"' \
		--extra-vars "cm_dir=${_cm_ubfoc_tvm_cm_main_cm_dir}" \
		--extra-vars "@${_cm_ubfoc_tvm_cm_main_cm_dir}/vars.yml" \
		--extra-vars "@${_cm_ubfoc_tvm_db_config_within_containers}" \
		--extra-vars "APP_ENV_DB_CONFIG_FILE=${_cm_ubfoc_tvm_db_config_within_containers}" \
		--extra-vars "secrets_dir=${_cm_ubfoc_tvm_cm_secrets_dir}" \
		--extra-vars do_reset_and_migrate='yes' \
		${_cm_ubfoc_tvm_cm_mkfile_dir}/playbook.yml

# TODO: not sure if needed, reset database from host
cm-ubuntu-focal-test-vm-cm-db-reset: build cm-ubuntu-focal-test-vm-generate-config-files
	SECRETS_FILE="${_cm_ubfoc_tvm_cm_secrets_file}" \
		APP_ENV_DB_CONFIG_FILE=${_cm_ubfoc_tvm_db_config_host_to_containers} \
		${_cm_ubfoc_tvm_cm_app_dir}/bin/db/setup/reset-and-migrate
