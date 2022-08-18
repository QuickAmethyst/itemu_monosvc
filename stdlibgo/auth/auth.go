package auth

import (
	"context"
	"crypto/rsa"
	"github.com/go-redis/redis/v9"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"strings"
	"time"
)

var TimeFunc = time.Now

type Auth interface {
	CreateTokenPair(subject string) (accessToken TokenDetail, refreshToken TokenDetail, err error)
	ParseClaimFromAccessToken(tokenStr string) (claims *jwt.RegisteredClaims, err error)
	RefreshAccessToken(rToken string) (accessToken TokenDetail, refreshToken TokenDetail, err error)
	Authenticate(bearer string) (*jwt.RegisteredClaims, error)
}

type Options struct {
	Redis                redis.UniversalClient
	PublicKeyPath        string
	PrivateKeyPath       string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type auth struct {
	redis                redis.UniversalClient
	publicKey            *rsa.PublicKey
	privateKey           *rsa.PrivateKey
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type TokenDetail struct {
	Raw    string
	Claims *jwt.RegisteredClaims
}

func (a *auth) CreateTokenPair(subject string) (accessToken TokenDetail, refreshToken TokenDetail, err error) {
	if accessToken.Raw, accessToken.Claims, err = a.createAccessToken(subject); err != nil {
		return
	}

	if refreshToken.Raw, refreshToken.Claims, err = a.createRefreshToken(subject); err != nil {
		return
	}

	return
}

func (a *auth) createAccessToken(subject string) (token string, claims *jwt.RegisteredClaims, err error) {
	var (
		id, _     = uuid.NewV4()
		expiresAt = TimeFunc().Add(a.accessTokenDuration)
	)

	claims = &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        id.String(),
		IssuedAt:  jwt.NewNumericDate(TimeFunc()),
		Subject:   subject,
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(a.privateKey)

	return
}

func (a *auth) createRefreshToken(subject string) (token string, claims *jwt.RegisteredClaims, err error) {
	var (
		id, _     = uuid.NewV4()
		expiresAt = TimeFunc().Add(a.refreshTokenDuration)
	)

	claims = &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        id.String(),
		IssuedAt:  jwt.NewNumericDate(TimeFunc()),
		Subject:   subject,
	}

	expiration := calcClaimsExpiredDuration(claims)

	token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(a.privateKey)
	if err != nil {
		return
	}

	redisKey := refreshTokenRedisKey(id, subject)
	if err = a.redis.Set(context.Background(), redisKey, true, expiration).Err(); err != nil {
		return
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(a.privateKey)

	return
}

func (a *auth) ParseClaimFromAccessToken(tokenStr string) (claims *jwt.RegisteredClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidSigningMethod(token.Header["alg"])
		} else if method != jwt.SigningMethodRS256 {
			return nil, ErrInvalidSigningMethod(token.Header["alg"])
		}

		return a.publicKey, nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return claims, ErrTokenExpired
		}

		return claims, ErrParseClaim
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return claims, ErrClaimTypeCasting
	}

	return claims, nil
}

func (a *auth) ParseClaimFromRefreshToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidSigningMethod(token.Header["alg"])
		} else if method != jwt.SigningMethodRS256 {
			return nil, ErrInvalidSigningMethod(token.Header["alg"])
		}

		return a.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claim := token.Claims.(*jwt.RegisteredClaims)

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return claim, ErrTokenExpired
		}

		return claim, ErrParseClaim
	}

	return claim, nil
}

func (a *auth) RefreshAccessToken(rToken string) (accessToken TokenDetail, refreshToken TokenDetail, err error) {
	refreshClaim, err := a.ParseClaimFromRefreshToken(rToken)
	if err != nil {
		return
	}

	refreshTokenID, err := uuid.FromString(refreshClaim.ID)
	if err != nil {
		return
	}

	redisKey := refreshTokenRedisKey(refreshTokenID, refreshClaim.Subject)
	_, err = a.redis.Get(context.Background(), redisKey).Result()
	if err != nil {
		return
	}

	// Delete old refresh token from redis
	_, err = a.redis.Del(context.Background(), refreshTokenRedisKey(refreshTokenID, refreshClaim.Subject)).Result()
	if err != nil {
		return
	}

	// create new access token
	accessToken.Raw, accessToken.Claims, err = a.createAccessToken(refreshClaim.Subject)
	if err != nil {
		return
	}

	// and create new refresh token
	refreshToken.Raw, refreshToken.Claims, err = a.createRefreshToken(refreshClaim.Subject)
	if err != nil {
		return
	}

	refreshID, err := uuid.FromString(refreshClaim.ID)
	if err != nil {
		return
	}

	redisKey = refreshTokenRedisKey(refreshID, refreshClaim.Subject)

	if err = a.redis.Set(context.Background(), redisKey, true, calcClaimsExpiredDuration(refreshClaim)).Err(); err != nil {
		return
	}

	return accessToken, refreshToken, err
}

func (a *auth) Authenticate(bearer string) (*jwt.RegisteredClaims, error) {
	tokenVal := strings.TrimPrefix(bearer, "Bearer ")

	if tokenVal == bearer {
		return nil, ErrInvalidBearerToken
	}

	claim, err := a.ParseClaimFromAccessToken(tokenVal)
	if err != nil {
		return nil, err
	}

	return claim, nil
}

func New(opt *Options) (Auth, error) {
	verifyBytes, err := ioutil.ReadFile(opt.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	decryptBytes, err := ioutil.ReadFile(opt.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, ErrParsePublicKey
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(decryptBytes)
	if err != nil {
		return nil, ErrParsePrivateKey
	}

	return &auth{
		redis:                opt.Redis,
		publicKey:            publicKey,
		privateKey:           privateKey,
		accessTokenDuration:  opt.AccessTokenDuration,
		refreshTokenDuration: opt.RefreshTokenDuration,
	}, nil
}
