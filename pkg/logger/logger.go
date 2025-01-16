package logger

import (
	"log"
	"net/http"
	"time"
)

type ContextKey string

const ContextUserKey ContextKey = "request_id"

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func WrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, req)
		request_id := req.URL.Query().Get("request_id")
		statusCode := lrw.statusCode
		log.Printf("time: %s, client IP: %s, status code: %d %s, request_id: %s", time.Now().String(), req.Host, statusCode, http.StatusText(statusCode), request_id)

	})
}
