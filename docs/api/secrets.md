# Secret Management API

This document outlines the CRUD (Create, Read, Update, Delete) operations for managing secrets, including opaque secrets like API keys and passwords, TLS secrets which include certificates and private keys.

## API Endpoints

### Create a Secret

- **POST** `/v1/namespaces/{namespace}/secrets`
  - **Description**: Create a new secret, either TLS, or opaque.
  - **Request Body**: (For a TLS secret)
    ```json
    {
      "type":      "string", // 'tls'
      "name":      "string", // Name of the secret
      "cert":      "string", // Certificate content
      "key":       "string"  // Private key content
    }
    ```
  - **Request Body**: (For an opaque secret)
    ```json
    {
      "type": "string",  // 'opaque'
      "name": "string",  // Name of the secret
      "data": [
        {
          "key": "string",  // Key
          "value": "string" // Value
        },
        {
          "key": "string",
          "value": "string"
        }
      ]
    }
    ```
  - **Response**: Returns the created secret with an ID, or an error message.
    ```json
    {
      "namespace": "string",
      "type":      "string",
      "name":      "string",
    }
    ```

### Read a Secret

- **GET** `/v1/namespaces/{namespace}/secrets`
  - **Description**: Retrieve details of a specific secret, TLS, or opaque.
  - **Path Parameter**: `namespace` - The namespace of the secrets.
  - **Query Params**:
    - `limit`: (optional) Specifies the maximum number of namespaces to return per page. Defaults to 9.
    - `continue`: (optional) Used for fetching the next set of results.
  - **Response**: Returns the secret details or an error message if not found.
    ```json
    {
      "count": "int",
      "secrets": [{
          "name": "string",
          "type": "string",
          "namespace": "string"
      }]
    }
    ```
- **GET** `/v1/namespaces/{namespace}/secrets/{secretName}`
  - **Description**: Get an existing secret.
  - **Path Parameter**:
    - `secretName` - The name of the secret.
    - `namespace` - The namespace of the secret.
  - **Response**: Returns the updated secret details or an error message.
    ```json
    {
     "secret": 
      {
      "id": "string",
      "type": "string",
      "name": "string",
      "data": [
        {
          "key": "string",
          "value": "string"
        }
      ]
    }}
    ```

### Update a Secret

- **PUT** `/v1/namespaces/{namespace}/secrets/{secretName}`
  - **Description**: Update an existing secret, either TLS or opaque.
  - **Path Parameter**:
    - `secretName` - The name of the secret.
    - `namespace` - The namespace of the secret.
  - **Request Body**: (For updating an opaque secret)
    ```json
    {
      "data": [
        {
          "key": "string",  // New Key
          "value": "string" // New Value
        }
      ]
    }
    ```
  - **Response**: Returns the updated secret details or an error message.
    ```json
    {
      "id": "string",
      "type": "string",
      "name": "string",
      "data": [
        {
          "key": "string",
          "value": "string"
        }
      ]
    }
    ```

### Delete a Secret

- **DELETE** `/v1/namespaces/{namespace}/secrets/{secretName}`
  - **Description**: Delete a specific secret, TLS, or opaque.
  - **Path Parameter**:
    `namespace` - The namespace of the secret to delete.
    `secretName` - The name of the secret to delete.
  - **Response**: Confirmation of deletion or an error message.
    ```json
    {
      "message": "string"  // Confirmation message
    }
    ```
