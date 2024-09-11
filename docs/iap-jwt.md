# Checking IAP JWT

[Google Clouds Identity-Aware Proxy
(IAP)](https://cloud.google.com/iap/docs/concepts-overview) provides a JWT in
the request header `x-goog-iap-jwt-assertion`. Here we show how to [verify this
JWT](https://cloud.google.com/iap/docs/signed-headers-howto).

## Ensure application is running behind IAP

Get your application deployed without the explicit IAP JWT checking, and
[setup and turn on IAP](./iap.md).

## Find the audience string

Part of JWT verification includes making sure the audience (`aud`)
field from the token is what we expect. You can find the audience
string in the Google Cloud console:

- `Navigation menu` &gt; `Security` &gt; `Identity-Aware Proxy`
- Find your backend service in the list, click the menu at the end of the row.
- Click `Get JWT audience code`
- Copy the value in the text input.

## Update the build trigger

There is a substition variable ready to go in `cloudbuild.yaml` called
`_GCP_IAP_JWT_AUDIENCE`. You just need to override this in the Cloud Build
Trigger settings. In the Google Cloud console:

- `Navigation menu` &gt; `Cloud Build` &gt; `Triggers`
- Click on the trigger for your project.
- Go to `Advanced` &gt; Substitution variables. Add the variable
  `_GCP_IAP_JWT_AUDIENCE` and paste the audience string copied previously.
- Save this.

## Add your validation code

There is a sample IAP JWT validator application in `cmd/cli/validator-iap.go`
which is based off [this
reference](https://cloud.google.com/iap/docs/signed-headers-howto#retrieving_the_user_identity).
Feel free to use this or roll your own.

## Rebuild, deploy and test

The [documentation has some notes on
testing](https://cloud.google.com/iap/docs/signed-headers-howto#validation_testing).

## References

- https://cloud.google.com/iap/docs/signed-headers-howto
- https://cloud.google.com/iap/docs/signed-headers-howto#retrieving_the_user_identity
