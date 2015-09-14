package main

import (
	"flag"
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"log"
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
}

type OneTimeListener struct {
	*net.TCPListener
	Count int
}

func (wh *wrapHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	t := time.Time{} // zero time
	http.ServeContent(rw, r, wh.Basename, t, wh.File)
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

	handler := &wrapHandler{
		Basename: basename,
		File:     fh,
	}

	http.Handle("/"+basename, handler)

	server := &graceful.Server{}

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

	server.ConnState = func(con net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:
			log.Println("New Connection!")
		case http.StateActive:
			log.Println("Active Connection!")
		case http.StateIdle:
			log.Println("Idle Connection!")
		case http.StateHijacked:
			log.Println("Hijacked Connection!")
		case http.StateClosed:
			log.Println("Closed Connection!")
		}
	}

	check(server.ListenAndServe())
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
