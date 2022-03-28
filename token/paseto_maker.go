package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	paseto *paseto.V2
	symKey []byte
}

func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return p.paseto.Encrypt(p.symKey, payload, nil)
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	//panic("implement me")
	payload := &Payload{}
	err := p.paseto.Decrypt(token, p.symKey, payload, nil)
	if err != nil {
		return nil, ErrInValid
	}
	err = payload.Valid()
	if err != nil {
		return nil, ErrExpiredAt
	}
	return payload, nil
}

func NewPasetoMaker(symKey string) (Maker, error) {
	if len(symKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid size must [%d]", chacha20poly1305.KeySize)
	}
	return &PasetoMaker{
		paseto: paseto.NewV2(),
		symKey: []byte(symKey),
	}, nil
}
