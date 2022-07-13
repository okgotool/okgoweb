package okauth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtTokenKey = "OkGoWebAuth20220808"
)

func SetJwtTokenKey(key string) error {
	if len(jwtTokenKey) < 2 {
		return errors.New("Length should max than 2")
	}

	jwtTokenKey = key

	return nil
}

func CreateJwtToken(tokenPairs map[string]interface{}) (string, error) {
	mapClaims := jwt.MapClaims{}
	for k, v := range tokenPairs {
		mapClaims[k] = v
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err := t.SignedString([]byte(jwtTokenKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func DecodeJwtToken(tokenStr string) (map[string]interface{}, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
		}
		return []byte(jwtTokenKey), nil
	})

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		tokenPairs := map[string]interface{}{}
		for k, v := range claims {
			tokenPairs[k] = v
		}
		return tokenPairs, nil
	} else {
		return nil, errors.New("token验证失败")
	}
}
