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
	packages := app.Group("/packages")
	{
		packages.Get("/", appWire.PackageHandler.GetPackages)
		packages.Get("/:id", appWire.PackageHandler.GetPackageByID)
	}

	// Invoice routes (public)
	invoices := app.Group("/invoices")
	{
		invoices.Post("/", appWire.InvoiceHandler.CreateInvoice)
		invoices.Get("/", appWire.InvoiceHandler.GetAllInvoices)
		invoices.Get("/:id", appWire.InvoiceHandler.GetInvoiceByID)
	}

	// Protected routes example
	logged := app.Group("/", appWire.AuthMiddleware.Authenticate())
	{
		// Example protected route
		logged.Get("/profile", func(c *fiber.Ctx) error {
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
