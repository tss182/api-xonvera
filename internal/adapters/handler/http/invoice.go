package http

import (
	"context"
	"time"

	"app/xonvera-core/internal/adapters/dto"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"
	"app/xonvera-core/internal/utils/validator"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type InvoiceHandler struct {
	service portService.InvoiceService
	rto     time.Duration
}

func NewInvoiceHandler(service portService.InvoiceService, rto time.Duration) *InvoiceHandler {
	return &InvoiceHandler{
		service: service,
		rto:     rto,
	}
}

// GetAllInvoices handles getting all invoices with pagination
// @Summary Get all invoices
// @Description Get all invoices with pagination
// @Tags Invoice
// @Accept json
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} Resp
// @Failure 400 {object} Resp
// @Router /api/v1/invoices [get]
func (h *InvoiceHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.rto)
	defer cancel()

	var req dto.PaginationRequest
	if err := validator.HandlerBindingError(c, &req, validator.HandlerQuery); err != nil {
		return BadRequest(c, []string{"invalid pagination parameters"})
	}

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		return NoAuth(c)
	}
	req.UserID = userID

	res, err := h.service.Get(ctx, &req)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return Page(c, res)
}

// Create handles invoice creation
// @Summary Create a new invoice
// @Description Create a new invoice with items
// @Tags Invoice
// @Accept json
// @Produce json
// @Param request body dto.CreateInvoiceRequest true "Create Invoice Request"
// @Success 200 {object} Resp
// @Failure 400 {object} Resp
// @Router /api/v1/invoices [post]
func (h *InvoiceHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.rto)
	defer cancel()

	var req dto.InvoiceRequest
	var ok bool

	userIDVal := c.Locals("userID")
	req.UserID, ok = userIDVal.(uint)
	if !ok || req.UserID == 0 {
		return NoAuth(c)
	}

	// Validate request
	if err := validator.HandlerBindingError(c, &req, validator.HandlerBody, "id"); err != nil {
		logger.Error("error when binding request in invoice service", zap.Strings("error validation body", err))
		return BadRequest(c, err)
	}

	err := h.service.Create(ctx, &req)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, nil)
}

// Update handles invoice update
// @Summary Update invoice
// @Description Update invoice details with items
// @Tags Invoice
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Param request body dto.InvoiceRequest true "Update Invoice Request"
// @Success 200 {object} Resp
// @Failure 400 {object} Resp
// @Failure 404 {object} Resp
// @Router /api/v1/invoices/{id} [put]
func (h *InvoiceHandler) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.rto)
	defer cancel()

	var req dto.InvoiceRequest
	var ok bool

	userIDVal := c.Locals("userID")
	req.UserID, ok = userIDVal.(uint)
	if !ok || req.UserID == 0 {
		return NoAuth(c)
	}

	if err := validator.HandlerBindingError(c, &req, validator.HandlerBody); err != nil {
		logger.Error("error when binding request in invoice service", zap.Strings("error validation body", err))
		return BadRequest(c, err)
	}

	if err := h.service.Update(ctx, &req); err != nil {
		return HandlerErrorGlobal(c, err)
	}

	return OK(c, nil)
}
