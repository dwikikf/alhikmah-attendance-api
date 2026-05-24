package response

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func Success(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func SuccessWithPagination(message string, data interface{}, pagination interface{}) PaginatedResponse {
	return PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}
}

func Error(message string) Response {
	return Response{
		Success: false,
		Message: message,
	}
}

// Optionally, we can have a helper for Validation errors if needed.
type ValidationErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

func ValidationError(message string, errors interface{}) ValidationErrorResponse {
	return ValidationErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}
