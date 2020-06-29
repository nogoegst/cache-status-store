package cachestatusstore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
)

type EncryptedStorage struct {
	storage *Storage
}

func NewEncryptedStorage(storage *Storage) *EncryptedStorage {
	return &EncryptedStorage{
		storage: storage,
	}
}

func derive(password []byte) (slug, iv, key []byte) {
	root := sha256.Sum256(password)
	slugArray := sha256.Sum256(append(root[:], 0x1))
	ivArray := sha256.Sum256(append(root[:], 0x2))
	keyArray := sha256.Sum256(append(root[:], 0x3))
	return slugArray[:], ivArray[:aes.BlockSize], keyArray[:]
}

func newCTR(iv, key []byte) cipher.Stream {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Errorf("create cipher: %w", err))
	}
	return cipher.NewCTR(block, iv)
}

func (es *EncryptedStorage) SetBytes(password, bytes []byte) error {
	slug, iv, key := derive(password)
	stream := newCTR(iv, key)
	ciphertext := make([]byte, len(bytes))
	stream.XORKeyStream(ciphertext, bytes)
	return es.storage.SetBytes(slug, ciphertext)
}

func (es *EncryptedStorage) GetBytes(password []byte, length int64) ([]byte, error) {
	slug, iv, key := derive(password)
	stream := newCTR(iv, key)
	ciphertext, err := es.storage.GetBytes(slug, length)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
