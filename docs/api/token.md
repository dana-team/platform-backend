# Token Management API

- **POST** `/v1/namespaces/{namespace}/serviceaccounts`
  - **Description**: Creates a new service account.
  - **Path Parameter**:
    - `namespace` - The namespace of the service account.
  - **Body**:
    ```json
    {
      "serviceAccountName": "string"
    }
    ```
  - **Response**: Confirmation of creation or an error message.
    ```json
    {
      "serviceAccountName": "string",
    }
    ```

- **GET** `/v1/namespaces/{namespace}/token`
  - **Description**: Gets a new token from the service account.
  - **Path Parameter**:
    - `namespace` - The namespace of the service account to create token from.
  - **Body**:
    ```json
    {
      "serviceAccountName": "string"
    }
  - **Response**: Token from the created service account.
    ```json
    {
      "token": "string",
    }
    ```
