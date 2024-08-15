# platform-backend

![Version: 0.0.0](https://img.shields.io/badge/Version-0.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

A Helm chart for platform-backend

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| config.cluster | object | `{"apiPort":6443,"domain":"domain-test.com","name":"cluster-test"}` | Configuration relating to the cluster where the backend is deployed |
| config.cluster.apiPort | int | `6443` | Port of the API Server of the cluster |
| config.cluster.domain | string | `"domain-test.com"` | Domain of the cluster where the code is deployed |
| config.cluster.name | string | `"cluster-test"` | Cluster name where the code is deployed |
| config.defaultPaginationLimit | int | `100` | Default pagination limit |
| config.insecureSkipVerify | bool | `true` | Flag to indicate whether to skip HTTPS verification |
| config.kubeClientID | string | `"openshift-challenging-client"` | The kube client ID to use |
| config.name | string | `"config"` | Name of the ConfigMap where authentication endpoints are stored |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"Always"` | The pull policy for the image. |
| image.repository | string | `"ghcr.io/dana-team/platform-backend"` | The repository of the manager container image. |
| image.tag | string | `""` | The tag of the manager container image. |
| livenessProbe.initialDelaySeconds | int | `5` |  |
| livenessProbe.path | string | `"/healthz"` |  |
| livenessProbe.periodSeconds | int | `10` |  |
| livenessProbe.port | int | `8080` |  |
| nameOverride | string | `""` |  |
| readinessProbe | object | `{"initialDelaySeconds":5,"periodSeconds":10,"port":8080}` | Readiness and Liveness Probes Configuration |
| scaleMetric | string | `"concurrency"` | Name of the scale metric to use for Capp |

