package main

import (
	db "alteroSmartTestTask/common/database"
	"alteroSmartTestTask/common/log"
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ms_persistence "alteroSmartTestTask/backend/services/MS_Persistence"
	api "alteroSmartTestTask/backend/services/MS_Persistence/common/api"
	db_client "alteroSmartTestTask/backend/services/MS_Persistence/database"
	"alteroSmartTestTask/common/flagenv"
	log_context "alteroSmartTestTask/common/log/context"
)

var gRpcPortEnvName = "GRPC_PORT"
var gRpcPortFlag = flag.Int(
	"groc_port",
	0,
	"This is the port from which server will listen grpc.",
)

var restApiPortEnvName = "REST_PORT"
var restApiPortFlag = flag.Int(
	"rest_port",
	0,
	"This is the port for REST API (grpc mirror) listening.",
)

func main() {
	flag.Parse()

	logger := logrus.NewEntry(
		log.ProvideLogrusLoggerUseFlags(),
	)

	logger.Info("Starting grpc")
	gRpcStarting := make(chan struct{})
	go runGRpcListener(logger, gRpcStarting)
	<-gRpcStarting
	logger.Info("GRpc listen")

	logger.Info("Starting http REST Api middleware.")
	httpStarting := make(chan struct{})
	go func(httpStarting <-chan struct{}) {
		<-httpStarting
		logger.Info("Http REST API listen")
	}(httpStarting)
	runHttpRestListener(logger, httpStarting)

	return
}

func runGRpcListener(logger *logrus.Entry, done chan<- struct{}) {
	gRpcListenAddress := getAddressFromPortFlag(
		gRpcPortFlag, gRpcPortEnvName,
	)

	listener, err := net.Listen("tcp", gRpcListenAddress)
	if err != nil {
		logger.Fatalf(
			"failed to open port %s for listen: %s",
			gRpcListenAddress, err,
		)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(18*1024*1024), // 18 Mb
		grpc.MaxSendMsgSize(18*1024*1024), // 18 Mb
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_validator.UnaryServerInterceptor(),
				log_context.ProvideLogContextInterceptor(
					log.ProvideLogrusLoggerUseFlags(),
				).LogContextUnaryServerInterceptor,
			),
		),
	)
	api.RegisterMsPersistenceServer(
		grpcServer,
		ms_persistence.NewService(
			db_client.NewClient(
				db.MustGetNewPostgresConnectionUseFlags(),
			),
		),
	)
	reflection.Register(grpcServer)

	done <- struct{}{}
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatalf("failed to start server: %s\n", err.Error())
	}
}

func runHttpRestListener(logger *logrus.Entry, done chan<- struct{}) {
	muxServer := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
	grpcAddress := getAddressFromPortFlag(
		gRpcPortFlag,
		gRpcPortEnvName,
	)
	err := api.RegisterMsPersistenceHandlerFromEndpoint(
		contextWithLogger(),
		muxServer,
		grpcAddress,
		dialOptions,
	)
	if err != nil {
		logger.Fatalf("failed to register http handler: %s\n", err.Error())
	}
	listenAddress := getAddressFromPortFlag(
		restApiPortFlag,
		restApiPortEnvName,
	)
	logger.WithField("rest_port", listenAddress).
		WithField("grpc_port", grpcAddress).
		Info("start REST")
	server := &http.Server{
		Addr: listenAddress,
		// TODO: create flag for turn on cors and flag for allowed IP
		Handler: allowCORS(muxServer),
	}
	done <- struct{}{}
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("failed to start http endpoint: %s\n", err.Error())
	}
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Grpc-Metadata-auth-token, Grpc-Metadata-app-name, Host, Origin")
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

func getAddressFromPortFlag(portFlag *int, portEnvName string) string {
	return fmt.Sprintf(
		":%d",
		flagenv.MustParseInt(
			portFlag,
			portEnvName,
		),
	)
}

func contextWithLogger() context.Context {
	return log_context.WithLogger(
		context.Background(),
		logrus.NewEntry(
			logrus.New(),
		),
	)
}