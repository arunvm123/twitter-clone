package auth

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//CustomClaims contains custom claims like UserName as well as standard claims
type CustomClaims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

//JWTTokenGeneration generates a jwt token with the given UserName and returns a JWT token
func JWTTokenGeneration(userName string) (string, error) {
	//Temporary signing key
	mySigningKey := []byte("TempSigningKey")

	// Create the Claims
	claims := &CustomClaims{
		userName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	//Signing the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {

		return "", err
	}
	return ss, nil
}

func ValidateJWT(t string) (interface{}, error) {
	mySigningKey := []byte("TempSigningKey")

	if t == "" {
		return nil, errors.New("Authorization token must be present")
	}
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return mySigningKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_name"].(string), nil
	} else {
		return nil, errors.New("Invalid authorization token")
	}
}
