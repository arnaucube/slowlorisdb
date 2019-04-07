package node

import (
	"crypto/ecdsa"
	"time"

	"github.com/arnaucube/slowlorisdb/core"
)

type Node struct {
	PrivK      *ecdsa.PrivateKey
	Addr       core.Address
	Bc         *core.Blockchain
	Miner      bool // indicates if the node is running as a miner
	PendingTxs []core.Tx
}

func NewNode(privK *ecdsa.PrivateKey, bc *core.Blockchain, isMiner bool) (*Node, error) {
	addr := core.AddressFromPrivK(privK)

	node := &Node{
		PrivK:      privK,
		Addr:       addr,
		Bc:         bc,
		Miner:      isMiner,
		PendingTxs: []core.Tx{},
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
	block := node.NewBlock(node.PendingTxs)
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

func (node *Node) NewBlock(txs []core.Tx) *core.Block {
	block := &core.Block{
		Height:    node.Bc.GetHeight() + 1,
		PrevHash:  node.Bc.LastBlock.Hash,
		Txs:       txs,
		Miner:     core.Address{},
		MinerPubK: &node.PrivK.PublicKey,
		Timestamp: time.Now(),
		Nonce:     0,
		Hash:      core.Hash{},
		Signature: []byte{},
	}
	return block
}
