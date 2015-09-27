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

func ConnStateListener(c net.Conn, cs http.ConnState) {
  log.Printf("CONN STATE: %v, %v\n", cs, c)
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

  server := &http.Server{
    Addr: ":" + *port,
    ConnState: ConnStateListener,
  }
  check(server.ListenAndServe())
}
