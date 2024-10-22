package controllers

import (
	"net/http"
	"opachat/tools"
	"strings"

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
		"deb":        route{"GET", "/deb/:uq", Deb},
		"di":         route{"POST", "/di", Di},
	}
}

// GetRouters returns routers
func GetRouters() (router *httprouter.Router) {
	router = httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(tools.Env().Static))

	for _, r := range list {
		router.Handle(r.method, r.pattern, r.handle)
	}

	return
}

func ro(alias string, pars ...string) string {
	pat := list[alias].pattern
	pata := strings.Split(pat, ":")

	if len(pata) == 1 {
		return pat

	}

	if len(pars) == 0 {
		return pat
	}

	purl := ""
	for _, par := range pars {
		if len(purl) > 0 {
			purl += "/"
		}
		purl += par
	}

	ret := pata[0] + purl

	return ret
}
