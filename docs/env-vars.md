# Adding environment variables

We can add environment variables for consumption by any of our sidecars. Here
we outline how to add new variables to the validator application for use at
**run** time.

For an example of adding an environment variable for use at **build** time
rather than run time, look for `OBSERVABLE_TELEMETRY_DISABLE`. Note that with
[interpolation](https://docs.docker.com/compose/environment-variables/set-environment-variables/#additional-information),
`.env` in the current directory will be read automatically.

In the steps below `MY_ENV_VARIABLE` is the name of your new environment
variable.

## Local development environment

- Add `MY_ENV_VARIABLE` to your `.env`.
- Update your code and make sure it acts how you like.
- Add your new environment variable to `.env-sample` and make any notes
  in there that may be useful.

## Update ./README.md

Add notes on your new environment variable to `README.md` in the section
`.env`.

## Docker environment

In `compose.yaml`

- Add your environment variable to `environment` in the validator service
  section:

  ```
  services:
    validator:
      environment:
        - MY_ENV_VARIABLE=${MY_ENV_VARIABLE}
  ```

## GCP environment

In `cloudbuild.yaml`:

- Add a user-defined substitution of your environment variable name to
  the `substitutions` section with its default value:

  ```
  substitutions:
    ...
    _MY_ENV_VARIABLE: "this is the default value"
    ...
  ```

- Since this is a run-time variable, we need to pass it along to `gcloud run
  services`. Add it to the `env` section of `SET_SERVICE`, and do a
  substitution in the `script` section as well:

  ```
  - name: 'alpine'
    id: SET_SERVICE
    env:
    ...
    - 'MY_ENV_VARIABLE=${_MY_ENV_VARIABLE}'
    script: |
      ...
      sed -i s@%GCP_JWT_AUDIENCE%@${GCP_JWT_AUDIENCE}@g run-service.yaml
      ...
  ```

In `run-service.yaml`:

- Add to the `env` section of the validator container:

  ```
  spec:
    template:
      spec:
        containers:
          - image: "%IMAGE_VALIDATOR%"
            ...
            env:
              - name: MY_ENV_VARIABLE
                value: "%MY_ENV_VARIABLE%"
  ```

For variables that shouldn't be committed, the values of substitution variables
can be overridden in the Google Cloud console in the Cloud Build trigger or the
Cloud Run settings.
