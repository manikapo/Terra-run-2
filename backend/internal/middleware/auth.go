package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserEmailKey contextKey = "user_email"

func AuthRequired(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == auth {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization format"})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid claims"})
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing sub claim"})
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user id"})
		}

		c.Locals("user_id", userID)
		if email, ok := claims["email"].(string); ok {
			c.Locals("user_email", email)
		}
		return c.Next()
	}
}

func AuthOrGuest(jwtSecret string, allowGuest bool) fiber.Handler {
	jwtAuth := AuthRequired(jwtSecret)
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") && len(auth) > 7 {
			return jwtAuth(c)
		}

		if allowGuest {
			guestID := c.Get("X-Guest-User")
			if guestID != "" {
				userID, err := uuid.Parse(guestID)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid guest user id"})
				}
				c.Locals("user_id", userID)
				c.Locals("user_email", "guest@"+userID.String()+".local")
				c.Locals("is_guest", true)
				return c.Next()
			}
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "sign in required — use guest header X-Guest-User or Bearer token",
		})
	}
}

func InternalJobAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if secret == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "internal jobs disabled"})
		}
		if c.Get("X-Internal-Secret") != secret {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) uuid.UUID {
	if v := c.Locals("user_id"); v != nil {
		if id, ok := v.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}
