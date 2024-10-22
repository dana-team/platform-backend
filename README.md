# Platform Backend

A Go-based API server built using the `Gin` framework. It serves as the backend for a serverless platform, simplifying the consumption of `Capp` (a Kubernetes Custom Resource).

It is part of the [`rcs-ocm-deployer`](https://github.com/dana-team/rcs-ocm-deployer) and [`container-app-operator`](https://github.com/dana-team/container-app-operator) ecosystem.

## Documentation

The API is documented on the `/docs` route, for example at `http://localhost:8080/docs`. It uses the [`Huma` project](https://huma.rocks/) to generate the docs. 

## Set Up

In order to fully use the `platform-backend`, it is needed to have at least:

1. One OpenShift cluster with `Advanced Cluster Management` installed, which acts as the `Hub Cluster`.
2. One OpenShift cluster which is a `Managed Cluster` and is connected to the `Hub Cluster`.

### Configure the Managed Cluster

To install all the prerequisites needed for the `Managed Cluster` to run a fully-functioning `platform-backend` on it as a `Capp`, use the following `Makefile` target, which configures the `Capp` environment to work properly using an actual DNS Server:

```bash
$ make install-prereq PROVIDER_DNS_REALM=<value> PROVIDER_DNS_KDC=<value> PROVIDER_DNS_POLICY=<value> PROVIDER_DNS_NAMESERVER=<value> PROVIDER_DNS_USERNAME=<value> PROVIDER_DNS_PASSWORD=<value> CAPP_RELEASE=<value>
```

| Value Name              | Value Default                            | Explanation                                                                                                  |
|-------------------------|------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| PROVIDER_DNS_REALM      | `DANA-DEV.COM`                           | Defines the name of the Kerberos Realm to use in the provider.                                               |
| PROVIDER_DNS_KDC        | `dana-wdc-1.dana-dev.com`                | Defines the name of the Kerberos Key Distribution Center server.                                             |
| PROVIDER_DNS_POLICY     | `ClusterFirst`                           | Defines the `dnsPolicy` of the `provider-dns` deployment. If used then it should be set to `None`.           |
| PROVIDER_DNS_NAMESERVER | `8.8.8.8`     | The nameserver to use in the `dnsConfig` of the `provider-dns` deployment if `dnsPolicy` is set to `None`.   |
| PROVIDER_DNS_USERNAME   | `dana`                                   | Defines the username to connect to the KDC with.                                                             |
| PROVIDER_DNS_PASSWORD   | `passw0rd`                               | Defines the password to connect to the KDC with.     
| CAPP_RELEASE   | `main`                               | The image tag with which to deploy the `container-app-operator`.     

### Configure the Hub Cluster

To install all the needed prerequisites on the `Hub Cluster`, use the following `Makefile` target:

```bash
$ make setup-hub PLACEMENT_NAME=<value> PLACEMENTS_NAMESPACE=<value> MANAGED_CLUSTER_NAME=<value>
```

### Deploy the backend

Use the provided Helm Chart in this repository in order to deploy the `platform-backend` to an OpenShift cluster using a `Capp` CR.

```bash
$ helm upgrade --install platform-backend --namespace platform-backend-system --create-namespace oci://ghcr.io/dana-team/helm-charts/platform-backend --version <release>
```

Alternatively, deploy the latest, non-released version using the Chart which is at `charts/platform-backend` on this repo. Use the `Makefile` target:

```bash
$ make deploy CLUSTER_NAME=<value> CLUSTER_DOMAIN=<value>
```

### Cleanup

To undeploy the `platform-backend`, run:

```bash
$ make undeploy
```

To cleanup the prerequisites from the `Hub Cluster`, run:

```bash
$ make cleanup-hub PLACEMENT_NAME=<value> PLACEMENTS_NAMESPACE=<value> MANAGED_CLUSTER_NAME=<value>
```

To cleanup the prerequisites from the `Managed Cluster`, run:

```bash
$ make uninstall-prereq
```

## Local Development

To run the backend locally, a `.env` file is used where several environment variables are set.

```bash
$ make env-file CLUSTER_NAME=<CLUSTER_NAME> CLUSTER_DOMAIN=<CLUSTER_DOMAIN>
```

For example:

```bash
$ make env-file CLUSTER_NAME=ocp-test CLUSTER_DOMAIN=dana.com
```

### Adding a new variable

To create a new environment variable, both locally and for the production deployment to use, create a new entry in the `charts/platform-backend/_config-data.yaml` file. Then, set its value in the `charts/platform-backend/values.yaml` file:

For example:

`config-data.yaml`:

```diff
INSECURE_SKIP_VERIFY: "{{ .Values.config.insecureSkipVerify }}"
+ NEW_VARIABLE: "{{ .Values.config.newVariable }}"
```

`values.yaml`:

```diff
config:
  # -- Name of the ConfigMap where authentication endpoints are stored
  name: config
+  # -- Indicates a value
+  newVariable: <VALUE>
```

### Running locally

Then, use the following `Makefile` target in order to run the backend:

```bash
$ make run
```

### Build an image

To build the backend as a Docker image, use the following `Makefile` targets:

```bash
$ make docker-build docker-push IMG=<registry>/platform-backend:<tag>
```

## E2E Testing

To ensure the backend works correctly, run end-to-end (`e2e`) tests.

To run `e2e` tests, use the `make test-e2e` target defined in the `Makefile`. You can pass the `PLATFORM_URL` flag to specify the URL of the platform backend.

### Run tests with a Specific Platform URL

```bash
$ make test-e2e PLATFORM_URL=<your-platform-url>
```

Replace <your-platform-url> with the actual URL of the platform backend.

Example:

If your platform URL is http://localhost:8080, you would run:

```bash
$ make test-e2e PLATFORM_URL=http://localhost:8080
```

If running in debug mode, add the `platformUrl` environment variable to your IDE.

### Run tests with a Default URL

```bash
$ make test-e2e
```

If you leave the `PLATFORM_URL` flag empty, the tests will automatically use the default URL of the backend deployed in `OpenShift`.