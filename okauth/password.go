package okauth

import (
	"golang.org/x/crypto/bcrypt"
)

func CreatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // 相同密码原文，每次调用产生的字符串都不同，需要使用CompareHashAndPassword才能判断是否相同
	if err != nil {
		return password, err
	}

	return string(hash), nil
}

// passwordHashInDb: pasword hash str stored in db
// passwordLogin: password from login input
func IsPasswordRight(passwordHashInDb string, passwordLogin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHashInDb), []byte(passwordLogin))
	if err != nil {
		return false
	} else {
		return true
	}
}
