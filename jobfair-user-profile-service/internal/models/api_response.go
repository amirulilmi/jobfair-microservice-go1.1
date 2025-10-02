package models

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

// Helper function for success response
func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// Helper function for error response
func ErrorResponse(message string, code string, error interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Code:    code,
		Error:   error,
	}
}

// Helper function for paginated response
func PaginatedSuccessResponse(message string, data interface{}, meta PaginationMeta) PaginatedResponse {
	return PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}
