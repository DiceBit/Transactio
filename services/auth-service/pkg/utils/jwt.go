package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Claims struct {
	Email string
	Roles []string
	jwt.RegisteredClaims
}

/*var jwtPrivateKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
var jwtPublicKey = jwtPrivateKey.Public()*/

var jwtPrivateKey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	jwtPrivateKey = []byte(os.Getenv("JWT_SECRET"))
}

func GenerateJWT(email string, roles []string) (string, error) {
	log.Println(string(jwtPrivateKey))

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Email: email,
		Roles: roles,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	//token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtPrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		//return jwtPublicKey, nil
		return jwtPrivateKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
