package random

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"sync"
)

var letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

type RandomStringGenerator struct {
	sync.RWMutex
	rnd   *rand.Rand
	runes string
}

func NewRandomStringGenerator() (RandomStringGenerator, error) {
	nBig, err := crand.Int(crand.Reader, big.NewInt(9223372036854775807))
	if err != nil {
		return RandomStringGenerator{}, err
	}
	res := RandomStringGenerator{
		rnd:   rand.New(rand.NewSource(nBig.Int64())),
		runes: letterRunes,
	}
	return res, nil
}

func (r *RandomStringGenerator) UpdateSeed() error {
	nBig, err := crand.Int(crand.Reader, big.NewInt(9223372036854775807))
	if err != nil {
		return err
	}
	r.Lock()
	defer r.Unlock()

	r.rnd = rand.New(rand.NewSource(nBig.Int64()))
	return nil
}

func (r *RandomStringGenerator) RandStringRunes(n int) string {
	r.RLock()
	defer r.RUnlock()

	result := ""

	for i := 0; i < n; i++ {
		result += string(r.runes[r.rnd.Intn(len(letterRunes))])
	}
	return result
}
