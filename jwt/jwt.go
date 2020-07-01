package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID     string            `json:"userId"`
	Account    string            `json:"account"`
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	Nickname   string            `json:"nickname"`
	GivenName  string            `json:"given_name"`
	FamilyName string            `json:"family_name"`
	Picture    string            `json:"picture"`
	Context    string            `json:"context,omitempty"`
	Groups     map[string]string `json:"groups,omitempty"`
	jwtgo.StandardClaims
}

func BuildSignedToken(privateKey *rsa.PrivateKey, keyID string, claims Claims) (string, error) {
	// Declare the token with the algorithm used for signing, and the claims
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, claims)
	token.Header["kid"] = keyID

	// Create the JWT string
	return token.SignedString(privateKey)
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	return string(pubkeyPem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}
