include deploy/docker/Makefile
include deploy/aws/cm/test-vm/Makefile

# "frontend" seems to be a reserved make target
js:
	cd ./frontend && BROWSER=none npm start

domain: reset
	./bin/db/setup/migrate && \
	./bin/db/testseed/testseed "$(domain)" && \
	./bin/debug/debug

add-domain: build
	./bin/db/testseed/testseed "$(domain)" && \
	./bin/debug/debug

build:
	go build -o bin/exec main.go

http: build
	./bin/serv/http

reset-and-migrate: build
	./bin/db/setup/reset-and-migrate

reset-and-seed: build
	./bin/db/setup/reset-and-migrate && \
		./bin/db/testseed/testseed && \
		./bin/debug/debug

reset: build
	./bin/db/setup/reset

debug: build
	./bin/debug/debug

live-reload:
	npx nodemon --signal SIGTERM --ext go --exec 'make http' && \
		notify-send --expire-time 12000 '... live reload exited'
