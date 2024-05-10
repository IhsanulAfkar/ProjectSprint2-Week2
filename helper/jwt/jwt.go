package jwt

import (
	"Week2/models"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type StaffClaims struct {
	jwt.RegisteredClaims
	Id 	string `json:"id"`
	PhoneNumber    string `json:"phoneNumber"`
}

func SignJWT(user models.Staff) string {
	// expiredIn := 28800 // 8 hours
	exp := time.Now().Add(time.Hour * 8)
	claims := StaffClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer: "Cat Socials",
		},
		Id: user.Id,
		PhoneNumber: user.Phone,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return ""
	}
	return signedToken

}
func ParseToken(jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", err
	}
	parsedToken, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		return "", errors.New("unable to parse claims")
	}
	id:=fmt.Sprint(parsedToken["pkId"])
	return id, nil
}