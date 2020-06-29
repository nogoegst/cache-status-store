package cachestatusstore

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

type Cache interface {
	Touch(key string) (hit bool, err error)
}

type Storage struct {
	cache          Cache
	PrintDebugBits bool
}

func NewStorage(cache Cache) *Storage {
	return &Storage{
		cache: cache,
	}
}

func cacheKey(slug []byte, offset int64) string {
	h := sha256.New()
	h.Write(slug)
	binary.Write(h, binary.BigEndian, offset)
	hash := h.Sum(nil)[:16]
	return base64.RawURLEncoding.EncodeToString(hash)
}

func (s *Storage) GetBit(slug []byte, offset int64) (bool, error) {
	bit, err := s.cache.Touch(cacheKey(slug, offset))
	if s.PrintDebugBits {
		if err != nil {
			fmt.Print("x")
		} else {
			if bit {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
		}
	}
	return bit, err
}

func (s *Storage) SetBit(slug []byte, offset int64, value bool) error {
	var err error
	if value {
		hit, e := s.cache.Touch(cacheKey(slug, offset))
		if hit {
			return fmt.Errorf("cache was already toggled")
		}
		err = e
	}
	if s.PrintDebugBits {
		if err != nil {
			fmt.Print("x")
		}
		if value {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	}
	return err
}

func (s *Storage) GetByte(key []byte, offset int64) (byte, error) {
	var x byte
	for i := 0; i < 8; i++ {
		bit, err := s.GetBit(key, offset*8+int64(i))
		if err != nil {
			return 0x0, fmt.Errorf("get bit #%v: %w", i, err)
		}
		if bit {
			x |= 1 << i
		}
	}
	return x, nil
}

func (s *Storage) SetByte(slug []byte, offset int64, x byte) error {
	for i := 0; i < 8; i++ {
		err := s.SetBit(slug, offset*8+int64(i), x&(1<<i) != 0)
		if err != nil {
			return fmt.Errorf("set bit #%v: %w", i, err)
		}
	}
	return nil
}

func (s *Storage) SetBytes(slug, bytes []byte) error {
	for i, b := range bytes {
		if err := s.SetByte(slug, int64(i), byte(b)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetBytes(slug []byte, length int64) ([]byte, error) {
	var ret []byte
	for i := int64(0); i < length; i++ {
		b, err := s.GetByte(slug, i)
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return ret, nil
}
