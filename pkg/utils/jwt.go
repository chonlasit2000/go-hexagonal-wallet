package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ของจริงควรดึงจาก Config แต่นี่เพื่อความง่าย Hardcode ไปก่อน
var jwtSecret = []byte("my-super-secret-key")

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID, // UUID เป็น string
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
