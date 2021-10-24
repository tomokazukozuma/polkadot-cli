package bip39

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/unicode/norm"
)

func GenerateMnemonic() ([]string, []byte, error) {
	entropy := make([]byte, 16)
	if _, err := rand.Read(entropy); err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed generate entropy: %s", err.Error()))
	}
	if len(entropy) < 16 || len(entropy) > 32 || len(entropy)%4 != 0 {
		return nil, nil, errors.New(fmt.Sprintf("Invalid entropy: %x", entropy))
	}
	entropyBits := bytesToBits(entropy)
	checksumBits := generateChecksum(entropy)
	bits := entropyBits + checksumBits
	chunks := regexp.MustCompile(".{1,11}").FindAllString(bits, -1)
	var words []string
	for _, chunk := range chunks {
		index, err := strconv.ParseInt(chunk, 2, 64)
		if err != nil {
			return nil, nil, err
		}
		words = append(words, englishWordList[index])
	}
	return words, entropy, nil
}

func MnemonicToSeed(mnemonic, passphrase string) ([]byte, error) {
	entropy, err := MnemonicToEntropy(mnemonic)
	if err != nil {
		return nil, err
	}
	return pbkdf2.Key(norm.NFKD.Bytes(entropy), norm.NFKD.Bytes([]byte("mnemonic"+passphrase)), 2048, 64, sha512.New), nil
}
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

func bytesToBits(bytes []byte) string {
	var bits string
	for _, byte := range bytes {
		bits += fmt.Sprintf("%08b", byte)
	}

	return bits
}

func generateChecksum(entropy []byte) string {
	checksumSize := len(entropy) * 8 / 32
	hashedEntropy := sha256.Sum256(entropy)
	return fmt.Sprintf("%08b", hashedEntropy[0])[:checksumSize]
}
