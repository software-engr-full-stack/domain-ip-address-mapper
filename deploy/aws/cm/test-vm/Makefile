ssh_pub_key_copy := ./ssh-pub-key
priv_file := ../../../../priv/init.env.sh
test_vm_name = test-vm

not_in_working_dir = [ "$$(basename "$$(readlink -f .)")" != "${test_vm_name}" ]
working_dir := ./deploy/aws/cm/${test_vm_name}

test-vm-build:
	if $(call not_in_working_dir); then cd "${working_dir}"; fi && \
	cp "${HOME}/.ssh/id_ed25519.pub" "${ssh_pub_key_copy}" && \
	. "${priv_file}" && \
	export $$(cut -d= -f1 ${priv_file}) && \
	docker build --tag "${test_vm_name}:latest" \
		--build-arg PRIV_REMOTE_USER="$$PRIV_REMOTE_USER" \
		--build-arg PRIV_REMOTE_PW="$$PRIV_REMOTE_PW" \
		. && \
	rm "${ssh_pub_key_copy}"

test-vm-run:
	if $(call not_in_working_dir); then cd "${working_dir}"; fi && \
	docker-compose up

test-vm-nuke:
	docker system prune --force && \
	docker rmi "${test_vm_name}"
