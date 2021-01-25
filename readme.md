# **Linux Job Worker**

A worker service for scheduling and running linux processes concurrently from multiple clients. The server provides a REST API secured with mutual TLS for authentication. It offers the ability to start, stop and query the status and logs for all jobs using a simple client cli.

self signed certificates are stored in the project's ssl directory for sole purpose of running locally and as proof of concept.

## Required
- Install Go 1.14.6 or greater: [Download Go](https://golang.org/dl/)
## Project Setup
1. **Download the project**
    ```
   git clone git@github.com:bradyfontenot/ljw.git
   cd ljw
   go get -u ./...
   ```

2. **Build the packages**
   ```
   go build -o /bin/server cmd/server/main.go
   go build -o /bin/client cmd/client/main.go
   ```

   You should now have 2 binaries in your /bin directory named `server` and `client`

**Important Note:** \
Do not relocate the ssl directory or the application will not be able to find the certificates.

## Run the Server
1. You should still be in the project's root directory.
2. **Start the server:**
   - `./bin/server`

## Run the Client
  1. Start a new terminal session in another window.
  2. cd into project's root directory if not already there.

### **Usage:**

**Quick Start** \
prefix all commands with: `./bin/client` 
- `start <linux command>`
- `stop <job id>`
- `list`
- `status <job id>`
- `log <job id>`

<br>

**START**
```bash
# Start a job with a command
./bin/client start ls # where ls is your linux command
```

**LIST**
```bash
# List jobs will retrieve a list of job ids for all jobs 
./bin/client list
```

**STATUS**
```bash
# Check Status of a specific job
./bin/client status <id>  
```

**STOP**
```bash
# Stop a job
./bin/client stop <id> 
```

**LOG**
```bash
# Get a job log
./bin/client log <id> 

```

## Tests

There is currently one test file named `server_test.go` located in `internal/server` directory.

### Run the tests
- **From the project root directory:**
  - `go test ./internal/server`

OR
- **From the server package directory `internal/server`:**
  - `go test`