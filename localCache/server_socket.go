package main

import (
    "net"
    "sync"

    "github.com/paoqi1997/pqgb/codec"
)

type IServerSocket interface {
    AddClient(net.Conn)
    DelClient(*ServerSocketClient)
    GetClientById(uint32) *ServerSocketClient
    HandlePacket(uint32, *codec.Packet)
}

type UnixServerSocket struct {
    address       string
    listener      *net.UnixListener
    clients       map[uint32]*ServerSocketClient
    clientCount   uint
    clientCounter uint32
    lockOfClients *sync.RWMutex
    handlePacket  func(uint32, *codec.Packet)
}

func NewUnixServerSocket(address string) *UnixServerSocket {
    return &UnixServerSocket{
        address:       address,
        listener:      nil,
        clients:       make(map[uint32]*ServerSocketClient),
        clientCount:   0,
        clientCounter: 0,
        lockOfClients: &sync.RWMutex{},
        handlePacket:  func(uint32, *codec.Packet){},
    }
}

func (us *UnixServerSocket) Start() error {
    unixAddr, err := net.ResolveUnixAddr("unix", us.address)
    if err != nil {
        return err
    }

    us.listener, err = net.ListenUnix("unix", unixAddr)
    if err != nil {
        return err
    }

    go us.Run()

    return nil
}

func (us *UnixServerSocket) Close() {
    if err := us.listener.Close(); err != nil {
        Printf("[UnixServerSocket][Close] Close err: %v", err)
    }
}

func (us *UnixServerSocket) Run() bool {
    for {
        unixConn, err := us.listener.AcceptUnix()
        if err != nil {
            Printf("[UnixServerSocket][Run] AcceptUnix err: %v", err)
            return false
        }

        us.HandleConn(unixConn)
    }
}

func (us *UnixServerSocket) HandleConn(conn net.Conn) {
    us.AddClient(conn)
}

func (us *UnixServerSocket) AddClient(conn net.Conn) {
    client := &ServerSocketClient{
        ss:   us,
        conn: conn,
    }

    us.clientCounter++

    client.clientId = us.clientCounter

    us.lockOfClients.Lock()
    us.clients[us.clientCounter] = client
    us.clientCount++
    us.lockOfClients.Unlock()

    client.Start()
}

func (us *UnixServerSocket) DelClient(client *ServerSocketClient) {
    us.lockOfClients.Lock()
    defer us.lockOfClients.Unlock()
    delete(us.clients, client.clientId)
    us.clientCount--
}

func (us *UnixServerSocket) GetClientById(clientId uint32) *ServerSocketClient {
    us.lockOfClients.RLock()
    client, existed := us.clients[clientId]
    us.lockOfClients.RUnlock()

    if existed {
        return client
    }

    return nil
}

func (us *UnixServerSocket) HandlePacket(clientId uint32, packet *codec.Packet) {
    us.handlePacket(clientId, packet)
}

func (us *UnixServerSocket) Send(clientId uint32, buff []byte) {
    client := us.GetClientById(clientId)
    if client != nil {
        client.Send(buff)
    }
}
