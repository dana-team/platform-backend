# Container Logs API

## Overview
This API provides endpoints to fetch logs from individual pods or from Capps.

## API Endpoints

### GET `/v1/logs/pod/{namespace}/{name}`
  - **Description**: Fetches logs from the specified pod.
  - **Path Parameters**:
    - `namespace` - The namespace of the pod.
    - `name` - The name of the pod.
  - **Query Parameters**:
    - `container` - The name of the container to fetch logs from. Defaults to the pod name if not specified.

### GET `/v1/logs/capp/{namespace}/{name}`
  - **Description**: Fetches logs from the specified Container Application (Capp).
  - **Path Parameters**:
    - `namespace` - The namespace of the Capp.
    - `name` - The name of the Capp.
  - **Query Parameters**:
    - `container` - The name of the container to fetch logs from. Defaults to the Capp name if not specified.
    - `podName` - The name of the pod to fetch logs from. Defaults to the first pod in the list if not specified.