package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
)

// Notes about packfiles:
// - The first 4 bytes spell "PACK".
// - The next 4 bytes are the version number.
// - The next 4 bytes are the number of objects.
// - The final 20 bytes are a SHA-1 checksum.
// - The rest is the data.

type pack struct {
	Version    int
	NumObjects int
	Checksum   [20]byte
	Data       []byte
}

func parsePack(b []byte) (pack, error) {
	p := pack{
		Data: make([]byte, len(b)-32),
	}

	checksum := sha1.Sum(b[:len(b)-20])
	if !bytes.Equal(checksum[:], b[len(b)-20:]) {
		return p, errors.New("invalid checksum")
	}
	p.Checksum = checksum

	v, ok := binary.Varint(b[4:8])
	if ok <= 0 {
		return p, errors.New("cannot read version number")
	}
	p.Version = int(v)

	n, ok := binary.Varint(b[8:12])
	if ok <= 0 {
		return p, errors.New("cannot read number of objects")
	}
	p.NumObjects = int(n)

	p.Data = b[12 : len(b)-20]
	return p, nil
}

func unpack(dir string, p *pack) error {
	// TODO: implement unpack function
	return nil
}
