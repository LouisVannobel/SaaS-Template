package middleware

import (
	"fmt"
	"strings"

	"github.com/LouisVannobel/SaaS-Template/backend/auth"
	"github.com/gofiber/fiber/v2"
)

// JWTProtected is a middleware that checks for a valid JWT token
func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authorization := c.Get("Authorization")
		fmt.Printf("Authorization header: %s\n", authorization)

		// Check if the header is empty or doesn't start with "Bearer "
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			fmt.Println("Authorization header is empty or doesn't start with 'Bearer '")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid or missing token",
			})
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authorization, "Bearer ")
		fmt.Printf("Token: %s\n", tokenString)

		// Validate the token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			fmt.Printf("Token validation error: %v\n", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Unauthorized: %v", err),
			})
		}

		// Store user ID in context for later use
		c.Locals("userID", claims.UserID)
		fmt.Printf("User ID: %d\n", claims.UserID)

		// Continue to the next middleware/handler
		return c.Next()
	}
}
