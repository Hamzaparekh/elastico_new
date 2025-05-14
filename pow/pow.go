package pow

import (
	"crypto/sha256"
	"encoding/binary"
)

func GetNonce(data []byte, difficulty int) []byte {
	var nonce int64
	for {
		combined := append(data, intToBytes(nonce)...)
		hash := sha256.Sum256(combined)
		if isValid(hash[:], difficulty) {
			return intToBytes(nonce)
		}
		nonce++
	}
}

func Fulfill(data []byte, nonce []byte, difficulty int) bool {
	combined := append(data, nonce...)
	hash := sha256.Sum256(combined)
	return isValid(hash[:], difficulty)
}

func Hash(data []byte, nonce int64) []byte {
	combined := append(data, intToBytes(nonce)...)
	hash := sha256.Sum256(combined)
	return hash[:]
}

func isValid(hash []byte, difficulty int) bool {
	// naive check: count leading 0s
	for i := 0; i < difficulty; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}

func intToBytes(n int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(n))
	return buf
}
