package monitoring

import "net/http"

type Prometheus struct{}

func (p *Prometheus) Monitor() {
	srv := &http.Server{
		Addr: ":9090",
	}
}
