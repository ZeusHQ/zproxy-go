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
		// req.URL.Scheme = "http"
		// req.URL.Host = proxy

		// client := &http.Client{}

		// //http: Request.RequestURI can't be set in client requests.
		// //http://golang.org/src/pkg/net/http/client.go
		// req.RequestURI = ""

		// delHopHeaders(req.Header)

		// if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// 	appendHostToXForwardHeader(req.Header, clientIP)
		// }

		// resp, err := client.Do(req)
		// if err != nil {
		// 	http.Error(wr, "Server Errorf", http.StatusInternalServerError)
		// 	log.Fatal("ServeHTTP:", err)
		// }
		// defer resp.Body.Close()

		// // log.Println(req.RemoteAddr, " ", resp.Status)

		// delHopHeaders(resp.Header)

		// copyHeader(wr.Header(), resp.Header)
		// wr.WriteHeader(resp.StatusCode)
		// io.Copy(wr, resp.Body)
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

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
// var hopHeaders = []string{
// 	"Connection",
// 	"Keep-Alive",
// 	"Proxy-Authenticate",
// 	"Proxy-Authorization",
// 	"Te", // canonicalized version of "TE"
// 	"Trailers",
// 	"Transfer-Encoding",
// 	"Upgrade",
// 	"Sec-WebSocket-Key",
// }

// func copyHeader(dst, src http.Header) {
// 	for k, vv := range src {
// 		for _, v := range vv {
// 			dst.Add(k, v)
// 		}
// 	}
// }

// func delHopHeaders(header http.Header) {
// 	for _, h := range hopHeaders {
// 		header.Del(h)
// 	}
// }

// func appendHostToXForwardHeader(header http.Header, host string) {
// 	// If we aren't the first proxy retain prior
// 	// X-Forwarded-For information as a comma+space
// 	// separated list and fold multiple headers into one.
// 	if prior, ok := header["X-Forwarded-For"]; ok {
// 		host = strings.Join(prior, ", ") + ", " + host
// 	}
// 	header.Set("X-Forwarded-For", host)
// }
