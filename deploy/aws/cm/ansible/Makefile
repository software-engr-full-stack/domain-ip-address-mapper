inventory_file := ./inventory
remote_ip="$$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' 'test-vm-container')"

build-inventory:
	export makefile_remote_ip="${remote_ip}" && \
		cat ./inventory.template | envsubst | tee "${inventory_file}"

run: build-inventory
	ansible-playbook \
		--inventory-file "${inventory_file}" \
		--vault-password-file=../../../../priv/ansible-vault-pw \
		./playbook.yml

ssh:
	ssh -o StrictHostKeyChecking=no "remote@${remote_ip}"; \
		c--change-terminal-title "${title}"
