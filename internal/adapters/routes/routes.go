package routes

import (
	"app/xonvera-core/internal/adapters/handler/http"
	"app/xonvera-core/internal/adapters/middleware"
	"app/xonvera-core/internal/dependencies"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	r *dependencies.Application,
) {

	// Auth routes (public) with rate limiting
	auth := app.Group("/auth", middleware.AuthRateLimiter(r.Redis))
	{
		auth.Post("/register", r.AuthHandler.Register)
		auth.Post("/login", r.AuthHandler.Login)
		auth.Post("/refresh", r.AuthHandler.RefreshToken)
		auth.Post("/logout", r.AuthMiddleware.Authenticate(), r.AuthHandler.Logout)
	}

	// Package routes (public)
	app.Get("/packages", r.PackageHandler.GetPackages)
	app.Get("/packages/:id", r.PackageHandler.GetPackageByID)

	// Protected routes example
	appLogged := app.Use("/", r.AuthMiddleware.Authenticate())

	// Example protected route
	appLogged.Get("/profile", func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return http.NoAuth(c)
		}
		return http.OK(c, fiber.Map{
			"user_id": userID,
		})
	})

	//invoice
	invoice := appLogged.Group("/invoices")
	{
		invoice.Post("", r.InvoiceHandler.Create)
		invoice.Get("", r.InvoiceHandler.GetAllInvoices)
		invoice.Get("/:id", r.InvoiceHandler.GetInvoiceByID)
	}
}
