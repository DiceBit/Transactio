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
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	jwtPrivateKey = []byte(os.Getenv("JWT_SECRET"))
}

func Int64ToNumericDate(ts int64) *jwt.NumericDate {
	tm := time.Unix(ts, 0)
	return jwt.NewNumericDate(tm)
}

func GenerateJWT(email string, roles []string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	nowTime := time.Now().Unix()

	claims := &Claims{
		Email: email,
		Roles: roles,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  Int64ToNumericDate(nowTime),
			NotBefore: Int64ToNumericDate(nowTime),
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
