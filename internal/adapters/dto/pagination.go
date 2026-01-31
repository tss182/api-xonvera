package dto

type (
	PaginationRequest struct {
		Page   uint8  `query:"page"`
		Limit  uint8  `query:"limit"`
		Offset uint64 `query:"-"`
		Search string `query:"search"`
		UserID uint   `query:"user_id"`
	}

	PaginationResponse struct {
		Meta PaginationMetaResponse `json:"meta"`
		Data []any                  `json:"data"`
	}

	PaginationMetaResponse struct {
		Page      uint8  `json:"page"`
		Limit     uint8  `json:"limit"`
		TotalData uint64 `json:"total_data"`
		TotalPage uint8  `json:"total_page"`
	}
)
