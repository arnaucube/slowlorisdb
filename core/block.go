package core

import (
	"encoding/json"
	"time"
)

type Input struct {
}
type Output struct {
}

// Tx holds the data structure of a transaction
type Tx struct {
	From       Address
	To         Address
	InputCount uint64
	Inputs     []Input
	Outputs    []Output
}

func NewTx(from, to Address, in []Input, out []Output) *Tx {
	tx := &Tx{
		From:       from,
		To:         to,
		InputCount: uint64(len(in)),
		Inputs:     in,
		Outputs:    out,
	}
	return tx
}

// Block holds the data structure for the block
type Block struct {
	Height    uint64
	PrevHash  Hash
	NextHash  Hash
	Txs       []Tx
	Miner     Address
	Timestamp time.Time
	Nonce     uint64
	Hash      Hash
	Signature []byte
}

// Bytes outputs a byte array containing the data of the Block
func (blk Block) Bytes() []byte {
	b, _ := json.Marshal(blk)
	return b
}

func (blk *Block) GetNonce() uint64 {
	return blk.Nonce
}

func (blk *Block) IncrementNonce() {
	blk.Nonce++
}
func (block *Block) CalculatePoW(difficulty int) error {
	hash := HashBytes(block.Bytes())
	for !CheckPoW(hash, difficulty) {
		block.IncrementNonce()
		hash = HashBytes(block.Bytes())
	}
	block.Hash = hash
	return nil
}

func BlockFromBytes(b []byte) (*Block, error) {
	var block *Block
	err := json.Unmarshal(b, &block)
	return block, err
}
