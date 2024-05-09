package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/asm/base64"
)

// Protected protect routes
func Protected() func(*fiber.Ctx) error {
	secret := os.Getenv("JWT_SECRET")

	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(secret)},
		ErrorHandler: jwtError,
		ContextKey:   "userData",
	})

}

func SetEmailAndUserID() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("userData").(*jwt.Token).Claims.(jwt.MapClaims)
		emailByte, err := base64.RawStdEncoding.DecodeString(user["ea"].(string))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
		}
		c.Locals("email", string(emailByte))

		userIDByte, err := base64.RawStdEncoding.DecodeString(user["ui"].(string))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
		}
		c.Locals("userID", string(userIDByte))

		return c.Next()
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})

	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}
