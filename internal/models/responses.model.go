package models

type Response struct {
	IsSuccess bool `json:"is_success"`
	Code      int  `json:"code,omitempty"`
	Page      int  `json:"page,omitempty"`
}

type ErrorResponse struct {
	Response
	Err string `json:"error" example:"Error message..."`
}

type TokenResponse struct {
	Response
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Miwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzU3MDM4NjQyfQ.J2MAUbAZvFpQl18BkSSyZOSMnbZxPziyZ6q6Bsuj8GU"`
}

type ProfileResponse struct {
	Response
	Data User `json:"data"`
}

type MoviesResponse struct {
	Response
	Data []MovieList `json:"data"`
}

type ScheduleResponse struct {
	Response
	Data []Schedule `json:"data"`
}

type SeatResponse struct {
	Response
	Data []BookedSeatBySchedule `json:"data"`
}

type DetailsMovieResponse struct {
	Response
	Data MovieDetails `json:"data"`
}

type UserDetailResponse struct {
	Response
	Data User `json:"data"`
}

type HistoiesResponse struct {
	Response
	Data []History `json:"data"`
}
