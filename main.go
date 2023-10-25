package main

import (
	"code.8labs.io/rdemeritt/dnm-searcher/app/modules/silkroad4"
	"code.8labs.io/rdemeritt/dnm-searcher/app/modules/slilpp"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/publicsuffix"
)

// Specify Tor proxy ip and port
var torProxyServiceHost string = os.Getenv("TORPROXY_SERVICE_HOST")
var torProxyServicePort string = os.Getenv("TORPROXY_SERVICE_PORT")
var torProxy string = fmt.Sprintf("socks5://%s:%s", torProxyServiceHost, torProxyServicePort)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

  client := setupClient()

	if !slilpp.Start(client) {
		log.Println("something went wrong")
	}

	if !silkroad4.Start(client) {
		log.Println("something went wrong")
	}
}

func setupClient() *http.Client {
  // make sure our proxy url is sane
  torProxyUrl, err := url.Parse(torProxy)
	if err != nil {
		log.Fatal("Error parsing Tor proxy URL:", torProxy, ".", err)
  }

	// setup cookiejar
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal("Error setting up cookiejar:", err)
  }

  // Set up a custom HTTP transport to use the proxy and create the client
  torTransport := &http.Transport{Proxy: http.ProxyURL(torProxyUrl)}
  client := &http.Client{Transport: torTransport, Timeout: time.Second * 120, Jar: jar}
	// client := &http.Client{Transport: torTransport, Timeout: time.Second * 60}


  return client
}
