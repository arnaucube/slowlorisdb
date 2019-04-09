package core

import "crypto/ecdsa"

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
	From       *ecdsa.PublicKey
	To         *ecdsa.PublicKey
	InputCount uint64
	Inputs     []Input
	Outputs    []Output
	Signature  []byte
}

func NewTx(from, to *ecdsa.PublicKey, in []Input, out []Output) *Tx {
	tx := &Tx{
		From:       from,
		To:         to,
		InputCount: uint64(len(in)),
		Inputs:     in,
		Outputs:    out,
	}
	return tx
}

// CheckTx checks if the transaction is consistent
func CheckTx(tx *Tx) bool {
	// TODO check that inputs and outputs are not empty

	// check that inputs == outputs
	totalIn := 0
	for _, in := range tx.Inputs {
		totalIn = totalIn + int(in.Value)
	}
	totalOut := 0
	for _, out := range tx.Outputs {
		totalOut = totalOut + int(out.Value)
	}
	if totalIn < totalOut {
		return false
	}

	// TODO check signature
	return true
}
