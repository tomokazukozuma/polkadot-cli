package bip39

import (
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
	//var checksumBits = bits[int(dividerIndex):]
	entropy := bitsToBytes(entropyBits)
	log.Printf("entropy: %x", entropy)
	return entropy, nil
}

func getWordIndex(wordList []string, word string) int64 {
	for i, w := range wordList {
		if w == word {
			return i
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
