package dto

import (
	"app/xonvera-core/internal/core/domain"
	"testing"
	"time"
)

func TestToUserResponse(t *testing.T) {
	if ToUserResponse(nil) != nil {
		t.Fatalf("expected nil response for nil user")
	}

	createdAt := time.Date(2026, 1, 29, 10, 0, 0, 0, time.UTC)
	user := &domain.User{
		ID:        10,
		Name:      "Jane",
		Email:     "jane@example.com",
		Phone:     "62812345678",
		CreatedAt: createdAt,
	}

	resp := ToUserResponse(user)
	if resp.ID != user.ID || resp.Email != user.Email || resp.Phone != user.Phone || resp.Name != user.Name {
		t.Fatalf("unexpected user response: %+v", resp)
	}
	if resp.CreatedAt != createdAt.Format("2006-01-02T15:04:05Z07:00") {
		t.Fatalf("unexpected created_at: %s", resp.CreatedAt)
	}
}

func TestToAuthResponse(t *testing.T) {
	if ToAuthResponse(nil) != nil {
		t.Fatalf("expected nil response for nil auth response")
	}

	user := &domain.User{ID: 1, Name: "Jane"}
	resp := ToAuthResponse(&domain.AuthResponse{
		User:         user,
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresAt:    123,
	})

	if resp.AccessToken != "access" || resp.RefreshToken != "refresh" || resp.ExpiresAt != 123 {
		t.Fatalf("unexpected auth response: %+v", resp)
	}
	if resp.User == nil || resp.User.ID != user.ID || resp.User.Name != user.Name {
		t.Fatalf("unexpected nested user response: %+v", resp.User)
	}
}

func TestToRegisterRequest(t *testing.T) {
	if ToRegisterRequest(nil) != nil {
		t.Fatalf("expected nil for nil register request")
	}

	req := &RegisterRequest{Name: "Jane", Email: "jane@example.com", Phone: "628123", Password: "secret"}
	mapped := ToRegisterRequest(req)
	if mapped.Name != req.Name || mapped.Email != req.Email || mapped.Phone != req.Phone || mapped.Password != req.Password {
		t.Fatalf("unexpected register request mapping: %+v", mapped)
	}
}

func TestToLoginAndRefreshRequests(t *testing.T) {
	if ToLoginRequest(nil) != nil || ToRefreshTokenRequest(nil) != nil {
		t.Fatalf("expected nil for nil requests")
	}

	login := &LoginRequest{Username: "jane", Password: "secret"}
	mappedLogin := ToLoginRequest(login)
	if mappedLogin.Username != login.Username || mappedLogin.Password != login.Password {
		t.Fatalf("unexpected login request mapping: %+v", mappedLogin)
	}

	refresh := &RefreshTokenRequest{RefreshToken: "token"}
	mappedRefresh := ToRefreshTokenRequest(refresh)
	if mappedRefresh.RefreshToken != refresh.RefreshToken {
		t.Fatalf("unexpected refresh request mapping: %+v", mappedRefresh)
	}
}

func TestToInvoiceResponse(t *testing.T) {
	if ToInvoiceResponse(nil, nil) != nil {
		t.Fatalf("expected nil response for nil invoice")
	}

	createdAt := time.Date(2026, 1, 29, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 1, 29, 13, 0, 0, 0, time.UTC)
	invoice := &domain.Invoice{
		ID:          100,
		AddTo:       "Client",
		InvoiceFor:  "Services",
		InvoiceFrom: "Company",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	items := []domain.InvoiceItem{
		{
			ID:          1,
			InvoiceID:   100,
			Description: "Work",
			Quantity:    2,
			UnitPrice:   50,
			Total:       100,
			CreatedAt:   createdAt,
		},
	}

	resp := ToInvoiceResponse(invoice, items)
	if resp.ID != invoice.ID || resp.AddTo != invoice.AddTo || resp.InvoiceFor != invoice.InvoiceFor || resp.InvoiceFrom != invoice.InvoiceFrom {
		t.Fatalf("unexpected invoice response: %+v", resp)
	}
	if len(resp.Items) != 1 || resp.Items[0].Description != "Work" {
		t.Fatalf("unexpected invoice items: %+v", resp.Items)
	}
	if resp.CreatedAt != createdAt.Format("2006-01-02T15:04:05Z07:00") || resp.UpdatedAt != updatedAt.Format("2006-01-02T15:04:05Z07:00") {
		t.Fatalf("unexpected timestamps: %+v", resp)
	}
}

func TestToInvoiceListResponse(t *testing.T) {
	invoices := []domain.Invoice{
		{
			ID:        1,
			AddTo:     "A",
			InvoiceFor: "Service",
			InvoiceFrom: "Company",
			CreatedAt: time.Date(2026, 1, 29, 8, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 1, 29, 9, 0, 0, 0, time.UTC),
		},
	}

	resp := ToInvoiceListResponse(invoices, 10, 20)
	if resp.Total != 1 || resp.Limit != 10 || resp.Offset != 20 {
		t.Fatalf("unexpected list response metadata: %+v", resp)
	}
	if len(resp.Invoices) != 1 || resp.Invoices[0].ID != invoices[0].ID {
		t.Fatalf("unexpected invoices list: %+v", resp.Invoices)
	}
}
