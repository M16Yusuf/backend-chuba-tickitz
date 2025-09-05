package models

type Response struct {
	IsSuccess bool `json:"is_success"`
	Code      int  `json:"code,omitempty"`
}

type ErrorResponse struct {
	Response
	Err string `json:"error" example:"Error message..."`
}

type NoRouteResponse struct {
	ErrorResponse
	Message string
}

type TokenResponse struct {
	Response
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Miwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzU3MDM4NjQyfQ.J2MAUbAZvFpQl18BkSSyZOSMnbZxPziyZ6q6Bsuj8GU"`
}

type ProfileResponse struct {
	Response
	Data User
}
