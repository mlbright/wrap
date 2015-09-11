package main

import (
    "errors"
    "flag"
    "log"
    "net"
    "net/http"
    "os"
    "path/filepath"
    "time"
)

type OneTimeListener struct {
    *net.TCPListener //Wrapped listener
    Accepted         bool
}

func New(l net.Listener) (*OneTimeListener, error) {
    tcpL, ok := l.(*net.TCPListener)

    if !ok {
        return nil, errors.New("Cannot wrap listener")
    }

    retval := &OneTimeListener{}
    retval.TCPListener = tcpL
    retval.Accepted = false

    return retval, nil
}

var AcceptedError = errors.New("Listener stopped")

func (otl *OneTimeListener) Accept() (net.Conn, error) {
    if otl.Accepted {
        return nil, AcceptedError
    }

    newConn, err := otl.TCPListener.Accept()
    return newConn, err
}

type wrapHandler struct {
    Basename string
    File     *os.File
    OTL      *OneTimeListener
}

func (wh *wrapHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    http.ServeContent(rw, r, wh.Basename, time.Now(), wh.File)
    wh.OTL.Accepted = true
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
    log.Println(basename)

    originalListener, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }

    otl, err := New(originalListener)
    if err != nil {
        panic(err)
    }

    handler := &wrapHandler{
        Basename: basename,
        File:     fh,
        OTL:      otl,
    }

    http.Handle("/"+basename, handler)

    server := &http.Server{}
    err = server.Serve(otl)
    if err != nil {
        log.Println(err)
    }
}
