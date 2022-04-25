package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type proxy struct {
	proxies map[string]*httputil.ReverseProxy
}

func (p *proxy) proxy(wr http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.URL.RequestURI(), "_next") && !strings.Contains(req.URL.RequestURI(), "favicon.ico") {
		log.Println("[Proxy]", req.Method, " ", req.Host, req.URL)
	}

	if proxy, ok := p.proxies[req.Host]; ok {
		proxy.ServeHTTP(wr, req)
	} else {
		fmt.Println("Invalid host", req.Host)
		wr.Write([]byte("Invalid host"))
		return
	}

}

func (p *proxy) handleZeus(wr http.ResponseWriter, req *http.Request) {
	log.Println("[Zeus]", req.RemoteAddr, " ", req.Method, " ", req.Host, " ", req.URL)
	wr.Write([]byte("Zeus"))

}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	if req.Host == "dev.z" {
		p.handleZeus(wr, req)
	} else {
		p.proxy(wr, req)
	}
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}
