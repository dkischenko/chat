package middleware

import (
	"github.com/dkischenko/chat/pkg/logger"
	"net/http"
	"time"
)

func Logging(next http.Handler, l *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		l.Entry.Logger.Infof("Method: %s | Reqest: %s | Latency: %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func PanicAndRecover(next http.Handler, l *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				l.Entry.Logger.Errorf("panic: %+v", err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
