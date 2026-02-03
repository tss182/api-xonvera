package http

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"app/xonvera-core/internal/adapters/dto"
	portService "app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/logger"
	"app/xonvera-core/internal/utils/validator"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type InvoiceHandler struct {
	service    portService.InvoiceService
	pdfService portService.PDFService
	rto        time.Duration
}

func NewInvoiceHandler(service portService.InvoiceService, pdfService portService.PDFService, rto time.Duration) *InvoiceHandler {
	return &InvoiceHandler{
		service:    service,
		pdfService: pdfService,
		rto:        rto,
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
func (h *InvoiceHandler) Get(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
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
func (h *InvoiceHandler) Create(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
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
func (h *InvoiceHandler) Update(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
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
// @Router /api/v1/invoice/{id}/pdf [get]
func (h *InvoiceHandler) GetInvoicePDF(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), h.rto)
	defer cancel()

	// Extract invoice ID from route parameter
	invoiceIDStr := c.Params("id")
	if invoiceIDStr == "" {
		return BadRequest(c, []string{"invoice ID is required"})
	}

	var invoiceID int64
	_, err := fmt.Sscanf(invoiceIDStr, "%d", &invoiceID)
	if err != nil || invoiceID <= 0 {
		return BadRequest(c, []string{"invalid invoice ID format"})
	}

	// Get user ID from context
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		return NoAuth(c)
	}

	// Step 1: Check if PDF already exists
	if h.pdfService.PDFExists(ctx, invoiceID) {
		logger.StdContextInfo(ctx, "retrieving existing invoice PDF", zap.Int64("invoice_id", invoiceID))
		pdfData, err := h.pdfService.GetInvoicePDF(ctx, invoiceID)
		if err == nil {
			c.Set("Content-Type", "application/pdf")
			c.Set("Content-Disposition", fmt.Sprintf("inline; filename=invoice_%d.pdf", invoiceID))
			return c.SendStream(bytes.NewReader(pdfData))
		}
		// Log error but continue to regenerate
		logger.StdContextWarn(ctx, "failed to retrieve existing PDF, will regenerate", zap.Error(err))
	}

	// Step 2: If PDF doesn't exist, fetch invoice data
	logger.StdContextInfo(ctx, "generating new invoice PDF", zap.Int64("invoice_id", invoiceID))

	invoiceData, err := h.service.GetByID(ctx, invoiceID, userID)
	if err != nil {
		logger.StdContextError(ctx, "failed to fetch invoice data", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return HandlerErrorGlobal(c, err)
	}

	if invoiceData == nil {
		return NotFound(c, []string{"invoice not found"})
	}

	// Step 3: Convert invoice items for PDF generation
	pdfItems := make([]dto.InvoiceItemDTO, len(invoiceData.Items))
	for i, item := range invoiceData.Items {
		pdfItems[i] = dto.InvoiceItemDTO{
			ID:          item.ID,
			InvoiceID:   item.InvoiceID,
			Description: item.Description,
			Qty:         item.Qty,
			Price:       item.Price,
			Total:       item.Total,
		}
	}

	// Step 4: Generate PDF
	pdfData, err := h.pdfService.GenerateInvoicePDF(ctx, invoiceData, pdfItems)
	if err != nil {
		logger.StdContextError(ctx, "failed to generate invoice PDF", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		return InternalServerError(c, []string{"failed to generate PDF"}, false)
	}

	// Step 5: Save PDF to filesystem
	pdfPath, err := h.pdfService.SaveInvoicePDF(ctx, invoiceID, pdfData)
	if err != nil {
		logger.StdContextError(ctx, "failed to save invoice PDF", zap.Error(err), zap.Int64("invoice_id", invoiceID))
		// Don't fail here - still return the PDF content even if save fails
		logger.StdContextWarn(ctx, "proceeding to return PDF despite save failure", zap.String("path", pdfPath))
	} else {
		logger.StdContextInfo(ctx, "invoice PDF saved successfully", zap.String("path", pdfPath))
	}

	// Step 6: Return PDF content for preview
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=invoice_%d.pdf", invoiceID))
	return c.SendStream(bytes.NewReader(pdfData))
}

// fiber:context-methods migrated
