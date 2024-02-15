package shortener

import (
	"crypto/sha256"
	"fmt"
	"github.com/itchyny/base58-go"
	"math/big"
	"os"
)

const UrlLength = 8

// convert str to hash value (256-bit binary value)
func sha256f(input string) []byte {
	algo := sha256.New()
	algo.Write([]byte(input))
	return algo.Sum(nil)
}

// convert hash value to text
func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

// GenerateShortLink userId here is increasing integer to differentiate same original URL
// also helps avoid collision
func GenerateShortLink(originalUrl string, userId string) string {
	hashBytes := sha256f(originalUrl + userId)
	// convert the hash into a big integer, giving you a unique numerical representation of the hash
	genNum := new(big.Int).SetBytes(hashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", genNum)))
	return finalString[:UrlLength]
}
