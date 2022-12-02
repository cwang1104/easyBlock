package block_chain

import (
	"block/transaction"
	"block/utils"
	"bytes"
	"crypto/sha256"
	"time"
)

// Block 单一区块
type Block struct {
	Timestamp int64  //时间戳
	Hash      []byte //当前区块的hash值
	PrevHash  []byte //前一个区块的hash值

	Target []byte //目标难度值
	Nonce  int64  // 节点寻找到的作为卷王的证据

	//Data []byte // 数据
	Transactions []*transaction.Transaction
}

func (b *Block) SetHash() {
	//将区块各项属性串联起来的字符串
	information := bytes.Join([][]byte{
		utils.ToHexInt(b.Timestamp),
		b.PrevHash,
		b.Target,
		utils.ToHexInt(b.Nonce),
		b.BackTransactionSummary(),
	}, []byte{})

	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

// CreateBlock 创建区块
func CreateBlock(prevHash []byte, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevHash, []byte{}, 0, txs}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

// GenesisBlock 创世区块，（链表头节点）
func GenesisBlock() *Block {
	tx := transaction.BaseTx([]byte("start"))
	return CreateBlock([]byte{}, []*transaction.Transaction{tx})
}

// BackTransactionSummary 交易信息的序列化
func (b *Block) BackTransactionSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}

	summary := bytes.Join(txIDs, []byte{})
	return summary
}
