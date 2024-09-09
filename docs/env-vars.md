# Adding environment variables

We can add environment variables for consumption by our application. Here we
outline how to add new variables to the ingress application our sidecar setup.
Steps are the same for addingn to any other sidecar; just use the appropriate
build section to receive the environment variable.

In the steps below `MY_ENV_VARIABLE` is the name of your new environment
variable.

## Local development environment

- Add `MY_ENV_VARIABLE` to your `.env`.
- Update your code and make sure it acts how you like.
- Add your new environment variable to `.env-sample` and make any notes
  that may be helpful by your peers.

## Update ./README.md

Add notes on your environment to `README.md` in the section `.env`.

## Docker environment

In `Dockerfile.ingress-local`:

- Add an `ARG` named for your environment variable in `stage-build-of-app`:
 
  ```
  FROM node:22-alpine AS stage-build-of-app
  ...
  ARG MY_ENV_VARIABLE
  ...
  ```

In `compose.yaml`

- Add your environment variable to `args` in the `services` section that references
  the same name:

  ```
  services:
    ingress:
      build:
      ...
        args:
          MY_ENV_VARIABLE: ${MY_ENV_VARIABLE}
      ...
  ```

## GCP environment

In `Dockerfile.ingress-gcp`:

- Add an `ARG` named for your environment variable in `stage-build-of-app`:

  ```
  FROM node:22-alpine AS stage-build-of-app
  ...
  ARG MY_ENV_VARIABLE
  ...
  ```

In `cloudbuild.yaml`:

- Add a user-defined substitution of your environment variable to
  `substitutions` with it's default value:

  ```
  substitutions:
    ...
    _MY_ENV_VARIABLE: "this is the default value"
    ...
  ```

- Add `--build-arg` to step `BUILD_INGRESS` in `steps` that references your
  environment variable and it's substitution:

  ```
  steps:
  - name: 'gcr.io/cloud-builders/docker'
    id: BUILD_INGRESS
    ...
    '--build-arg', 'MY_ENV_VARIABLE=${_MY_ENV_VARIABLE}',
    ...
  ```

