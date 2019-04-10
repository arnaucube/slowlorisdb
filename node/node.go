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

func (node *Node) SignBlock(block *core.Block) error {
	block.CalculateHash()
	sig, err := core.Sign(node.PrivK, block.Hash[:])
	if err != nil {
		return err
	}
	block.Signature = sig.Bytes()
	return nil
}

func (node *Node) AddToPendingTxs(tx core.Tx) {
	node.PendingTxs = append(node.PendingTxs, tx)
}

func (node *Node) BlockFromPendingTxs() (*core.Block, error) {
	block, err := node.NewBlock(node.PendingTxs)
	if err != nil {
		return nil, err
	}
	block.PrevHash = node.Bc.LastBlock.Hash
	err = block.CalculatePoW(node.Bc.Difficulty)
	if err != nil {
		return nil, err
	}
	err = node.SignBlock(block)
	if err != nil {
		return nil, err
	}
	block.CalculateHash()

	return block, nil
}

func (node *Node) NewBlock(txs []core.Tx) (*core.Block, error) {
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
	block.CalculateHash()
	err := node.SignBlock(block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (node *Node) CreateGenesis() (*core.Block, error) {
	block := &core.Block{
		Height:    node.Bc.LastBlock.Height + 1,
		PrevHash:  node.Bc.LastBlock.Hash,
		Txs:       []core.Tx{},
		Miner:     node.Addr,
		MinerPubK: &node.PrivK.PublicKey,
		Timestamp: time.Now(),
		Nonce:     uint64(0),
		Hash:      core.Hash{},
		Signature: []byte{},
	}

	block.CalculateHash()
	err := node.SignBlock(block)
	if err != nil {
		return nil, err
	}
	return block, nil
}
