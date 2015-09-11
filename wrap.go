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
	Done     chan struct{}
}

func (wh *wrapHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	http.ServeContent(rw, r, wh.Basename, time.Now(), wh.File)
	wh.Done <- struct{}{}
	close(wh.Done)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error

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

	done := make(chan struct{})

	handler := &wrapHandler{
		Basename: basename,
		File:     fh,
		Done:     done,
	}

	http.Handle("/"+basename, handler)

	server := &graceful.Server{
		Server: &http.Server{
			Addr:    ":8080",
			Handler: handler,
		},
	}

	ip := getIP()

	u, err := url.Parse("http://" + ip + server.Addr + "/" + basename)
	check(err)
	fmt.Println(u)

	go func() {
		<-done
		server.Stop(1)
	}()

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
