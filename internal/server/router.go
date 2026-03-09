package server

import (
	"bytes"
	"io"
	"net/http"
	"simple_twitter/internal/logger"
	"simple_twitter/internal/monitoring"
	"time"

	"github.com/rs/zerolog/hlog"
)

type responseRecorder struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body.Write(b)                  // Capture the data
	return rr.ResponseWriter.Write(b) // Send to client
}

func BodyLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		rr := &responseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		next.ServeHTTP(rr, r)

		hlog.FromRequest(r).Info().
			Bytes("request_body", reqBody).
			Str("response_body", rr.body.String()).
			Msg("Body Details")
	})
}

func StartRouter() http.Handler {
	mx := http.NewServeMux()

	mx.HandleFunc("GET /home", ListPost)
	mx.HandleFunc("GET /user/login", LoginHandler)
	mx.Handle("POST /user/login", BodyLogger(http.HandlerFunc(LoginHandler)))

	mx.Handle("POST /user/logout", BodyLogger(http.HandlerFunc(LogoutHandler)))

	mx.HandleFunc("GET /user/register", CreateUserHandler)
	mx.Handle("POST /user/register", BodyLogger(http.HandlerFunc(CreateUserHandler)))

	mx.HandleFunc("POST /posts/create", CreatePostHandler)

	mx.HandleFunc("/post", DetailPostHandler)
	mx.Handle("/posts/update", BodyLogger(http.HandlerFunc(UpdatePostHandler)))
	mx.Handle("/posts/delete/{id}", BodyLogger(http.HandlerFunc(DeletePostHandler)))
	var chain http.Handler = mx

	chain = hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
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
