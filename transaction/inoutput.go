package transaction

import (
	"bytes"
)

/*
TxOutput包含Value与ToAddress两个属性，前者是转出的资产值，后者是资产的接收者的地址。
TxInput包含的TxID用于指明支持本次交易的前置交易信息，
OutIdx是具体指明是前置交易信息中的第几个Output，
FromAddress就是资产转出者的地址。
*/

// TxOutput 交易信息的Output
type TxOutput struct {
	Value     int
	ToAddress []byte
}

// TxInput 交易信息的Input
type TxInput struct {
	TxID        []byte
	OutIdx      int
	FromAddress []byte
}

func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.FromAddress, address)
}

func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.ToAddress, address)
}
