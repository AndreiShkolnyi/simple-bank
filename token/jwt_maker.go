package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)
	claims := jwt.MapClaims{
		"sub": maker.secretKey,               // subject (например, id пользователя)
		"name": payload.Username,                // произвольные данные
		"iat": payload.IssuedAt.Unix(),          // issued at (время выпуска)
		"exp": payload.ExpiredAt.Unix(), // время истечения токена (через 1 час)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	if !token.Valid {
		return nil, fmt.Errorf("недействительный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("не получилось преобразовать claims")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("поле exp отсутствует или неверного типа")
	}
	iatFloat, ok := claims["iat"].(float64)
	if !ok {
		return nil, fmt.Errorf("поле iat отсутствует или неверного типа")
	}

	// Переводим во время:
	exp := time.Unix(int64(expFloat), 0)
	iat := time.Unix(int64(iatFloat), 0)

	duration := exp.Sub(iat)

	payload := NewPayload(claims["name"].(string), duration)

	return payload, nil
}