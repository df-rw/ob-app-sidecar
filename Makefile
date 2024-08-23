APP=ob-nginx-test

docker.build:
	npm run build
	docker build \
		--tag ${APP} \
		-f Dockerfile.local .

docker.run:
	docker run \
		-p 8080:80 \
		${APP}:latest

# Run a shell so we can manually inspect the generated image.
docker.inspect:
	docker run -it ${APP}:latest sh
