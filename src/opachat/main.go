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

	if e.Acme {
		startAcme(e)
		return
	}

	startSelf(e)
}

func startSelf(e *tools.Configuration) {
	fmt.Printf("\n[%s] %s\ntime\t\t%s\ncrt\t\t%s\nkey\t\t%s\naddress\t\t%s:%d\nturn\t\t%s\n",
		"self", e.Appname,
		time.Now().Format("2006-01-02 15:04:05"),
		e.Crt, e.Key,
		e.Address, e.Port,
		e.IceList[0]["urls"],
	)

	mux := controllers.GetRouters()

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", e.Address, e.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(ct["ReadTimeout"] * int64(time.Second)),
		WriteTimeout:   time.Duration(ct["WriteTimeout"] * int64(time.Second)),
		IdleTimeout:    time.Duration(ct["IdleTimeout"] * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatalln(server.ListenAndServeTLS(e.Crt, e.Key))
}

func startAcme(e *tools.Configuration) {
	fmt.Printf("\n[%s] %s\ntime\t\t%s\nacmehost\t%s\ndirCache\t%s\naddress\t\t%s:%d\nturn\t\t%s\n",
		"acme", e.Appname,
		time.Now().Format("2006-01-02 15:04:05"),
		e.Acmehost, e.DirCache, e.Address, e.Port,
		e.IceList[0]["urls"],
	)

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(e.Acmehost...),
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
