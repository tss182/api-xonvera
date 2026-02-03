# Invoice PDF API Documentation

## Overview

The Invoice PDF API endpoint provides functionality to retrieve or generate PDF representations of invoices. The endpoint automatically handles PDF generation, storage, and retrieval with intelligent caching.

## Endpoint Specification

### Get Invoice PDF

**Endpoint:** `GET /api/v1/invoice/{id}/pdf`

**Authentication:** Required (Bearer Token)

**Content-Type:** `application/pdf`

#### Path Parameters

| Parameter | Type    | Required | Description        |
|-----------|---------|----------|--------------------|
| `id`      | integer | Yes      | The invoice ID      |

#### Response

- **Status Code 200:** PDF file content is returned with proper headers for browser preview
- **Status Code 400:** Bad request (invalid invoice ID format)
- **Status Code 401:** Unauthorized (missing or invalid authentication token)
- **Status Code 404:** Invoice not found or user does not have access
- **Status Code 500:** Internal server error

#### Response Headers

```
Content-Type: application/pdf
Content-Disposition: inline; filename=invoice_{id}.pdf
```

The `inline` disposition tells the browser to preview the PDF directly instead of downloading it.

## Workflow

The endpoint performs the following operations:

```
1. Validate invoice ID format
   ↓
2. Check if PDF already exists in assets/pdf directory
   ├─→ YES: Retrieve and return existing PDF
   └─→ NO: Continue to step 3
   ↓
3. Fetch invoice data from database
   ├─→ Verify invoice belongs to authenticated user
   └─→ Retrieve all invoice items
   ↓
4. Generate PDF from invoice data
   ↓
5. Save PDF to assets/pdf/{invoice_id}.pdf
   ↓
6. Return PDF content for browser preview
```

## Design Considerations

### 1. Reusability & Modularity

The implementation leverages:

- **PDF Service Interface** (`portService.PDFService`): Decoupled PDF operations
- **Repository Pattern** (`InvoiceRepository`): Database operations isolation
- **Dependency Injection**: Wire framework manages all dependencies
- **DTOs**: Clean separation between domain and API layers

### 2. PDF Storage Strategy

- **Location:** `assets/pdf/` directory (auto-created if missing)
- **Naming Convention:** `invoice_{invoiceID}.pdf`
- **Caching:** PDFs are cached after first generation for performance
- **File Permissions:** `0644` (readable by all, writable by owner)

### 3. Error Handling

| Error Scenario                          | Status Code | Action                                   |
|-----------------------------------------|-------------|------------------------------------------|
| Invalid invoice ID format               | 400         | Return error message                     |
| Invoice not found                       | 404         | Return not found error                   |
| User unauthorized (different user)      | 404         | Return not found (security: no leak)     |
| PDF generation failure                  | 500         | Return internal error                    |
| PDF save failure                        | Partial     | Return PDF but log warning               |
| Existing PDF retrieval failure          | Regenerate  | Generate new PDF instead of failing      |

### 4. Security

- **Authentication:** Required via Bearer token middleware
- **Authorization:** User can only access their own invoices
- **Validation:** Invoice ID format validated before processing
- **Error Messages:** Do not leak information about other users' invoices

## API Testing with Postman

### Setup Instructions

#### 1. Import Authentication

```
1. Create a new request or use existing collection
2. Go to "Tests" tab
3. Add the following script to extract and save tokens:

var jsonData = pm.response.json();
if (jsonData.data && jsonData.data.access_token) {
    pm.environment.set("access_token", jsonData.data.access_token);
}
```

#### 2. Configure Bearer Token

```
1. Go to "Authorization" tab
2. Select "Bearer Token" type
3. Enter: {{access_token}}
4. This uses the token from your environment variable
```

#### 3. Test the Endpoint

**Request:**
```
GET http://localhost:3000/api/v1/invoice/1/pdf
Authorization: Bearer {{access_token}}
```

**Expected Response:**
- Status: `200 OK`
- Binary PDF content in response body
- Postman will display a "Preview" button to view the PDF

### Step-by-Step Testing Guide

#### Step 1: Login to Get Access Token

```bash
Method: POST
URL: http://localhost:3000/auth/login
Body (JSON):
{
    "email": "user@example.com",
    "password": "your_password"
}
```

Response:
```json
{
    "success": true,
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIs...",
        "refresh_token": "...",
        "expires_in": 3600
    }
}
```

**In Postman:**
- Copy the `access_token` value
- Set it in your collection variables or environment as `access_token`

#### Step 2: Create an Invoice

```bash
Method: POST
URL: http://localhost:3000/api/v1/invoice
Authorization: Bearer {{access_token}}
Content-Type: application/json

Body:
{
    "id": 2025010001,
    "issuer": "Company Name",
    "customer": "Customer Name",
    "issue_date": "2025-01-15",
    "due_date": "2025-02-15 23:59:59",
    "note": "Payment terms: Net 30",
    "items": [
        {
            "description": "Web Development Services",
            "qty": 10,
            "price": 5000
        },
        {
            "description": "Consulting",
            "qty": 5,
            "price": 2000
        }
    ]
}
```

Note the `id` field - you'll use this in the next step.

#### Step 3: Generate/Retrieve Invoice PDF

```bash
Method: GET
URL: http://localhost:3000/api/v1/invoice/2025010001/pdf
Authorization: Bearer {{access_token}}
```

**Expected Result:**
- Status: `200 OK`
- Response body contains binary PDF data
- In Postman, click "Preview" tab to view the PDF in the browser

### Advanced Testing Scenarios

#### Test 1: PDF Caching Verification

```
1. Generate PDF for invoice #2025010001 (note response time)
2. Generate PDF again for same invoice (should be faster)
3. Verify both responses are identical binary content
```

#### Test 2: Unauthorized Access

```
Method: GET
URL: http://localhost:3000/api/v1/invoice/2025010001/pdf
(No Authorization header or invalid token)

Expected Result: 401 Unauthorized
```

#### Test 3: Invalid Invoice ID

```
Method: GET
URL: http://localhost:3000/api/v1/invoice/invalid-id/pdf
Authorization: Bearer {{access_token}}

Expected Result: 400 Bad Request
```

#### Test 4: Non-existent Invoice

```
Method: GET
URL: http://localhost:3000/api/v1/invoice/9999999999/pdf
Authorization: Bearer {{access_token}}

Expected Result: 404 Not Found
```

#### Test 5: Other User's Invoice (Security Test)

```
1. Login as User A, get their invoice ID (e.g., 2025010001)
2. Logout/Switch to User B
3. Try to access User A's invoice:
   GET /api/v1/invoice/2025010001/pdf
   Authorization: Bearer {{user_b_token}}

Expected Result: 404 Not Found (secure: doesn't reveal existence)
```

## Configuration

### Environment Variables

```bash
# PDF storage directory (created automatically if not exists)
PDF_STORAGE_PATH="assets/pdf"

# Request timeout for PDF operations (default: 5 seconds)
REQUEST_TIMEOUT="5s"
```

### Directory Structure

```
project-root/
├── assets/
│   └── pdf/
│       ├── invoice_2025010001.pdf
│       ├── invoice_2025010002.pdf
│       └── invoice_2025010003.pdf
├── internal/
│   ├── adapters/
│   │   ├── handler/http/
│   │   │   └── invoice.go          (PDF endpoint handler)
│   │   ├── repositories/
│   │   └── routes/
│   │       └── routes.go            (PDF route registration)
│   ├── core/
│   │   ├── ports/
│   │   │   └── service/
│   │   │       └── pdf.go           (PDF service interface)
│   │   └── services/
│   │       └── pdf.go               (PDF service implementation)
│   └── dependencies/
│       └── wire.go                  (Dependency injection setup)
└── docs/
    └── INVOICE_PDF_API.md           (This file)
```

## Implementation Details

### Handler: `InvoiceHandler.GetInvoicePDF()`

Located in: `internal/adapters/handler/http/invoice.go`

**Responsibilities:**
- Extract and validate invoice ID from URL parameters
- Verify user authentication
- Orchestrate PDF retrieval or generation
- Return PDF content with proper headers

**Key Methods Called:**
- `pdfService.PDFExists()` - Check if PDF is cached
- `pdfService.GetInvoicePDF()` - Retrieve cached PDF
- `invoiceService.GetByID()` - Fetch invoice from database
- `pdfService.GenerateInvoicePDF()` - Generate new PDF
- `pdfService.SaveInvoicePDF()` - Store PDF on disk

### Service: `PDFService`

Located in: `internal/core/services/pdf.go`

**Key Methods:**
- `GenerateInvoicePDF()` - Creates PDF from invoice data
- `GetInvoicePDF()` - Reads cached PDF from filesystem
- `SaveInvoicePDF()` - Persists PDF to disk
- `PDFExists()` - Checks if PDF file exists

**PDF Format:** Simple text-based PDF with:
- Invoice header and metadata
- Invoice items table
- Total amount calculation
- Generation timestamp

### Repository Enhancement

Added to `InvoiceRepository` interface:
- `GetItemsByInvoiceID()` - Retrieves all items for a specific invoice

Implementation in: `internal/adapters/repositories/sql/invoice.go`

## Performance Considerations

1. **PDF Caching:** Generated PDFs are cached in `assets/pdf/` directory
   - First access: ~100-200ms (includes generation + save)
   - Subsequent accesses: ~10-30ms (filesystem read only)

2. **Database Optimization:** 
   - Invoice fetched via `GetByID()` 
   - Items fetched via `GetItemsByInvoiceID()`
   - Both operations use indexed lookups

3. **Disk Space:**
   - Each PDF is ~2-10KB depending on item count
   - Cleanup strategy: Manual deletion or scheduled maintenance job

## Troubleshooting

### Issue: "PDF directory permission denied"

**Solution:**
```bash
# Ensure assets/pdf directory is writable
chmod 755 assets/pdf
```

### Issue: "404 on valid invoice"

**Possible Causes:**
1. Invoice belongs to different user (check authentication)
2. Invoice ID doesn't exist (create test invoice first)
3. Incorrect URL format (verify path parameter)

### Issue: "PDF generation very slow"

**Debug Steps:**
1. Check if PDF already exists (`assets/pdf/invoice_*.pdf`)
2. Monitor database query times
3. Check system disk space
4. Verify no file system issues

### Issue: "Postman shows blank PDF preview"

**Solutions:**
1. Verify PDF was actually generated (check `assets/pdf/` directory)
2. Try downloading and opening in external PDF reader
3. Check response headers (should be `application/pdf`)
4. Clear Postman cache and try again

## Future Enhancements

1. **PDF Generation:**
   - Upgrade to professional PDF library (e.g., `go-pdf`, `gofpdf`)
   - Add company logo and branding
   - Support for custom templates

2. **Storage:**
   - Move PDFs to cloud storage (S3, GCS)
   - Implement automatic cleanup of old PDFs
   - Add backup strategy

3. **Features:**
   - Email PDF directly to customer
   - Generate PDF for multiple invoices
   - PDF digital signature support
   - Watermark for draft/preview invoices

4. **Performance:**
   - Async PDF generation with background jobs
   - CDN integration for PDF delivery
   - Compression for large PDFs

## Related Documentation

- [Invoice API Reference](./INVOICE_API.md) - Full invoice CRUD operations
- [Authentication Guide](./AUTH.md) - Bearer token authentication
- [Error Handling](./ERROR_HANDLING.md) - Standard error responses
