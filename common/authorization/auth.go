package authorization

import "github.com/dgrijalva/jwt-go"

// NOTE ERROR ------------------------------
var (
	ErrInvalidSigningMethod  = "invalid signing method"
	ErrInvalidToken          = "invalid token"
	ErrInvalidClaims         = "invalid claims"
	ErrInvalidExpirationTime = "invalid expiration time"
	ErrInvalidIssuerInToken  = "invalid issuer in token"
	ErrTokenHasExpired       = "token has expired"
)

// NOTE Port -------------------------------
type AppAuthorization interface {
	// สำหรับ Cenerate JWT Tokan
	GenerateToken(payload AppAuthorizationClaim) (token string, err error)

	// สำหรับ Validate JWT Tokan
	ValidateToken(tokenString string, paserTo interface{}) (err error)
}

type authCustomClaims struct {
	Name    string `json:"name,omitempty"`
	Channel string `json:"channel,omitempty"`
	jwt.StandardClaims
}

type AppAuthorizationClaim struct {
	UserId   string `json:"sub,omitempty"`
	Name     string `json:"name,omitempty"`
	Audience string `json:"aud,omitempty"`
	Issuer   string `json:"issuer,omitempty"`
	Channel  string `json:"channel,omitempty"`
}
