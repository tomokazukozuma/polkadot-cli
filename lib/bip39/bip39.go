package bip39

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func MnemonicToEntropy(mnemonic string) ([]byte, error) {
	wordList := strings.Split(mnemonic, " ")
	var indices []string
	for _, word := range wordList {
		i := getWordIndex(englishWordList, word)
		indices = append(indices, fmt.Sprintf("%011s", strconv.FormatInt(i, 2)))
	}
	bits := strings.Join(indices, "")

	var dividerIndex = math.Floor(float64(len(bits))/33) * 32
	var entropyBits = bits[:int(dividerIndex)]
	// TODO confirm checksum
	var checksumBits = bits[int(dividerIndex):]
	entropy := bitsToBytes(entropyBits)
	checksum := generateChecksum(entropy)
	if checksum != checksumBits {
		log.Fatalf("mismatch checksum. checksum: %b, checksumBits: %b", checksum, checksumBits)
	}
	log.Printf("entropy: %x", entropy)
	return entropy, nil
}

func getWordIndex(wordList []string, word string) int64 {
	for i, w := range wordList {
		if w == word {
			return int64(i)
		}
	}
	return -1
}

func bitsToBytes(bitString string) []byte {
	lenB := len(bitString)/8 + 1
	bs := make([]byte, lenB)

	count, i := 0, 0
	var now byte
	for _, v := range bitString {
		if count == 8 {
			bs[i] = now
			i++
			now, count = 0, 0
		}
		now = now<<1 + byte(v-'0')
		count++
	}
	if count != 0 {
		bs[i] = now << (8 - byte(count))
		i++
	}

	bs = bs[:i:i]
	return bs
}

func generateChecksum(entropy []byte) string {
	checksumSize := len(entropy) * 8 / 32
	hashedEntropy := sha256.Sum256(entropy)
	return fmt.Sprintf("%08b", hashedEntropy[0])[:checksumSize]
}
