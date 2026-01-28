package dto

import (
	"app/xonvera-core/internal/core/domain"
)

// ToUserResponse converts domain.User to UserResponse
func ToUserResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToAuthResponse converts domain.AuthResponse to dto.AuthResponse
func ToAuthResponse(domainResp *domain.AuthResponse) *AuthResponse {
	if domainResp == nil {
		return nil
	}
	return &AuthResponse{
		User:         ToUserResponse(domainResp.User),
		AccessToken:  domainResp.AccessToken,
		RefreshToken: domainResp.RefreshToken,
		ExpiresAt:    domainResp.ExpiresAt,
	}
}

// ToRegisterRequest converts dto.RegisterRequest to domain.RegisterRequest
func ToRegisterRequest(req *RegisterRequest) *domain.RegisterRequest {
	if req == nil {
		return nil
	}
	return &domain.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}
}

// ToLoginRequest converts dto.LoginRequest to domain.LoginRequest
func ToLoginRequest(req *LoginRequest) *domain.LoginRequest {
	if req == nil {
		return nil
	}
	return &domain.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

// ToRefreshTokenRequest converts dto.RefreshTokenRequest to domain.RefreshTokenRequest
func ToRefreshTokenRequest(req *RefreshTokenRequest) *domain.RefreshTokenRequest {
	if req == nil {
		return nil
	}
	return &domain.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}
}

// ToInvoiceResponse converts domain.Invoice to InvoiceResponse
func ToInvoiceResponse(invoice *domain.Invoice, items []domain.InvoiceItem) *InvoiceResponse {
	if invoice == nil {
		return nil
	}
	
	var itemResponses []InvoiceItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, InvoiceItemResponse{
			ID:          item.ID,
			InvoiceID:   item.InvoiceID,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Total:       item.Total,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	
	return &InvoiceResponse{
		ID:          invoice.ID,
		AddTo:       invoice.AddTo,
		InvoiceFor:  invoice.InvoiceFor,
		InvoiceFrom: invoice.InvoiceFrom,
		Items:       itemResponses,
		CreatedAt:   invoice.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   invoice.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToInvoiceListResponse converts list of invoices to InvoiceListResponse
func ToInvoiceListResponse(invoices []domain.Invoice, limit, offset int) *InvoiceListResponse {
	var invoiceResponses []InvoiceResponse
	for _, invoice := range invoices {
		invoiceResponses = append(invoiceResponses, InvoiceResponse{
			ID:          invoice.ID,
			AddTo:       invoice.AddTo,
			InvoiceFor:  invoice.InvoiceFor,
			InvoiceFrom: invoice.InvoiceFrom,
			CreatedAt:   invoice.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   invoice.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	
	return &InvoiceListResponse{
		Invoices: invoiceResponses,
		Total:    len(invoices),
		Limit:    limit,
		Offset:   offset,
	}
}
