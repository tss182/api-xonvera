package http

import (
	"context"
	"strconv"
	"time"

	"app/xonvera-core/internal/adapters/dto"
	"app/xonvera-core/internal/core/domain"
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

// CreateInvoice handles invoice creation
// @Summary Create a new invoice
// @Description Create a new invoice with items
// @Tags Invoice
// @Accept json
// @Produce json
// @Param request body dto.CreateInvoiceRequest true "Create Invoice Request"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/invoices [post]
func (h *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.rto)
	defer cancel()

	var req dto.CreateInvoiceRequest

	// Validate request
	if err := validator.HandlerBindingError(c, &req, validator.HandlerBody); err != nil {
		logger.Error("error when binding request in invoice service", zap.Strings("error validation body", err))
		return BadRequest(c, err)
	}

	// Convert DTO to domain
	invoice := &domain.Invoice{
		AddTo:       req.AddTo,
		InvoiceFor:  req.InvoiceFor,
		InvoiceFrom: req.InvoiceFrom,
	}

	items := make([]domain.InvoiceItem, len(req.Items))
	for i, itemReq := range req.Items {
		items[i] = domain.InvoiceItem{
			Description: itemReq.Description,
			Quantity:    itemReq.Quantity,
			UnitPrice:   itemReq.UnitPrice,
			Total:       itemReq.Quantity * itemReq.UnitPrice,
		}
	}

	err := h.service.CreateInvoice(ctx, invoice, items)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	// Get the created invoice with items
	invoiceItems, err := h.service.GetInvoiceItems(ctx, invoice.ID)
	if err != nil {
		logger.Error("failed to get invoice items after creation", zap.Int64("invoice_id", invoice.ID), zap.Error(err))
		// Return invoice without items rather than failing the request
		invoiceItems = []domain.InvoiceItem{}
	}
	response := dto.ToInvoiceResponse(invoice, invoiceItems)

	return Created(c, response)
}

// GetInvoiceByID handles getting an invoice by ID
// @Summary Get invoice by ID
// @Description Get invoice details by ID
// @Tags Invoice
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /api/v1/invoices/{id} [get]
func (h *InvoiceHandler) GetInvoiceByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.rto)
	defer cancel()

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return BadRequest(c, []string{"invalid invoice ID"})
	}

	invoice, err := h.service.GetInvoiceByID(ctx, id)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	// Get invoice items
	items, err := h.service.GetInvoiceItems(ctx, id)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	response := dto.ToInvoiceResponse(invoice, items)
	return OK(c, response)
}

// GetAllInvoices handles getting all invoices with pagination
// @Summary Get all invoices
// @Description Get all invoices with pagination
// @Tags Invoice
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/invoices [get]
func (h *InvoiceHandler) GetAllInvoices(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), h.rto)
	defer cancel()

	// Get pagination parameters
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	invoices, err := h.service.GetAllInvoices(ctx, limit, offset)
	if err != nil {
		return HandlerErrorGlobal(c, err)
	}

	response := dto.ToInvoiceListResponse(invoices, limit, offset)
	return OK(c, response)
}
