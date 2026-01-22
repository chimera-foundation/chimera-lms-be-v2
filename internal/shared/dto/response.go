package dto

// WebResponse is the generic template for all API responses
type WebResponse struct {
	Code   int         `json:"code"`   // e.g., 200, 201, 400
	Status string      `json:"status"` // e.g., "OK", "BAD_REQUEST", "INTERNAL_SERVER_ERROR"
	Data   interface{} `json:"data,omitempty"`  // The actual payload
	Errors interface{} `json:"errors,omitempty"` // Validation or error details
}