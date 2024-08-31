package authorization

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// NOTE Adapter -----------------------------
type jwtHS256 struct {
	Signature string        `json:"signature"`
	Duration  time.Duration `json:"duration"`
}

// JWT แบบ HS256
func NewJWT_HS256(signature string, duration time.Duration) AppAuthorization {
	return jwtHS256{Signature: signature, Duration: duration}
}

func (c jwtHS256) GenerateToken(payload AppAuthorizationClaim) (tokenString string, err error) {
	// FIX EDIT PAYLOAD HERE ----------------------------
	claim := &authCustomClaims{
		payload.Name,
		payload.Channel,
		jwt.StandardClaims{
			Audience:  payload.Audience,                  // aud Audience (who or what the token intended for)
			ExpiresAt: time.Now().Add(c.Duration).Unix(), // exp Expiration time (seconds since Unix epoch)
			Id:        "",                                // jti JWT ID (unique identifier for this token)
			IssuedAt:  time.Now().Unix(),                 // iat isused at (seconds since Unix epoch)
			Issuer:    payload.Issuer,                    // iss issuer (who created and signed this token)
			NotBefore: 0,                                 // nbf No valid before (seconds since Unix epoch)
			Subject:   payload.UserId,                    // sub Subject (whom the token reference to)
		},
	}
	// FIX EDIT PAYLOAD HERE ----------------------------

	// NOTE Create a new JWT token & Set the claims for the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// NOTE Sign the token with the key
	tokenString, err = token.SignedString([]byte(c.Signature))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (c jwtHS256) ValidateToken(tokenString string, data interface{}) (err error) {
	// NOTE Parse the token string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// NOTE Check the signing method of the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(ErrInvalidSigningMethod)
		}

		// NOTE Return the key for verifying the signature
		return []byte(c.Signature), nil
	})
	if err != nil {
		return err
	}
	// NOTE Check if the token is valid
	if !token.Valid {
		return errors.New(ErrInvalidToken)
	}

	// NOTE Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New(ErrInvalidClaims)
	}

	// NOTE Check expired from the token
	expirationTime, ok := claims["exp"].(float64)
	if !ok {
		return errors.New(ErrInvalidExpirationTime)
	}

	if time.Now().Unix() > int64(expirationTime) {
		return errors.New(ErrTokenHasExpired)
	}

	// NOTE Check issuer from the token
	_, ok = claims["iss"].(string)
	if !ok {
		return errors.New(ErrInvalidIssuerInToken)
	}

	// NOTE Convert the data to the specified type
	dataBytes, err := json.Marshal(claims)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return err
	}

	return nil
}
