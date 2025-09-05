package models

type NoRouteResponse struct {
	Message string
	Status  string
}

type TokenResponse struct {
	IsSuccess bool   `json:"is_success" example:"true"`
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Miwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzU3MDM4NjQyfQ.J2MAUbAZvFpQl18BkSSyZOSMnbZxPziyZ6q6Bsuj8GU"`
}

type ErrorResponse struct {
	IsSuccess bool   `json:"is_success" example:"false"`
	Err       string `json:"error" example:"Error message..."`
	Code      int    `json:"code,omitempty" example:"400"`
}

type SuccessResponse struct {
	IsSuccess bool `json:"is_success" example:"true"`
	Data      any  `json:"data"`
	Code      int  `json:"code,omitempty" example:"200"`
}
