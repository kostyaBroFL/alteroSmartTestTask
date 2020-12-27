package context

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type logContextInterceptor struct {
	Logger *logrus.Logger
}

// LogContextUnaryServerInterceptor is interceptor that adds standard logrus logger into the request's context.
func (l *logContextInterceptor) LogContextUnaryServerInterceptor(
	requestContext context.Context,
	request interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	requestContextWithLogger := WithLogger(requestContext, logrus.NewEntry(l.Logger))
	return handler(requestContextWithLogger, request)
}

func (l *logContextInterceptor) ProvideServerOption() grpc.ServerOption {
	return grpc.UnaryInterceptor(l.LogContextUnaryServerInterceptor)
}

func ProvideLogContextInterceptor(
	logger *logrus.Logger,
) *logContextInterceptor {
	return &logContextInterceptor{Logger: logger}
}
