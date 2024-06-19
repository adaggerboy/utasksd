package jwt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/adaggerboy/utasksd/models/config"
	"github.com/adaggerboy/utasksd/models/context"
	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	signingKey                       = "somekey"
	passPassphrase                   = "6368616e676520746869732070617373"
	iv                               = "my16digitIvKey12"
	validatingKey                    = signingKey
	issuer                           = "issuer"
	validIssuer                      = "issuer"
	method         jwt.SigningMethod = jwt.SigningMethodHS384
	isHMAC         bool              = true
	isECDSA        bool              = false
)

//https://medium.com/insiderengineering/aes-encryption-and-decryption-in-golang-php-and-both-with-full-codes-ceb598a34f41

func decrypt(encrypted string, key string) (string, error) {

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)

	length := len(ciphertext)
	unpadding := int(ciphertext[length-1])

	ciphertext = ciphertext[:(length - unpadding)]

	return fmt.Sprintf("%s", ciphertext), nil
}

func encrypt(plaintext string, key string) (string, error) {

	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

type CustomClaims struct {
	Container context.Session `json:"container"`
	jwt.RegisteredClaims
}

func ValidateIssuer(iss string) bool {
	return iss == validIssuer
}

func IsRevoked(string) bool {
	return false
}

func VerifyToken(tokenStr string) (claims context.Session, valid bool, err error) {
	fullClaims := &CustomClaims{}
	valid = false
	token, err := jwt.NewParser().ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if isECDSA {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
		} else if isHMAC {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
		}
		return []byte(validatingKey), nil
	})
	if err != nil {
		return
	}
	fullClaims, ok := token.Claims.(*CustomClaims)
	if !ok {
		err = fmt.Errorf("invalid claims type")
		return
	}
	valid = true
	if !fullClaims.ExpiresAt.After(time.Now()) {
		valid = false
	}
	if !ValidateIssuer(fullClaims.Issuer) {
		valid = false
	}
	if IsRevoked(fullClaims.ID) {
		valid = false
	}
	claims = fullClaims.Container
	encPass := claims.EncryptedPassword
	pass, err := decrypt(encPass, passPassphrase)
	if err != nil {
		return claims, false, err
	}
	claims.EncryptedPassword = pass
	return
}

func SignToken(claims context.Session, subject string, until time.Time) (tkn string, err error) {

	pass := claims.EncryptedPassword
	encPass, err := encrypt(pass, passPassphrase)
	if err != nil {
		return "", err
	}
	claims.EncryptedPassword = encPass

	signingKey := []byte(signingKey)

	fullClaims := CustomClaims{
		claims,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(until),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   subject,
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(method, fullClaims)
	tkn, err = token.SignedString(signingKey)
	return
}

func InitJWT(config config.JWTConfig) error {
	signingKey = config.SigningKey
	validatingKey = config.SigningKey
	passPassphrase = config.EncryptPassphrase
	issuer = config.Issuer
	validIssuer = config.Issuer

	switch config.Method {
	case "hmac":
		method = jwt.SigningMethodHS384
		isHMAC = true
		isECDSA = false
	case "ecdsa":
		method = jwt.SigningMethodES384
		isHMAC = false
		isECDSA = true
	default:
		return fmt.Errorf("invalid jwt method: %s", config.Method)
	}
	return nil
}
