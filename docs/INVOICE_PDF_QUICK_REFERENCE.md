# Invoice PDF API - Quick Reference Guide

## Overview

The Invoice PDF API endpoint allows authenticated users to retrieve or automatically generate PDF versions of their invoices.

## Endpoint

```
GET /api/v1/invoice/{invoiceID}/pdf
```

## Authentication

Bearer token authentication required (JWT token from login endpoint).

## Request Example

```bash
curl -X GET "http://localhost:3000/api/v1/invoice/2025010001/pdf" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Accept: application/pdf"
```

## Response

### Success (200)

Returns PDF file content as binary stream. Postman will automatically preview the PDF.

```
Headers:
Content-Type: application/pdf
Content-Disposition: inline; filename=invoice_2025010001.pdf
```

### Error Responses

| Status | Message               | Cause                                      |
|--------|----------------------|--------------------------------------------|
| 400    | Invalid invoice ID    | Malformed invoice ID format                |
| 401    | Unauthorized          | Missing or invalid authentication token    |
| 404    | Not found             | Invoice doesn't exist or user unauthorized |
| 500    | Internal server error | PDF generation failure                     |

## Workflow

```
┌─ Request Received
│
├─ Validate authentication (401 if missing)
│
├─ Validate invoice ID format (400 if invalid)
│
├─ Check if PDF exists in cache
│  ├─ YES → Return cached PDF (fast)
│  └─ NO  → Continue
│
├─ Fetch invoice from database
│  ├─ Not found → Return 404
│  └─ Unauthorized user → Return 404
│
├─ Generate PDF from invoice data
│
├─ Save PDF to assets/pdf/invoice_{id}.pdf
│
└─ Return PDF content to user
```

## Quick Start

### 1. Login

```bash
curl -X POST "http://localhost:3000/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your_password"
  }'
```

Response:
```json
{
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "...",
    "expires_in": 3600
  }
}
```

### 2. Create Invoice

```bash
curl -X POST "http://localhost:3000/api/v1/invoice" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 2025010001,
    "issuer": "My Company",
    "customer": "John Doe",
    "issue_date": "2025-01-15",
    "due_date": "2025-02-15 23:59:59",
    "note": "Payment due within 30 days",
    "items": [
      {
        "description": "Service",
        "qty": 1,
        "price": 50000
      }
    ]
  }'
```

### 3. Get PDF

```bash
curl -X GET "http://localhost:3000/api/v1/invoice/2025010001/pdf" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -o invoice.pdf
```

## Key Features

✅ **Automatic PDF Generation** - Generates from invoice data on first access  
✅ **Intelligent Caching** - Subsequent requests return cached PDF (fast)  
✅ **Browser Preview** - PDF served inline for direct preview (no download prompt)  
✅ **Authorization** - Users can only access their own invoices  
✅ **Error Handling** - Comprehensive error messages  
✅ **Async Ready** - Non-blocking PDF save operations  

## Postman Setup

### Import Collection

1. Open Postman
2. Click "Import" → "Upload Files"
3. Select `docs/postman_invoice_pdf.json`
4. Set environment variables:
   - `base_url`: http://localhost:3000
   - `access_token`: (auto-populated after login)

### Test PDF Endpoint

1. Execute "Login" request first
2. Execute "Create Invoice" to create test data
3. Execute "Get Invoice PDF - First Time (Generate)"
4. Click "Preview" tab to view PDF
5. Execute same request again to test caching

## File Storage

PDFs are stored in: `assets/pdf/invoice_{id}.pdf`

```
assets/
└── pdf/
    ├── invoice_2025010001.pdf
    ├── invoice_2025010002.pdf
    └── ...
```

## Common Issues & Solutions

### "PDF preview shows blank"
- Check `assets/pdf/` directory for PDF files
- Try downloading and opening in external PDF reader
- Check response headers for `application/pdf`

### "404 Not Found"
- Verify invoice ID exists (create one first)
- Verify invoice belongs to your user (not another user)
- Check correct ID format (numeric)

### "401 Unauthorized"
- Token may have expired (login again)
- Token not included in Authorization header
- Invalid token format (should be "Bearer TOKEN")

### "Slow PDF generation"
- First request includes generation + save time (~100-200ms)
- Subsequent requests use cache (~10-30ms)
- Check system disk space

## Development Notes

### Adding Custom PDF Styling

Edit `internal/core/services/pdf.go` - `buildContentStream()` method

### Changing PDF Storage Location

Update `internal/dependencies/wire.go` - `ProvidePDFService()` function:

```go
func ProvidePDFService(cfg *config.Config) (portService.PDFService, error) {
    pdfDir := "your/custom/path"  // Change this
    return services.NewPDFService(pdfDir)
}
```

### Adding PDF to Email

In the future, you can add email functionality to `InvoiceHandler.GetInvoicePDF()`:

```go
// Send PDF via email
pdfPath, _ := h.pdfService.SaveInvoicePDF(ctx, invoiceID, pdfData)
h.emailService.SendInvoicePDF(ctx, invoiceData.Customer, pdfPath)
```

## Related Endpoints

- `POST /api/v1/invoice` - Create invoice
- `GET /api/v1/invoice` - List invoices
- `PUT /api/v1/invoice` - Update invoice
- `POST /auth/login` - Get access token

## API Documentation

For complete API documentation, see [INVOICE_PDF_API.md](./INVOICE_PDF_API.md)
