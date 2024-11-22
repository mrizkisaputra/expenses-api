package model

type ApiResponse struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	TotalItems int         `json:"total_items,omitempty"`
	TotalPages int         `json:"total_pages,omitempty"`
}
