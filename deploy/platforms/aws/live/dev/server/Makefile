_aws_live_dev_server_mkfile_dir := $(dir $(abspath $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))
_aws_live_dev_server_mkfile_dir := $(_aws_live_dev_server_mkfile_dir:/=)

_aws_live_dev_server_app_dir := ${_aws_live_dev_server_mkfile_dir}/../../../../../..

include ${_aws_live_dev_server_app_dir}/lib/Makefile

-aws-live-dev-server-action-validation:
	@:$(call check_defined, action)

aws-live-dev-server-localstack: -aws-live-dev-server-action-validation
	cd "${_aws_live_dev_server_mkfile_dir}" && \
	/bin/cp './provider-localstack' './provider.tf' && \
	./run.sh "${action}"

-aws-live-dev-server-tfvars-file-validation:
	@:$(call check_defined, tfvars_file)

aws-live-dev-server-aws: -aws-live-dev-server-action-validation -aws-live-dev-server-tfvars-file-validation
	cd "${_aws_live_dev_server_mkfile_dir}" && \
	/bin/cp './provider-aws' './provider.tf' && \
	${_aws_live_dev_server_app_dir}/lib/terraform.sh "${action}" "${tfvars_file}"

-aws-live-dev-server-pem-file-path-validation:
	@:$(call check_defined, pem_file)

# SSH tasks: 1st arg is AWS tag name

# Using elastic IP
aws-ssh-server: -aws-live-dev-server-pem-file-path-validation
	${_aws_live_dev_server_mkfile_dir}/ssh.sh terrainfra "${pem_file}"


aws-ssh-bastion: -aws-live-dev-server-pem-file-path-validation
	${_aws_live_dev_server_mkfile_dir}/ssh.sh bastion "${pem_file}"
