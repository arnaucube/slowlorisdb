package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKey(t *testing.T) {
	_, err := NewKey()
	assert.Nil(t, err)
}
func TestAddress(t *testing.T) {
	privK, err := NewKey()
	assert.Nil(t, err)

	addr := AddressFromPrivK(privK)
	fmt.Println(addr.String())
}

func TestSignAndVerify(t *testing.T) {
	privK, err := NewKey()
	assert.Nil(t, err)

	// Sign
	m := []byte("test")
	sig, err := Sign(privK, m)
	assert.Nil(t, err)

	// Verify
	verified := VerifySignature(&privK.PublicKey, m, sig)
	assert.True(t, verified)
}
