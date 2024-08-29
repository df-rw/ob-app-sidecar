include .env

APP=ob-app-sidecar

APP_INGRESS=${APP}-ingress
APP_BACKEND=${APP}-backend

.PHONY: \
	docker.build \
	docker.build.ingress \
	docker.up \
	docker.down \
	docker.clean \
	clean \
	cloudbuild \
	help

.DEFAULT_GOAL=help

docker.build: docker.build.ingress # Build front and back ends.

docker.build.ingress: # Build just the ingress container.
	docker compose build ingress

docker.up: # Run docker compose up
	docker compose up

docker.down: # Run docker compose down
	docker compose down

docker.clean: docker.down # Clear out all the docker things.
	docker image rm -f ${APP_INGRESS}

clean: docker.clean

cloudbuild:
	gcloud builds submit \
		--region=${GOOGLE_REGION} \
		--service-account=projects/${GOOGLE_CLOUD_PROJECT}/serviceAccounts/${GOOGLE_CLOUD_SERVICE_ACCOUNT} \
		--config=cloudbuild.yaml

help: # me
	@grep '^[a-z]' Makefile | sed -e 's/^\(.*\): .*# \(.*\)/\1: \2/'
