# the-better-backend
A GoLang backend using Fiber and MongoDB for the Wholesome living Project

## Getting Started

### Prerequisites

- [GoLang](https://golang.org/doc/install)
- [MongoDB](https://docs.mongodb.com/manual/installation/)

### Installing

0. Install extra packages: 
    ```go install github.com/cosmtrek/air@latest```
    ```go install github.com/swaggo/swag/cmd/swag@latest```
1. Clone the repo
2. Add .env file
3. ```make dev```
4. view docs at http://localhost:8080/swagger

### Scripts

- ```make dev``` - runs the server in development mode
- ```make swagger``` - generates the swagger docs
- ```make test``` - runs the tests

### Testing
Run test ```make test``` or ```go test ./...``` 