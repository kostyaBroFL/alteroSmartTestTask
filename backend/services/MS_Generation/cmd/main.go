package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"
	"alteroSmartTestTask/backend/services/MS_Generation/generator"
	"alteroSmartTestTask/backend/services/MS_Generation/server"
	msperapi "alteroSmartTestTask/backend/services/MS_Persistence/common/api"
	"alteroSmartTestTask/common/flagenv"
	"alteroSmartTestTask/common/log"
	logcontext "alteroSmartTestTask/common/log/context"
)

var (
	gRpcPortEnvName = "GRPC_PORT"
	gRpcPortFlag    = flag.Int(
		"grpc_port",
		0,
		"This is the port from which server will listen grpc.",
	)

	restApiPortEnvName = "REST_PORT"
	restApiPortFlag    = flag.Int(
		"rest_port",
		0,
		"This is the port for REST API (grpc mirror) listening.",
	)

	msPersistenceGrpcHostEnvName = "MS_PERSISTENCE_GRPC_HOST"
	msPersistenceGrpcHostFlag    = flag.String(
		"ms_persistence_grpc_host",
		"",
		"This is the host for grpc connection to ms persistence server.",
	)

	pprofPortEnvName = "PPROF_PORT"
	pprofPortFlag    = flag.Int(
		"pprof_port",
		0,
		"This is the port for pprof",
	)
)

// TODO[#0]: Need to make graceful shutdown.
// TODO[#8]: DI using https://github.com/Monnoroch/go-inject

func main() {
	flag.Parse()

	logger := newLogger()

	go func() {
		if err := http.ListenAndServe(getAddressFromPortFlag(
			pprofPortFlag, pprofPortEnvName), nil); err != nil {
			logger.WithError(err).Error("can not start pprof")
		}
	}()

	logger.Info("starting grpc")
	gRpcStarting := make(chan struct{})
	go runGRpcListener(logger, gRpcStarting)
	<-gRpcStarting
	logger.Info("grpc listening started")

	logger.Info("starting http REST API middleware")
	httpStarting := make(chan struct{})
	go func() {
		<-httpStarting
		logger.Info("http REST API listening started")
	}()
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
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				grpcvalidator.UnaryServerInterceptor(),
				logcontext.ProvideLogContextInterceptor(
					log.ProvideLogrusLoggerUseFlags(),
				).LogContextUnaryServerInterceptor,
			),
		),
	)

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(flagenv.MustParseString(
		msPersistenceGrpcHostFlag, msPersistenceGrpcHostEnvName,
	), grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect to ms_persistence: %s\n", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Fatal("can not to close ms_persistence sonnection")
		}
	}()

	api.RegisterMsGenerationServer(
		grpcServer,
		server.NewServer(
			generator.NewGenerator(),
			msperapi.NewMsPersistenceClient(conn),
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
	err := api.RegisterMsGenerationHandlerFromEndpoint(
		newContextWithLogger(),
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
	// TODO[#1]: create flag for allowed origins coors. (* by default).
	restServer := &http.Server{
		Addr:    listenAddress,
		Handler: allowCORS(muxServer),
	}
	done <- struct{}{}
	if err := restServer.ListenAndServe(); err != nil {
		logger.Fatalf("failed to start http endpoint: %s\n", err.Error())
	}
}

// allowCORS allows Cross Origin Resources Sharing from any origin.
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
	(*w).Header().Set("Access-Control-Allow-Methods",
		"POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, Authorization, Host, Origin")
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

func newContextWithLogger() context.Context {
	return logcontext.WithLogger(
		context.Background(),
		newLogger(),
	)
}

func newLogger() *logrus.Entry {
	return logrus.NewEntry(
		log.ProvideLogrusLoggerUseFlags(),
	)
}
