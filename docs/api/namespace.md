
# Namespace Management API

This document outlines the CRUD (Create, Read, Update, Delete) operations for managing namespaces.

- **GET** `/v1/namespaces`
  - **Description**: Get all namespaces.
  - **Path Parameter**:
    - `namespace` - The namespace of the capp.
  - **Query Params**:
    - `limit`: (optional) Specifies the maximum number of namespaces to return per page. Defaults to 9.
    - `continue`: (optional) Used for fetching the next set of results.
    - `labels`: (optional) Used for filtering namespaces by labels.
  - **Response**: Namespaces names or an error message.
    ```json
    {
       "namespaces": []str,
       "count": int
    }
    ```

- **GET** `/v1/namespaces/{namespace}`
  - **Description**: Get specific namespace.
  - **Path Parameter**: 
    - `namespace` - The namespace of the capp.
  - **Response**: Namespace name or an error message.
    ```json
    {
       "namespace": str,
    }
    ```

- **GET** `/v1/namespaces/{namespace}/users`
  - **Description**: Get users of namespace.
  - **Path Parameter**:
    -  `namespace` - The namespace of the capp.
  - **Response**: The users of the namespace or an error message.
    ```json
    {
       "users": []str 
    }
    ```

- **POST** `/v1/namespaces`
  - **Description**: Create a new namespace.
  - **Path Parameter**:
  - **Body**:
    ```json
    {
      name: str
    }
    ```
  - **Response**: Confirmation of creation or an error message.
    ```json
    {
       "namespace": str,
    }
    ```

- **POST** `/v1/namespaces/users`
  - **Description**: Add users to namespace.
  - **Path Parameter**:
    - `userNames` - User names.
  - **Body**:
    ```json
    {
      userNames: []str
    }
    ```
  - **Response**: Confirmation of creation or an error message.
    ```json
    {
       "namespace": str,
       "users": []str
    }
    ```

- **DELETE** `/v1/namespaces/{namespaece}`
  - **Description**: Delete a specific namespace.
  - **Path Parameter**:
    - `namespace` - Namespace name.
  - **Response**: Confirmation of deletion an error message.
    ```json
    {
       "message": "string",
    }
    ```

- **GET** `/v1/users/{userName}`
  - **Description**: Get a specific user.
  - **Path Parameter**:
    - `userName` - The user name.
  - **Response**: The user or an error message if not found.
    ```json
    {
       "user": str
    }
    ```
