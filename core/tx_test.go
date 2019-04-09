package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTx(t *testing.T) {
	privK0, err := NewKey()
	assert.Nil(t, err)
	pubK0 := privK0.PublicKey
	privK1, err := NewKey()
	assert.Nil(t, err)
	pubK1 := privK1.PublicKey

	tx := NewTx(&pubK0, &pubK1, []Input{}, []Output{})

	assert.Equal(t, tx.From, &pubK0)
	assert.Equal(t, tx.To, &pubK1)

	assert.True(t, CheckTx(tx))
}
