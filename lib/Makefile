# Common

# https://stackoverflow.com/questions/11958626/make-file-warning-overriding-commands-for-target/49804748
# ifeq here is just saying "is the value empty"
# the value "defined" is just anything that isn't blank
ifeq ($(_COMMON_MAKEFILE_),)
_COMMON_MAKEFILE_ := defined

check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
        $(error Undefined '$1'$(if $2, ($2))$(if $(value @), \
                required by target '$@')))

# Long var names because make variable files are global
_lib_mkfile_dir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
_lib_app_root_dir := $$(dirname ${_lib_mkfile_dir})

build:
	cd "${_lib_app_root_dir}" && \
	go build -o bin/exec main.go

db-reset: build
	./bin/db/setup/reset

db-reset-and-migrate: build
	cd "${_lib_app_root_dir}" && \
	./bin/db/setup/reset-and-migrate

db-reset-and-seed: build
	cd "${_lib_app_root_dir}" && \
	./bin/db/setup/reset-and-migrate && \
		./bin/db/testseed/testseed && \
		./bin/debug/debug

endif
