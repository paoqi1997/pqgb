package util

import (
    "bytes"
    "encoding/binary"
)

type IntN interface {
    int8 | int16 | int32 | int64
}

func IntN2Bytes[T IntN](x T) []byte {
    buff := bytes.NewBuffer([]byte{})
    binary.Write(buff, binary.LittleEndian, x)
    return buff.Bytes()
}

func Bytes2IntN[T IntN](in []byte) T {
    buff := bytes.NewBuffer(in)
    var x T
    binary.Read(buff, binary.LittleEndian, &x)
    return x
}
