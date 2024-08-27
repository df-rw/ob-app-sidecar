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
prior to the request being handed off to the application. What the validation
does is application dependent, and could be anything like a basic auth check, a
cookie check or JWT validation.

The validator application exists _separately_ from the application and proxy;
however the proxy is (and in production, must be!) configured to send all
requests through the validator.

The purpose of having a validator application, aside from it's function, is to
provide a single method of validation that doesn't require changes to a backend
application.

## Prerequisites

To run this demo:

- Go (`brew install go` or [rtfm](https://go.dev/doc/install))
- nginx (`brew install nginx` or [rtfm](https://nginx.org/en/docs/install.html))

### Optional

- [air](https://github.com/air-verse/air) for rebuilding the Go backend
  application on file changes during development:

```shell
go install github.com/air-verse/air@latest
```

Configuration for air is in `./.air.toml`. Replace the `go run ./cmd/web/*.go
-p 8082` line (below in **Development**) with `air` to get application rebuild
when a backend file changes.

## Install

```shell
git clone https://github.com/df-rw/ob-app-sidecar
cd ob-app-sidecar
npm install # install modules for Observable Framework
```

### Development: write your application

To get your development environment up and running:

```shell
go run ./cmd/web/*.go -p 6082        # (or "air") Start backend server.
npm run dev -- --port 6081 --no-open # Start Observable framework (diff terminal).
nginx -p . -c ./nginx-dev.conf       # Start nginx (diff terminal).
```

How traffic moves through the development environment:

```
----------      ---------------
| client | <--> | nginx proxy |
----------      ---------------
                     ^    ^
                     |    |      ----------------------
                     |    -----> | backend app server |
                     |           ----------------------
                     |           -------------------------------
                     ----------> | Observable Framework server |
                                 -------------------------------
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

### Development: adding backend routes

`nginx-dev.conf` is setup to pass any requests starting with `/api/` to the
backend application. If there are specific paths you wish to forward to the
backend application, adjust `nginx-dev.conf` to suit.

### Development: test your containers

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
                -------------------------------------------------       validator container
                                ingress container
                                        ^
                                        |
                                        v
                -------------------------------------------------
                |             ----------------------            |
                |             | backend app server |            |
                |             ----------------------            |
                -------------------------------------------------
                              application container
```

- We setup the containers using [docker compose](https://docs.docker.com/compose/).
  Configuration is in `compose.yaml`.
- `nginx-docker.conf` is configured to listen on port `8080` for inbound
  requests, pass off application requests to the application server on
  `8082`, with the validation service listening on port `8081`.
- `make docker.build` will build the observable Framework frontend application
  and all backends.
- `make docker.up` will run everything.
- `make docker.down` will stop everything.
- `make docker.clean` will kill everything.

Open browser to <http://localhost:8080>. Click click click.

### Deploy to cloudbuild

TODO
