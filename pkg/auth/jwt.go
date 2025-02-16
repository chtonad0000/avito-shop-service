package auth

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtSecretKey = []byte(getSecretKey())

func GenerateToken(employeeUsername string) (string, error) {
	claims := Claims{
		EmployeeUsername: employeeUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func getSecretKey() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "pkg/auth/jwt_key/config.yaml"
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file ", err)
	}

	secret := viper.GetString("jwt_secret_key")
	if secret == "" {
		log.Fatal("jwt_secret_key is not set in config file")
	}
	return secret
}
