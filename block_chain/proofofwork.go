package block_chain

/*
	简单共识机制建立  源于比特币的PoW机制
*/

import (
	"block/constcoe"
	"block/utils"
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

func (b *Block) GetTarget() []byte {

	target := big.NewInt(1)

	//Lsh函数就是向左移位，移的越多目标难度值越大，哈希取值落在的空间就更多也就越容易找到符合条件的nonce
	target.Lsh(target, uint(256-constcoe.Difficulty))
	return target.Bytes()
}

// GetBase4Nonce 每次输入nonce 对应的区块hash值会变化
func (b *Block) GetBase4Nonce(nonce int64) []byte {

	//bytes.Join可以将多个字节串连接，第二个参数是将字节串连接时的分隔符，这里设置为[]byte{}即为空
	data := bytes.Join([][]byte{
		utils.ToHexInt(b.Timestamp),
		b.PrevHash,
		utils.ToHexInt(nonce),
		b.Target,
		b.BackTransactionSummary(),
	},
		[]byte{},
	)
	return data
}

// FindNonce 寻找合适的nonce值
// nonce不过是从0开始取的整数而已，随着不断尝试，
// 每次失败nonce就加1直到由当前nonce得到的区块哈希转化为数值小于目标难度值为止
func (b *Block) FindNonce() int64 {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	var nonce int64
	nonce = 0
	intTarget.SetBytes(b.Target)

	for nonce < math.MaxInt64 {
		data := b.GetBase4Nonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&intTarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce
}

// ValidatePow 验证选出来的节点
func (b *Block) ValidatePow() bool {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	intTarget.SetBytes(b.Target)
	data := b.GetBase4Nonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	if intHash.Cmp(&intTarget) == -1 {
		return true
	}
	return false
}
