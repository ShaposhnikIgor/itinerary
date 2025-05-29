package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

const (
	blocksize = 136
)

var (
	rc = []uint64{
		0x0000000000000001, 0x0000000000008082, 0x800000000000808a, 0x8000000080008000,
		0x000000000000808b, 0x0000000080000001, 0x8000000080008081, 0x8000000000008009,
		0x000000000000008a, 0x0000000000000088, 0x0000000080008009, 0x000000008000000a,
		0x000000008000808b, 0x800000000000008b, 0x8000000000008089, 0x8000000000008003,
		0x8000000000008002, 0x8000000000000080, 0x000000000000800a, 0x800000008000000a,
		0x8000000080008081, 0x8000000000008080, 0x0000000080000001, 0x8000000080008008,
	}
)

func padMessage(message []byte) []byte {
	l := len(message)
	message = append(message, 0x06)
	for len(message)%blocksize != blocksize-16 {
		message = append(message, 0x00)
	}
	message = append(message, byte(l>>8), byte(l))
	return message
}

func keccakF(state *[25]uint64) {
	var C [5]uint64
	var D [5]uint64
	for x := 0; x < 5; x++ {
		C[x] = state[x] ^ state[x+5] ^ state[x+10] ^ state[x+15] ^ state[x+20]
	}
	for x := 0; x < 5; x++ {
		D[x] = C[(x+4)%5] ^ ((C[(x+1)%5])<<1 | (C[(x+1)%5])>>(64-1))
		for y := 0; y < 25; y += 5 {
			state[y+x] ^= D[x]
		}
	}
}

func hashSHA3256(message []byte) []byte {

	state := [25]uint64{}

	message = padMessage(message)

	for len(message) > 0 {
		var block []byte
		if len(message) >= blocksize {
			block = message[:blocksize]
			message = message[blocksize:]
		} else {
			block = make([]byte, blocksize)
			copy(block, message)
			message = nil
		}

		// Absorb
		for i := 0; i < len(block); i += 8 {
			for j := 0; j < 8; j++ {
				if i+j < len(block) {
					state[i/8] ^= uint64(block[i+j]) << (8 * uint(j))
				}
			}
		}

		keccakF(&state)
	}

	result := make([]byte, 32)
	for i := range result {
		result[i] = byte((state[i/8] >> (8 * uint(i) % 64)) & 0xff)
	}

	return result
}

func startGeneretion() {
	signature := createSignature("IGOR GENNADIEVICH SHAPOSHNIK")

	signatureHex := hex.EncodeToString(signature)

	file, err := os.Create("signature.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(signatureHex)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func createSignature(data string) []byte {
	message := []byte(data)

	hashedData := hashSHA3256(message)

	return hashedData
}
