package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Test_BuildToken(t *testing.T) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		UserID:  "user id",
		Account: "acc name",
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(rsaPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	publicKeyBytes, _ := ExportRsaPublicKeyAsPemStr(&rsaPrivateKey.PublicKey)

	claims2 := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims2, func(token *jwt.Token) (interface{}, error) {
		return ParseRsaPublicKeyFromPemStr(publicKeyBytes)
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	if !tkn.Valid {
		t.Fatal("not valid")
		return
	}
	//js, _ := json.Marshal(claims2)
	//t.Log(string(js))
}
