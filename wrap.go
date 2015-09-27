package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

type singleListener struct {
	net.Listener
	active   chan struct{}
	serviced bool
}

type singleConn struct {
	net.Conn
	active chan struct{}
}

func NewSingleListener(l net.Listener) net.Listener {
	return &singleListener{l, make(chan struct{}, 1), false} // allow only 1 connection
}

type errDelivered struct{}

func (e *errDelivered) Error() string {
	return "Delivered Content."
}

func (l *singleListener) Accept() (net.Conn, error) {
	l.active <- struct{}{}
	if l.serviced {
		return nil, &errDelivered{}
	}
	l.serviced = true
	c, err := l.Listener.Accept()
	if err != nil {
		<-l.active
		return nil, err
	}
	return &singleConn{c, l.active}, err
}

func (l *singleConn) Close() error {
	err := l.Conn.Close()
	<-l.active
	return err
}

func main() {
	var err error
	var port = flag.String("port", "8080", "port to run wrap on")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("You must specify a file.")
	}
	path := flag.Args()[0]
	fh, err := os.Open(path)
	check(err)
	defer fh.Close()
	basename := filepath.Base(path)

	ip := getIP()

	u, err := url.Parse("http://" + ip + ":" + *port + "/" + basename)
	check(err)
	log.Println(u)

	hostname, err := os.Hostname()
	u, err = url.Parse("http://" + hostname + ":" + *port + "/" + basename)
	check(err)
	log.Println(u)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := time.Time{} // zero time
		http.ServeContent(w, r, basename, t, fh)
		log.Println("Delivered.")
	})

	var state http.ConnState
	connStateLog := func(c net.Conn, cs http.ConnState) {
		log.Printf("NEW CONN STATE: %v, %v\n", cs, c)
		log.Printf("OLD CONN STATE: %v\n", state)
		if state.String() == "active" && cs.String() == "idle" {
			log.Println("We can shutdown now")
		}
		state = cs
	}

	server := &http.Server{
		Addr:      ":" + *port,
		ConnState: connStateLog,
	}
	server.SetKeepAlivesEnabled(false)

	l, err := net.Listen("tcp", ":"+*port)
	check(err)
	err = server.Serve(NewSingleListener(l))
	if _, ok := err.(*errDelivered); ok {
		log.Println("shutting down safely")
	} else {
		log.Fatal(err)
	}
}
