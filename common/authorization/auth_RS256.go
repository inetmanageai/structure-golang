package authorization

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// NOTE Adapter -----------------------------
type jwtRS256 struct {
	PublicKey  string        `json:"public_key"`
	PrivateKey string        `json:"private_key"`
	Duration   time.Duration `json:"duration"`
}

// JWT แบบ RS256
func NewJWT_RS256(privateKey string, publicKey string, duration time.Duration) AppAuthorization {
	return jwtRS256{PrivateKey: privateKey, PublicKey: publicKey, Duration: duration}
}

func (c jwtRS256) GenerateToken(payload AppAuthorizationClaim) (tokenString string, err error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(c.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

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
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	// NOTE Sign the token with the key
	tokenString, err = token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c jwtRS256) ValidateToken(tokenString string, data interface{}) (err error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(c.PublicKey)
	if err != nil {
		return fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return fmt.Errorf("create: parse key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// NOTE Check the signing method of the token
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New(ErrInvalidSigningMethod)
		}

		// NOTE Return the key for verifying the signature
		return key, nil
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

	fmt.Println(claims)
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return err
	}

	return nil
}
