package middleware

import (
	"github.com/dkischenko/chat/pkg/logger"
	"net/http"
	"os"
	"runtime/pprof"
	"time"
)

func ProfilingCPU(next http.Handler, l *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Create("profiling/cpu.prof")
		if err != nil {
			l.Entry.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		next.ServeHTTP(w, r)
	})
}

func ProfilingMemory(next http.Handler, l *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m, err := os.Create("profiling/mem.prof")
		if err != nil {
			l.Entry.Fatal(err)
		}
		_ = pprof.WriteHeapProfile(m)
		next.ServeHTTP(w, r)
	})
}

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
