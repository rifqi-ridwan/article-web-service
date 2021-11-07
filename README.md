# Simple Article Web Service
This is simple article REST API implementation.

This project has 4 layer:
* Entity Layer (this layer will contain data structure)
* Repository Layer (this will handle any db or external request transaction)
* Service Layer (this will handle bussiness logic, you might want to check /internal/article/service.go)
* Delivery Layer (this will handle how to deliver the service)

## Project structure:

```tree
.
├───cmd                     // Main applications for this project.
│   └───web
└───internal                // Private application and library code.
    ├───article             // Article service implementation.
    ├───cache               // Cache implementation.
    ├───entity              // All entity layer.
    ├───pkg                 // Library code.
    │   └───loadenv
    └───test
        └───mocks           // Mocks implementation.
```

## Setup

1. Run
```
$ docker-compose up -d
```

## How to use

### Create Article

```http
POST /articles HTTP/1.1
Content-Type: application/json
```

Request Body:
```json
{
    "title": "test title",
    "content": "test content",
    "author": "test auhtor",
}
```

### Get Articles

```http
GET /articles HTTP/1.1
```

```http
GET /articles?keyword=yourkeyword&author=yourauthor HTTP/1.1
```

keyword and author optional

### Get Article by ID

```http
GET /articles/{id} HTTP/1.1
```
