package encoder

import (
	"crypto/rand"
	"math"
	"math/big"
)

const minGeneratedNumber int64 = 1_000_000_000
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// TODO: clean up + add comments
func Base58() string {
	id := generateID()
	// Number of digits in base58 representation of id,
	// 64 (64-bits) / log2(58) - the number to represent 1 digit in base n.
	buffer := make([]byte, 0, 11)

	elements := 0
	for id > 58 {
		mod := id % 58
		buffer = append(buffer, base58Alphabet[mod])
		id = id / 58
		elements += 1
	}
	// append last element
	buffer = append(buffer, base58Alphabet[id])

	// reverse buffer
	for i := 0; i < elements/2; i++ {
		buffer[i], buffer[elements-i-1] = buffer[elements-i-1], buffer[i]
	}

	return string(buffer[:elements])
}

func generateID() int64 {
	maxPossibleValue := math.MaxInt64 - (math.MaxInt64 % minGeneratedNumber)
	randomInt, err := rand.Int(rand.Reader, big.NewInt(maxPossibleValue))
	if err != nil {
		panic(err)
	}

	id := minGeneratedNumber + randomInt.Int64()
	return id
}
