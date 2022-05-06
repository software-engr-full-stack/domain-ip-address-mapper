_error_icon := /usr/share/icons/gnome/32x32/status/error.png

_main_mkfile_dir := $(dir $(abspath $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))
_main_mkfile_dir := $(_main_mkfile_dir:/=)

_main_app_dir := ${_main_mkfile_dir}

_main_platforms_aws_dir := ${_main_mkfile_dir}/deploy/platforms/aws
_main_platforms_do_dir := ${_main_mkfile_dir}/deploy/platforms/digitalocean

include ${_main_mkfile_dir}/lib/Makefile

include ${_main_platforms_aws_dir}/live/dev/server/Makefile
include ${_main_mkfile_dir}/deploy/cm/ubuntu-focal/platforms/aws/Makefile

include ${_main_platforms_do_dir}/live/dev/server/Makefile
include ${_main_mkfile_dir}/deploy/cm/ubuntu-focal/platforms/digitalocean/Makefile

_main_base_name := terrainfra
_main_front_end_dir := $(abspath ${_main_app_dir}/frontend)

_main_secrets_dir := $(abspath ${_main_app_dir}/secrets)
_main_secrets_file := ${_main_secrets_dir}/secrets.yml
_main_tfvars_file := ${_main_secrets_dir}/auto-generated_tfvars.json
_main_fetched_tls_files_dir := ${_main_secrets_dir}/tmp

_dir_to_list_of_domain_tls_files := ${_main_secrets_dir}/tls/letsencrypt

-prep:
	mkdir -p "${_main_fetched_tls_files_dir}"

-sleep:
	sleep 10

# "frontend" seems to be a reserved make target
js:
	cd "${_main_front_end_dir}" && BROWSER=none npm start

js-build:
	cd "${_main_front_end_dir}" && npm run build

domain: db-reset
	export SECRETS_FILE="${_main_secrets_file}" && \
	./bin/db/setup/migrate && \
	./bin/db/testseed/testseed "$(domain)" && \
	./bin/debug/debug

add-domain: build
	export SECRETS_FILE="${_main_secrets_file}" && \
	./bin/db/testseed/testseed "$(domain)" && \
	./bin/debug/debug

http: build
	SECRETS_FILE="${_main_secrets_file}" ./bin/serv/http

live-reload:
	npx nodemon --signal SIGTERM --ext go --exec 'make http' && \
		notify-send --expire-time 12000 '... live reload exited'

server-localstack:
	docker-compose --file=${_main_platforms_aws_dir}/localstack/server-docker-compose.yml up

_main_aws_app_env_db_config_file := ${_main_secrets_dir}/auto-generated_aws-db.yml
_main_aws_pem_file := "$$(./lib/secrets_data.py pem_file)"
_main_aws_remote_ip := $$(${_main_platforms_aws_dir}/lib/remote_ip.py "${_main_base_name}")
_main_aws_domain_config_file := ${_main_platforms_aws_dir}/domains.yml
_main_aws_tls := $$( \
	${_main_mkfile_dir}/lib/tls.py \
		--domain-config-file ${_main_aws_domain_config_file} \
		--dir-to-list-of-domains ${_dir_to_list_of_domain_tls_files} \
)

-aws-generate-db-config:
	${_main_mkfile_dir}/lib/aws-db-config-file-builder.py \
		--tag-name-value "${_main_base_name}" \
		--output-file "${_main_aws_app_env_db_config_file}"

aws-http: build -aws-generate-db-config
	SECRETS_FILE="${_main_secrets_file}" \
		APP_ENV_DB_CONFIG_FILE="${_main_aws_app_env_db_config_file}" \
		./bin/serv/http

aws-db: build -aws-generate-db-config
	SECRETS_FILE="${_main_secrets_file}" \
		APP_ENV_DB_CONFIG_FILE="${_main_aws_app_env_db_config_file}" \
		./bin/db/setup/reset-and-migrate

_aws_cm_args := remote_ip="${_main_aws_remote_ip}" \
		tls="${_main_aws_tls}" \
		pem_file="${_main_aws_pem_file}" \
		secrets_dir="${_main_secrets_dir}" \
		app_env_db_config_file="${_main_aws_app_env_db_config_file}" \
		fetched_tls_files_dir="${_main_fetched_tls_files_dir}" \
		domain_config_file="${_main_aws_domain_config_file}"

aws-cm-deploy-code: -prep -aws-generate-db-config
	make cm-ubuntu-focal-aws-deploy-code ${_aws_cm_args}; \
	notify-send --expire-time 12000 "... AWS deploy code \"done\" $$(date)"

aws-ssh:
	make aws-ssh-server pem_file="${_main_aws_pem_file}"

action :=
aws-prov:
	${_main_mkfile_dir}/lib/generate-tfvars-file.py "${_main_tfvars_file}" && \
	sync && \
	sleep 1 && \
	make aws-live-dev-server-aws action=${action} tfvars_file="${_main_tfvars_file}" || \
	( \
		notify-send --icon ${_error_icon} --expire-time 12000 "... AWS provisioning failed $$(date)" && \
		exit 1 \
	)

aws-cm: -prep -aws-generate-db-config
	make cm-ubuntu-focal-aws ${_aws_cm_args} \
		&& \
		${_main_mkfile_dir}/lib/extract-fetched-tls-archive.py \
			--fetched-tls-files-dir ${_main_fetched_tls_files_dir} \
			--dest-dir ${_main_mkfile_dir}/secrets/tls/letsencrypt \
			--domain-config-file ${_main_aws_domain_config_file} || \
		( \
			notify-send --icon ${_error_icon} --expire-time 12000 "... AWS CM failed $$(date)" && \
			exit 1 \
		)

-aws-prov-apply:
	make aws-prov action=apply

aws-zero-to-launch: \
	-aws-prov-apply \
	-sleep \
	aws-cm
	notify-send --expire-time 12000 "... AWS zero to launch passed $$(date)"

digitalocean-prov:
	make digitalocean-live-dev-server-run action=${action} || \
	( \
		notify-send --icon ${_error_icon} --expire-time 12000 "... Digital Ocean provisioning failed $$(date)" && \
		exit 1 \
	)

_do_domain_config_file := ${_main_platforms_do_dir}/domains.yml
_do_tls := $$( \
	${_main_mkfile_dir}/lib/tls.py \
		--domain-config-file ${_do_domain_config_file} \
		--dir-to-list-of-domains ${_dir_to_list_of_domain_tls_files} \
)

digitalocean-cm: -prep
	make cm-ubuntu-focal-digitalocean-run \
		tls="${_do_tls}" \
		fetched_tls_files_dir="${_main_fetched_tls_files_dir}" \
		domain_config_file="${_do_domain_config_file}" && \
		${_main_mkfile_dir}/lib/extract-fetched-tls-archive.py \
			--fetched-tls-files-dir ${_main_fetched_tls_files_dir} \
			--dest-dir ${_main_mkfile_dir}/secrets/tls/letsencrypt \
			--domain-config-file ${_main_platforms_do_dir}/domains.yml || \
		( \
			notify-send --icon ${_error_icon} --expire-time 12000 "... Digital Ocean CM failed $$(date)" && \
			exit 1 \
		)

-digitalocean-prov-apply:
	make digitalocean-prov action=apply

digitalocean-zero-to-launch: \
	-digitalocean-prov-apply \
	-sleep \
	digitalocean-cm
	notify-send --expire-time 12000 "... Digital Ocean zero to launch passed $$(date)"

digitalocean-ssh:
	make digitalocean-live-dev-server-ssh
