package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"

	"github.com/mr-tron/base58"
)

const (
	DbPath      = "../tmp/blocks"
	DbFile      = "../tmp/blocks/MANIFEST"
	GenesisData = "First Transaction from Genesis"
)

func ErrorHandling(err error) {
	if err != nil {
		log.Panic(err);
	}
}

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input);

	return []byte(encode);
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]));
	if err != nil {
		log.Panic(err);
	}

	return decode;
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer);
	err := binary.Write(buff, binary.BigEndian, num);
	ErrorHandling(err);

	return buff.Bytes();
}

func DBexists() bool {
	if _, err := os.Stat(DbFile); os.IsNotExist(err) {
		return false
	}

	return true
}