package http

import (
	portService "app/xonvera-core/internal/core/ports/service"

	"github.com/gofiber/fiber/v2"
)

type PackageHandler struct {
	service portService.PackageService
}

func NewPackageHandler(service portService.PackageService) *PackageHandler {
	return &PackageHandler{
		service: service,
	}
}

func (h *PackageHandler) GetPackages(c *fiber.Ctx) error {
	packages, err := h.service.GetPackages(c.Context())
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch packages: "+err.Error())
	}

	return SuccessResponse(c, fiber.StatusOK, "Packages fetched successfully", packages)
}

func (h *PackageHandler) GetPackageByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid package ID")
	}

	pkg, err := h.service.GetPackageByID(c.Context(), id)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch package")
	}

	if pkg == nil {
		return ErrorResponse(c, fiber.StatusNotFound, "Package not found")
	}

	return SuccessResponse(c, fiber.StatusOK, "Package fetched successfully", pkg)
}
