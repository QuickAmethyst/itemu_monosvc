package auth

import (
	"errors"
	"fmt"
)

var (
	ErrLoadPublicKey        = errors.New("error load public key")
	ErrLoadPrivateKey       = errors.New("error load private key")
	ErrParsePublicKey       = errors.New("error parse public key")
	ErrParsePrivateKey      = errors.New("error parse private key")
	ErrClaimTypeCasting     = errors.New("error claim type casting")
	ErrTokenExpired         = errors.New("error token expired")
	ErrParseClaim           = errors.New("error parse claims")
	ErrParseInfo            = errors.New("error parse info")
	ErrUnmarshalInfo        = errors.New("error unmarshal info")
	ErrorMarshalInfo        = errors.New("error marshal info")
	ErrEncryptInfo          = errors.New("error encrypt info")
	ErrInvalidBearerToken   = errors.New("invalid bearer token")
	ErrInvalidSigningMethod = func(sign interface{}) error {
		return fmt.Errorf("unexpected signing method: %v", sign)
	}
)
