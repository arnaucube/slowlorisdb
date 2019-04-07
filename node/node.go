package node

import (
	"crypto/ecdsa"

	"github.com/arnaucube/slowlorisdb/core"
	"github.com/arnaucube/slowlorisdb/db"
)

type Node struct {
	PrivK      *ecdsa.PrivateKey
	Addr       core.Address
	Bc         *core.Blockchain
	PendingTxs []core.Tx
}

func NewNode(db *db.Db, dif uint64) (*Node, error) {
	privK, err := core.NewKey()
	if err != nil {
		return nil, err
	}
	addr := core.AddressFromPrivK(privK)

	bc := core.NewBlockchain(db, dif)

	node := &Node{
		PrivK: privK,
		Addr:  addr,
		Bc:    bc,
	}
	return node, nil
}

func (node *Node) Sign(m []byte) (*core.Signature, error) {
	return core.Sign(node.PrivK, m)
}

func (node *Node) SignBlock(block *core.Block) (*core.Signature, error) {
	return core.Sign(node.PrivK, block.Hash[:])
}

func (node *Node) AddToPendingTxs(tx core.Tx) {
	node.PendingTxs = append(node.PendingTxs, tx)
}

func (node *Node) BlockFromPendingTxs() (*core.Block, error) {
	block := node.Bc.NewBlock(node.PendingTxs)
	err := block.CalculatePoW(node.Bc.Difficulty)
	if err != nil {
		return nil, err
	}
	sig, err := node.SignBlock(block)
	if err != nil {
		return nil, err
	}
	block.Signature = sig.Bytes()
	return block, nil
}
