package main

import (
	"flag"
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type wrapHandler struct {
	Basename string
	File     *os.File
	Server   *graceful.Server
}

func (wh *wrapHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	http.ServeContent(rw, r, wh.Basename, time.Now(), wh.File)
	wh.Server.Stop(1 * time.Nanosecond)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	var port = flag.String("port", "8080", "port to run wrap on")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("You must specify a file.")
		os.Exit(1)
	}
	path := flag.Args()[0]
	fh, err := os.Open(path)
	check(err)
	defer fh.Close()
	basename := filepath.Base(path)

	server := &graceful.Server{}
	server.ListenLimit = 0

	handler := &wrapHandler{
		Basename: basename,
		File:     fh,
		Server:   server,
	}

	http.Handle("/"+basename, handler)

	server.Server = &http.Server{
		Addr:    ":" + *port,
		Handler: handler,
	}

	ip := getIP()

	u, err := url.Parse("http://" + ip + server.Addr + "/" + basename)
	check(err)
	fmt.Println(u)

	hostname, err := os.Hostname()
	u, err = url.Parse("http://" + hostname + server.Addr + "/" + basename)
	check(err)
	fmt.Println(u)

	// This is actually a problem: there's an error here, and we're ignoring it.
	server.ListenAndServe()
}

func getIP() string {
	ifaces, err := net.Interfaces()
	check(err)
	var ip net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		check(err)
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
		}
	}
	return ip.String()
}
