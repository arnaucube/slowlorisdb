package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	block := &Block{
		PrevHash:  HashBytes([]byte("prevhash")),
		NextHash:  HashBytes([]byte("nextHash")),
		Txs:       []Tx{},
		Miner:     Address(HashBytes([]byte("addrfromminer"))),
		Timestamp: time.Now(),
		Nonce:     0,
		Hash:      HashBytes([]byte("blockhash")),
	}

	block, err := BlockFromBytes(block.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, block.Bytes(), block.Bytes())

	difficulty := 2
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
		NextHash:  HashBytes([]byte("nextHash")),
		Txs:       []Tx{},
		Miner:     Address(HashBytes([]byte("addrfromminer"))),
		Timestamp: time.Now(),
		Nonce:     0,
		Hash:      HashBytes([]byte("blockhash")),
	}

	block2, err := BlockFromBytes(block.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, block2.Bytes(), block.Bytes())

	difficulty := 2
	nonce, err := CalculatePoW(block, difficulty)
	assert.Nil(t, err)
	block.Nonce = nonce
	h := HashBytes(block.Bytes())

	// CheckPoW
	assert.True(t, CheckPoW(h, difficulty))
}

func TestTx(t *testing.T) {
	addr0 := Address(HashBytes([]byte("addr0")))
	addr1 := Address(HashBytes([]byte("addr1")))

	tx := NewTx(addr0, addr1, []Input{}, []Output{})

	assert.Equal(t, tx.From, addr0)
	assert.Equal(t, tx.To, addr1)
}
