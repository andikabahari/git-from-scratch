package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"io"
)

// Notes about packfiles:
// - The first 4 bytes spell "PACK".
// - The next 4 bytes are the version number.
// - The next 4 bytes are the number of objects.
// - The final 20 bytes are a SHA-1 checksum.
// - The rest is the data.

type ObjectType int

const (
	OBJ_COMMIT ObjectType = iota + 1
	OBJ_TREE
	OBJ_BLOB
	OBJ_TAG
	_
	OBJ_OFS_DELTA
	OBJ_REF_DELTA
)

type pack struct {
	Signature  []byte
	Version    int
	NumObjects int
	Checksum   [20]byte
	Data       []byte
}

func parsePack(b []byte) (*pack, error) {
	p := pack{}
	p.Signature = b[:4]

	v, ok := binary.Varint(b[4:8])
	if ok <= 0 {
		return nil, errors.New("cannot read version number")
	}
	p.Version = int(v)

	n, ok := binary.Varint(b[8:12])
	if ok <= 0 {
		return nil, errors.New("cannot read number of objects")
	}
	p.NumObjects = int(n)

	checksum := sha1.Sum(b[:len(b)-20])
	if !bytes.Equal(checksum[:], b[len(b)-20:]) {
		return nil, errors.New("invalid checksum")
	}
	p.Checksum = checksum

	p.Data = b[12 : len(b)-20]
	return &p, nil
}

func unpack(dir string, p *pack) error {
	r := bytes.NewReader(p.Data)

	for r.Len() > 0 {
		size, objType := parseHeader(r)
		_ = size

		// Remove this
		io.ReadAll(r)

		switch objType {
		case OBJ_COMMIT:
			// TODO: decompress commit

		case OBJ_TREE:
			// TODO: decompress tree

		case OBJ_BLOB:
			// TODO: decompress blob

		case OBJ_OFS_DELTA:
			// TODO: decompress ofs delta

		case OBJ_REF_DELTA:
			// TODO: decompress ref delta
		}
	}

	return nil
}

func parseHeader(r *bytes.Reader) (int, ObjectType) {
	b, err := r.ReadByte()
	if err != nil {
		panic(err)
	}

	objType := ObjectType((b & 0b01110000) >> 4)

	var val int
	haed := int(b) & 0b00001111
	val = haed
	if int(b)&128 != 0 {
		tail, err := parseVarint(r)
		if err != nil {
			panic(err)
		}
		val += (tail << 4) + haed
	}

	return val, objType
}

// Source: https://github.com/ChimeraCoder/gitgo/blob/master/delta.go
func parseVarint(r io.Reader) (int, error) {
	_bytes := make([]byte, 1)
	_, err := r.Read(_bytes)
	if err != nil {
		return 0, err
	}

	_byte := _bytes[0]
	MSB := (_byte & 128) // will be either 128 or 0
	var objectSize = int((uint(_byte) & 127))
	var shift uint
	for MSB > 0 {
		shift += 7
		_bytes := make([]byte, 1)
		_, err := r.Read(_bytes)
		if err != nil {
			return 0, err
		}
		_byte := _bytes[0]
		MSB = _byte & 128
		objectSize += int((uint(_byte) & 127) << shift)
	}

	return objectSize, nil
}
