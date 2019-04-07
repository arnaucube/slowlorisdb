package node

import (
	"crypto/ecdsa"

	"github.com/arnaucube/slowlorisdb/core"
	"github.com/arnaucube/slowlorisdb/db"
)

type Node struct {
	PrivK *ecdsa.PrivateKey
	Addr  core.Address
	Bc    *core.Blockchain
}

func NewNode(db *db.Db) (*Node, error) {
	privK, err := core.NewKey()
	if err != nil {
		return nil, err
	}
	addr := core.AddressFromPrivK(privK)

	bc := core.NewBlockchain(db)

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
