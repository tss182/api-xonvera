package routes

import (
	"app/xonvera-core/internal/adapters/handler/http"
	"app/xonvera-core/internal/adapters/middleware"
	"app/xonvera-core/internal/dependencies"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	appWire *dependencies.Application,
) {

	// Auth routes (public) with rate limiting
	auth := app.Group("/auth", middleware.AuthRateLimiter(appWire.Redis))
	{
		auth.Post("/register", appWire.AuthHandler.Register)
		auth.Post("/login", appWire.AuthHandler.Login)
		auth.Post("/refresh", appWire.AuthHandler.RefreshToken)
		auth.Post("/logout", appWire.AuthMiddleware.Authenticate(), appWire.AuthHandler.Logout)
	}

	// Package routes (public)
	app.Get("/packages", appWire.PackageHandler.GetPackages)
	app.Get("/packages/:id", appWire.PackageHandler.GetPackageByID)

	// Invoice routes (public)
	app.Post("/invoices", appWire.InvoiceHandler.CreateInvoice)
	app.Get("/invoices", appWire.InvoiceHandler.GetAllInvoices)
	app.Get("/invoices/:id", appWire.InvoiceHandler.GetInvoiceByID)

	// Protected routes example
	protected := app.Group("/protected", appWire.AuthMiddleware.Authenticate())
	{
		// Example protected route
		protected.Get("/profile", func(c *fiber.Ctx) error {
			userID, ok := c.Locals("userID").(uint)
			if !ok {
				return http.NoAuth(c)
			}
			return http.OK(c, fiber.Map{
				"user_id": userID,
			})
		})
	}

}
