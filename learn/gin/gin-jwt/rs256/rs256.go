package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type MyCustomClaims struct {
	UserID     int
	Username   string
	GrantScope string
	jwt.RegisteredClaims
}

// 随机字符串
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(str_len int) string {
	rand_bytes := make([]rune, str_len)
	for i := range rand_bytes {
		rand_bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(rand_bytes)
}

// pkcs1
func parsePriKeyBytes(buf []byte) (*rsa.PrivateKey, error) {
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse key error")
	}
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

const pri_key = `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBAL00QsML/ovZle3Lq3C7QBo9s00ivsLhG2xlamhHOZDrjTGJX4OA
H27qQbDREcYXpUt5JqOt+KzB4MA/vUKCbT0CAwEAAQJBAINbkS5RWXxGqCzcRj6S
AkM1qxJWmRI7rwpmrqWPLYxKiS1i/i3bwSA3H+NODWIk1p2BWtycWzx5s3cNLn4b
gIECIQD6WuNzXxZHRIxRJQDRyEeWLsrRv9nkZJXHde78DoIZuQIhAMF4ZOgQX2hV
+y9YZmca2tW7etwGPmVjFWQd6JFtjyGlAiBFR9GZo76uijGqYusPIrVswhYuZUEP
CybHw8MWzY0DQQIgc4DDDWCo9QtP+MYX7Lo1p6BUCwOXQMRUwv6wGBKGfxkCIQDn
EKF3Ee6bnLT5DMfrnGY20RNg1Yes+14KkEyYsx0++Q==
-----END RSA PRIVATE KEY-----
`

const pub_key = `-----BEGIN RSA PUBLIC KEY-----
MEgCQQC9NELDC/6L2ZXty6twu0AaPbNNIr7C4RtsZWpoRzmQ640xiV+DgB9u6kGw
0RHGF6VLeSajrfisweDAP71Cgm09AgMBAAE=
-----END RSA PUBLIC KEY-----
`

func generateTokenUsingRS256() (string, error) {
	claim := MyCustomClaims{
		UserID:     000001,
		Username:   "Tom",
		GrantScope: "read_user_info",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   // 签发者
			Subject:   "Tom",                                           // 签发对象
			Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
			ID:        randStr(10),                                     // jwt ID, 类似于盐值
		},
	}
	rsa_pri_key, err := parsePriKeyBytes([]byte(pri_key))
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claim).SignedString(rsa_pri_key)
	return token, err
}

func parsePubKeyBytes(pub_key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub_key)
	if block == nil {
		return nil, errors.New("block nil")
	}
	pub_ret, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("x509.ParsePKCS1PublicKey error")
	}

	return pub_ret, nil
}

func parseTokenRs256(token_string string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(token_string, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		pub, err := parsePubKeyBytes([]byte(pub_key))
		if err != nil {
			fmt.Println("err = ", err)
			return nil, err
		}
		return pub, nil
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

	token, err := generateTokenUsingRS256()
	if err != nil {
		panic(err)
	}
	fmt.Println("Token = ", token)

	time.Sleep(time.Second * 2)

	my_claim, err := parseTokenRs256(token)
	if err != nil {
		panic(err)
	}
	fmt.Println("my claim = ", my_claim)

}
