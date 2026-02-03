package routes

import (
	"app/xonvera-core/internal/adapters/middleware"
	"app/xonvera-core/internal/dependencies"

	"github.com/gofiber/fiber/v3"
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

	//invoice
	invoice := appLogged.Group("/invoice")
	{
		invoice.Post("", r.InvoiceHandler.Create)
		invoice.Get("", r.InvoiceHandler.Get)
		invoice.Put("", r.InvoiceHandler.Update)
	}
}
