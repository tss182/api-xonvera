package http

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"app/xonvera-core/internal/core/domain"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"
	"app/xonvera-core/internal/utils/validator"

	"github.com/gofiber/fiber/v3"
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

// GetAllinvoice handles getting all invoice with pagination
// @Summary Get all invoice
// @Description Get all invoice with pagination
// @Tags Invoice
// @Accept json
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} Resp
// @Failure 400 {object} Resp
// @Router /invoice [get]
func (h *InvoiceHandler) Get(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
	defer cancel()

	var req domain.PaginationRequest
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
// @Param request body domain.InvoiceRequest true "Create Invoice Request"
// @Success 200 {object} Resp
// @Failure 400 {object} Resp
// @Router /invoice [post]
func (h *InvoiceHandler) Create(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
	defer cancel()

	var req domain.InvoiceRequest
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

	if len(req.Items) == 0 {
		return BadRequest(c, []string{"at least one invoice item is required"})
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
// @Param request body domain.InvoiceRequest true "Update Invoice Request"
// @Success 200 {object} Resp
// @Failure 400 {object} Resp
// @Failure 404 {object} Resp
// @Router /invoice/{id} [put]
func (h *InvoiceHandler) Update(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
	defer cancel()

	var req domain.InvoiceRequest
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

// GetInvoicePDF handles retrieving or generating invoice PDF
// @Summary Get invoice PDF
// @Description Retrieve existing invoice PDF or generate new one based on invoice data
// @Tags Invoice
// @Accept json
// @Produce application/pdf
// @Param id path int true "Invoice ID"
// @Success 200 {file} application/pdf
// @Failure 400 {object} Resp
// @Failure 404 {object} Resp
// @Router /invoice/{id}/pdf [get]
func (h *InvoiceHandler) GetInvoicePDF(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
	defer cancel()

	invoiceID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || invoiceID <= 0 {
		return BadRequest(c, []string{"invalid invoice ID format"})
	}

	// Get user ID from context
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		return NoAuth(c)
	}

	res, err := h.service.GetPDF(ctx, invoiceID, userID)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=invoice_%d.pdf", invoiceID))
	return c.SendStream(bytes.NewReader(res))
}
