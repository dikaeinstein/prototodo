# prototodo

My first shot at gRPC with golang

Features:

* Health check
* Graceful shutdown

## Run Locally

Prerequisites:

* make
* Go 1.12+

### Available Commands

`build-proto` - Compiles todo.proto using protoc compiler for golang

`start-server` - Starts the gRPC server

`run-client` - Runs the gRPC client against the server
