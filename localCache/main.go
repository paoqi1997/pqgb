package main

import (
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/paoqi1997/pqgb/codec"
    "github.com/paoqi1997/pqgb/network"
)

func OnExit(us *network.UnixServerSocket) {
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
    us := network.NewUnixServerSocket(addr)

    us.PacketHandler = func(clientId uint32, packet *codec.Packet) {
        packetType := packet.Type
        packetDataLen := packet.DataLen
        pakDataLen := len(packet.Data)
        Printf("[main] packetType: %d packetDataLen: %d pakDataLen: %d", packetType, packetDataLen, pakDataLen)
        us.Send(clientId, packet.Data)
    }

    err := us.Start()
    if err != nil {
        Printf("[main] Start err: %v", err)
    }

    OnExit(us)

    for {
        time.Sleep(time.Second)
    }
}
