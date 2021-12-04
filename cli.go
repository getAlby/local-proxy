package main

import (
	"strings"
	"fmt"
	"encoding/base64"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httputil"
	"net/url"
	"flag"
	// "github.com/wailsapp/wails/v2/pkg/runtime"
	// "github.com/cretz/bine/tor"
	"golang.org/x/net/proxy"
)

func main() {
	address := flag.String("target", "", "Target URL")
	secret := flag.String("secret", "", "Local proxy secret")
	port := flag.String("port", "8181", "Local proxy port")
	certBase64 := flag.String("cert", "", "Optional request certificate Base64 encoded")

	flag.Parse()

	if (*address == "") {
		panic("Config missing. Check -help")
	}
	fmt.Printf("Starting Proxy Server for %s\n", *address)
	remoteUrl, err := url.Parse(*address)
	if (err != nil) {
		panic(err)
	}


	reverseProxy := httputil.NewSingleHostReverseProxy(remoteUrl)

	if strings.Contains(*address, ".onion") {
		// Start Tor
		// t, err := tor.Start(nil, nil)
		// if err != nil {
		// 	panic(err)
		// }
		// //defer t.Close()
		// torDialer, err := t.Dialer(context.Background(), nil)
		// if err != nil {
		// 	panic(err)
		// }
		// reverseProxy.Transport = &http.Transport{
		// 	DialContext: torDialer.DialContext,
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// }

		socksDialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9150", nil, nil)
		if err != nil {
			panic(err)
		}
		reverseProxy.Transport = &http.Transport{
			Dial: socksDialer.Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		if *certBase64 != "" {
			certDER, err := base64.RawURLEncoding.DecodeString(*certBase64)
			if err != nil {
				panic(err)
			}
			certPool := x509.NewCertPool()
			cert, err := x509.ParseCertificate(certDER)
			if err != nil {
				panic(err)
			}
			//certPool.AppendCertsFromPEM(cert) // FOR PEM
			certPool.AddCert(cert)

			reverseProxy.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: certPool,
				},
			}
		}
	}
	proxyHandler := func(w http.ResponseWriter, r *http.Request) {
		r.Host = remoteUrl.Host
		// remove secret form the URL
		r.URL.Path = strings.Replace(r.URL.Path, *secret + "/", "", 1)

		reverseProxy.ServeHTTP(w, r)
	}
	path := "/"
	if *secret != "" {
		path = path + *secret + "/"
	}
	http.HandleFunc(path , proxyHandler)
	proxy := &http.Server{Addr: "0.0.0.0:" + *port}

	fmt.Printf("Proxying  http://0.0.0.0:%s to %s\n\n", *port, remoteUrl.Host)
	fmt.Printf("Proxy running. Connect Alby to: http://localhost:%s/%s\n", *port, *secret)
	if err := proxy.ListenAndServe(); err != http.ErrServerClosed {
		// TODO: report back to the frontend
		fmt.Printf("Could not listen on %s: %v", port, err)
	}

}
