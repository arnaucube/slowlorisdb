package core

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashBytes(t *testing.T) {
	m := []byte("test")
	h := HashBytes(m)
	assert.Equal(t, hex.EncodeToString(h[:]), "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")
	assert.True(t, !h.IsZero())

	z := Hash{}
	assert.True(t, z.IsZero())
}

type testData struct {
	Data  []byte
	Nonce uint64
}

func (d *testData) Bytes() []byte {
	b, _ := json.Marshal(d)
	return b
}

func (d *testData) GetNonce() uint64 {
	return d.Nonce
}
func (d *testData) IncrementNonce() {
	d.Nonce++
}

func TestPoW(t *testing.T) {
	difficulty := uint64(2)
	data := &testData{
		Data:  []byte("test"),
		Nonce: 0,
	}
	nonce, err := CalculatePoW(data, difficulty)
	assert.Nil(t, err)
	data.Nonce = nonce

	h := HashBytes(data.Bytes())

	assert.Equal(t, hex.EncodeToString(h[:]), "0000020881c02f5171b978e74bb710242e95cc67b36416e382118a7ab2a69321")

	// CheckPoW
	assert.True(t, CheckPoW(h, difficulty))
}
