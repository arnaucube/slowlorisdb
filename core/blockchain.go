package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/arnaucube/slowlorisdb/db"
	log "github.com/sirupsen/logrus"
	lvldberrors "github.com/syndtr/goleveldb/leveldb/errors"
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
	addressdb  *db.Db
	walletsdb  *db.Db
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
	blockDb := database.WithPrefix([]byte("blockDb"))
	txDb := database.WithPrefix([]byte("txDb"))
	addressDb := database.WithPrefix([]byte("addressDb"))

	poa := PoA{
		AuthMiners: authNodes,
	}
	blockchain := &Blockchain{
		Id:         []byte{},
		Difficulty: uint64(0),
		Genesis:    HashBytes([]byte("genesis")), // tmp
		LastBlock:  &Block{},
		blockdb:    blockDb,
		txdb:       txDb,
		addressdb:  addressDb,
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
	if err != nil {
		return err
	}

	// add each tx to txDb & update addressDb balances
	for _, tx := range block.Txs {
		err = bc.txdb.Put(tx.TxId[:], tx.Bytes())
		if err != nil {
			return err
		}
		bc.UpdateWalletsWithNewTx(&tx)
	}
	return nil
}

func (bc *Blockchain) UpdateWalletsWithNewTx(tx *Tx) error {
	for _, in := range tx.Inputs {
		balanceBytes, err := bc.addressdb.Get(PackPubK(tx.From))
		if err != nil {
			return err
		}
		balance := Uint64FromBytes(balanceBytes)
		balance = balance - in.Value
		err = bc.addressdb.Put(PackPubK(tx.From), Uint64ToBytes(balance))
		if err != nil {
			return err
		}
		log.Info("sent-->: balance of " + hex.EncodeToString(PackPubK(tx.From)[:10]) + ": " + strconv.Itoa(int(balance)))
	}
	for _, out := range tx.Outputs {
		balanceBytes, err := bc.addressdb.Get(PackPubK(tx.To))
		if err != nil && err != lvldberrors.ErrNotFound {
			return err
		}
		if err == lvldberrors.ErrNotFound {
			balanceBytes = Uint64ToBytes(uint64(0))
		}
		balance := Uint64FromBytes(balanceBytes)
		balance = balance + out.Value
		err = bc.addressdb.Put(PackPubK(tx.To), Uint64ToBytes(balance))
		if err != nil {
			return err
		}
		log.Info("--> received: balance of " + hex.EncodeToString(PackPubK(tx.To)[:10]) + ": " + strconv.Itoa(int(balance)))
	}
	return nil

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

func (bc *Blockchain) GetBalance(pubK *ecdsa.PublicKey) (uint64, error) {
	balanceBytes, err := bc.addressdb.Get(PackPubK(pubK))
	if err != nil {
		return uint64(0), err
	}
	balance := Uint64FromBytes(balanceBytes)
	return balance, nil

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
		log.Error("signer is not miner")
		return false
	}

	// get the signature
	sig, err := SignatureFromBytes(block.Signature)
	if err != nil {
		log.Error("error parsing signature")
		return false
	}

	// check if the signature is by the miner
	return VerifySignature(block.MinerPubK, block.Hash[:], *sig)
}

func (bc *Blockchain) VerifyBlock(block *Block) bool {
	// verify block signature
	// for the moment just covered the case of PoA blockchain
	if !bc.verifyBlockSignature(block) {
		log.Error("signature verification error")
		return false
	}

	// verify timestamp
	if block.Timestamp.Unix() < bc.LastBlock.Timestamp.Unix() {
		return false
	}

	// verify prev hash
	// check that the block.PrevHash is the blockchain current last block
	if !bytes.Equal(block.PrevHash[:], bc.LastBlock.Hash[:]) {
		fmt.Println(block.PrevHash.String())
		fmt.Println(bc.LastBlock.Hash.String())
		log.Error("block.PrevHash not equal to last block hash")
		return false
	}

	// verify block height
	// check that the block height is the last block + 1
	if block.Height != bc.LastBlock.Height+1 {
		log.Error("block.Height error")
		return false
	}

	// verify block transactions (not if the block is the genesis block)
	if !bytes.Equal(block.Txs[0].TxId[:], GenesisHashTxInput[:]) {
		for _, tx := range block.Txs {
			txVerified := CheckTx(&tx)
			if !txVerified {
				log.Error("tx could not be verified")
				return false
			}
		}
	}

	// TODO in --> out0
	//          -> out1
	//          -> ...

	return true
}

// func (bc *Blockchain) Mint(toAddr Address, amount uint64) error {
//         fromAddr := Address(HashBytes([]byte("mint")))
//         out :=
//         tx := NewTx(fromAddr, toAddr, []Input, )
// }
