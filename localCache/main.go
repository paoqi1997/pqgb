package main

import (
    "os"
    "os/signal"
    "syscall"
    "time"
)

func OnExit(us *UnixServerSocket) {
    ch := make(chan os.Signal, 1)

    signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

    go func(ch chan os.Signal) {
        <-ch

        us.Close()

        os.Exit(0)
    }(ch)
}

func main() {
    addr := "/tmp/local-cache.sock"
    us := NewUnixServerSocket(addr)

    err := us.Start()
    if err != nil {
        Printf("[main] Start err: %v", err)
    }

    OnExit(us)

    for {
        time.Sleep(time.Second)
    }
}
