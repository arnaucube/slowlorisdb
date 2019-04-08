package core

import (
	"crypto/ecdsa"
	"encoding/json"
	"time"
)

// Block holds the data structure for the block
type Block struct {
	Height    uint64
	PrevHash  Hash
	Txs       []Tx
	Miner     Address
	MinerPubK *ecdsa.PublicKey // tmp meanwhile no ecrecover function
	Timestamp time.Time
	Nonce     uint64
	Hash      Hash
	Signature []byte
}

func (block Block) Copy() *Block {
	return &Block{
		Height:    block.Height,
		PrevHash:  block.PrevHash,
		Txs:       block.Txs,
		Miner:     block.Miner,
		MinerPubK: block.MinerPubK,
		Timestamp: block.Timestamp,
		Nonce:     block.Nonce,
		Hash:      block.Hash,
		Signature: block.Signature,
	}
}

// Bytes outputs a byte array containing the data of the Block
func (blk Block) Bytes() []byte {
	// TODO add parser, to use minimum amount of bytes
	b, _ := json.Marshal(blk)
	return b
}

func (block *Block) CalculateHash() {
	blockCopy := block.Copy()
	blockCopy.Hash = Hash{}
	blockCopy.Signature = []byte{}
	hash := HashBytes(blockCopy.Bytes())
	block.Hash = hash
}

func (blk *Block) GetNonce() uint64 {
	return blk.Nonce
}

func (blk *Block) IncrementNonce() {
	blk.Nonce++
}
func (block *Block) CalculatePoW(difficulty uint64) error {
	blockCopy := block.Copy()
	blockCopy.Hash = Hash{}

	hash := HashBytes(blockCopy.Bytes())
	for !CheckPoW(hash, difficulty) {
		blockCopy.IncrementNonce()
		hash = HashBytes(blockCopy.Bytes())
	}
	block.Hash = hash
	block.Nonce = blockCopy.Nonce
	return nil
}

func CheckBlockPoW(block *Block, difficulty uint64) bool {
	blockCopy := block.Copy()
	blockCopy.Hash = Hash{}
	blockCopy.Signature = []byte{}
	return CheckPoW(HashBytes(blockCopy.Bytes()), difficulty)
}

func BlockFromBytes(b []byte) (*Block, error) {
	var block *Block
	err := json.Unmarshal(b, &block)
	return block, err
}
