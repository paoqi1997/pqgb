package main

import (
    "io"
    "net"
)

type ServerSocketClient struct {
    ss       IServerSocket
    conn     net.Conn
    clientId uint32
    sendCh   chan []byte
}

func (sc *ServerSocketClient) Start() {
    sc.sendCh = make(chan []byte, 64)
    go sc.Run()
    go sc.SendLoop()
}

func (sc *ServerSocketClient) Close() {
    if err := sc.conn.Close(); err != nil {
        Printf("[ServerSocketClient][Close] err: %v", err)
    }

    if sc.ss != nil {
        sc.ss.DelClient(sc)
    }
}

func (sc *ServerSocketClient) Run() {
    buffer := make([]byte, 8192) // 8K
    done := false

    for !done {
        nBytes, err := sc.conn.Read(buffer)
        if err != nil {
            if err == io.EOF {
                Printf("[ServerSocketClient][Run] client %d close the connection.")
            } else {
                Printf("[ServerSocketClient][Run] client %d Read err: %v")
            }

            done = true
            continue
        }

        if nBytes > 0 {
            sc.HandlePacket(buffer[:nBytes])
        }
    }

    sc.Close()
}

func (sc *ServerSocketClient) HandlePacket(buff []byte) {

}

func (sc *ServerSocketClient) SendLoop() {
    for buff := range sc.sendCh {
        if buff == nil {
            return
        } else {
            _, err := sc.Send(buff)
            if err != nil {
                Printf("[ServerSocketClient][SendLoop] client %d Send err: %v")
            }
        }
    }
}

func (sc *ServerSocketClient) Send(buff []byte) (int, error) {
    return sc.conn.Write(buff)
}
