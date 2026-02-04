package server

import (
	"net/http"
	"simple_twitter/internal/logger"
	"simple_twitter/internal/monitoring"
	"time"

	"github.com/rs/zerolog/hlog"
)

func StartRouter() http.Handler {
	mx := http.NewServeMux()

	mx.HandleFunc("/", ListPost)
	mx.HandleFunc("/user/login", LoginHandler)
	mx.HandleFunc("/user/logout", LogoutHandler)
	mx.HandleFunc("/user/register", CreateUserHandler)

	mx.HandleFunc("/posts/create", CreatePostHandler)
	mx.HandleFunc("/post", DetailPostHandler)
	mx.HandleFunc("/posts/update", UpdatePostHandler)
	var chain http.Handler = mx

	chain = hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("HTTP Request")
	})(chain)

	chain = hlog.RemoteIPHandler("remote_ip")(chain)
	chain = hlog.UserAgentHandler("user_agent")(chain)
	chain = hlog.RefererHandler("referer")(chain)
	chain = hlog.RequestIDHandler("req_id", "Request-Id")(chain)

	chain = hlog.NewHandler(*logger.Logger.Logger)(chain)

	return monitoring.PrometheusMiddleware(chain)
}
