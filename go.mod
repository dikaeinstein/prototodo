module github.com/dikaeinstein/prototodo

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/jinzhu/gorm v1.9.11
	github.com/joho/godotenv v1.3.0
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107
	google.golang.org/grpc v1.19.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.24.0
