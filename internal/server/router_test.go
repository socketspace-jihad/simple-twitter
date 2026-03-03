package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Route(t *testing.T) {
	router := StartRouter()
	t.Run("GET /user/login", func(t *testing.T) {
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/user/login", nil)
		if err != nil {
			t.Error(err)
			return
		}

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Error(errors.New("wrong response status code"))
		}

	})

	t.Run("GET /", func(t *testing.T) {
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Error(err)
			return
		}

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Error(errors.New("wrong response status code"))
			return
		}

		t.Run("GET /user/register", func(t *testing.T) {
			rr := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/user/register", nil)
			if err != nil {
				t.Error(err)
				return
			}

			router.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				t.Error(errors.New("wrong response status code"))
				return
			}
		})

	})
}
