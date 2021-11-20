package http

import (
	"hw9/internal/cache"
	"hw9/internal/store"
)

type ServerOption func(srv *Server)

func WithAddress(address string) ServerOption {
	return func(srv *Server) {
		srv.Address = address
	}
}

func WithStore(store store.Store) ServerOption {
	return func(srv *Server) {
		srv.store = store
	}
}

func WithCache(cache cache.Cache) ServerOption {
	return func(srv *Server) {
		srv.cache = cache
	}
}
