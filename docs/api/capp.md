### Capp Namespaced API

## API Endpoints

- **GET** `/v1/namespaces/{namespace}/capps`
  - **Description**: Get all capp of namespace.
  - **Path Parameter**:
    - `namespace` - The namespace of the capp.
  - **Query Params**:
    ```json
    {
      "lables": {
        "key": str,
        "key1": str
      }
    }
    ```
  - **Response**: Capp names or an error message.
    ```json
    {
       "capps": []str,
       "count": int
    }
    ```

- **GET** `/v1/namespaces/{namespace}/capps/{cappName}`
  - **Description**: Get all capps of namespace.
  - **Path Parameter**:
    - `namespace` - The namespace of the capp.
    - `cappName` - The capp name to fetch.
  - **Query Params**: 
    ```json
    {
      "lables": {
        "key": str,
        "key1": str
      }
    }
    ```
  - **Response**: Capp info or an error message.

```json
{
  "capp": {
    "metadata": {
      "name": "string",  // max length 53 char
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
      // [CappSpec] (https://github.com/dana-team/container-app-operator/blob/v0.2.0/api/v1alpha1/capp_types.go#L29
    }
  }
}
```

- **GET** `/v1/namespaces/{namespace}/capps`
  - **Description**: Get all capps of namespace.
  - **Path Parameter**:
    -  `namespace` - The namespace of the capp.
  - **Query Params**:
    ```json
    {
      "lables": {
        "key": str,
        "key1": str
      }
    }
    ```
  - **Response**: Confirmation of deletion or an error message.
    ```json
    {
       "capps": []Capp,
       "count": int
    }
    ```

- **POST** `/v1/namespaces/{namespace}/capps`
  - **Description**: Create capp in a namespace.
  - **Path Parameter**:
    - `namespace` - The namespace of the capp.
  - **Body**:

```json
{
  "capp": {
    "metadata": {
      "name": "string",  // max length 53 char
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
      // [CappSpec] (https://github.com/dana-team/container-app-operator/blob/v0.2.0/api/v1alpha1/capp_types.go#L29
    }
  }
}
```
    ```
  - **Response**: Confirmation of creation or an error message.
```json
{
  "capp": {
    "metadata": {
      "name": "string",  // max length 53 char
      "namespace": "string"
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
      // [CappSpec] (https://github.com/dana-team/container-app-operator/blob/v0.2.0/api/v1alpha1/capp_types.go#L29
    },
    "status": {
      // [CappStatus] (https://github.com/dana-team/container-app-operator/blob/v0.2.0/api/v1alpha1/capp_types.go#L202
    }
  }
}
```

- **PUT** `/v1/namespaces/{namespace}/capps/{capa_name}`
  - **Description**: Update capp in a namespace.
  - **Path Parameter**:
    - `namespace` - The namespace of the capp.
    - `capp_name` - The name of the capp you want to update.
  - **Body**:

```json
{
  "capp": {
    "metadata": {
      "name": "string"   // max length 53 char
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
      // [CappSpec] (https://github.com/dana-team/container-app-operator/blob/v0.2.0/api/v1alpha1/capp_types.go#L29
    }
  }
}
```

- **Response**: Confirmation of update or an error message.
```json
{
  "capp": {
    "metadata": {
      "name": "string",  // max length 53 char
      "namespace": "string"
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
      // [CappSpec] (https://github.com/dana-team/container-app-operator/blob/v0.2.0/api/v1alpha1/capp_types.go#L29
    },
    "status": {
      // [CappStatus] (https://github.com/dana-team/container-app-operator/blob/v0.2.0/api/v1alpha1/capp_types.go#L202
    }
  }
}
```

- **DELETE** `/v1/namespaces/{namespace}/capps/{cappName}`
  - **Description**: Get all capps of namespace.
  - **Path Parameter**: 
    - `namespace` - The namespace of the capp.
    - `cappName` - The capp name to fetch.
  - **Response**: Confirmation of deletion an error message.
    ```json
    {
       "message": "string",
    }
    ```
