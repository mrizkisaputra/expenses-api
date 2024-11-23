package model

type ApiResponse struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	TotalItems int64       `json:"total_items,omitempty"`
	TotalPages int64       `json:"total_pages,omitempty"`
}

type SearchExpenseRequestQueryParam struct {
	Page      int    `form:"page" validate:"omitempty,numeric,min=1"`
	Limit     int    `form:"limit" validate:"omitempty,numeric,min=10"`
	Filter    string `form:"filter" validate:"omitempty,oneof=last_week last_month last_3_month custom"`
	StartDate string `form:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" validate:"omitempty,datetime=2006-01-02"`
}
