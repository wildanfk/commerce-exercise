package librest

import (
	"net/http"
	"time"
	"user-service/internal/util/liberr"

	"go.uber.org/zap"
)

type loggingMiddleware struct {
	path    string
	logger  *zap.Logger
	handler GatewayHandler
}

func WithLoggingMiddleware(path string, logger *zap.Logger) GatewayMiddleware {
	return func(handle GatewayHandler) GatewayHandler {
		lm := loggingMiddleware{path: path, handler: handle, logger: logger}

		return lm.handle
	}
}

func (lm *loggingMiddleware) handle(w http.ResponseWriter, r *http.Request) error {
	rw := WrapResponseWriter(w)

	start := time.Now()
	err := lm.handler(rw, r)

	duration := time.Since(start)
	fields := []zap.Field{
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
		zap.Int("status_code", rw.StatusCode()),
		zap.Duration("duration", duration),
	}

	if err != nil {
		fields = liberr.AppendErrorLogField(fields, err)

		lm.logger.Error("Request failed", fields...)
	} else {
		lm.logger.Info("Request completed", fields...)
	}

	return err
}
