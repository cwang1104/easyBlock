package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

// ToHexInt 将int64转为[]byte
func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	Handle(err)
	return buff.Bytes()
}

// Handle 错误处理
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
