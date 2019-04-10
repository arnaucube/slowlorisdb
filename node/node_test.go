package node

import (
	"crypto/ecdsa"
	"io/ioutil"
	"testing"

	"github.com/arnaucube/slowlorisdb/core"
	"github.com/arnaucube/slowlorisdb/db"
	"github.com/stretchr/testify/assert"
)

func newTestPoABlockchain() (*ecdsa.PrivateKey, *core.Blockchain, error) {
	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		return nil, nil, err
	}
	db, err := db.New(dir)
	if err != nil {
		return nil, nil, err
	}

	privK, err := core.NewKey()
	if err != nil {
		return nil, nil, err
	}

	var authNodes []*ecdsa.PublicKey
	authNodes = append(authNodes, &privK.PublicKey)

	bc := core.NewPoABlockchain(db, authNodes)
	return privK, bc, nil
}
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
	privK, bc, err := newTestPoABlockchain()
	assert.Nil(t, err)

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
	assert.True(t, core.CheckBlockPoW(block, node.Bc.Difficulty))
	assert.True(t, node.Bc.VerifyBlock(block))
}

func TestBlockFromPendingTxsIteration(t *testing.T) {
	privK, bc, err := newTestPoABlockchain()
	assert.Nil(t, err)

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
	privK, bc, err := newTestPoABlockchain()
	assert.Nil(t, err)

	node, err := NewNode(privK, bc, true)
	assert.Nil(t, err)

	// create the genesis block
	genesisBlock, err := node.CreateGenesis(&privK.PublicKey, uint64(100))
	assert.Nil(t, err)
	assert.NotEqual(t, genesisBlock.Signature, core.Signature{})
	assert.NotEqual(t, genesisBlock.Hash, core.Hash{})
	assert.True(t, node.Bc.VerifyBlock(genesisBlock))

	// add the genesis block into the blockchain
	err = node.Bc.AddBlock(genesisBlock)
	assert.Nil(t, err)
	assert.NotEqual(t, genesisBlock.Hash, core.Hash{})
	assert.Equal(t, genesisBlock.Hash, node.Bc.LastBlock.Hash)

	// add another tx sending coins to the pubK0
	privK0, err := core.NewKey()
	assert.Nil(t, err)
	pubK0 := privK0.PublicKey

	var ins []core.Input
	in := core.Input{
		TxId:  genesisBlock.Txs[0].TxId,
		Vout:  0,
		Value: 100,
	}
	ins = append(ins, in)
	var outs []core.Output
	out0 := core.Output{
		Value: 10,
	}
	out1 := core.Output{
		Value: 90,
	}
	outs = append(outs, out0)
	outs = append(outs, out1)
	tx := core.NewTx(&privK.PublicKey, &pubK0, ins, outs)

	// verify tx
	txVerified := core.CheckTx(tx)
	assert.True(t, txVerified)

	// then create a new block with the tx and add it to the blockchain
	var txs []core.Tx
	txs = append(txs, *tx)
	block, err := node.NewBlock(txs)
	assert.Nil(t, err)
	err = node.Bc.AddBlock(block)
	assert.Nil(t, err)

	// add another tx sending coins to the pubK1
	privK1, err := core.NewKey()
	assert.Nil(t, err)
	pubK1 := privK1.PublicKey

	ins = []core.Input{}
	in = core.Input{
		TxId:  block.Txs[0].TxId,
		Vout:  0,
		Value: 10,
	}
	ins = append(ins, in)
	outs = []core.Output{}
	out0 = core.Output{
		Value: 10,
	}
	outs = append(outs, out0)
	tx = core.NewTx(&pubK0, &pubK1, ins, outs)

	// verify tx
	txVerified = core.CheckTx(tx)
	assert.True(t, txVerified)

	// then create a new block with the tx and add it to the blockchain
	txs = []core.Tx{}
	txs = append(txs, *tx)
	block, err = node.NewBlock(txs)
	assert.Nil(t, err)
	err = node.Bc.AddBlock(block)
	assert.Nil(t, err)
}
