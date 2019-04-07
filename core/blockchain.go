package core

import (
	"errors"
	"time"

	"github.com/arnaucube/slowlorisdb/db"
)

type Blockchain struct {
	Id        []byte // Id allows to have multiple blockchains
	Genesis   Hash
	LastBlock *Block
	db        *db.Db
}

func NewBlockchain(database *db.Db) *Blockchain {
	blockchain := &Blockchain{
		Id:        []byte{},
		Genesis:   Hash{},
		LastBlock: &Block{},
		db:        database,
	}
	return blockchain
}

func (bc *Blockchain) NewBlock(txs []Tx) *Block {
	block := &Block{
		Height:    bc.GetHeight() + 1,
		PrevHash:  bc.LastBlock.Hash,
		Txs:       txs,
		Miner:     Address{}, // TODO put the node address
		Timestamp: time.Now(),
		Nonce:     0,
		Hash:      Hash{},
		Signature: []byte{},
	}
	return block
}

func (bc *Blockchain) GetHeight() uint64 {
	return bc.LastBlock.Height
}

func (bc *Blockchain) GetLastBlock() *Block {
	return bc.LastBlock
}

func (bc *Blockchain) AddBlock(block *Block) error {
	bc.LastBlock = block
	err := bc.db.Put(block.Hash[:], block.Bytes())
	return err
}

func (bc *Blockchain) GetBlock(hash Hash) (*Block, error) {
	v, err := bc.db.Get(hash[:])
	if err != nil {
		return nil, err
	}
	block, err := BlockFromBytes(v)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (bc *Blockchain) GetPrevBlock(hash Hash) (*Block, error) {
	currentBlock, err := bc.GetBlock(hash)
	if err != nil {
		return nil, err
	}
	if currentBlock.PrevHash.IsZero() {
		return nil, errors.New("This was the oldest block")
	}
	prevBlock, err := bc.GetBlock(currentBlock.PrevHash)
	if err != nil {
		return nil, err
	}

	return prevBlock, nil
}
