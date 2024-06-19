
# Namespace Management CRUD API

This document outlines the CRUD (Create, Read, Update, Delete) operations for managing Users.

- **GET** `/v1/namespaces/{namespace}/users`
  - **Description**: Get users of namespace.
  - **Path Parameter**: 
    - `namespace` - Namespace of the capp.
  - **Query Params**:
    - `limit`: (optional) Specifies the maximum number of namespaces to return per page. Defaults to 9.
    - `continue`: (optional) Used for fetching the next set of results.
  - **Response**: User of the namespace or an error message.
    ```json
    {
       "users": []str 
    }
    ```

- **POST** `/v1/namespaces/{namespace}/users`
  - **Description**: Add users to namespace.
  - **Path Parameter**:
    - `userNames` - The users of the namespace.
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

- **GET** `/v1/users/{userName}`
  - **Description**: Get user.
  - **Path Parameter**: 
    - `userName` - The user name.
  - **Response**: User or an error message if not found.
    ```json
    {
       "user": str 
    }
    ```
