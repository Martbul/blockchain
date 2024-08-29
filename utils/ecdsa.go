package utils

import (
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int // public key X coordinate
	S *big.Int // public key Y coordinate
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}