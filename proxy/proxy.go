package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type InterceptingTransport struct {
	Target      *url.URL
	Broadcaster *Broadcaster
	ProxyName   string
}

func (t *InterceptingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqBodyBytes []byte
	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	var respBodyBytes []byte
	if resp.Body != nil {
		respBodyBytes, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(respBodyBytes))
		
		go t.Broadcaster.Broadcast(map[string]interface{}{
			"proxy":    t.ProxyName,
			"url":      req.URL.String(),
			"method":   req.Method,
			"reqBody":  string(reqBodyBytes),
			"respBody": string(respBodyBytes),
			"status":   resp.StatusCode,
		})
	}

	return resp, nil
}

func setupReverseProxy(listenAddr, targetURL, proxyName string, broadcaster *Broadcaster) *http.Server {
	target, _ := url.Parse(targetURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	
	proxy.Transport = &InterceptingTransport{
		Target:      target,
		Broadcaster: broadcaster,
		ProxyName:   proxyName,
	}
	
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	go func() {
		fmt.Printf("[%s] Listening on %s, forwarding to %s\n", proxyName, listenAddr, targetURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[%s] error: %v", proxyName, err)
		}
	}()

	return srv
}
