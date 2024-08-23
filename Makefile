APP=ob-test

APP_INGRESS=${APP}-ingress
APP_BACKEND=${APP}-backend
APP_VALIDATOR=${APP}-validator

.PHONY: \
	docker.build-ingress \
	docker.build-backend \
	docker.build-validator \
	docker.run-ingress \
	docker.run-backend \
	docker.run-validator \
	docker.inspect-ingress \
	docker.inspect-backend \
	docker.inspect-validator

.DEFAULT_GOAL=help

docker.build-ingress: # build the ingress container
	npm run build # build observable framework static site
	docker build --tag ${APP_INGRESS} -f Dockerfile.local-ingress .

docker.build-backend: # build the backend application container
	# docker build --tag ${APP_BACKEND} -f Dockerfile.local-backend .

docker.build-validator: # build the validator container
	# docker build --tag ${APP_VALIDATOR} -f Dockerfile.local-validator .

docker.run-ingress: # run the ingress container
	docker run -p 8080:8080 ${APP_INGRESS}:latest

docker.run-backend: # run the backend container
	# docker run -p 8081:8081 ${APP_BACKEND}:latest

docker.run-validator: # run the validator container
	# docker run -p 8082:8082 ${APP_VALIDATOR}:latest

docker.inspect-ingress: # inspect the ingress container
	# docker run -it ${APP_INGRESS}:latest sh

docker.inspect-backend: # inspect the backend container
	# docker run -it ${APP_BACKEND}:latest sh

docker.inspect-validator: # inspect the validator container
	# docker run -it ${APP_VALIDATOR}:latest sh

cloudbuild: # do a deploy onto cloudbuild
	# @echo "TODO send everything to cloud build"

help: # me
	@grep '^[a-z]' Makefile | sed -e 's/ #//'
