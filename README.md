# GO ToDo List Task Tracker - REST API

# 1. Architecture of the server

## 1.1. Clean Architecture
        - Handler Layer
        - Use Case Layer (Service/Business Logic Layer)
        - Adapter Layer (Postgres Database)
        - Domain Layer

```plaintext
+-----------------------+
|   Handler Layer       |  ---> Receives all of the HTTP requests & calls Use Case Layer
+-----------------------+
            |
            v
+-----------------------+
|  Use Case Layer       |  ---> Contains business logic, calls Adapter Layer
+-----------------------+
            |
            v
+-------------------------+
|   Adapter Layer         |  ---> Communicates with external systems (PostgresDB)
|     (Postgres)          |
+-------------------------+
            |
            v
+-----------------------+
|    Domain Layer       |  ---> Core domain entities & rules
+-----------------------+
```

## 1.2. Design decisions and overview
        - Handler Layer - All of the incoming HTTP request are processed through the handler layer. As a routing framework, the application uses 'chi', becuase of it's optimization, speed and middleware support.
        - Use Case Layer - All of the business logic is stored in this layer. Everything related to tasks CRUD operations. It makes the proper requests to the repository.
        - Adapter Layer (Postgres Database) - This layer makes all the requests to our database.
        - Domain Layer - Every entity struct is kept here, as well as the interfaces that are used to loosely couple the adapter layer.
        - Database schema - In the Postgres db we have 1 table - tasks. All of the information about the tasks is kept in the 'tasks' table.

# 2. How to run the application

    API_PORT=8080
    DB_CONNECTION_URL=postgres://postgres:password@localhost:5432/postgres?sslmode=disable

## 2.1. Locally
    - Firstly you need a running postgres connection. A db creation service is provided inside the docker-compose.yaml. Then open the root directory terminal of the project and run the following command:
    'API_PORT=${PORT} DB_CONNECTION_URL=${DB_URL} go run .' . After providing the needed env vars, the application will work correctly.


## 2.2. Docker
    - Inside the application root dir via the terminal, run 'docker build -t go-task-tracker .'. This command is used to build the application using the provided Dockerfile inside the project. After that to start the application, you need to run 'docker run go-task-tracker'. It is very important before running the application that there is a running postgres database one the provided connection string for the DB_CONNECTION_URL env var. All of the obligatory env vars are found inside the Dockerfile's ENV sections. For the purposes of the demo a .env file isn't provided, so all the values for the env vars are directly found isnide the Dockerfile.

    - The preferred way to run the application via Docker is using the provided docker-compose.yaml file. It has 2 services defined - one for the postgres db and another for our API. The API section is dependant on the db, so only when the db is created and running, only then we build and run our API. To build the project run the following command inside the roots terminal 'docker compose up --build -d'. This will build and run the application in the background. To stop it - use 'docker compose down'. If the app is already built, then just use 'docker compose up -d' to start it again.

# 3. Endpoints

## 3.1. /api/task/{id} (GET)
        - Takes an id a a URL param called 'id'
        - Fetches data from the postgres db. If no data is found for a specific task then it returns HTTP 404 StatusNotFound

        Request:
            (GET) ${apiUrl}/api/task/1461ec84-ccff-4f3c-af34-65d0856ac3ce

```jsx
        Response: 
            (OK - 200):
                {
                    "id": "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
                    "title": "Do unit tests",
                    "description": "Create extensive unit tests for all layers",
                    "status": "PENDING",
                    "due_date": "2025-05-12T00:00:00Z",
                    "created_at": "2025-04-10T22:12:23.273317Z"
                }
                
            (Bad Request - 400):
                {
                    "code": 400,
                    "message": "wrong id format provided"
                }

            (Not Found - 404):
                {
                    "code": 404,
                    "message": "task not found"
                }

            (Internal Server Error - 500):
                {
                    "code": 500,
                    "message": "error occurred"
                }
```

## 3.2. /api/tasks (GET)
        - Fetches all tasks from the postgres db. If no data is found then it returns an empty array

        Request:
            (GET) ${apiUrl}/api/tasks

```jsx
        Response: 
            (OK - 200):
                [
                    {
                        "id": "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
                        "title": "Do unit tests",
                        "description": "Create extensive unit tests for all layers",
                        "status": "PENDING",
                        "due_date": "2025-05-12T00:00:00Z",
                        "created_at": "2025-04-10T22:12:23.273317Z"
                    },
                    {
                        "id": "1461ec84-ccff-4f3c-af34-65d0856ac3cd",
                        "title": "Do unit tests",
                        "description": "Create extensive unit tests for all layers",
                        "status": "PENDING",
                        "due_date": "2025-05-12T00:00:00Z",
                        "created_at": "2025-04-10T22:12:23.273317Z"
                    }
                ]
                
            (Internal Server Error - 500):
                {
                    "code": 500,
                    "message": "error occurred"
                }
```

## 3.3. /api/task (POST)
        - Takes an id a a URL param called 'id'
        - Fetches data from the postgres db. If no data is found for a specific task then it returns HTTP 404 StatusNotFound

        Request:
            (GET) ${apiUrl}/api/task/1461ec84-ccff-4f3c-af34-65d0856ac3ce
    
        Body:
```jsx
            {
                "id": "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
                "title": "Do unit tests",
                "description": "Create extensive unit tests for all layers",
                "status": "PENDING",
                "due_date": "2025-05-12T00:00:00Z"
            }
```

```jsx
        Response: 
            (OK - 200):
                {
                    "id": "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
                    "title": "Do unit tests",
                    "description": "Create extensive unit tests for all layers",
                    "status": "PENDING",
                    "due_date": "2025-05-12T00:00:00Z",
                    "created_at": "2025-04-10T22:12:23.273317Z"
                }

            (Bad Request - 400):
                {
                    "code": 400,
                    "message": "wrong id format provided"
                }

            (Internal Server Error - 500):
                {
                    "code": 500,
                    "message": "error occurred"
                }
```

# 4. Others

## 4.1. Testing
        - Use the command 'go test ./...' to run all tests across all paths in the project
        - Unit tests for the handler layer
        - Unit tests for the use case layer
        - Unit/Integration tests for the adapter (repo) layer
        - Used https://pkg.go.dev/github.com/golang/mock/gomock for the mocking
        - Added mockgen commands into the Makefile, so that mock generation is simplified

## 4.2. Tools used for the api
        - chi - The router framework used to build the HTTP services - https://go-chi.io/#/
        - sqlc - Used to auto generate the sql code into the application. Can be used via the 'sqlc generate' terminal command - https://sqlc.dev/
        - go-migrate - Used to read and execute all db migrations script on project startup. Also used in the integration tests for the db to generate all the needed tables for the tests - https://pkg.go.dev/github.com/golang-migrate/migrate/v4
        - gomock - Used to generate all the mocked objects needed for all the unit tests - https://pkg.go.dev/github.com/golang/mock/gomock
        - Makefile - Used for the mockgen commands
        - docker-compose - Has 2 services defined - db & app. To run the application with docker-compose, the user must run first the command 'docker-compose up --build -d'. This will build everything and run the app. If everything is already built, then the command 'docker compose up -d' is enough to run both the db and the application.

# 5. Workflows

## 5.1. Continuous Integration Workflow - Go CI
    This repository includes a GitHub Actions CI workflow that automatically builds and tests the Go project on every push to the main branch.

    Workflow Steps:

        - Code Checkout – Uses actions/checkout@v4 to pull the latest code.
    
        - Go Setup – Uses actions/setup-go@v5 to install Go version 1.24.
    
        - Build Step – Runs go build ./... to compile the project.
    
        - Test Step – Runs go test ./... to execute all tests.

## 5.2. Continuous Delivery Workflow - Go CD (Release Build)
    This GitHub Actions workflow is triggered on the creation of a new release. It automates the process of building and preparing a release asset for distribution.

    Workflow Steps:

        - Code Checkout – Uses actions/checkout@v4 to fetch the latest code.
        
        - Go Setup – Uses actions/setup-go@v5 to set up Go version 1.24.
        
        - Download Go Dependencies – Runs go mod download to download the project dependencies.
        
        - Install Staticcheck – Installs staticcheck for linting the Go code.
        
        - Lint Code with Staticcheck – Runs static analysis with staticcheck to ensure code quality.
        
        - Run Tests – Executes tests with go test ./... to verify correctness.
        
        - Build Binary – Runs go build -o myapp to compile the Go application into a binary.
        
        - Upload Release Asset – Uses softprops/action-gh-release@v2 to upload the built binary as a release asset.

# 6. Which GitHub Actions from the Marketplace were used
    - CI - Go by GitHub Actions
