package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func encrypt(data any, key []byte) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	padding := aes.BlockSize - len(jsonData)%aes.BlockSize
	padtext := append(jsonData, bytes.Repeat([]byte{byte(padding)}, padding)...)

	ciphertext := make([]byte, len(padtext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padtext)

	return base64.StdEncoding.EncodeToString(append(iv, ciphertext...)), nil
}

func decrypt(encrypted string, key []byte) ([]Storage, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	padding := plaintext[len(plaintext)-1]
	plaintext = plaintext[:len(plaintext)-int(padding)]

	var result []Storage
	err = json.Unmarshal(plaintext, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 33->125
func generatePass(seed *Seed) string {
	pass := ""
	for i := 0; i < seed.Length; i++ {

		newRand, err := rand.Int(rand.Reader, big.NewInt(time.Now().UnixNano()))
		if err != nil {
			log.Fatal(err)
		}
		offset := newRand.Int64() % int64(seed.Key)

		newRand, err = rand.Int(rand.Reader, big.NewInt(time.Now().UnixNano()))
		if err != nil {
			log.Fatal(err)
		}

		ind := 33 + ((newRand.Int64()%93)+offset)%93
		pass += string(rune(ind))
	}
	return pass
}

func hashKey(key []byte) ([]byte, error) {
	hashedKey, err := bcrypt.GenerateFromPassword(key, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedKey, nil
}

func compareHashKey(key, hashedKey []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedKey, key)
	return err == nil
}

func generateKey() string {
	seed := Seed{1, 16, ""}
	return generatePass(&seed)
}
