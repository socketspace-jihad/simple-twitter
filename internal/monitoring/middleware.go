package monitoring

import (
	"net/http"
	"time"
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues(path, r.Method).Observe(duration)
	})
}
