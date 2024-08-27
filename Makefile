APP=ob-app-sidecar

APP_INGRESS=${APP}-ingress
APP_BACKEND=${APP}-backend
APP_VALIDATOR=${APP}-validator

.PHONY: \
	up \
	down

.DEFAULT_GOAL=help

docker.build: build.of docker.build.ingress docker.build.backend docker.build.validator # Build front and back ends.

docker.build.ingress: # Build just the ingress container.
	docker compose build ingress

docker.build.validator: # Build just the validator container.
	docker compose build validator

docker.build.backend: # Build just the backend container.
	docker compose build backend

docker.up: # Run docker compose up
	docker compose up

docker.down: # Run docker compose down
	docker compose down

docker.clean: docker.down # Clear out all the docker things.
	docker image rm -f ${APP_INGRESS}
	docker image rm -f ${APP_BACKEND}
	docker image rm -f ${APP_VALIDATOR}

clean: docker.clean
	rm -rf ./dist

build.of: # Build the Observable Framework application.
	npm run build

cloudbuild: # do a deploy onto cloudbuild
	# @echo "TODO send everything to cloud build"

help: # me
	@grep '^[a-z]' Makefile | sed -e 's/^\(.*\): .*# \(.*\)/\1: \2/'
