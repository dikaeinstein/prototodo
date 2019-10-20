package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/dikaeinstein/prototodo/pkg/config"
	"github.com/dikaeinstein/prototodo/pkg/logger"
	pb "github.com/dikaeinstein/prototodo/pkg/proto"
	g "github.com/dikaeinstein/prototodo/pkg/protocol/grpc"
	"github.com/dikaeinstein/prototodo/pkg/protocol/grpc/interceptor"
	"github.com/dikaeinstein/prototodo/pkg/todo"
	"github.com/dikaeinstein/prototodo/pkg/todo/service"
	"github.com/dikaeinstein/prototodo/pkg/todo/storage"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func connectToDatabase(dbURI string, l *zap.Logger) *gorm.DB {
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		l.Fatal("failed to open database connection", zap.Error(err))
	}
	db.Debug().AutoMigrate(&todo.Todo{})

	return db
}

func main() {
	cfg := config.New()
	flag.BoolVar(&cfg.TLS, "tls", cfg.TLS, "Connection uses TLS if true, else plain TCP")
	flag.StringVar(&cfg.DBName, "db_name", cfg.DBName, "The database name")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "The server port")
	flag.StringVar(&cfg.AppEnv, "app_env", cfg.AppEnv, "The app environment")
	flag.IntVar(&cfg.LogLevel, "log_level", cfg.LogLevel, "Global log level")
	flag.StringVar(&cfg.CertFile, "cert_file", cfg.CertFile, "The TLS cert file")
	flag.StringVar(&cfg.KeyFile, "key_file", cfg.KeyFile, "The TLS key file")

	flag.Parse()

	zapLogger := logger.NewZapLogger(cfg.LogLevel, "2006-01-02T15:04:05Z07:00", cfg.AppEnv)
	defer zapLogger.Sync()

	dbURI := fmt.Sprintf("host=localhost user=Dikaeinstein dbname=%s sslmode=disable", cfg.DBName)
	db := connectToDatabase(dbURI, zapLogger)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.Port))
	if err != nil {
		zapLogger.Fatal("failed to listen", zap.Error(err))
	}
	defer lis.Close()

	r := storage.NewPostgresStore(db)
	s := service.New(r)
	srv := g.NewGRPCTodoHandler(s)

	var opts []grpc.ServerOption
	if cfg.TLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			zapLogger.Fatal("Failed to generate credentials", zap.Error(err))
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		interceptor.LogRPCCalls(zapLogger),
	))

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTodoServiceServer(grpcServer, srv)
	h := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, h)
	h.SetServingStatus("TodoService", grpc_health_v1.HealthCheckResponse_SERVING)

	msg := fmt.Sprintf("gRPC server listening on %d...", cfg.Port)
	zapLogger.Info(msg)
	if err := grpcServer.Serve(lis); err != nil {
		zapLogger.Fatal("failed to serve", zap.Error(err))
	}

	h.SetServingStatus("TodoService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	grpcServer.GracefulStop()
}
