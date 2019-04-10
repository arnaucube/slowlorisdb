package core

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/arnaucube/slowlorisdb/db"
)

type PoA struct {
	AuthMiners []*ecdsa.PublicKey
}

type Blockchain struct {
	Id         []byte // Id allows to have multiple blockchains
	Difficulty uint64
	Genesis    Hash
	LastBlock  *Block
	blockdb    *db.Db
	txdb       *db.Db
	PoA        PoA
}

func NewBlockchain(database *db.Db, dif uint64) *Blockchain {
	blockchain := &Blockchain{
		Id:         []byte{},
		Difficulty: dif,
		Genesis:    HashBytes([]byte("genesis")),
		LastBlock:  &Block{},
		blockdb:    database,
		PoA:        PoA{},
	}
	return blockchain
}

func NewPoABlockchain(database *db.Db, authNodes []*ecdsa.PublicKey) *Blockchain {
	poa := PoA{
		AuthMiners: authNodes,
	}
	blockchain := &Blockchain{
		Id:         []byte{},
		Difficulty: uint64(0),
		Genesis:    HashBytes([]byte("genesis")), // tmp
		LastBlock:  &Block{},
		blockdb:    database,
		PoA:        poa,
	}
	return blockchain
}

func (bc *Blockchain) GetHeight() uint64 {
	return bc.LastBlock.Height
}

func (bc *Blockchain) GetLastBlock() *Block {
	return bc.LastBlock
}

func (bc *Blockchain) AddBlock(block *Block) error {
	if !bc.VerifyBlock(block) {
		return errors.New("Block could not be verified")
	}
	bc.LastBlock = block.Copy()
	err := bc.blockdb.Put(block.Hash[:], block.Bytes())
	return err
}

func (bc *Blockchain) GetBlock(hash Hash) (*Block, error) {
	v, err := bc.blockdb.Get(hash[:])
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

func (bc *Blockchain) verifyBlockSignature(block *Block) bool {
	// check if the signer is one of the blockchain.AuthMiners
	signerIsMiner := false
	for _, pubK := range bc.PoA.AuthMiners {
		if bytes.Equal(PackPubK(pubK), PackPubK(block.MinerPubK)) {
			signerIsMiner = true
		}
	}
	if !signerIsMiner && len(bc.PoA.AuthMiners) > 0 {
		fmt.Println("signer is not miner")
		return false
	}

	// get the signature
	sig, err := SignatureFromBytes(block.Signature)
	if err != nil {
		fmt.Println("error parsing signature")
		return false
	}

	// check if the signature is by the miner
	return VerifySignature(block.MinerPubK, block.Hash[:], *sig)
}

func (bc *Blockchain) VerifyBlock(block *Block) bool {
	// verify block signature
	// TODO for the moment just covered the case of PoA blockchain
	if !bc.verifyBlockSignature(block) {
		fmt.Println("signature verification error")
		return false
	}

	// verify timestamp

	// verify prev hash
	// check that the block.PrevHash is the blockchain current last block
	if !bytes.Equal(block.PrevHash[:], bc.LastBlock.Hash[:]) {
		fmt.Println("block.PrevHash not equal to last block hash")
		return false
	}

	// verify block height
	// check that the block height is the last block + 1
	if block.Height != bc.LastBlock.Height+1 {
		fmt.Println("block.Height error")
		return false
	}

	// verify block transactions

	return true
}

// func (bc *Blockchain) Mint(toAddr Address, amount uint64) error {
//         fromAddr := Address(HashBytes([]byte("mint")))
//         out :=
//         tx := NewTx(fromAddr, toAddr, []Input, )
// }
