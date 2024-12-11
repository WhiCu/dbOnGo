package types

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/WhiCu/mdb/config"
)

type User struct {
	Email          string `json:"email" bson:"email"`
	Phone          string `json:"phone" bson:"phone"`
	UserName       string `json:"username" bson:"username"`
	Login          string `json:"login" bson:"login"`
	Password       string `json:"password" bson:"password"`
	Token          string `json:"token" bson:"token"`
	TemporaryToken string `json:"temporaryToken" bson:"temporaryToken"`
}

var tokenLength = config.MustGetInt("TOKEN_LENGTH")

func New(
	Email string,
	Phone string,
	UserName string,
	Login string,
	Password string,
) *User {
	u := &User{
		Email:    Email,
		Phone:    Phone,
		UserName: UserName,
		Login:    Login,
		Password: Password,
	}

	u.GenerateTokens()

	return u
}

func (u *User) GenerateTokens() {
	u.generateToken()
	u.generateTemporaryToken()
}

func (u *User) generateToken() {

	bytes := make([]byte, tokenLength)
	_, err := rand.Read(bytes)
	if err != nil {
		//TODO: add logging
		log.Fatal(err)
	}

	u.Token = base64.RawURLEncoding.EncodeToString(bytes)
}

// generateTemporaryToken создает временный токен, который зависит от текущего времени и основного токена.
func (u *User) generateTemporaryToken() {
	// Получаем текущее время
	currentTime := time.Now().Unix()

	// Создаем строку, которая будет использоваться для хеширования
	dataToHash := fmt.Sprintf("%s%d", u.Token, currentTime)

	// Хэшируем строку с помощью SHA-256
	hash := sha256.Sum256([]byte(dataToHash))

	// Кодируем хэш в base64 для удобства
	u.TemporaryToken = base64.RawURLEncoding.EncodeToString(hash[:])
}

func (u *User) Json() []byte {
	data, err := json.Marshal(u)

	if err != nil {
		//TODO: add logging
		log.Fatal("User.Json error: ", err.Error())
	}

	return data
}
