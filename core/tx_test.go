package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTx(t *testing.T) {
	addr0 := Address(HashBytes([]byte("addr0")))
	addr1 := Address(HashBytes([]byte("addr1")))

	tx := NewTx(addr0, addr1, []Input{}, []Output{})

	assert.Equal(t, tx.From, addr0)
	assert.Equal(t, tx.To, addr1)

	assert.True(t, CheckTx(tx))
}
