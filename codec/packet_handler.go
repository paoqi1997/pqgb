package codec

import (
    "fmt"

    "github.com/paoqi1997/pqgb/util"
)

const (
    PACKET_FIELD_TYPE_LEN     = 1 // Packet.Type 字段的长度
    PACKET_FIELD_DATA_LEN_LEN = 4 // Packet.DataLen 字段的长度
    PACKET_FIELD_DATA_MIN_LEN = 1 // Packet.Data 字段的最小长度
    PACKET_MIN_LEN            = 6 // Packet 大小的最小值
)

type Packet struct {
    Type    uint8
    DataLen uint32
    Data    []byte
}

type PacketHandler struct {
    readIndex  uint32
    writeIndex uint32
    buffer     []byte
}

func NewPacketHandler() *PacketHandler {
    return &PacketHandler{
        readIndex:  0,
        writeIndex: 0,
        buffer:     make([]byte, 8192),
    }
}

func Packet2Bytes(packet *Packet) []byte {
    packetBytes := append([]byte{}, packet.Type)
    dataLenField := util.IntN2Bytes(int32(packet.DataLen))
    packetBytes = append(packetBytes, dataLenField...)
    packetBytes = append(packetBytes, packet.Data...)
    return packetBytes
}

func (ph *PacketHandler) vacate() {
    halfCap := uint32(cap(ph.buffer) / 2)
    if ph.readIndex >= halfCap {
        ph.buffer = ph.buffer[ph.readIndex:]
        ph.writeIndex -= ph.readIndex
        ph.readIndex = 0
    }
}

func (ph *PacketHandler) Len() int {
    return int(ph.writeIndex - ph.readIndex)
}

func (ph *PacketHandler) Cap() int {
    return cap(ph.buffer)
}

func (ph *PacketHandler) ParsePacket() *Packet {
    buffLen := ph.Len()

    if buffLen < PACKET_MIN_LEN {
        return nil
    }

    dataLenField := ph.buffer[PACKET_FIELD_TYPE_LEN:(PACKET_FIELD_TYPE_LEN + PACKET_FIELD_DATA_LEN_LEN)]
    dataLen := util.Bytes2IntN[int32](dataLenField)

    buffDataLen := buffLen - (PACKET_FIELD_TYPE_LEN + PACKET_FIELD_DATA_LEN_LEN)

    if buffDataLen < int(dataLen) {
        return nil
    }

    packetType := ph.UnpackInt8()
    packetDataLen := ph.UnpackInt32()
    packetData := ph.Unpack(uint32(packetDataLen))

    fmt.Printf("[PacketHandler][ParsePacket] dataLen: %d, packetDataLen: %d\n", dataLen, packetDataLen)

    return &Packet{
        Type:    uint8(packetType),
        DataLen: uint32(packetDataLen),
        Data:    packetData,
    }
}

func (ph *PacketHandler) Pack(inData []byte) {
    ph.buffer = append(ph.buffer[:ph.writeIndex], inData...)
    dataLen := len(inData)
    ph.writeIndex += uint32(dataLen)
}

func (ph *PacketHandler) PackInt8(x int8) {
    ph.Pack(util.IntN2Bytes(x))
}

func (ph *PacketHandler) PackInt16(x int16) {
    ph.Pack(util.IntN2Bytes(x))
}

func (ph *PacketHandler) PackInt32(x int32) {
    ph.Pack(util.IntN2Bytes(x))
}

func (ph *PacketHandler) PackInt64(x int64) {
    ph.Pack(util.IntN2Bytes(x))
}

func (ph *PacketHandler) PackString(s string) {
    ph.Pack([]byte(s))
}

func (ph *PacketHandler) Unpack(size uint32) []byte {
    buffLen := ph.writeIndex - ph.readIndex
    sz := min(size, buffLen)

    outData := make([]byte, sz)
    copy(outData, ph.buffer[ph.readIndex:(ph.readIndex + sz)])

    ph.readIndex += sz
    ph.vacate()

    return outData
}

func (ph *PacketHandler) UnpackInt8() int8 {
    return util.Bytes2IntN[int8](ph.Unpack(1))
}

func (ph *PacketHandler) UnpackInt16() int16 {
    return util.Bytes2IntN[int16](ph.Unpack(2))
}

func (ph *PacketHandler) UnpackInt32() int32 {
    return util.Bytes2IntN[int32](ph.Unpack(4))
}

func (ph *PacketHandler) UnpackInt64() int64 {
    return util.Bytes2IntN[int64](ph.Unpack(8))
}

func (ph *PacketHandler) UnpackString(size uint32) string {
    return string(ph.Unpack(size))
}
