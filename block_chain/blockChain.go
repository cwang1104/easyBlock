package block_chain

import (
	"block/transaction"
	"block/utils"
	"encoding/hex"
	"fmt"
)

// BlockChain 区块链
type BlockChain struct {
	Blocks []*Block
}

// CreateBlockChain 创建区块链
func CreateBlockChain() *BlockChain {
	blockChain := BlockChain{}
	blockChain.Blocks = append(blockChain.Blocks, GenesisBlock())
	return &blockChain
}

// AddBlock 增加区块
func (bc *BlockChain) AddBlock(txs []*transaction.Transaction) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// FindUnspentTransactions 寻找前置可用交易信息OutPut
func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {

	//返回包含指定地址的可用交易信息的切片
	var unSpentTxs []transaction.Transaction

	//记录遍历区块链时那些已经被使用的交易信息的Output,
	// key值为交易信息的ID值（需要转成string），value值为Output在该交易信息中的序号。
	spentTxs := make(map[string][]int)

	for idx := len(bc.Blocks) - 1; idx >= 0; idx-- { //遍历所有区块，从后往前
		block := bc.Blocks[idx]
		for _, tx := range block.Transactions { //遍历每一个区块的交易信息
			txID := hex.EncodeToString(tx.ID)
		IterOutputs:
			for outIdx, out := range tx.Outputs {
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIdx {
							continue IterOutputs //如果改OutPut在spentTxs中就被跳过，说明已经被消费
						}
					}
				}

				//确认ToAddress正确与否，正确就是我们要找的可用交易信息
				if out.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
			}

			//检查当前交易信息是否为Base Transaction
			//如果不是就检查当前交易信息的input中是否包含目标地址，
			//有的话就将指向的Output信息加入到spentTxs中
			if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if in.FromAddressRight(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
					}
				}
			}
		}
	}
	return unSpentTxs
}

// FindUTXOs 找到一个地址的所有UTXO以及该地址对应的资产总和
func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unSpentOuts := make(map[string]int)
	unSpentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unSpentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) {
				accumulated += out.Value
				unSpentOuts[txID] = outIdx
				continue Work
			}
		}
	}
	return accumulated, unSpentOuts
}

// FindSpendableOutputs 找到资产总量大于本次交易转账额的一部分UTXO就行
func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				if accumulated >= amount {
					break Work
				}
				continue Work
			}
		}
	}
	return accumulated, unspentOuts
}

// CreateTransaction 创建交易信息
func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		fmt.Println("Not enough coins")
		return &transaction.Transaction{}, false
	}

	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)

		input := transaction.TxInput{TxID: txID, OutIdx: outidx, FromAddress: from}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, transaction.TxOutput{Value: amount, ToAddress: to})
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{Value: acc - amount, ToAddress: from})
	}

	tx := transaction.Transaction{ID: nil, Inputs: inputs, Outputs: outputs}
	tx.SetID()
	return &tx, true
}

// Mine 不希望再使用AddBlock直接添加区块进入到区块链中，而是预留一个函数模拟一下从交易信息池中获取交易信息打包并挖矿这个过程。
func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
	bc.AddBlock(txs)
}
