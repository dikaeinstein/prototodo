## compile protobuf using protoc
build-proto:
	protoc --go_out=plugins=grpc:. pkg/proto/*.proto

## start gRPC server
start-server:
	go run cmd/server/*.go

## run the gRPC client
run-client:
	go run cmd/client/*.go
