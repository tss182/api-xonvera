package http

import (
	portService "app/xonvera-core/internal/core/ports/service"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type PackageHandler struct {
	service portService.PackageService
}

func NewPackageHandler(service portService.PackageService) *PackageHandler {
	return &PackageHandler{
		service: service,
	}
}

func (h *PackageHandler) GetPackages(c fiber.Ctx) error {
	resp, err := h.service.GetPackages(c.Context())
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, resp)
}

func (h *PackageHandler) GetPackageByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return HandlerErrorGlobal(c, fmt.Errorf("400:invalid package"))
	}

	pkg, err := h.service.GetPackageByID(c.Context(), id)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, pkg)
}

// fiber:context-methods migrated
