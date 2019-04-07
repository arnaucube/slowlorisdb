package node

import (
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

	node, err := NewNode(db, uint64(1))
	assert.Nil(t, err)

	assert.Equal(t, node.Addr, core.AddressFromPrivK(node.PrivK))
}

func TestNodeSignature(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := db.New(dir)
	assert.Nil(t, err)

	node, err := NewNode(db, uint64(1))
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

	node, err := NewNode(db, uint64(1))
	assert.Nil(t, err)

	addr0 := core.Address(core.HashBytes([]byte("addr0")))
	addr1 := core.Address(core.HashBytes([]byte("addr1")))
	tx := core.NewTx(addr0, addr1, []core.Input{}, []core.Output{})
	node.AddToPendingTxs(*tx)
	block, err := node.BlockFromPendingTxs()
	assert.True(t, core.CheckBlockPoW(block, node.Bc.Difficulty))
	// TODO add VerifyBlockSignature
}

// func TestBlockFromPendingTxsIteration(t *testing.T) {
//         dir, err := ioutil.TempDir("", "db")
//         assert.Nil(t, err)
//         db, err := db.New(dir)
//         assert.Nil(t, err)
//
//         node, err := NewNode(db, uint64(1))
//         assert.Nil(t, err)
//
//         addr0 := core.Address(core.HashBytes([]byte("addr0")))
//         addr1 := core.Address(core.HashBytes([]byte("addr1")))
//         for i := 0; i < 10; i++ {
//                 tx := core.NewTx(addr0, addr1, []core.Input{}, []core.Output{})
//
//         }
// }
