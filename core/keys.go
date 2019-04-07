package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
)

// Address is the type data for addresses
type Address [32]byte

func (addr Address) String() string {
	return hex.EncodeToString(addr[:])
}

func NewKey() (*ecdsa.PrivateKey, error) {
	curve := elliptic.P256()

	privatekey := new(ecdsa.PrivateKey)
	privatekey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return privatekey, err
}

func PackPubK(pubK *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(pubK.Curve, pubK.X, pubK.Y)
}
func UnpackPubK(b []byte) *ecdsa.PublicKey {
	x, y := elliptic.Unmarshal(elliptic.P256(), b)
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
}

func AddressFromPrivK(privK *ecdsa.PrivateKey) Address {
	h := HashBytes(PackPubK(&privK.PublicKey))
	return Address(h)
}

func (sig *Signature) Bytes() []byte {
	b := sig.R.Bytes()
	b = append(b, sig.S.Bytes()...)
	return b
}

func SignatureFromBytes(b []byte) (*Signature, error) {
	if len(b) != 64 {
		return nil, errors.New("Invalid signature")
	}
	rBytes := b[:32]
	sBytes := b[32:]

	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	sig := &Signature{
		R: r,
		S: s,
	}
	return sig, nil
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func Sign(privK *ecdsa.PrivateKey, m []byte) (*Signature, error) {
	r := big.NewInt(0)
	s := big.NewInt(0)

	hashMsg := HashBytes(m)

	r, s, err := ecdsa.Sign(rand.Reader, privK, hashMsg[:])
	if err != nil {
		return nil, err
	}
	sig := &Signature{
		R: r,
		S: s,
	}

	return sig, nil

}

func VerifySignature(pubK *ecdsa.PublicKey, m []byte, sig Signature) bool {
	hashMsg := HashBytes(m)
	verified := ecdsa.Verify(pubK, hashMsg[:], sig.R, sig.S)
	return verified
}
