package network

import (
    "io"
    "net"

    "github.com/paoqi1997/pqgb/codec"
    "github.com/paoqi1997/pqgb/util"
)

type ServerSocketClient struct {
    ss       IServerSocket
    conn     net.Conn
    clientId uint32
    sendCh   chan []byte
    inPakHd  *codec.PacketHandler
    outPakHd *codec.PacketHandler
}

func (sc *ServerSocketClient) Start() {
    sc.sendCh = make(chan []byte, 64)
    sc.inPakHd = codec.NewPacketHandler()
    sc.outPakHd = codec.NewPacketHandler()

    go sc.Run()
    go sc.SendLoop()
}

func (sc *ServerSocketClient) Close() {
    sc.sendCh <- nil

    if err := sc.conn.Close(); err != nil {
        util.Printf("[ServerSocketClient][Close] err: %v", err)
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
                util.Printf("[ServerSocketClient][Run] client %d close the connection.", sc.clientId)
            } else {
                util.Printf("[ServerSocketClient][Run] client %d Read err: %v", sc.clientId, err)
            }

            done = true
            continue
        }

        if nBytes > 0 {
            sc.HandleRead(buffer[:nBytes])
        }
    }

    sc.Close()
}

func (sc *ServerSocketClient) HandleRead(buff []byte) {
    sc.inPakHd.Pack(buff)

    packet := sc.inPakHd.ParsePacket()
    if packet != nil {
        sc.ss.HandlePacket(sc.clientId, packet)
    }
}

func (sc *ServerSocketClient) Send(buff []byte) {
    go func() {
        sc.sendCh <- buff
    }()
}

func (sc *ServerSocketClient) SendLoop() {
    for buff := range sc.sendCh {
        if buff == nil {
            return
        } else {
            sc.HandleWrite(buff)
        }
    }
}

func (sc *ServerSocketClient) HandleWrite(buff []byte) {
    dataLen := len(buff)

    nBytes, err := sc.conn.Write(buff)
    if err != nil {
        util.Printf("[ServerSocketClient][HandleWrite] client %d Write err: %v", sc.clientId, err)
    }

    if nBytes < dataLen {
        go func() {
            sc.sendCh <- buff[nBytes:]
        }()
    }

    remain := dataLen - nBytes

    util.Printf("[ServerSocketClient][HandleWrite] client %d write %d bytes, remain %d bytes.", sc.clientId, nBytes, remain)
}
