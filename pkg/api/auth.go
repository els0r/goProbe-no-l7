package api

import (
	"context"
	"net/http"
	"strings"

	log "github.com/els0r/log"
)

func (k Keys) exists(key string) bool {
	_, exists := k[key]
	return exists
}

// authenticator implements the middleware http.Handler
type authenticator struct {
	h      http.Handler
	keys   Keys
	logger log.Logger
}

// ServeHTTP checks a request for valid authorization parameters and calls the next handler if successful
func (a *authenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		authCtx context.Context
		userKey string
	)

	ctx := r.Context()
	select {
	case <-ctx.Done():
		http.Error(w, errContextCanceled, http.StatusServiceUnavailable)
		return
	default:
		// get the authorization header of the request
		authHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(authHeader) != 2 || authHeader[0] != "digest" {
			a.logger.Infof("user key denied: invalid auth header")

			ReturnStatus(w, http.StatusUnauthorized)
			return
		}
		userKey = authHeader[1]

		if !a.keys.exists(userKey) {
			a.logger.Debugf("user key '%s' denied: not registered")

			ReturnStatus(w, http.StatusUnauthorized)
			return
		}

		a.logger.Debugf("user successfully authenticated")

		// set key in request context
		authCtx = context.WithValue(r.Context(), apiKeyCtxKey, userKey)

		// call next handler
		a.h.ServeHTTP(w, r.WithContext(authCtx))
	}
}

// AuthenticationHandler registers API authentication keys and returns a middleware that checks them
func (s *Server) AuthenticationHandler(keys Keys) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &authenticator{h: next, keys: keys, logger: s.logger}
	}
}
