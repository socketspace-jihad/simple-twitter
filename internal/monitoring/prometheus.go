package monitoring

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	mx *http.ServeMux
}

func NewPrometheus(configs ...PrometheusConfig) *Prometheus {
	p := &Prometheus{
		mx: http.NewServeMux(),
	}
	for _, config := range configs {
		config(p)
	}
	return p
}

type PrometheusConfig func(*Prometheus)

func WithPromHTTP() PrometheusConfig {
	return func(p *Prometheus) {
		p.mx.Handle("/metrics", promhttp.Handler())
	}
}

func (p *Prometheus) Monitor() {
	srv := &http.Server{
		Addr:    ":9091",
		Handler: p.mx,
	}
	log.Println("prometheus metrics is listening on port 9091...")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err.Error())
	}
}
