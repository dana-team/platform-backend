# Kubernetes Pod Containers API

## Overview

This API allowes retrieving the containers within a specified pod in a given namespace.

### capp

- **GET** `/v1/namespaces/{namespace}/capps/{cappName}/pods`
  - **Description**: Retrieve a list of all pods associated with a specific capp in the given namespace.
  - **Path Parameter**:
    - `namespace` - The namespace of the capp. This is a required parameter that scopes the request to a specific namespace.
    - `cappName` - The name of the capp for which to retrieve the pods.
  - **Query Params**:
    - `limit`: (optional) Specifies the maximum number of pods to return per page. Defaults to 9.
    - `continue`: (optional) Used for fetching the next set of results.
    - `labelSelector`: (optional) Used for filtering by labels.
  - **Response**: A JSON object containing the list of pods for the specified capp or an error message if the request fails.
    ```json
    {
       "pods": [{
                    "name": "string"   
                }, ...],
       "count": int
    }
    ```

- **GET** `/v1/namespaces/{namespace}/pods/{podName}/containers`
  - **Description**: Retrieve a list of containers within a specific pod. This endpoint provides details about the containers running inside the specified pod.
  - **Path Parameter**:
    - `namespace` - The namespace of the pod.
    - `podName` -  The name of the pod whose containers are being listed.
  - **Response**: A JSON object containing the list of containers within the specified pod or an error message if the request fails.
    ```json
    {
       "containers": [{
                    "name": "string"   
                }, ...],
       "count": int
    }
    ```
