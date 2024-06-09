
# Namespace Management CRUD API

## Overview

This document outlines the CRUD (Create, Read, Update, Delete) operations for managing namespaces.

- **GET** `/v1/namespaces`
  - **Description**: get all namespaces
 - **Path Parameter**: 
    -  `namespace` - namespace of the capp.
  - **Query Params**: 
    ```json
    {
      "lables": {
        "key": str,  
        "key1": str
      }
    }
    ```
  - **Response**:  namespaces names  or an error message.
    ```json
    {
       "namespaces": []str,
       "count": int
    }
    ```

- **GET** `/v1/namespaces/{namespace}`
  - **Description**: get specific namespace
 - **Path Parameter**: 
    -  `namespace` - namespace of the capp.

  - **Response**:  namespace name and users for the namespaces or an error message.
    ```json
    {
       "namespace": str ,
    }
    ```

- **GET** `/v1/namespaces/{namespace}/users`
  - **Description**: get user of namespace
 - **Path Parameter**: 
    -  `namespace` - namespace of the capp.

  - **Response**:  user of the namespace or an error message.
    ```json
    {
       "users": []str 
    }
    ```

        


- **POST** `/v1/namespaces`
  - **Description**: Create namespace
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


- **DELETE** `/v1/namespaces/{namespaece}`
  - **Description**: delete namespace
 - **Path Parameter**: 
    -  `namespace` - namespace name.
  - **Response**: Confirmation of deletion an error message.
    ```json
    {
       "message": "string",
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

        