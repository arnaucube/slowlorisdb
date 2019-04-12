package core

import (
	"crypto/ecdsa"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

var GenesisHashTxInput = HashBytes([]byte("genesis"))

type Input struct {
	TxId  Hash
	Vout  int // index of the output from the TxId
	Value uint64
}

type Output struct {
	Value uint64
}

// Tx holds the data structure of a transaction
type Tx struct {
	TxId       Hash
	From       *ecdsa.PublicKey
	To         *ecdsa.PublicKey
	InputCount uint64
	Inputs     []Input
	Outputs    []Output
	Signature  []byte
}

func (tx *Tx) Bytes() []byte {
	// TODO add parser, to use minimum amount of bytes
	b, _ := json.Marshal(tx)
	return b
}

func (tx *Tx) CalculateTxId() {
	h := HashBytes(tx.Bytes())
	tx.TxId = h
}

func NewTx(from, to *ecdsa.PublicKey, in []Input, out []Output) *Tx {
	tx := &Tx{
		From:       from,
		To:         to,
		InputCount: uint64(len(in)),
		Inputs:     in,
		Outputs:    out,
		Signature:  []byte{},
	}
	tx.CalculateTxId()
	return tx
}

// CheckTx checks if the transaction is consistent
func CheckTx(tx *Tx) bool {
	// check that inputs == outputs
	totalIn := 0
	for _, in := range tx.Inputs {
		// check that inputs are not empty, to avoid spam tx
		if in.Value == uint64(0) {
			return false
		}
		totalIn = totalIn + int(in.Value)
	}
	totalOut := 0
	for _, out := range tx.Outputs {
		totalOut = totalOut + int(out.Value)
	}
	if totalIn != totalOut {
		log.Info("totalIn != totalOut")
		return false
	}

	// TODO check signature
	return true
}
