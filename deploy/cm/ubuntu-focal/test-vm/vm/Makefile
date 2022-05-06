# https://ftp.gnu.org/old-gnu/Manuals/make-3.80/html_node/make_17.html
_cm_ubfoc_tvm_vm_mkfile_dir := $(dir $(abspath $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))
_cm_ubfoc_tvm_vm_mkfile_dir := $(_cm_ubfoc_tvm_vm_mkfile_dir:/=)

_cm_ubfoc_tvm_vm_app_dir := $(abspath ${_cm_ubfoc_tvm_vm_mkfile_dir}/../../../../..)

_cm_ubfoc_tvm_vm_init_env := ${_cm_ubfoc_tvm_vm_mkfile_dir}/config/init.env

_cm_ubfoc_tvm_vm_docker_compose := ${_cm_ubfoc_tvm_vm_mkfile_dir}/docker-compose.yml
_cm_ubfoc_tvm_vm_ssh_pub_key_copy := ${_cm_ubfoc_tvm_vm_mkfile_dir}/ssh-pub-key

include ${_cm_ubfoc_tvm_vm_app_dir}/lib/Makefile
include ${_cm_ubfoc_tvm_vm_mkfile_dir}/../lib/Makefile

cm-ubuntu-focal-test-vm-pw:
	cat "${_cm_ubfoc_tvm_vm_init_env}" | /bin/grep REMOTE_PW | cut -f 2 -d = | xsel

cm-ubuntu-focal-test-vm-nuke:
	-cd "${_cm_ubfoc_tvm_vm_mkfile_dir}" && \
	docker kill "${_cm_ubfoc_tvm_container_name}"; \
	docker system prune --force; \
	docker rmi "${_cm_ubfoc_tvm_image_name}"

cm-ubuntu-focal-test-vm-db-up: cm-ubuntu-focal-test-vm-generate-config-files
	cd "${_cm_ubfoc_tvm_vm_mkfile_dir}" && \
	docker-compose --file ${_cm_ubfoc_tvm_vm_docker_compose} --env-file ${_cm_ubfoc_tvm_docker_db_env} up --detach db

# TODO: a more robust way of checking if the database is alive, pg_isready probably.
-cm-ubuntu-focal-test-vm-delay:
	sleep 5

cm-ubuntu-focal-test-vm-db-seed: cm-ubuntu-focal-test-vm-generate-config-files
	cd "${_cm_ubfoc_tvm_vm_mkfile_dir}" && \
	APP_ENV_DB_CONFIG_FILE="${_cm_ubfoc_tvm_db_config_host_to_containers}" make db-reset-and-seed

cm-ubuntu-focal-test-vm-db-down: cm-ubuntu-focal-test-vm-generate-config-files
	cd "${_cm_ubfoc_tvm_vm_mkfile_dir}" && \
	docker-compose --file ${_cm_ubfoc_tvm_vm_docker_compose} --env-file ${_cm_ubfoc_tvm_docker_db_env} down

cm-ubuntu-focal-test-vm-build:
	cd "${_cm_ubfoc_tvm_vm_mkfile_dir}" && \
	cp "${HOME}/.ssh/id_ed25519.pub" "${_cm_ubfoc_tvm_vm_ssh_pub_key_copy}" && \
	. "${_cm_ubfoc_tvm_vm_init_env}" && \
	export $$(cut -d= -f1 ${_cm_ubfoc_tvm_vm_init_env}) && \
	docker build \
		--file "${_cm_ubfoc_tvm_vm_mkfile_dir}/Dockerfile" \
		--tag "${_cm_ubfoc_tvm_image_name}:latest" \
		--build-arg REMOTE_USER="$$REMOTE_USER" \
		--build-arg REMOTE_PW="$$REMOTE_PW" \
		"${_cm_ubfoc_tvm_vm_mkfile_dir}" && \
	rm "${_cm_ubfoc_tvm_vm_ssh_pub_key_copy}"

cm-ubuntu-focal-test-vm-up: cm-ubuntu-focal-test-vm-generate-config-files
	cd "${_cm_ubfoc_tvm_vm_mkfile_dir}" && \
	docker-compose --file ${_cm_ubfoc_tvm_vm_docker_compose} --env-file ${_cm_ubfoc_tvm_docker_db_env} up

cm-ubuntu-focal-test-vm-down: cm-ubuntu-focal-test-vm-generate-config-files
	cd "${_cm_ubfoc_tvm_vm_mkfile_dir}" && \
	docker-compose --file ${_cm_ubfoc_tvm_vm_docker_compose} --env-file ${_cm_ubfoc_tvm_docker_db_env} down

-cm-ubuntu-focal-test-vm-notify-build-done:
	notify-send '... test VM build done'

cm-ubuntu-focal-test-vm-db-init: \
	cm-ubuntu-focal-test-vm-db-up \
	-cm-ubuntu-focal-test-vm-delay \
	cm-ubuntu-focal-test-vm-db-seed \
	cm-ubuntu-focal-test-vm-db-down

test-vm-full-reset: \
	cm-ubuntu-focal-test-vm-nuke \
	cm-ubuntu-focal-test-vm-db-init \
	cm-ubuntu-focal-test-vm-build \
	-cm-ubuntu-focal-test-vm-delay \
	-cm-ubuntu-focal-test-vm-notify-build-done \
	cm-ubuntu-focal-test-vm-up
