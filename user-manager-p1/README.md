# User Manager

## Setup Instructions

1. Clone the repository:
    ```bash
    git clone https://github.com/rajaarslan95/golang-pst.git
    cd user-manager-p1
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. If Postgres is already Installed then skip this step otherwise Download and Run Postgre Using Below Command:
    ```docker run -p 5432:5432 --name postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
    ```

4. Updating ENV var
- Update these env var according to Postgres running in your machine 
    ```DATABASE_HOST, DATABASE_PORT, DATABASE_USERNAME, DATABASE_PASSWORD, DATABASE_DBNAME
    ```

5. Run the application:
    ```bash
    go run main.go
    ```

6. Run the unit tests:
    ```bash
    go test ./tests/unit_test.go
    ```

## Endpoints

- POST `/users`: Adds a new user.
    ```curl http://127.0.0.1:9001/users -v -X POST -d '{"name": "arslan", "email": "arslan@gmail.com", "age": 23}'
    ```
- GET `/users/{id}`: Retrieves a user by ID.
    ```curl http://127.0.0.1:9001/users/1 -v
    ```
- PUT `/users/{id}`: Updates an existing user.
    ```curl http://127.0.0.1:9001/users/1 -v -X PUT -d '{"name": "Ali", "email": "ali@gmail.com", "age": 25}'
    ```
- DELETE `/users/{id}`: Deletes a user by ID.
    ```curl http://127.0.0.1:9001/users/1 -v -X DELETE
    ```

