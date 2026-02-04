package main

import (
	"encoding/binary"
	"io"
)

func encodeVarInt(val uint64) []byte {
	switch {
	case val < 0xfd:
		return []byte{byte(val)}

	case val <= 0xffff:
		b := make([]byte, 3)
		b[0] = 0xfd
		binary.LittleEndian.PutUint16(b[1:], uint16(val))
		return b

	case val <= 0xffffffff:
		b := make([]byte, 5)
		b[0] = 0xfe
		binary.LittleEndian.PutUint32(b[1:], uint32(val))
		return b

	default:
		b := make([]byte, 9)
		b[0] = 0xff
		binary.LittleEndian.PutUint64(b[1:], val)
		return b
	}
}

func decodeVarInt(r io.Reader) (uint64, error) {
	var prefix [1]byte
	if _, err := r.Read(prefix[:]); err != nil {
		return 0, err
	}

	switch prefix[0] {
	case 0xfd:
		var v uint16
		err := binary.Read(r, binary.LittleEndian, &v)
		return uint64(v), err
	case 0xfe:
		var v uint32
		err := binary.Read(r, binary.LittleEndian, &v)
		return uint64(v), err
	case 0xff:
		var v uint64
		err := binary.Read(r, binary.LittleEndian, &v)
		return v, err
	default:
		return uint64(prefix[0]), nil
	}
}

func appendVarInt(buf []byte, val uint64) []byte {
	switch {
	case val < 0xfd:
		return append(buf, byte(val))
	case val <= 0xffff:
		buf = append(buf, 0xfd)
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(val))
		return append(buf, b...)
	case val <= 0xffffffff:
		buf = append(buf, 0xfe)
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(val))
		return append(buf, b...)
	default:
		buf = append(buf, 0xff)
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, val)
		return append(buf, b...)
	}
}

func appendVarBytes(buf []byte, data []byte) []byte {
	buf = appendVarInt(buf, uint64(len(data)))
	return append(buf, data...)
}
