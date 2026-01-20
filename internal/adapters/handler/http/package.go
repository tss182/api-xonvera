package http

import (
	portRepository "app/xonvera-core/internal/core/ports/repository"

	"github.com/gofiber/fiber/v2"
)

type PackageHandler struct {
	packageRepo portRepository.PackageRepository
}

func NewPackageHandler(packageRepo portRepository.PackageRepository) *PackageHandler {
	return &PackageHandler{
		packageRepo: packageRepo,
	}
}

func (h *PackageHandler) GetPackages(c *fiber.Ctx) error {
	packages, err := h.packageRepo.GetAll(c.Context())
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

	pkg, err := h.packageRepo.GetByID(c.Context(), id)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch package")
	}

	if pkg == nil {
		return ErrorResponse(c, fiber.StatusNotFound, "Package not found")
	}

	return SuccessResponse(c, fiber.StatusOK, "Package fetched successfully", pkg)
}

type CreatePackageRequest struct {
	ID           string `json:"id" validate:"required,min=1,max=255"`
	Name         string `json:"name" validate:"required,min=1,max=255"`
	Price        int    `json:"price" validate:"required,min=0"`
	DiscountType string `json:"discount_type" validate:"required,oneof=PERCENTAGE AMOUNT"`
	Discount     int    `json:"discount" validate:"required,min=0"`
}
