package main

import (
	"strings"
	"context"
	// "fmt"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httputil"
	"net/url"
	// "github.com/wailsapp/wails/v2/pkg/runtime"
	// "github.com/cretz/bine/tor"
  "golang.org/x/net/proxy"
)

// App struct
type App struct {
	ctx context.Context
	proxy   *http.Server
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	b.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	if b.proxy != nil {
		if err := b.proxy.Shutdown(context.Background()); err != nil {
			// runtime.LogWarning(ctx, "Failed to stop proxy server. Not running?")
		}
	}
}

func (b *App) StartProxy(address string, cert string, port string) (string, error) {
	// runtime.LogInfo(b.ctx, fmt.Sprintf("Starting Proxy Server for %s", address))
	remoteUrl, err := url.Parse(address)
	if (err != nil) {
		return "", err
	}


	reverseProxy := httputil.NewSingleHostReverseProxy(remoteUrl)

	if strings.Contains(address, ".onion") {
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
		if cert != "" {
			// runtime.LogInfo(b.ctx, "Setting certificate")
			certPool := x509.NewCertPool()
			certPool.AppendCertsFromPEM([]byte(cert))

			reverseProxy.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: certPool,
				},
			}
		}
	}
	proxyHandler := func(w http.ResponseWriter, r *http.Request) {
		// runtime.LogInfo(b.ctx, fmt.Sprintf("%s, %q", r.Method, r.URL.Path))
		reverseProxy.ServeHTTP(w, r)
	}
	http.HandleFunc("/", proxyHandler)
	b.proxy = &http.Server{Addr: "0.0.0.0:" + port}

	go func() {
		// runtime.LogInfo(b.ctx, fmt.Sprintf("Proxying  http://0.0.0.0:%s\n to %s", port, remoteUrl.Host))
		if err := b.proxy.ListenAndServe(); err != http.ErrServerClosed {
			// TODO: report back to the frontend
			// runtime.LogWarning(b.ctx, fmt.Sprintf("Could not listen on %s: %v", port, err))
		}
	}()

	return "Proxy active", nil // TODO: error handling
}
