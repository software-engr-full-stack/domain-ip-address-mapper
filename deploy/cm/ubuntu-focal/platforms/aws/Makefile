_cm_ubuntu_focal_aws_mkfile_dir := $(dir $(abspath $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))
_cm_ubuntu_focal_aws_mkfile_dir := $(_cm_ubuntu_focal_aws_mkfile_dir:/=)

_cm_ubuntu_focal_aws_app_dir := ${_cm_ubuntu_focal_aws_mkfile_dir}/../../../../..
_cm_ubuntu_focal_aws_cm_dir := ${_cm_ubuntu_focal_aws_mkfile_dir}/../..

_cm_ubuntu_focal_aws_inventory_file := /tmp/ansible/ubuntu-focal_aws.inventory

include ${_cm_ubuntu_focal_aws_app_dir}/lib/Makefile

-cm-ubuntu-focal-remote-ip-validation:
	@:$(call check_defined, remote_ip)

-cm-ubuntu-focal-build-inventory: -cm-ubuntu-focal-remote-ip-validation
	cd "${_cm_ubuntu_focal_aws_mkfile_dir}" && \
	${_cm_ubuntu_focal_aws_app_dir}/lib/create-inventory-file.py \
		--remote-ips "[\"${remote_ip}\"]" \
		--inventory-file "${_cm_ubuntu_focal_aws_inventory_file}"

-cm-ubuntu-focal-aws-validations:
	@:$(call check_defined, tls) && \
	 	:$(call check_defined, pem_file) && \
	 	:$(call check_defined, secrets_dir) && \
		:$(call check_defined, app_env_db_config_file) && \
		:$(call check_defined, fetched_tls_files_dir) && \
		:$(call check_defined, domain_config_file)

_cm_ubfoc_common_args_aws := --inventory-file "${_cm_ubuntu_focal_aws_inventory_file}" \
  --extra-vars "@${_cm_ubuntu_focal_aws_cm_dir}/vars.yml" \
	--extra-vars "ansible_ssh_private_key_file=${pem_file}" \
	--extra-vars 'make_task_name_to_launch_http="aws-http"' \
	--extra-vars 'make_task_name_to_reset_and_migrate_db=aws-db'

# do_reset_and_migrate=yes: run on one server first, then
#   do_reset_and_migrate=no on the rest of the servers
cm-ubuntu-focal-aws: -cm-ubuntu-focal-build-inventory -cm-ubuntu-focal-aws-validations
	cd "${_cm_ubuntu_focal_aws_mkfile_dir}" && \
	ansible-playbook \
		--extra-vars do_reset_and_migrate='yes' \
		--extra-vars '{"tls_container":${tls}}' \
		--extra-vars fetched_tls_files_dir="${fetched_tls_files_dir}" \
		--extra-vars secrets_dir="${secrets_dir}" \
		--extra-vars "@${app_env_db_config_file}" \
		--extra-vars '@${domain_config_file}' \
		--extra-vars "cm_dir=${_cm_ubuntu_focal_digitalocean_cm_dir}" \
		${_cm_ubfoc_common_args_aws} \
		${_cm_ubuntu_focal_aws_mkfile_dir}/playbook.yml

-cm-ubuntu-focal-aws-deploy-code-validations:
	 	:$(call check_defined, secrets_dir) && \
	 	:$(call check_defined, pem_file) && \
		:$(call check_defined, app_env_db_config_file)

cm-ubuntu-focal-aws-deploy-code: -cm-ubuntu-focal-build-inventory -cm-ubuntu-focal-aws-deploy-code-validations
	cd "${_cm_ubuntu_focal_aws_mkfile_dir}" && \
	ansible-playbook \
		--extra-vars secrets_dir="${secrets_dir}" \
		--extra-vars "@${app_env_db_config_file}" \
		${_cm_ubfoc_common_args_aws} \
		${_cm_ubuntu_focal_aws_mkfile_dir}/playbooks/deploy-code.yml
