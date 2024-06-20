# ConfigMaps Management API

- **GET** `/v1/namespaces/{namespace}/{configMapName}`
  - **Description**: Gets data from the specified config map.
  - **Path Parameter**:
    - `namespace` - The namespace of the config map.
    - `configMapName` - The name of the config map.
  - **Response**: Data of the config map.
    ```json
    {
      "data": {
        "key": "string",
        "key1": "string"
      }
    }
    ```
