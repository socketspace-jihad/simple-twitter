package monitoring

import (
	"net/http"
	"os"
	"simple_twitter/internal/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Prometheus struct {
	mx  *http.ServeMux
	srv *http.Server
}

func NewPrometheus(configs ...PrometheusConfig) *Prometheus {
	router := http.NewServeMux()
	p := &Prometheus{
		mx: router,
	}
	p.srv = &http.Server{
		Handler: p.mx,
		Addr:    ":9091",
	}
	for _, config := range configs {
		if config != nil {
			config(p)
		}
	}
	return p
}

type PrometheusConfig func(*Prometheus)

func WithPromHTTP() PrometheusConfig {
	return func(p *Prometheus) {
		p.mx.Handle("/metrics", promhttp.Handler())
	}
}

func WithHostEnv() PrometheusConfig {
	return func(p *Prometheus) {
		val := os.Getenv("MONITORING_HOST")
		if val != "" {
			p.srv.Addr = val
		}
	}
}

func (p *Prometheus) Monitor() {
	logger.Logger.Info().Str("addr", p.srv.Addr).Msg("monitoring service is listening")
	if err := p.srv.ListenAndServe(); err != nil {
		log.Err(err)
	}
}
