APP=ob-app-sidecar

APP_INGRESS=${APP}-ingress
APP_BACKEND=${APP}-backend
APP_VALIDATOR=${APP}-validator

.PHONY: \
	docker.build \
	docker.build.ingress \
	docker.build.backend \
	docker.build.validator \
	docker.up \
	docker.down \
	docker.clean \
	clean \
	cloudbuild \
	help \
	backend.dev \
	frontend.dev \
	ingress.dev

.DEFAULT_GOAL=help

help: # me
	@grep '^[a-z]' Makefile | sed -e 's/^\(.*\): .*# \(.*\)/\1: \2/'

docker.build: docker.build.ingress docker.build.backend docker.build.validator # Build front and back ends.

docker.build.ingress: # Build just the ingress container.
	docker compose build ingress

docker.build.backend: # Build just the backend container.
	docker compose build backend

docker.build.validator: # Build just the validator container.
	docker compose build validator

docker.up: # Run docker compose up
	docker compose up

docker.down: # Run docker compose down
	docker compose down

docker.clean: docker.down # Clear out all the docker things.
	docker image rm -f ${APP_INGRESS}
	docker image rm -f ${APP_BACKEND}
	docker image rm -f ${APP_VALIDATOR}

clean: docker.clean # Clean things.

backend.dev: # run the backend server
	make -C backend

frontend.dev: # run the frontend server
	make -C frontend

ingress.dev: # run nginx dev
	nginx -p ./nginx -c nginx-dev.conf
