package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"photofinish/internal/domain/user"
	"strings"
	"time"
)

func CreateTokenWithDuration(username string, role user.Role, duration time.Duration) (string, error) {
	var err error
	claims := jwt.MapClaims{}
	claims["role"] = role.Values()
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(duration).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := getSecret()
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func CreateToken(username string, role user.Role) (string, error) {
	return CreateTokenWithDuration(username, role, time.Second*60)
}

func CheckToken(tokenString string) (*jwt.Token, bool) {
	token, err := parseToken(tokenString)
	if err != nil {
		return nil, false
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, false
	}

	return token, true
}

func parseToken(tokenString string) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getSecret()), nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func ParseToken(header string) (string, bool) {
	if header == "" {
		return "Authorization header not set", false
	}
	bearerToken := strings.Split(header, " ")
	if len(bearerToken) != 2 {
		return "Cannot read token", false
	}
	if bearerToken[0] != "Bearer" {
		return "Error in authorization token. it needs to be in form of 'Bearer <token>'", false
	}

	return bearerToken[1], true
}

func getSecret() string {
	secret := os.Getenv("ACCESS_SECRET")

	if secret == "" {
		secret = "sdmalncnjsdsmf"
	}

	return secret
}
