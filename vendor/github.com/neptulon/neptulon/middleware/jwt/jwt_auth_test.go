package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	pass = "pass"
	now  = time.Now().Unix()
)

type testMsg struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func TestJWTLib(t *testing.T) {
	tokenStr := genToken(t)

	jt, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(pass), nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if !jt.Valid {
		t.Fatal("somehow generated invalid token with the library itself!")
	}
}

func TestMiddleware(t *testing.T) {

	//
	// mid := HMAC(pass)
}

func genToken(t *testing.T) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["userid"] = 1
	token.Claims["created"] = now
	tokenString, err := token.SignedString([]byte(pass))
	if err != nil {
		t.Fatal(err)
	}

	return tokenString
}
