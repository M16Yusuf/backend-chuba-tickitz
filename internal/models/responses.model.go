package models

type Response struct {
	IsSuccess bool   `json:"is_success"  example:"true"`
	Code      int    `json:"code,omitempty"  example:"200"`
	Page      int    `json:"page,omitempty"  example:"1"`
	Msg       string `json:"message,omitempty"  example:"Example message success..."`
}

type ResponseData struct {
	Response
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Response
	Err string `json:"error" example:"Error message..."`
}

type TokenResponse struct {
	Response
	Token     string  `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ... "`
	Role      string  `json:"role" example:"user"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Avatar    *string `json:"avatar_path"`
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

// ===================== { Example error response for swagger } =====================

type BadRequestResponse struct {
	IsSuccess bool   `json:"is_success" example:"false"`
	Code      int    `json:"code,omitempty" example:"400"`
	Err       string `json:"error" example:"Example bad request error..."`
}

type InternalErrorResponse struct {
	IsSuccess bool   `json:"is_success" example:"false"`
	Code      int    `json:"code,omitempty" example:"500"`
	Err       string `json:"error" example:"Example Internal server error..."`
}
