# ginkvtest

Repository for testing Key-Value datasources using the Gin framework and load testing with K6.

## Overview
This project uses the Gin Web Framework for building HTTP endpoints in Go. Load testing is performed using K6, enabling performance evaluation for KV datasources.

## Features
- Gin HTTP Framework: High-performance and lightweight server setup.
- Key-Value Data Source Tests: Benchmarks for KV store performance.
- Load Testing with K6: Simulates real-world usage scenarios for endpoints.
## Project Structure

```plaintext
ginkvtest/
├── cmd/                 # Main application
├── internal/            # Core business logic
├── load_test/           # Load testing scripts using K6
│   ├── load_test.js   # K6 script for performance testing 
├── Dockerfile           # Docker image configuration
├── docker-compose.yml   # Multi-service setup
└── go.mod               # Go module dependencies
```

## Pre-requisites

- Go 1.21+
- Docker and Docker Compose
- K6 CLI installed for load testing
  Setup and Usage
  Clone the Repository:

```bash
git clone https://github.com/Miguelme/ginkvtest.git
cd ginkvtest
```
## Run the Application:

```bash
docker-compose up --build
```

## Load Testing with K6 
From the `load_test` directory, run:
```bash
k6 run load_test.js
```


## Considerations 
Right now, the whole set-up is done to be tested locally. 
In order to do the tests pointing to proper dynamor, redis and aurora instances 
we should work with the environment variables and setup of the credentials in the respective services:

- Redis: `service/redis_service.go`
- Dynamo: `service/dynamo_service.go`
- Aurora: `service/aurora_service.go`

Also, now we are doing 3 goroutines (one for each of the datasources) and populating automatically the databases if needed. 
In case we don't like the approach of one endpoint testing 3 datasources with go routines, 
we can create only the layers for specific endpoints depending on the datasource and only call one datasource from each of the endpoints in a matter of 1-2 hours.