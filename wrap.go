package main

import (
    "flag"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"
    "net"
)

type wrapHandler struct {
    Basename string
    File     *os.File
}

func (wh *wrapHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    http.ServeContent(rw, r, wh.Basename, time.Now(), wh.File)
}

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

    handler := &wrapHandler{
        Basename: basename,
        File:     fh,
    }

    http.Handle("/"+basename, handler)

    Server: &http.Server{
        Addr: ":8080",
        Handler: handler,
    }
    server.ListenAndServe()
}
