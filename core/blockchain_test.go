package core

import (
	"io/ioutil"
	"testing"

	"github.com/arnaucube/slowlorisdb/db"
	"github.com/stretchr/testify/assert"
)

func TestBlockchainDataStructure(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	bc := NewBlockchain(db, uint64(1))
	block := bc.NewBlock([]Tx{})

	block2, err := BlockFromBytes(block.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, block2.Bytes(), block.Bytes())
}

func TestGetBlock(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	bc := NewBlockchain(db, uint64(1))

	block := bc.NewBlock([]Tx{})
	assert.Equal(t, block.Height, uint64(1))

	err = bc.AddBlock(block)
	assert.Nil(t, err)

	block2, err := bc.GetBlock(block.Hash)
	assert.Nil(t, err)
	assert.Equal(t, block.Bytes(), block2.Bytes())
}

func TestGetPrevBlock(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	bc := NewBlockchain(db, uint64(1))

	for i := 0; i < 10; i++ {
		block := bc.NewBlock([]Tx{})
		block.CalculatePoW(bc.Difficulty)
		assert.Equal(t, block.Height, uint64(i+1))

		err = bc.AddBlock(block)
		assert.Nil(t, err)
		assert.Equal(t, bc.LastBlock.Height, block.Height)
	}
	block9, err := bc.GetPrevBlock(bc.LastBlock.Hash)
	assert.Nil(t, err)
	assert.Equal(t, block9.Height, uint64(9))

	block8, err := bc.GetPrevBlock(block9.Hash)
	assert.Nil(t, err)
	assert.Equal(t, block8.Height, uint64(8))

	currentBlock := bc.LastBlock
	for err == nil {
		currentBlock, err = bc.GetPrevBlock(currentBlock.Hash)
	}
	assert.Equal(t, err.Error(), "This was the oldest block")
}

func TestAddBlockWithTx(t *testing.T) {
	addr0 := Address(HashBytes([]byte("addr0")))
	addr1 := Address(HashBytes([]byte("addr1")))

	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	bc := NewBlockchain(db, uint64(1))

	var txs []Tx
	tx := NewTx(addr0, addr1, []Input{}, []Output{})
	txs = append(txs, *tx)
	block := bc.NewBlock(txs)

	block2, err := BlockFromBytes(block.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, block2.Bytes(), block.Bytes())
}
