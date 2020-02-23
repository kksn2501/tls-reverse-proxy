package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
)

var (
	ListenPort string
	Upstream   string
	SslCert    string
	SslKey     string
)

func init() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	log.SetFlags(log.Lshortfile)

	ListenPort = os.Getenv("LISTEN_PORT")
	if len(ListenPort) == 0 {
		log.Fatal(`Please set environment variable "LISTEN_PORT"`)
	}
	log.Println(fmt.Sprintf(`LISTEN_PORT=[%s]`, ListenPort))

	Upstream = os.Getenv("UPSTREAM")
	if len(Upstream) == 0 {
		log.Fatal(`Please set environment variable "UPSTREAM"`)
	}
	log.Println(fmt.Sprintf(`UPSTREAM=[%s]`, Upstream))

	SslCert = os.Getenv("SSL_CERT")
	if len(SslCert) == 0 {
		log.Fatal(`Please set environment variable "SSL_CERT"`)
	}
	log.Println(fmt.Sprintf(`SSL_CERT=[%s]`, SslCert))

	SslKey = os.Getenv("SSL_KEY")
	if len(SslKey) == 0 {
		log.Fatal(`Please set environment variable "SSL_KEY"`)
	}
	log.Println(fmt.Sprintf(`SSL_KEY=[%s]`, SslKey))
}

func main() {
	target, err := url.Parse(Upstream)
	if err != nil {
		log.Println(`Invalid environment variable "UPSTREAM"`)
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(`X-Forwarded-Host`, r.Header.Get(`Host`))
		r.Header.Set(`X-Forwarded-Proto`, r.URL.Scheme)
		// r.Host = target.Host
		// r.URL.Host = target.Host
		// r.URL.Scheme = target.Scheme
		proxy.ServeHTTP(w, r)
	})

	err = http.ListenAndServeTLS(":"+ListenPort, SslCert, SslKey, nil)
	if err != nil {
		log.Fatal(err)
	}
}
