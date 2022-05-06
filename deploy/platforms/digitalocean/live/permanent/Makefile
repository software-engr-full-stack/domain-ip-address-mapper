_digitalocean_live_permanent_mkfile_dir := $(dir $(abspath $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))
_digitalocean_live_permanent_mkfile_dir := $(_digitalocean_live_permanent_mkfile_dir:/=)

_digitalocean_live_permanent_app_dir := ${_digitalocean_live_permanent_mkfile_dir}/../../../../..

include ${_digitalocean_live_permanent_app_dir}/lib/Makefile

action :=
-digitalocean-live-permanent-action-validation:
	@:$(call check_defined, action)

digitalocean-live-permanent-run: -digitalocean-live-permanent-action-validation
	cd "${_digitalocean_live_permanent_mkfile_dir}" && \
	${_digitalocean_live_permanent_app_dir}/lib/terraform.sh \
		"${action}" \
		"${_digitalocean_live_permanent_app_dir}/secrets/digitalocean.json" \
		"${_digitalocean_live_permanent_mkfile_dir}/../config.json"
