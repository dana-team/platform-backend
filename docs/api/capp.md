### Capp Namespaced
- **GET** `/v1/namespaces/{namespace}/capps`
  - **Description**: get all capp of namespace
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
  - **Response**:  capp names  or an error message.
    ```json
    {
       "capps": []str,
       "count": int
    }
    ```


- **GET** `/v1/namespaces/{namespace}/capps/{cappName}`
  - **Description**: get all capps of namespace
 - **Path Parameter**: 
    -  `namespace` - namespace of the capp.
    -  `cappName` - The capp name to fetch.
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
       "capps": []str,
       "count": int
    }
    ```


- **GET** `/v1/namespaces/{namespace}/capps`
  - **Description**: get all capps of namespace
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
  - **Response**: Confirmation of deletion or an error message.
    ```json
    {
       "capps": []Capp,
       "count": int
    }
    ```


- **POST** `/v1/namespaces/{namespace}/capps`
  - **Description**: Create Cpp
  - **Path Parameter**: 
    -  `namespace` - namespace of the capp.
  - **Body**: 
    ```json
    {
        Capp
    }
    ```
  - **Response**: Confirmation of creation or an error message.
    ```json
    {
       Capp,
    }
    ```


- **PUT** `/v1/namespaces/{namespace}/capps`
  - **Description**: Update capp
  - **Path Parameter**: 
    -  `namespace` - namespace of the capp.
  - **Body**: 
```json
{
 "capp": {

    "metadata": {
      "name": "string", // max length 53 char
      "namespace": "string"
    }
    ,
    "spec": {
      "scaleMetric": "concurrency",  // Available options: "cpu", "memory", "rps", "concurrency"
      "site": "example-cluster",     // Optional: specific cluster or placement name
      "state": "enabled",            // Optional: "enabled" (default) or "disabled"
      "configurationSpec": {
        // Configuration details specific to the Capp
      },
      "routeSpec": {
        // Routing specifications for the Capp
      },
      "logSpec": {
        // Log configuration for the Capp
      },
      "volumesSpec": {
        // Volume specifications for the Capp
      }
    }
  }
}

  - **Response**: Confirmation of update or an error message.
```json
{
  "capp": {
    "apiVersion": "yourdomain.com/v1",
    "kind": "Capp",
    "metadata": {
      "name": "example-capp", // max length 53 char
      "namespace": "default"
    }
    ,
    "spec": {
      "scaleMetric": "concurrency",  // Available options: "cpu", "memory", "rps", "concurrency"
      "site": "example-cluster",     // Optional: specific cluster or placement name
      "state": "enabled",            // Optional: "enabled" (default) or "disabled"
      "configurationSpec": {
        // Configuration details specific to the Capp
      },
      "routeSpec": {
        // Routing specifications for the Capp
      },
      "logSpec": {
        // Log configuration for the Capp
      },
      "volumesSpec": {
        // Volume specifications for the Capp
      }
    }
  }
}
```


- **DELETE** `/v1/namespaces/{namespace}/capps/{cappName}`
  - **Description**: get all capps of namespace
 - **Path Parameter**: 
    -  `namespace` - namespace of the capp.
    -  `cappName` - The capp name to fetch.
  - **Response**: Confirmation of deletion an error message.
    ```json
    {
       "message": "string",
    }
    ```