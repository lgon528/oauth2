package generates

import (
	cryptoRand "crypto/rand"
	"math"
	"math/big"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func Secret(n int) string {
	b := make([]byte, n)

	randInt, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(math.MaxInt64))
	i, cache, remain := n-1, randInt.Int64(), letterIdxMax
	for i >= 0 {
		if remain == 0 {
			randInt, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(math.MaxInt64))
			cache, remain = randInt.Int64(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
