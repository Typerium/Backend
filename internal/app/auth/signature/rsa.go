package signature

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"sync"

	"github.com/pkg/errors"
)

func NewRSACreator(bits int) KeyCreator {
	return &rsaCreator{
		bits: bits,
		keyPool: &sync.Pool{
			New: func() interface{} {
				return new(rsaKey)
			},
		},
	}
}

type rsaCreator struct {
	bits    int
	keyPool *sync.Pool
}

const rsaBlockType = "RSA PUBLIC KEY"

type rsaKey struct {
	secret   *rsa.PrivateKey
	isPublic bool
}

func (r *rsaKey) Bytes() []byte {
	if r == nil || r.secret == nil {
		return nil
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  rsaBlockType,
		Bytes: x509.MarshalPKCS1PublicKey(&r.secret.PublicKey),
	})
}

func (r *rsaKey) Key() interface{} {
	if r.isPublic {
		return &r.secret.PublicKey
	}

	return r.secret
}

func (r *rsaCreator) Acquire() (Key, error) {
	key, err := rsa.GenerateKey(rand.Reader, r.bits)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := r.keyPool.Get().(*rsaKey)
	result.secret = key

	return result, nil
}

func (r *rsaCreator) Release(key Key) {
	rsaKey, ok := key.(*rsaKey)
	if !ok {
		return
	}

	rsaKey.secret = nil
	rsaKey.isPublic = false
	r.keyPool.Put(rsaKey)
}

func (r *rsaCreator) Create(key []byte) (Key, error) {
	block, _ := pem.Decode(key)

	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := r.keyPool.Get().(*rsaKey)
	result.secret = &rsa.PrivateKey{
		PublicKey: *pubKey,
	}
	result.isPublic = true

	return result, nil
}
