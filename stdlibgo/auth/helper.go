package auth

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func calcClaimsExpiredDuration(claims *jwt.RegisteredClaims) time.Duration {
	var (
		expireTime = claims.ExpiresAt
		issueTime  = claims.IssuedAt
		expiration = expireTime.Sub(issueTime.Time)
	)

	return expiration
}

func refreshTokenRedisKey(id uuid.UUID, userID string) string {
	return fmt.Sprintf("%s:%s:%s", "refresh_token", userID, id.String())
}
