package handler

import "net/http"

type AppRouter interface {
	GetRouter() http.Handler
}
