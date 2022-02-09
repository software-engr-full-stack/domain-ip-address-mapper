dc_init := ./deploy/docker/docker-compose-init.yml
dc_main := ./deploy/docker/docker-compose-main.yml
priv_env := ./priv/priv.env

docker-db-up:
	docker-compose --file ${dc_init} --env-file ${priv_env} up --detach

# TODO: a more robust way of checking if the database is alive, pg_isready probably.
docker-db-delay:
	sleep 10

docker-db-seed:
	APP_ENV_DB=./priv/docker-db-host-to-container.yml make reset-and-seed

docker-db-down:
	docker-compose --file ${dc_init} --env-file ${priv_env} down

docker-build-go-base:
	docker build --tag go-base --file deploy/docker/go-base.Dockerfile .

docker-build-app:
	mkdir ./tmp && \
		cp ./priv/priv.yml ./tmp && \
		cp ./priv/docker-db-within-containers.yml ./tmp && \
		docker build --tag app --no-cache --file ./deploy/docker/app.Dockerfile . && \
		rm -rf ./tmp

docker-init-up: docker-db-up docker-db-delay docker-db-seed docker-db-down docker-build-go-base docker-build-app
	docker-compose --file ${dc_main} --env-file ${priv_env} up

docker-up:
	docker-compose --file ${dc_main} --env-file ${priv_env} up

docker-down:
	docker-compose --file ${dc_main} --env-file ${priv_env} down

docker-nuke:
	docker image prune --all --force; \
	docker kill $$(docker ps -q); \
	docker system prune --all --force

docker-debug:
	APP_ENV_DB=./priv/docker-db-host-to-container.yml make debug
