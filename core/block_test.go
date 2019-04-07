package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	block := &Block{
		PrevHash:  HashBytes([]byte("prevhash")),
		Txs:       []Tx{},
		Miner:     Address(HashBytes([]byte("addrfromminer"))),
		Timestamp: time.Now(),
		Nonce:     0,
		Hash:      HashBytes([]byte("blockhash")),
	}

	blockParsed, err := BlockFromBytes(block.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, blockParsed.Bytes(), block.Bytes())

	blockCopy := block.Copy()
	assert.Equal(t, blockCopy.Bytes(), block.Bytes())

	difficulty := uint64(2)
	nonce, err := CalculatePoW(block, difficulty)
	assert.Nil(t, err)
	block.Nonce = nonce
	h := HashBytes(block.Bytes())

	// CheckPoW
	assert.True(t, CheckPoW(h, difficulty))
}

func TestNewBlock(t *testing.T) {
	block := &Block{
		PrevHash:  HashBytes([]byte("prevhash")),
		Txs:       []Tx{},
		Miner:     Address(HashBytes([]byte("addrfromminer"))),
		Timestamp: time.Now(),
		Nonce:     0,
		Hash:      HashBytes([]byte("blockhash")),
	}

	block2, err := BlockFromBytes(block.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, block2.Bytes(), block.Bytes())

	difficulty := uint64(2)
	nonce, err := CalculatePoW(block, difficulty)
	assert.Nil(t, err)
	block.Nonce = nonce
	h := HashBytes(block.Bytes())

	// CheckPoW
	assert.True(t, CheckPoW(h, difficulty))
}
