package encode

import (
	"bytes"
	"fmt"

	"github.com/akamensky/base58"
	"golang.org/x/crypto/blake2b"
)

var (
	prefix = []byte("SS58PRE")
)

func EncodeAddress(pubKey []byte, ss58Prefix int8) string {
	var raw []byte
	raw = append([]byte{byte(ss58Prefix)}, pubKey...)
	checksum := blake2b.Sum512(append(prefix, raw...))
	address := base58.Encode(append(raw, checksum[0:2]...))
	return address
}

func DecodeAddress(address string) (publicKey []byte, ss58Prefix int8, err error) {
	raw, err := base58.Decode(address)
	if err != nil {
		return nil, 0, err
	}
	actualChecksum := raw[len(raw)-4:]
	checksum := blake2b.Sum512(raw[:len(raw)])
	if bytes.Equal(actualChecksum, checksum[:]) {
		return nil, 0, fmt.Errorf("Invalid checksum. actualChecksum: %s, checksum: %s", actualChecksum, checksum)
	}
	return raw[1:33], int8(raw[0]), nil
}
