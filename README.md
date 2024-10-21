# ob-app-sidecar

This repository shows how to setup and deploy an [Observable
Framework](https://observablehq.com/framework) frontend application with a
backend application written in [Go](https://go.dev), as well as a validator
sample application, onto [Google Cloud Run](https://cloud.google.com/run) in a
[sidecar](https://cloud.google.com/run/docs/deploying#sidecars) configuration.
[nginx](https://nginx.org) is used as the ingress container and hosts the
statically built Observable Framework application; the backend application and
the validator application are separate sidecars.

```
                ----------------------------------------------------------------------------------
                |                                                                                |
                | -------------------------------------------------    ------------------------- |
                | |              -------------------------------- |    |                       | |
----------      | |              |  Observable Framework static | |    |       validator       | |
| client | <--> | | --------------- site served from the proxy  | |<-->|      application      | |
----------      | | | nginx proxy |------------------------------ |    |         server        | |
                | | ---------------                               |    |                       | |
                | -------------------------------------------------    ------------------------- |
                |               ingress container                           sidecar container    |
                |                       ^                                                        |
                |                       |                                                        |
                |                       v                                                        |
                | -----------------------------------------------                                |
                | |                                             |                                |
                | |        backend application server           |                                |
                | |                                             |                                |
                | -----------------------------------------------                                |
                |                sidecar container                                               |
                |                                                                                |
                ----------------------------------------------------------------------------------
                                                   Cloud Run instance



```

The application is based on <https://github.com/df-rw/ob-app>.

## Aside: What is a validator application?

The validator application accepts requests from the client and validates them
prior to the request being handed off to the backend application. What the
validation does is application dependent, and could be anything like a basic
auth check, a cookie check or JWT validation.

The validator application exists _separately_ from the application and proxy;
however the proxy is (and in production, must be!) configured to send all
requests through the validator.

The purpose of having a validator application, aside from it's function, is to
provide a single method of validation that doesn't require changes to the
backend application.

It can also provide a second layer of security. For example, when using
Google's Identity-Aware Proxy (IAP), the header `X-Goog-IAP-JWT-Assertion` is
added to a client request by IAP before reaching the backend application. This
can be checked in the validator to ensure that IAP is enabled.

## Prerequisites

For this demo:

- Go (`brew install go` or [rtfm](https://go.dev/doc/install))
- nginx (`brew install nginx` or [rtfm](https://nginx.org/en/docs/install.html))

### Optional

- [air](https://github.com/air-verse/air) for rebuilding the Go backend
  application on file changes during development:

  ```shell
  go install github.com/air-verse/air@latest
  ```

  Configuration for air is in `./.air.toml`.

- [Docker](https://docker.com) for testing all applications in their own
  containers.

  ```shell
  brew install --cask docker
  ```

## Install

```shell
git clone https://github.com/df-rw/ob-app-sidecar
cd ob-app-sidecar
npm install # install modules for Observable Framework
```

### Optional

You may wish to set the email address of a valid user account while developing
application if you need to differentiate between users. See <a
href="#why-google-header-nginx-conf">Why do nginx-dev.conf and
nginx-docker.conf set a Google header?</a> for why you may want to do this:

```shell
sed -i '' -e 's/foo@bar.com/<a valid email address/>' nginx-dev.conf
sed -i '' -e 's/foo@bar.com/<a valid email address/>' nginx-docker.conf
```

## Local development

For ease of development, each of the nginx, Observable Framework and backend
servers should be started in separate terminals. This allows easy log viewing
of each individual service, as well as restarting only individual services if
required:

```shell
cd backend && go run ./cmd/web/*.go -p 6082         # (or "air") Start backend server.
cd frontend && npm run dev -- --port 6081 --no-open # Start Observable framework (diff terminal).
nginx -p ./nginx -c ./nginx-dev.conf                # Start nginx (diff terminal).
```

How traffic moves through the development environment:

```
----------      ---------------
| client | <--> | nginx proxy |
----------      ---------------
              :6080  ^    ^
                     |    |
                     |    |      ----------------------
                     |    -----> | backend app server |
                     |           ----------------------
                     |        :6082
                     |
                     |           -------------------------------
                     ----------> | Observable Framework server |
                                 -------------------------------
                              :6081
```

- `nginx-dev.conf` is configured to listen on port `6080` for inbound requests,
  and pass off to Observable on port `6081` and the application server on
  `6082`.
- nginx proxies requests from the client to either Observable Framework or the
  application server based on the URL of the request.
- nginx also proxies the Observable Framework websocket connection for live-reloading
  of the frontend. Changes you make to Observable Framework code will be automatically
  reloaded in the client.
- The validator service is [stubbed
  out](https://nginx.org/en/docs/http/ngx_http_rewrite_module.html#return) in
  `nginx-dev.conf` and is set to let all requests pass through. You can change this as
  and when required, but the idea is for validation to not get in the way while writing
  your application.

Open browser to <http://localhost:6080>. Click click click, hack hack hack.

## Adding backend routes

`nginx-dev.conf` is setup to pass any requests starting with `/api/` to the
backend application. If there are specific paths you wish to forward to the
backend application, adjust `nginx-*.conf` to suit.

## Testing containers

Once you're happy with your application, you may wish to test it locally with
each part of the whole in a separate container. We can do this with
[Docker](https://docker.com).

```
                -------------------------------------------------
                |              -------------------------------- |    --------------------------
----------      |              |  Observable Framework static | |    |  --------------------  |
| client | <--> | --------------- site served from the proxy  | |<-->|  | validator server |  |
----------      | | nginx proxy |------------------------------ |    |  --------------------  |
                | ---------------                               |    --------------------------
                -------------------------------------------------  :8081  validator container
             :8080              ingress container
                                        ^
                                        |
                                        v
                -------------------------------------------------
                |             ----------------------            |
                |             | backend app server |            |
                |             ----------------------            |
                -------------------------------------------------
             :8082           application container
```

- We setup the containers using [docker compose](https://docs.docker.com/compose/).
  Configuration is in `compose.yaml`.
- `nginx-docker.conf` is configured to listen on port `8080` for inbound
  requests, pass off application requests to the application server on
  `8082`, with the validation service listening on port `8081`.
- `make docker.build` will build all containers.
- `make docker.up` will run all containers.
- `make docker.down` will stop everything.
- `make docker.clean` will kill everything.

Open a browser to <http://localhost:8080>. Click click click.

## Deploy to Cloud Build

Configuration for Cloud Build can be found in `cloudbuild.yaml`. There is also
configuration for the Cloud Run service in `run-service.yaml`. The intended
deployment method is a push to a branch, so a Cloud Build trigger should be
setup to run the deploy on this event.

### Why two Cloud Build files?

New Cloud Run services can be created and updated from Cloud Build using
`gcloud run deploy`. However:

- sidecar deployments using `gcloud run deploy` are [currently in
  preview](https://cloud.google.com/run/docs/deploying#multicontainer-yaml);
- sidecar start up order, which is required for our layout, requires specifying
  both container dependencies and startup healthcheck probes.

Container dependencies [can be specified using `gcloud run
deploy`](https://cloud.google.com/run/docs/configuring/services/containers#gcloud_2)
but [healthchecks cannot, and can only be specified via the console, Terraform
or from a .yaml](https://cloud.google.com/run/docs/configuring/healthchecks).
We use the service configuration in `run-service.yaml` to specify the startup
probes.

We can't specify the service configuration with `gcloud run deploy`, so have to
use `gcloud services replace` instead to get a new service revision going. This
will only run if the configuration supplied (ie. from `run-services.yaml`)
changes between runs. This is good, as there is no point creating a new Cloud
Run revision if nothing changed. This is bad however, if we make application
changes (ie. to our images) and not to our service; we will build and upload
new images, but the service won't deploy a new revision.

`cloudbuild.yaml` is set to use the `COMMIT_SHA` as the tag for our images. This
value will be whatever the commit checksum is given on the push to repo. This
value is written in `run-sevice.yaml` on each deployment, so our service will
get a new revision each deploy as the configuration has changed.

- [Noted in this SO
question](https://stackoverflow.com/questions/77550717/gcloud-run-replace-not-deploying-a-new-revision-when-no-change-in-the-file).
- Worth noting that using `:latest` won't work; since this tag will never be
  changed in `run-service.yaml` we end up in the same
  uploaded-new-images-but-no-new-service-revision deployed basket.

- Also worth noting that we don't get any errors if `gcloud services replace`
  didn't start a new revision. You need to check the checksums on the images
  via `gcloud run revisions describe...` to see what images are being used.

## Notes

### Why are there so many Dockerfiles?

| Filename | Porpoise |
| --- | --- |
| `Dockerfile.backend` | Builds the container for the backend application. |
| `Dockerfile.ingress-gcp` | Builds the ingress container on GCP. |
| `Dockerfile.ingress-local` | Builds the ingress container for local (Docker) use. |
| `Dockerfile.validator-dummy` | Builds the container for the dummy validator application. |
| `Dockerfile.validator-iap` | Builds the container for the IAP validator application. |

`Dockerfile.ingress-gcp` differs from `Dockerfile.ingress-local` as they
reference (slightly) different nginx configurations (see below).

`Dockerfile.validator-dummy` builds a dummy validator that returns success for
every connection. `Dockerfile.validator-iap` uses `cmd/cli/validator-iap.go`
to validate connections through IAP.

### Why are there so many nginx configurations?

nginx configuration files are in `./nginx`:

| Filename | Porpoise |
| --- | --- |
| `nginx-dev.conf` | nginx configuration for local development environment. |
| `nginx-docker.conf` | nginx configuration for local Docker environment. |
| `nginx-gcp.conf` | nginx configuration for GCP environment. |

Docker compose [creates a network for the
container](https://docs.docker.com/compose/networking/). This allows internal
applications to refer to each other by service name as specified in
`compose.yaml`. For instance, we can refer to the backend application with the
service name `backend`. `nginx-docker.conf` is setup like this.

This differs to GCP which uses [localhost and unique port
numbers](https://cloud.google.com/run/docs/internet-proxy-nginx-sidecar#configure_the_nginx_ingress_container)
to refer to services. For instance, we refer to the backend application on
server `127.0.0.1`, with the port diferrentiating services. `nginx-gcp.conf` is
setup like this.

<a name="why-google-header-nginx-conf"></a>
### Why do nginx-dev.conf and nginx-docker.conf set a Google header?

The sample application in this repository is ready to be deployed onto Google
Cloud with IAP sitting in front of it. Two headers are provided by IAP to
backend applications:

- [`X-Goog-IAP-JWT-Assertion`](https://cloud.google.com/iap/docs/signed-headers-howto): 
    a JWT supplied by IAP that should be verified by the application.
- [`X-Goog-Authenticated-User-Email`](https://cloud.google.com/iap/docs/identity-howto#getting_the_users_identity_with_signed_headers): the email address of the authenticated user. 

The validator application `cmd/cli/validator-iap.go` validates the JWT, and
also checks the claim in the JWT for the user by using
`X-Goog-Authenticated-User-Email`. If both these conditions pass, the validator
lets the request through to the backend application.

The backend application should know who the user is, and while it could
parse the JWT again, it's easier to just use `X-Goog-Authenticated-User-Email`
since this has been verified by the validator application. From [the
docs](https://cloud.google.com/iap/docs/identity-howto#getting_the_users_identity_with_signed_headers):

> If you use these headers, you must compare them against the identity
> information from the authenticated JWT header listed above.

Since we don't use IAP in local development, we fake
`X-Goog-Authenticated-User-Email` in `nginx-dev.conf` and `nginx-docker.conf`.
This also allows easy testing across different accounts locally.

<a name="dot-env"></a>
### .env

See `.env-sample` for a... sample.

| Name | Value(s) | Notes |
|---|---|---|
| `OBSERVABLE_TELEMETRY_DISABLE` | `true` or `false` | [Observable Framework Telemetry](https://observablehq.com/framework/telemetry) off or on. |

## Other documentations

- [Adding environment variables](docs/env-vars.md)
- [Enabling IAP on a Google Cloud Run service](docs/iap.md)
- [Validating IAP JWT tokens](docs/iap-jwt.md)

## TODO

- Rewrite nginx configuration based on local (Docker) deploy or GCP deploy.
