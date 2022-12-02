package transaction

import (
	"block/constcoe"
	"block/utils"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Transaction struct {
	ID      []byte     //自身的ID值
	Inputs  []TxInput  //标记支持我们本次转账的前置的交易信息的TxOutput
	Outputs []TxOutput //记录我们本次转账的amount和Reciever
}

// BaseTx 初始交易信息
func BaseTx(toaddress []byte) *Transaction {
	txIn := TxInput{[]byte{}, -1, []byte{}}
	txOut := TxOutput{constcoe.InitCoin, toaddress}
	tx := Transaction{[]byte("this is the base tx"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}

// SetID 设置交易信息的ID值
func (tx *Transaction) SetID() {
	tx.ID = tx.TxHash()
}

// TxHash 设置生成hash功能
func (tx *Transaction) TxHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}
