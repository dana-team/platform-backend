
# Kubernetes Pod Containers API

## Overview

This API provides access to Kubernetes pod information, specifically retrieving the names of containers within a specified pod in a given namespace. It is intended for use in environments where Kubernetes pod management is essential. This document details the usage of the `/{namespace}/{pod_name}/containers` endpoint.

## API Endpoint

### GET `/v1/{namespace}/{capp_name}/{pod_name}/containers`

This endpoint retrieves a list of container names within the specified pod in the specified namespace and the count of these containers.

#### Path Parameters

##### Required

- **namespace**: The Kubernetes namespace in which the pod resides.
- **capp_name**: The capp name is required for us to know which cluster to query.
- **pod_name**: The name of the pod from which to retrieve container names. This string should accurately match the pod name within the specified namespace.

#### Response

The response is a JSON object that contains a list of container names and the total number of containers within the specified pod.

##### Response Object

```json
{
  "containerNames": ["container1", "container2", ...],
  "count": 2
}
```

##### Response Fields

- **containerNames**: An array of strings, each representing a container name within the specified pod.
- **count**: An integer representing the total number of containers listed in the response.

#### Response Status Codes

- **200 OK**: The request was successful, and the container names have been returned.
- **400 Bad Request**: The request was malformed. Please check the `namespace` and `pod_name` parameters.
- **404 Not Found**: The specified pod was not found in the specified namespace. Please verify the `namespace` and `pod_name` and try again.
- **500 Internal Server Error**: An error occurred on the server. Please contact the system administrator.

## Examples

### Request

```
GET /v1/my-namespace/sample-capp/sample-pod/containers
```

### Successful Response

```json
{
  "containerNames": ["nginx", "redis"],
  "count": 2
}
```
