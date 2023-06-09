package security

import (
	"errors"
	"fmt"
	"log"
	"os"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"strings"

	"time"

	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte("")
var jwtSingingMethod = jwt.SigningMethodHS256.Name

func EncryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func VerifyPassword(hashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func NewToken(userId string) (string, error) {
	claims := jwt.StandardClaims{
		Id:        userId,
		Issuer:    userId,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
	}

	tokenStr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenStr.SignedString(jwtSecretKey)
}

func ValidateSignedMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return jwtSecretKey, nil
}

func ParseToken(tokenStr string) (*jwt.StandardClaims, error) {

	log.Println(tokenStr)
	claims := new(jwt.StandardClaims)
	token, err := jwt.ParseWithClaims(strings.Split(tokenStr, " ")[1], claims, ValidateSignedMethod)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ok bool
	claims, ok = token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid auth token")
	}

	return claims, nil
}

func CheckRoleFromToken(token string, userRepo repository.UsersRepository, rc []string) (*models.User, error) {
	if tryApikey(token) {
		log.Println("authorized by apikey")
		user := models.User{
			Role: "server",
		}
		return &user, nil
	}

	if len(token) <= 7 || token[0:7] != "Bearer " {
		return nil, errors.New("invalid auth token")
	}

	payload, err := ParseToken(token)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// log.Println(payload)
	user, err := userRepo.GetById(payload.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(user)

	check := true
	for _, v := range rc {
		if v == user.Role || v == "all" {
			check = false
			break
		}
	}

	if check {
		return nil, errors.New("not permission")
	}

	return user, nil
}

func CheckRoleFromTokenInGet(token string, userRepo repository.UsersRepository, event string, role string) error {
	payload, err := ParseToken(token)
	if err != nil {
		log.Println(err)
		return err
	}
	user, err := userRepo.GetById(payload.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	if event == "get_all_profile" {
		if user.Role != "admin" && user.Role != "teacher" {
			return errors.New("You are not permission")
		}
	} else if event == "get_courses_id" {
		if user.Role != "admin" {
			if user.Role == "teacher" {

			}
		}
	}
	return nil
}

func tryApikey(authorization string) bool {
	apikey := os.Getenv("APIKEY")

	return authorization == apikey
}
