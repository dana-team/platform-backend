
# Namespace Management CRUD API

## Overview

This document outlines the CRUD (Create, Read, Update, Delete) operations for managing Users.



- **GET** `/v1/namespaces/{namespace}/users`
  - **Description**: get users of namespace
 - **Path Parameter**: 
    -  `namespace` - namespace of the capp.

  - **Response**:  user of the namespace or an error message.
    ```json
    {
       "users": []str 
    }
    ```



- **POST** `/v1/namespaces/{namespace}/users`
  - **Description**: Add members to namespace
  - **Path Parameter**: 
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
  - **Description**: get user
 - **Path Parameter**: 
    -  `userName` - user name.

  - **Response**:  user or an error message if not found.
    ```json
    {
       "user": str 
    }
    ```

        