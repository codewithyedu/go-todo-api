# Todo API in Go
This project is a simple RESTful API for managing a todo list, built using Go (Golang). It includes routes to create, read, update, and delete todos, as well as unit tests for each route.

## Features
- `Create Todo`: Add a new todo item with a title and completion status.
- `Get Todo`: Retrieve a single todo item by its unique ID.
- `Get Todos`: Retrieve all todos.
- `Update Todo`: Update the title or completion status of an existing todo item.
- `Delete Todo`: Remove a todo item by ID.
 
## Prerequisites
To run this project, you need:

- Go 1.18 or later
- Google UUID package (to generate unique IDs)
- Validator

Install the UUID package with:
```
go get github.com/google/uuid
```

Install the validator package with:
```
go get github.com/go-playground/validator/v10
```

## Getting Started
1. Clone the repository:
    ```
    git clone https://github.com/codewithyedu/go-todo-api.git
    cd go-todo-api
    ```

2. Run the API server:
    ```
    go run .
    ```
   
3. Access the API: The server will run on ``http://localhost:8080/api/v1/todos``

## Endpoints
- POST /api/v1/todos: Create a new todo
- GET /api/v1/todos: Get all todos
- GET /api/v1/todos/{id}: Get a todo by ID
- PUT /api/v1/todos/{id}: Update a todo by ID
- DELETE /api/v1/todos/{id}: Delete a todo by ID

## Running Tests
The project includes tests for each API endpoint to verify functionality. Run the tests with:
```
go test -v ./...
```
These tests use Go’s net/http/httptest package to simulate API requests and check responses for correctness.