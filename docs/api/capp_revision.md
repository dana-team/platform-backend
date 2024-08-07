# CappRevision Management API

This document outlines the CRUD (Create, Read, Update, Delete) on CappRevision. Capp revision is represents version of capp.

## API Endpoints

### Capp Revision

- **GET** `/v1/namespaces/{namespace}/cappRevisions`
  - **Description**: Get all cappRevision of namespace.
  - **Query Params**:
    - `limit`: (optional) Specifies the maximum number of namespaces to return per page.
    - `page`: (optional) Used for setting the current pge.
  - **Response**: CappRevision names or an error message.
    ```json
    {
       "cappRevisions": []str,
       "count": int
    }
    ```

- **GET** `/v1/namepaces/{namespace}/cappRevisions/{cappRevisionName}`
  - **Description**: Gets a specific cappRevision in a namespace.
  - **Response**: CappRevision info or an error message.
    ```json
    {
      "cappRevision": {
        "metadata": {
          "name": "string",
          "namespace": "string",
          "creationTimestamp": "string"
        },
        "annotations": [
          {
            "key": "string",
            // Key
            "value": "string"
            // Value
          },
          {
            "key": "string",
            "value": "string"
          }
        ],
        "labels": [
          {
            "key": "string",
            // Key
            "value": "string"
            // Value
          },
          {
            "key": "string",
            "value": "string"
          }
        ],
        "spec": {
          // [CappRevisionSpec] https://github.com/dana-team/container-app-operator/blob/main/api/v1alpha1/capprevision_types.go#L23-L30
        },
        "status": {
          // [CappRevisionSpec] https://github.com/dana-team/container-app-operator/blob/main/api/v1alpha1/capprevision_types.go#L32-L34
        }
      }
    }
    ```