package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var err error
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("You must specify a file.")
	}
	path := flag.Args()[0]
	fh, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	basename := filepath.Base(path)
	log.Println(basename)
	http.HandleFunc("/" + basename, func(rw http.ResponseWriter, r *http.Request) {
		http.ServeContent(rw, r, basename, time.Now(), fh)
	})
    log.Fatal(http.ListenAndServe(":80", nil))
}
