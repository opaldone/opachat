package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"opachat/controllers"
	"opachat/tools"

	"golang.org/x/crypto/acme/autocert"
)

var ct = map[string]int64{
	"ReadTimeout":  10,
	"WriteTimeout": 120,
	"IdleTimeout":  120,
}

func main() {
	e := tools.Env(false)

	if e.Debug {
		startDevTLS(e)
		return
	}

	startTLS(e)
}

func startDevTLS(e *tools.Configuration) {
	fmt.Printf("\n[%s] %s started\ncrt\t\t%s\nkey\t\t%s\naddress\t\t%s:%d\ncsrf\t\t%s\n",
		"debug", e.Appname,
		e.Crt, e.Key,
		e.Address, e.Port,
		tools.GetKeyCSRF(),
	)

	mux := controllers.GetRouters()

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", e.Address, e.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(ct["ReadTimeout"] * int64(time.Second)),
		WriteTimeout:   time.Duration(ct["WriteTimeout"] * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatalln(server.ListenAndServeTLS(e.Crt, e.Key))
}

func startTLS(e *tools.Configuration) {
	fmt.Printf("\n[%s] %s started\nacmehost\t%s\ndirCache\t%s\naddress\t\t%s:%d\ncsrf\t\t%s\n",
		"prod", e.Appname,
		e.Acmehost, e.DirCache, e.Address, e.Port,
		tools.GetKeyCSRF(),
	)

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(e.Acmehost),
		Cache:      autocert.DirCache(e.DirCache),
	}

	mux := controllers.GetRouters()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", e.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(ct["ReadTimeout"] * int64(time.Second)),
		WriteTimeout:   time.Duration(ct["WriteTimeout"] * int64(time.Second)),
		IdleTimeout:    time.Duration(ct["IdleTimeout"] * int64(time.Second)),
		TLSConfig:      certManager.TLSConfig(),
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatalln(server.ListenAndServeTLS("", ""))
}
