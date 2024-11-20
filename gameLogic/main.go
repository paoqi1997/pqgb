package main

import (
    "net"
    "time"

    "github.com/paoqi1997/pqgb/codec"
)

func OnRead(conn net.Conn) {
    buffer := make([]byte, 1024) // 1K

    for {
        nBytes, err := conn.Read(buffer)
        if err != nil {
            Printf("[OnRead] Read err: %v", err)
            break
        }

        reply := buffer[:nBytes]

        Printf("[OnRead] reply: %s", reply)
    }
}

func main() {
    addr := "/tmp/local-cache.sock"

    conn, err := net.Dial("unix", addr)
    if err != nil {
        Printf("[main] net.Dial err: %v", err)
    }

    defer conn.Close()

    go OnRead(conn)

    for {
        time.Sleep(3 * time.Second)

        msg := "Hello!"

        packet := &codec.Packet{
            Type:    1,
            DataLen: uint32(len(msg)),
            Data:    []byte(msg),
        }

        packetBytes := codec.Packet2Bytes(packet)

        _, err := conn.Write(packetBytes)
        if err != nil {
            Printf("[main] Write err: %v", err)
        }
    }
}
