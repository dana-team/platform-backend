# CappRevision Management API

This document outlines the CRUD (Create, Read, Update, Delete) on CappRevision. Capp revision is represents version of capp.

## API Endpoints

### Capp Revision

- **GET** `/v1/namespaces/{namespace}/cappRevisions`
  - **Description**: Get all cappRevision of namespace.
  - **Query Params**:
    ```json
    {
      "lables": {
        "key": str,
        "key1": str
      }
    }
    ```
  - **Response**: Capp revisions or an error message.
    ```json
    {
       "cappRevisions": []CappRevision,
       "count": int
    }
    ```

- **GET** `/v1/namepaces/{namespace}/cappRevisions/{cappRevisionName}`
  - **Response**: CappRevision info or an error message.
    ```json
    {
      CappRevision
    }
    ```
