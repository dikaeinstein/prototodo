package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/dikaeinstein/prototodo/pkg/config"

	"go.uber.org/zap"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/dikaeinstein/prototodo/pkg/logger"
	pb "github.com/dikaeinstein/prototodo/pkg/proto"
	g "github.com/dikaeinstein/prototodo/pkg/protocol/grpc"
	"github.com/dikaeinstein/prototodo/pkg/protocol/grpc/interceptor"
	"github.com/dikaeinstein/prototodo/pkg/todo"
	"github.com/dikaeinstein/prototodo/pkg/todo/service"
	"github.com/dikaeinstein/prototodo/pkg/todo/storage"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func connectToDatabase(dbURI string, l *zap.Logger) *gorm.DB {
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		l.Fatal("failed to open database connection", zap.Error(err))
	}
	db.Debug().AutoMigrate(&todo.ToDo{})

	return db
}

func getCertFile(cfg config.Config, certFileFlag string) string {
	if certFileFlag != "" {
		return certFileFlag
	}
	return cfg.CertFile
}

func getKeyFile(cfg config.Config, keyFileFlag string) string {
	if keyFileFlag != "" {
		return keyFileFlag
	}
	return cfg.KeyFile
}

func main() {
	cfg := config.New()
	flag.BoolVar(&cfg.TLS, "tls", false, "Connection uses TLS if true, else plain TCP")
	flag.StringVar(&cfg.DBName, "db_name", "prototodos", "The database name")
	flag.IntVar(&cfg.Port, "port", 10000, "The server port")
	flag.StringVar(&cfg.AppEnv, "app_env", "development", "The app environment")
	flag.IntVar(&cfg.LogLevel, "log_level", 0, "Global log level")
	certFileFlag := flag.String("cert_file", "", "The TLS cert file")
	keyFileFlag := flag.String("key_file", "", "The TLS key file")

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
	srv := g.NewGRPCToDoServiceServer(s)

	var opts []grpc.ServerOption
	if cfg.TLS {
		cFile := getCertFile(cfg, *certFileFlag)
		kFile := getKeyFile(cfg, *keyFileFlag)

		creds, err := credentials.NewServerTLSFromFile(cFile, kFile)
		if err != nil {
			zapLogger.Fatal("Failed to generate credentials", zap.Error(err))
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		interceptor.LogRPCCalls(zapLogger),
	))

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterToDoServiceServer(grpcServer, srv)

	zapLogger.Info("Starting gRPC server ...")
	if err := grpcServer.Serve(lis); err != nil {
		zapLogger.Fatal("failed to serve", zap.Error(err))
	}
}
