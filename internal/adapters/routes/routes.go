package routes

import (
	"app/xonvera-core/internal/adapters/handler/http"
	"app/xonvera-core/internal/adapters/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(
	app *fiber.App,
	authHandler *http.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	redisClient *redis.Client,
) {

	// Auth routes (public) with rate limiting
	auth := app.Group("/auth", middleware.AuthRateLimiter(redisClient))
	{
		auth.Post("/register", authHandler.Register)
		auth.Post("/login", authHandler.Login)
		auth.Post("/refresh", authHandler.RefreshToken)
		auth.Post("/logout", authMiddleware.Authenticate(), authHandler.Logout)
	}

	// Protected routes example
	logged := app.Group("/", authMiddleware.Authenticate())
	{
		// Example protected route
		logged.Get("/profile", func(c *fiber.Ctx) error {
			userID := c.Locals("userID").(uint)
			return http.OK(c, fiber.Map{
				"user_id": userID,
			})
		})
	}

}
