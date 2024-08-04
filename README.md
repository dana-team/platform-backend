# Platform Backend

A Go-based API server built using the `Gin` framework. It serves as the backend for a serverless platform, simplifying the consumption of the `Capp` (a Kubernetes Custom Resource).

It is part of the [`rcs-ocm-deployer`](https://github.com/dana-team/rcs-ocm-deployer) and [`container-app-operator`](https://github.com/dana-team/container-app-operator) ecosystem, and is a complementary project to the [`platform-frontend`](https://github.com/dana-team/platform-frontend) `React.js` application.

## Documentation

The API is documented in the [docs/api](./docs/api) directory of this repository. Refer to:

- [ContainerApp API](./docs/api/capp.md)
- [ContainerApp Revisions API](./docs/api/capp_revision.md)
- [Containers API](./docs/api/containers.md)
- [Namespace API](./docs/api/namespace.md)
- [Secrets API](./docs/api/secrets.md)
- [Users API](./docs/api/users.md)
- [Token API](./docs/api/token.md)

## Quickstart

The backend uses `OpenShift`'s built-in `OAuth` Server to authenticate users. To run it, configure environment variables in a new `.env` file:

```bash
# Define the allowed origin regex pattern for WebSocket connections from the frontend. If not specified, all origins will be approved.
ALLOWED_ORIGIN_REGEX="http:localhost:8080|https:example.com.*"

# Cluster configuration
CLUSTER_NAME=ocp-rcs-example    # Change to match your own
CLUSTER_DOMAIN=dana.com         # Change to match your own
KUBE_AUTH_BASE_URL="https://oauth-openshift.apps.${CLUSTER_NAME}.${CLUSTER_DOMAIN}"
KUBE_API_BASE_URL="https://api.${CLUSTER_NAME}.${CLUSTER_DOMAIN}:6443"

INSECURE_SKIP_VERIFY=true
KUBE_CLIENT_ID="openshift-challenging-client"
KUBE_AUTH_URL="${KUBE_AUTH_BASE_URL}/oauth/authorize"
KUBE_TOKEN_URL="${KUBE_AUTH_BASE_URL}/oauth/token"
KUBE_USERINFO_URL="${KUBE_API_BASE_URL}/apis/user.openshift.io/v1/users/~"
KUBE_API_SERVER="${KUBE_API_BASE_URL}"

cat <<EOF > .env
INSECURE_SKIP_VERIFY=${INSECURE_SKIP_VERIFY}
KUBE_CLIENT_ID=${KUBE_CLIENT_ID}
KUBE_AUTH_URL=${KUBE_AUTH_URL}
KUBE_TOKEN_URL=${KUBE_TOKEN_URL}
KUBE_USERINFO_URL=${KUBE_USERINFO_URL}
KUBE_API_SERVER=${KUBE_API_SERVER}
ALLOWED_ORIGIN_REGEX=${ALLOWED_ORIGIN_REGEX}
EOF
```

Then, use the following `Makefile` target in order to run the backend:

```bash
$ make run
```

Alternatively, to build the backend as a Docker image, use the following `Makefile` targets:

```bash
$ make docker-build docker-push IMG=<registry>/platform-backend:<tag>
```

## Routes Examples

### Login Example

#### Request

```bash
$ curl -X POST \
-H "Content-Type: application/json" \
-H "Authorization: Basic <base64('username:password')>" \
"localhost:8080/v1/login/"
```

#### Response

```bash
{
  "token": "<TOKEN_VALUE>"
}
```

### Get All Namespaces Example

#### Request

```bash
$ curl -H "Authorization: Bearer <TOKEN_VALUE>" \
"localhost:8080/v1/namespaces/"
```

#### Response

```bash
{
  "namespaces": [
    {
      "name": "project1"
    },
    {
      "name": "project2"
    }
  ],
  "count": 2
}
```

## Testing
To ensure your backend is working correctly, you can run end-to-end (e2e) tests.

### Running End-to-End Tests
To run e2e tests, use the test-e2e target defined in the Makefile. You can pass the platformUrl flag to specify the URL of the platform backend.

1. Run Tests with Specific Platform URL:
    ```bash
    make test-e2e platformUrl=<your-platform-url>
    ```
    Replace <your-platform-url> with the actual URL of the platform backend.
    Example:
    If your platform URL is http://localhost:8080, you would run:
    ```bash
    make test-e2e platformUrl=http://localhost:8080
    ```
2. Run Tests with Default URL:
   ```bash
   make test-e2e
   ```
   If you leave the platformUrl flag empty, the tests will automatically use the default URL of the backend deployed in OpenShift. 

