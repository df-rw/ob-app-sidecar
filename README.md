# ob-app

This document describes how to setup a development environment for a web
application using:

- [Observable Framework](https://observablehq.com/framework) frontend;
- [Go](https://go.dev) web server backend;
- [htmx](https://htmx.org) for hypertext transactions;
- [nginx](https://nginx.org) for tying everything together.

## tl;dr

```shell
git clone https://github.com/df-rw/ob-app
cd ob-app
npm install
go run ./cmd/web/*.go -p 8082        # Start backend server (in one terminal).
npm run dev -- --port 8081 --no-open # Start Observable framework (diff terminal).
nginx -p . -c ./nginx-dev.conf       # Start nginx (diff terminal).
```

Open browser to http://localhost:8080. Click click click.

## Prerequisites

- Go (`brew install go` or [rtfm](https://go.dev/doc/install))
- nginx (`brew install nginx` or [rtfm](https://nginx.org/en/docs/install.html))

## Optional

- [air](https://github.com/air-verse/air) for live-reloading the Go backend
  application:

```shell
go install github.com/air-verse/air@latest
```

Configuration for `air` is in the repo as `./.air.toml`. Replace the `go run
./cmd/web/*.go -p 8082` line with `air` to get live reloading when a backend
file changes.

## Why is nginx in there?

Observable Framework uses it's own web server to not only serve itself to the
client, but also provide for hot reloading when frontend content or backend
data changes. This occurs only in development; with production deploys,
Observable Framework builds a complete static site.

Writing general purpose applications with Observable Framework on localhost can
be a little kludgy, as the frontend must be able to make calls to the backend.
Since the backend server for the application lives outside of Observable
Framework server, the calls would have to be CORS and the application would
need to be CORS aware.

It's possible to workaround this by passing environment variables to Observable
Framework's pages that rewrite URLs for backend calls. However this ends up
being messy when integrating with other components. For instance: Observable
Framework doesn't rewrite environment variables inside DOM elements, so tools
like htmx will lose their neat DOM syntax, and have to be constructed within a
JavaScript code block at runtime. This is messy for developers and inefficient
when run.

We can also potentially end up with the development environment not having the
same code paths as production due to the environment wrangling, making problem
tracking more difficult.

In the development environment, nginx is used as a proxy between the client
(browser) and both the Observable Framework server and the application server:

```
----------      ---------------
| client | ---> | nginx proxy |
----------      ---------------
                     |    |      ----------------------
                     |    -----> | application server |
                     |           ----------------------
                     |           -------------------------------
                     ----------> | Observable Framework server |
                                 -------------------------------
```

Production layout will depend on how you wish to carve up your application.

## How this all works

- Frontend is a Observable Framework application.
- Backend is a Go application, that serves a couple of routes hanging off `/api`/.
- nginx proxies requests from the client browser off to the appropriate backend
  based on URL (see
  [`location`](https://nginx.org/en/docs/http/ngx_http_core_module.html#location)).

## Notes

- Check the ports you run your servers on. The defaults, as specified in
  `nginx-dev.conf` are:
  - `8080` for the nginx proxy;
  - `8081` for the Observable Framework server;
  - `8082` for the backend application.
- Requests from your client browser should come through
  `http://localhost:8080`. If nothing seems to be working correctly, check
  your browser URL and the console to make sure you're using the proxy.
- If you are only doing work on the front end, of course you don't need to go
  through the proxy.

## Production deploys

How the application is deployed to production will differ based on target. At a
minimum:

- `npm run build` will be needed to build the frontend Observable Framework
  website. This creates a static site under `./dist`.
- Either:
  - the backend server application will need to serve the contents of this
  directory in addition to handling API calls; or
  - the application is fronted by a proxy which will route client requests.

In this example, we'll front our application with `nginx`. The static site will
be hosted by `nginx`, and API requests are routed to the backend application:

```
                             --------------------------------
                             |  Observable Framework static |
----------      --------------- site mounted on the proxy   |
| client | ---> | nginx proxy |------------------------------
----------      ---------------
                     |           ----------------------
                     |---------> | application server |
                                 ----------------------
```

This layout is specified in `nginx-prod.conf`. To see this working:

```shell
# Stop the development Observable Framework server and the development nginx
# server. The application server can keep running.
# ^C ^C

# Build the Observable Framework application:
npm run build

# Run an nginx proxy with the prod configuration:
nginx -p . -c ./nginx-prod.conf

# Open your client browser to http://localhost:8080. Click click click.
```

## Todo

- Clean URLs in Observable Framework are turned off. This is so our sample
  application deployment is correct. Routing and clean URLs will depend on
  how the target webserver is setup.
