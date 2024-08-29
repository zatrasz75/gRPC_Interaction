package gRPCLogger

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
)

type GRPCLoggersInterface interface {
	Error(message string, err error, args ...interface{})
	Info(message string, args ...interface{})
}
type MyLogger struct {
	logger  *log.Logger
	console *os.File
}

func NewGRPCLogger() GRPCLoggersInterface {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	myLogger := &MyLogger{
		logger: logger,
	}

	return myLogger
}

func LogUnaryServerInterceptor(logger GRPCLoggersInterface) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		bodyBytes, ok := req.([]byte)
		if !ok {
			bodyBytes = []byte(fmt.Sprintf("%#v", req))
		}

		logger.Info("Получен запрос на %s: %v, Size: %d bytes", info.FullMethod, req, len(bodyBytes))
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("ошибка при запросе %v: %s", err, info.FullMethod)
			return resp, err
		}
		logger.Info("Успешно обработанный запрос на %s: %v", info.FullMethod, resp)

		return resp, nil
	}
}

func logWithCallerInfo(level string, message string, args ...interface{}) {
	var str strings.Builder
	str.WriteString("< ")
	str.WriteString(level)
	str.WriteString(" >")
	str.WriteString(" ")
	str.WriteString(" ")
	str.WriteString(message)

	formattedMessage := fmt.Sprintf(str.String(), args...)

	log.Println(formattedMessage)
}

// Error ...
func (l *MyLogger) Error(message string, err error, args ...interface{}) {
	if l.logger != nil {
		logWithCallerInfo("GRPC ERROR", message, err, args)
	} else {
		log.Println("No gRPCLogger available.")
	}
}

// Info ...
func (l *MyLogger) Info(message string, args ...interface{}) {
	if l.logger != nil {
		logWithCallerInfo("GRPC INFO", message, args...)
	} else {
		log.Println("No gRPCLogger available.")
	}
}
