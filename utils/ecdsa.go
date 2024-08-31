package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)


type Signature struct {
	R *big.Int // public key X coordinate
	S *big.Int // public key Y coordinate
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

func String2BigIntTuple(s string) (big.Int, big.Int) {
	bx,_ := hex.DecodeString(s[:64])
	by,_ := hex.DecodeString(s[64:])

	var bix big.Int
	var biy big.Int

	_ = bix.SetBytes(bx)
	_ = biy.SetBytes(by)

	return bix, biy
}

func PublicKeyFromString(s string) *ecdsa.PublicKey {
	x, y := String2BigIntTuple(s)

	return &ecdsa.PublicKey{elliptic.P256(), &x, &y}
}

func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey{
	b, _ := hex.DecodeString(s[:]) //hex.DecodeString() is decoding s from hexadecimal format into a byte slice. s[:] creates a slice of thr string 's', in this case the entire string
	var bi big.Int
	_ = bi.SetBytes(b) //converting byte slice 'b' into integer that is assigned to big.Int value 'bi'
	return &ecdsa.PrivateKey{*publicKey, &bi}
}

func SignatureFromString(s string) *Signature{ 
	x,y := String2BigIntTuple(s)
	return &Signature{&x, &y}
}