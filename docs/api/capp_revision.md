
# CappRevision Management CRUD API

## Overview

This document outlines the CRUD (Create, Read, Update, Delete) on CappRevision. Capp revision is represents version of capp.

## API Endpoints

### Capp Revision  
- **GET** `/v1/namespaces/{namespaece}/cappRevisions`
  - **Description**: get all cappRevision of namespace
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


- **GET** `/v1/namepaces/{namespaece}/cappRevisions/{cappRevisionName}`
  - **Response**: CappRevision info or an error message.`
    ```json
    {
      CappRevision 
    }
    ```