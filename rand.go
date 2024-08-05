package main

import (
	"crypto/rand"
	"math/big"
	"unicode"
)

func genRoomCode() Code {
	var codes [6]byte
	for i := 0; i < 6; i++ {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(26))
		codes[i] = uint8(65 + nBig.Int64())
	}
	return Code(codes[:])
}

func validateRoomCode(s Code) bool {
	if len(s) != 6 {
		return false
	}
	for _, c := range s {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}
