package node

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/arnaucube/slowlorisdb/core"
	"github.com/arnaucube/slowlorisdb/db"
	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	privK, err := core.NewKey()
	assert.Nil(t, err)

	dif := uint64(1)
	bc := core.NewBlockchain(db, dif)
	node, err := NewNode(privK, bc, true)
	assert.Nil(t, err)

	assert.Equal(t, node.Addr, core.AddressFromPrivK(node.PrivK))
}

func TestNodeSignature(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	privK, err := core.NewKey()
	assert.Nil(t, err)

	dif := uint64(1)
	bc := core.NewBlockchain(db, dif)
	node, err := NewNode(privK, bc, true)
	assert.Nil(t, err)

	m := []byte("test")
	sig, err := node.Sign(m)
	assert.Nil(t, err)
	pubK := node.PrivK.PublicKey
	assert.True(t, core.VerifySignature(&pubK, m, *sig))
}

func TestBlockFromPendingTxs(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	privK, err := core.NewKey()
	assert.Nil(t, err)

	dif := uint64(1)
	bc := core.NewBlockchain(db, dif)
	node, err := NewNode(privK, bc, true)
	assert.Nil(t, err)

	privK0, err := core.NewKey()
	assert.Nil(t, err)
	pubK0 := privK0.PublicKey
	privK1, err := core.NewKey()
	assert.Nil(t, err)
	pubK1 := privK1.PublicKey
	tx := core.NewTx(&pubK0, &pubK1, []core.Input{}, []core.Output{})
	node.AddToPendingTxs(*tx)
	block, err := node.BlockFromPendingTxs()
	assert.Nil(t, err)
	fmt.Println("h", block.Hash)
	assert.True(t, core.CheckBlockPoW(block, node.Bc.Difficulty))
	assert.True(t, node.Bc.VerifyBlock(block))
}

func TestBlockFromPendingTxsIteration(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	privK, err := core.NewKey()
	assert.Nil(t, err)

	dif := uint64(1)
	bc := core.NewBlockchain(db, dif)
	node, err := NewNode(privK, bc, true)
	assert.Nil(t, err)

	privK0, err := core.NewKey()
	assert.Nil(t, err)
	pubK0 := privK0.PublicKey
	privK1, err := core.NewKey()
	assert.Nil(t, err)
	pubK1 := privK1.PublicKey

	for i := 0; i < 10; i++ {
		tx := core.NewTx(&pubK0, &pubK1, []core.Input{}, []core.Output{})
		node.AddToPendingTxs(*tx)
	}
	block, err := node.BlockFromPendingTxs()
	assert.Nil(t, err)
	assert.True(t, core.CheckBlockPoW(block, node.Bc.Difficulty))
	assert.True(t, node.Bc.VerifyBlock(block))
}

func TestFromGenesisToTenBlocks(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	privK, err := core.NewKey()
	assert.Nil(t, err)

	dif := uint64(1)
	bc := core.NewBlockchain(db, dif)
	node, err := NewNode(privK, bc, true)
	assert.Nil(t, err)

	// create the genesis block
	genesisBlock, err := node.CreateGenesis()
	assert.Nil(t, err)
	assert.NotEqual(t, genesisBlock.Signature, core.Signature{})
	assert.NotEqual(t, genesisBlock.Hash, core.Hash{})

	// add the genesis block into the blockchain
	err = node.Bc.AddBlock(genesisBlock)
	assert.Nil(t, err)
	assert.NotEqual(t, genesisBlock.Hash, core.Hash{})
	assert.Equal(t, genesisBlock.Hash, node.Bc.LastBlock.Hash)

	// TODO add another block
	block := node.NewBlock([]core.Tx{})
	err = node.Bc.AddBlock(block)
	assert.Nil(t, err)
}
