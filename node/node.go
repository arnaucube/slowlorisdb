package node

import (
	"crypto/ecdsa"
	"time"

	"github.com/arnaucube/slowlorisdb/core"
)

// Node
type Node struct {
	PrivK      *ecdsa.PrivateKey
	Addr       core.Address
	Bc         *core.Blockchain
	Miner      bool // indicates if the node is running as a miner
	PendingTxs []core.Tx
}

// NewNode creates a new node
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

// SignBlock performs a signature of a byte array with the node private key
func (node *Node) Sign(m []byte) (*core.Signature, error) {
	return core.Sign(node.PrivK, m)
}

// SignBlock performs a signature of a block with the node private key
func (node *Node) SignBlock(block *core.Block) error {
	block.CalculateHash()
	sig, err := core.Sign(node.PrivK, block.Hash[:])
	if err != nil {
		return err
	}
	block.Signature = sig.Bytes()
	return nil
}

// AddToPendingTxs adds a transaction the the node.PendingTxs
func (node *Node) AddToPendingTxs(tx core.Tx) {
	node.PendingTxs = append(node.PendingTxs, tx)
}

// BlockFromPendingTxs creates a new block from the pending transactions
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

// NewBlock creates a new block with the given txs
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

// CreateGenesis creates the genesis block
// pubK is the wallet where the first coins will be created
// amount is the amount of coins that will be created
func (node *Node) CreateGenesis(pubK *ecdsa.PublicKey, amount uint64) (*core.Block, error) {
	in := core.Input{
		TxId:  core.GenesisHashTxInput,
		Vout:  0,
		Value: amount,
	}
	var ins []core.Input
	ins = append(ins, in)

	out := core.Output{
		Value: amount,
	}
	var outs []core.Output
	outs = append(outs, out)

	tx := core.Tx{
		From:       &ecdsa.PublicKey{},
		To:         pubK,
		InputCount: uint64(0),
		Inputs:     []core.Input{},
		Outputs:    outs,
		Signature:  []byte{},
	}

	// calculate TxId
	// tx.CalculateTxId()
	tx.TxId = core.GenesisHashTxInput
	// sign transaction

	var txs []core.Tx
	txs = append(txs, tx)

	block := &core.Block{
		Height:    node.Bc.LastBlock.Height + 1,
		PrevHash:  node.Bc.LastBlock.Hash,
		Txs:       txs,
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

// ParseReceivedBlock is just a caller of node.Bc.AddBlock() at the Node level
func (node *Node) ParseReceivedBlock(block *core.Block) error {
	err := node.Bc.AddBlock(block)
	return err
}
