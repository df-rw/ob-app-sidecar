# Enabling IAP on a Google Cloud Run service

Google Cloud Identity-Aware Proxy (IAP) provides simple authentication to
various Google services, including Google Cloud Run. Here we show the steps
required in setting up IAP to front a Google Cloud Run service.

More information about IAP is available in the [Identity-Aware Proxy
overview](https://cloud.google.com/iap/docs/concepts-overview).

## Overview

Adding IAP to a Google Cloud Run application involves a number of steps across
multiple Google Cloud services, and it can be a little hard to keep track of.
Here we show the minimum steps required and in the correct order to get IAP up
and running in front of Google Cloud Run.

The steps below are done using the `gcloud` CLI for brevity. All commands have
console equivalents, but are easier to explain in the context of the CLI.

Majority of this is taken from a [Google Codelabs
exercise](https://codelabs.developers.google.com/secure-serverless-application-with-identity-aware-proxy#0)
intermixed with other documentation, notes from previous attempts and some
experience. See [references](#references) for some links.

## Pre-requisites

A Google Cloud Run service that is up and running in the correct Google
project. To see what you've got:

```
gcloud config get project       # your currently active project
gcloud projects list            # projects you have access to
gcloud config set project <...> # set your current project
gcloud run services list        # see the Google Cloud Run services in this project
```

Make sure you have the correct project set, and you know the right service
that you want to front with IAP.

## Setup

Setting up some base environment variables will make it easier to construct
our `gcloud` commands. These are the base variables; make sure they are set
correctly in your environment:

```
export REGION=australia-southeast1
export NETWORK_TIER=PREMIUM
export PRINCIPAL_DOMAIN=domain.that.users.belong.to
export SERVICE_NAME=ob-app-sidecar
export SERVICE_TITLE="Observable Framework and app backend in sidecar conf"
```

Setup some derived variables - safe to copypasta this:

```
export PROJECT_ID=$(gcloud config get-value project)
export PROJECT_NUMBER=$(gcloud projects describe ${PROJECT_ID} --format='value(projectNumber)')
export USER_EMAIL=$(gcloud config list account --format "value(core.account)")
export IAP_NEG_NAME=${SERVICE_NAME}-iap-neg
export IAP_BACKEND_NAME=${SERVICE_NAME}-iap-backend
export IAP_URL_MAP_NAME=${SERVICE_NAME}-iap-url-map
export IAP_IP_NAME=${SERVICE_NAME}-iap-ip
export IAP_CERT_NAME=${SERVICE_NAME}-iap-cert
export IAP_HTTP_PROXY_NAME=${SERVICE_NAME}-iap-http-proxy
export IAP_FORWARDING_RULE=${SERVICE_NAME}-iap-forwarding-rule
export IAP_OAUTH_DISPLAY_NAME=${SERVICE_NAME}-oauth-client
```

Double check everything:

```
echo "\$REGION                : ${REGION}"
echo "\$NETWORK_TIER          : ${NETWORK_TIER}"
echo "\$PRINCIPAL_DOMAIN      : ${PRINCIPAL_DOMAIN}"
echo "\$SERVICE_NAME          : ${SERVICE_NAME}"
echo "\$SERVICE_TITLE         : ${SERVICE_TITLE}"
echo "\$PROJECT_ID            : ${PROJECT_ID}"
echo "\$PROJECT_NUMBER        : ${PROJECT_NUMBER}"
echo "\$USER_EMAIL            : ${USER_EMAIL}"
echo "\$IAP_NEG_NAME          : ${IAP_NEG_NAME}"
echo "\$IAP_BACKEND_NAME      : ${IAP_BACKEND_NAME}"
echo "\$IAP_URL_MAP_NAME      : ${IAP_URL_MAP_NAME}"
echo "\$IAP_IP_NAME           : ${IAP_IP_NAME}"
echo "\$IAP_CERT_NAME         : ${IAP_CERT_NAME}"
echo "\$IAP_HTTP_PROXY_NAME   : ${IAP_HTTP_PROXY_NAME}"
echo "\$IAP_FORWARDING_RULE   : ${IAP_FORWARDING_RULE}"
echo "\$IAP_OAUTH_DISPLAY_NAME: ${IAP_OAUTH_DISPLAY_NAME}"
```

## Enable APIs

By default, minimal APIs are enabled on Google Cloud accounts. Turn on the APIs
required:

```
gcloud services enable \
    iap.googleapis.com \
    cloudresourcemanager.googleapis.com \
    cloudidentity.googleapis.com \
    compute.googleapis.com
```

| Service | Purpose | Reference |
|---|---|---|
| `iap.googleapis.com` | IAP API. | [Docs](https://cloud.google.com/iap/docs/reference/rest) |
| `cloudresourcemanager.googleapis.com` | Update metadata for GCP resource containers. | [Docs](https://cloud.google.com/resource-manager/reference/rest) |
| `cloudidentity.googleapis.com` | Manage identity resources. | [Docs](https://cloud.google.com/identity/docs/reference/rest) |
| `compute.googleapis.com` | Manage VMs on GCP. | [Docs](https://cloud.google.com/compute/docs/reference/rest/v1) |

## Create Network Endpoint Group

```
gcloud compute network-endpoint-groups create ${IAP_NEG_NAME} \
    --project ${PROJECT_ID} \
    --region=${REGION} \
    --network-endpoint-type=serverless  \
    --cloud-run-service=${SERVICE_NAME}
```

To verify:

```
gcloud compute network-endpoint-groups describe --region=${REGION} ${IAP_NEG_NAME}
```

In console:

* https://console.cloud.google.com/compute/networkendpointgroups/list?project=${PROJECT_ID}

## Create backend, add the NEG

```
gcloud compute backend-services create ${IAP_BACKEND_NAME} \
    --global 
gcloud compute backend-services add-backend ${IAP_BACKEND_NAME} \
    --global \
    --network-endpoint-group=${IAP_NEG_NAME} \
    --network-endpoint-group-region=${REGION}
```

To verify:

```
gcloud compute backend-services list
gcloud compute backend-services describe --global ${IAP_BACKEND_NAME}
```

In console:

* https://console.cloud.google.com/net-services/loadbalancing/list/backends?project=${PROJECT_ID}

## Create URL map

```
gcloud compute url-maps create ${IAP_URL_MAP_NAME} \
    --default-service ${IAP_BACKEND_NAME}
```

To verify:

```
gcloud compute url-maps list
gcloud compute url-maps describe ${IAP_URL_MAP_NAME}
```

In console:

* https://console.cloud.google.com/net-services/loadbalancing/list/loadBalancers?project=${PROJECT_ID}

## Reserve static IP address

The [load balancer requires an IP address](https://cloud.google.com/load-balancing/docs/https/setup-global-ext-https-serverless#ip-address):

```
gcloud compute addresses create ${IAP_IP_NAME} \
    --network-tier=${NETWORK_TIER} \
    --ip-version=IPV4 \
    --global
```

To verify:

```
gcloud compute addresses list
gcloud compute addresses describe --global ${IAP_IP_NAME}
```

In console:

* https://console.cloud.google.com/networking/addresses/list?project=${PROJECT_ID}

| Flag | Purpose | Reference |
|---|---|---|
| --network-tier | RTFM; use `PREMIUM` here | [Docs](https://cloud.google.com/sdk/gcloud/reference/compute/addresses/create#--network-tier) |

## Get the domain name

We'll need a hostname to reference the application (through the load balancer).
We also have to setup a certificate which requires the hostname. For
demonstration purposes, we will use <nip.io> to make us a name that we can use
for testing that is registered on DNS. IRL, contact your domain manager and
grab a proper domain.

```
export DOMAIN=$(gcloud compute addresses list --filter ${IAP_IP_NAME} --format='value(ADDRESS)').nip.io
```

## Provision a certificate

Google [recommends a Google-managed
certificate](https://cloud.google.com/load-balancing/docs/https/setup-global-ext-https-serverless#ssl_certificate_resource)
as management of the cert is done automatically.

```
gcloud compute ssl-certificates create ${IAP_CERT_NAME} \
    --description=${IAP_CERT_NAME} \
    --domains=${DOMAIN} \
    --global
```

Provisioning the cert may [take up to 60
minutes](https://cloud.google.com/load-balancing/docs/ssl-certificates/google-managed-certs).
To periodically check it's status:

```
gcloud compute ssl-certificates list
gcloud compute ssl-certificates describe --global ${IAP_CERT_NAME}
```

Some Google documentation, including the codelabs example, states that you need
to wait for the certificate state to move from [`PROVISIONING` to `ACTIVE`
before
proceeding](https://codelabs.developers.google.com/secure-serverless-application-with-identity-aware-proxy#3).
Experience has shown that `PROVISIONING` will transition to
`FAILED_NOT_VISIBLE` and stay there. The [troubleshooting
documentation](https://cloud.google.com/load-balancing/docs/ssl-certificates/troubleshooting#domain-status)
states that one reason for this is because the certificate isn't attached to a
load balancer proxy.

Experience has shown that it is safe to continue without waiting for the
certificate to be `ACTIVE`. YMMV. One approach which has worked more than once
is to wait for the certificate to move to `FAILED_NOT_VISIBLE`, then attach to
the load balancer and wait for it to move to `ACTIVE` which doesn't take long.

Whatever your approach, feel free to check certificate state with `gcloud
compute ssl-certificates ...`.

## Create the proxy

```
gcloud compute target-https-proxies create ${IAP_HTTP_PROXY_NAME} \
    --ssl-certificates ${IAP_CERT_NAME} \
    --url-map ${IAP_URL_MAP_NAME}
```

To verify:

```
gcloud compute target-https-proxies list
gcloud compute target-https-proxies describe ${IAP_HTTP_PROXY_NAME}
```

## Setup forwarding rules

```
gcloud compute forwarding-rules create ${IAP_FORWARDING_RULE} \
    --load-balancing-scheme=EXTERNAL \
    --network-tier=${NETWORK_TIER} \
    --address=${IAP_IP_NAME} \
    --global \
    --ports=443 \
    --target-https-proxy ${IAP_HTTP_PROXY_NAME}
```

To verify:

```
gcloud compute forwarding-rules list
gcloud compute forwarding-rules describe --global ${IAP_FORWARDING_RULE}
```

## Inbound only through the LB

```
gcloud run services update ${SERVICE_NAME} \
    --ingress internal-and-cloud-load-balancing \
    --region ${REGION}
```

To verify:

```
gcloud run services list
gcloud run services describe ${SERVICE_NAME}
```

## Configure OAuth consent screen

Create a ["Cloud OAuth
brand"](https://cloud.google.com/sdk/gcloud/reference/alpha/iap/oauth-brands/create).
This will possibly / probably require review and alteration after your project
is up and running:

```
gcloud iap oauth-brands create \
    --application_title="${SERVICE_TITLE}" \
    --support_email=${USER_EMAIL}
```

To verify:

```
gcloud iap oauth-brands list
```

## Create an IAP Oauth Client

Create the IAP Oauth client. Grab the bits we need to enable IAP on the
backend:

```
gcloud iap oauth-clients create \
    projects/${PROJECT_ID}/brands/${PROJECT_NUMBER} \
    --display_name=${IAP_OAUTH_DISPLAY_NAME}

export CLIENT_NAME=$(gcloud iap oauth-clients list \
    projects/${PROJECT_ID}/brands/${PROJECT_NUMBER} --format='value(name)' \
    --filter="displayName:${IAP_OAUTH_DISPLAY_NAME}")

export CLIENT_ID=${CLIENT_NAME##*/}

export CLIENT_SECRET=$(gcloud iap oauth-clients describe $CLIENT_NAME --format='value(secret)')
```

## Enable IAP on the backend

```
gcloud iap web enable --resource-type=backend-services \
    --oauth2-client-id=${CLIENT_ID} \
    --oauth2-client-secret=${CLIENT_SECRET} \
    --service=${IAP_BACKEND_NAME}
```

## Add the principal's domain

```
gcloud iap web add-iam-policy-binding \
    --resource-type=backend-services \
    --service=${IAP_BACKEND_NAME} \
    --member=domain:${PRINCIPAL_DOMAIN} \
    --role="roles/iap.httpsResourceAccessor"
```

## Provision the IAP account

As per: https://cloud.google.com/iap/docs/enabling-cloud-run#troubleshooting_errors

```
gcloud beta services identity create \
    --service=iap.googleapis.com \
    --project=${PROJECT_ID}
```

## I think that's it

At this point, we should have a load balancer fronting a Google Cloud Run
service and using IAP for authentication. There are a few more steps from the
[codelabs
example](https://codelabs.developers.google.com/secure-serverless-application-with-identity-aware-proxy#0)
that require some attention:

- [Setting the publishing status to
  external](https://codelabs.developers.google.com/secure-serverless-application-with-identity-aware-proxy#4).
  We are only handling internal traffic so this doesn't need to be made `EXTERNAL`.

## Notes

### Why isn't this in a script / Terraform / Ansible / whatever?

To avoid blind copypasta. I might / probably have something wrong in the above.
Might be a good idea to understand what each command does before you commit. By
all means, feel free to put this into something more easily digestible.

<a name="references"></a>
## References

- https://cloud.google.com/iap/docs/concepts-overview
- https://codelabs.developers.google.com/secure-serverless-application-with-identity-aware-proxy#0
- https://cloud.google.com/iap/docs/enabling-cloud-run
- https://gist.github.com/df-rw/7af214d3ee37d5fcec7d14c0b4763871
- https://gist.github.com/df-rw/62168720124236dd61c304cd7ca8cb32
