// The MIT License
//
// Copyright (c) 2019-2020, Cloudflare, Inc. and Apple, Inc. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// HTTP constants. Fill in your proxy and target here.
	defaultPort    = "8080"
	proxyEndpoint  = "/proxy"
	queryEndpoint  = "/dns-query"
	healthEndpoint = "/health"

	// Environment variables
	targetNameEnvironmentVariable = "TARGET_INSTANCE_NAME"
)

type odohServer struct {
	endpoints map[string]string
	Verbose   bool
	proxy     *proxyServer
}

func (s odohServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s Handling %s\n", r.Method, r.URL.Path)
	fmt.Fprint(w, "Proxy service\n")
	fmt.Fprint(w, "----------------\n")
	fmt.Fprintf(w, "Proxy endpoint: https://%s%s{?targethost,targetpath}\n", r.Host, s.endpoints["Proxy"])
	fmt.Fprint(w, "----------------\n")
}

func (s odohServer) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s Handling %s\n", r.Method, r.URL.Path)
	fmt.Fprint(w, "ok")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var serverName string
	if serverNameSetting := os.Getenv(targetNameEnvironmentVariable); serverNameSetting != "" {
		serverName = serverNameSetting
	} else {
		serverName = "server_localhost"
	}
	log.Printf("Setting Server Name as %v", serverName)

	endpoints := make(map[string]string)
	endpoints["Proxy"] = proxyEndpoint

	proxy := &proxyServer{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 1024,
				TLSHandshakeTimeout: 0 * time.Second,
			},
		},
	}

	server := odohServer{
		endpoints: endpoints,
		proxy:     proxy,
	}

	http.HandleFunc(proxyEndpoint, server.proxy.proxyQueryHandler)
	http.HandleFunc(healthEndpoint, server.healthCheckHandler)
	http.HandleFunc("/", server.indexHandler)

	log.Printf("Listening on port %v without enabling TLS\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
