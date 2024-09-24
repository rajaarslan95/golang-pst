# Buggy Project

## Setup Instructions

1. Clone the repository:
    ```bash
    git clone https://github.com/rajaarslan95/golang-pst.git
    cd buggy-project-p2
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Update Postgres credentials in main.go file

4. Run the application:
    ```bash
    go run main.go
    ```

## Potential Bottlenecks

- Incorrect use of goroutines

- Improper error handling

- Missing logging

- Potential race condition

## Brief explanation of issues and fixes

### Incorrect use of goroutines
The original code used goroutines in the getUsers function to execute database queries asynchronously. However, this was unnecessary for a simple task like fetching data and could introduce complexity. 
The wg.Wait() calls in both getUsers and createUser block the main thread, defeating the purpose of goroutines for asynchronous execution. 

#### Fix
By removing the goroutines, the code becomes simpler and easier to understand. It's often better to keep things as straightforward as possible, especially for relatively small functions
The wg.Wait() calls were removed as they were preventing the main thread from handling other requests.

### Improper error handling
The original code lacked proper error handling for various operations like database connectivity, database queries and JSON decoding. This could lead to unexpected behavior or server crashes.

#### Fix
Specific error checks have been added.

### Missing logging
The original code lacked comprehensive logging for database operations and errors. This made it difficult to diagnose issues.

#### Fix
Logging statements have been added

### Potential race condition
If multiple goroutines were accessing the database simultaneously, there could be a risk of race conditions where data is modified or read inconsistently

#### Fix
Removed goroutines when quering database