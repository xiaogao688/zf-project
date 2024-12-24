package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

// MyCustomClaims 自定义的 payload.
type MyCustomClaims struct {
	UserID               int
	Username             string
	GrantScope           string
	jwt.RegisteredClaims // 标准的 jwt payload
}

// 签名密钥
const sign_key = "12345"

// 随机字符串
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(str_len int) string {
	rand_bytes := make([]rune, str_len)
	for i := range rand_bytes {
		rand_bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(rand_bytes)
}

func generateTokenUsingHs256() (string, error) {
	claim := MyCustomClaims{
		UserID:     000001,
		Username:   "Tom",
		GrantScope: "read_user_info",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   // 签发者，表示谁签发的这个jwt
			Subject:   "Tom",                                           // 这个jwt的作用对象
			Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      // 签发受众，表示谁在用
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   // 过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), // 最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  // 签发时间
			ID:        randStr(10),                                     // wt ID, 类似于盐值
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(sign_key))
	return token, err
}

func parseTokenHs256(token_string string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(token_string, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sign_key), nil //返回签名密钥
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("claim invalid")
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}

	return claims, nil
}

func main() {

	token, err := generateTokenUsingHs256()
	if err != nil {
		panic(err)
	}
	fmt.Println("Token = ", token)

	time.Sleep(time.Second * 2)

	my_claim, err := parseTokenHs256(token)
	if err != nil {
		panic(err)
	}
	fmt.Println("my claim = ", my_claim)

}
