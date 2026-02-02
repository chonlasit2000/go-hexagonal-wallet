package middleware

import (
	"strings"

	"github.com/chonlasit2000/e-wallet-hexagonal/internal/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. ดึง Header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing Token", "header is empty")
		}

		// 2. ตัดคำว่า "Bearer " ออก
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 3. Parse Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ใช้ secret ที่รับเข้ามาในการ verify
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid Token", err.Error())
		}

		// 4. ดึงข้อมูลใน Token (Claims)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid Token Claims", "cannot cast claims")
		}

		// 5. ดึง user_id และใส่ Context ไว้ให้ Handler อื่นใช้ต่อ
		// สำคัญ: ตอน Gen เราใช้ userID เป็น string (UUID) ดังนั้นตอนดึงต้อง cast เป็น string
		userID, ok := claims["user_id"].(string)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid User ID in Token", "user_id is not string")
		}

		c.Locals("user_id", userID) // ฝากไว้ที่ Locals

		return c.Next()
	}
}
