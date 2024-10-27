package controllers

import (
	"github.com/julienschmidt/httprouter"
)

type route struct {
	method  string
	pattern string
	handle  httprouter.Handle
}

type routes = map[string]route

var list routes

func init() {
	list = routes{
		"ws_connect": route{"GET", "/ws/:roomuq/:useruq/:perroom", Ws},
		"di":         route{"GET", "/di/:uq", Di},
	}
}

// GetRouters returns routers
func GetRouters() (router *httprouter.Router) {
	router = httprouter.New()

	for _, r := range list {
		router.Handle(r.method, r.pattern, r.handle)
	}

	return
}
