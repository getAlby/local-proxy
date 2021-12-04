# Alby Proxy

If you have problems connecting Alby to your node (for example because your node is behind TOR)
this little proxy application helps you.

* Start this application
* Enter your node details (host e.g. http://youronionurl.onion:9090)
* Configure Alby to connect to localhost:8181 (or whatever port you used)


## Installation

You can download an executable from the [release assets](https://github.com/getAlby/local-proxy/releases)


## NOTE!

Please know that this opens a http server listening on localhost. Other apps could also try to connect to it. To secure the Proxy we use a custom "secret" behind which the proxy runs. 

## TOR

If you enter a `.onion` URL make sure you have a Tor socks proxy running on `127.0.0.1:9150`. All traffic will be sent through this Tor proxy.

## Live Development

To run in live development mode, run `wails dev` in the project directory. The frontend dev server will run
on http://localhost:34115. Open this in your browser to connect to your application.

## Building

For a production build, use `wails build`.


## Screenshot
![screenshot](https://user-images.githubusercontent.com/318/144709052-9105cf80-02e9-42b8-a934-666e6d5babcc.png)

