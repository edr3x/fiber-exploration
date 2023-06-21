package model

type Response struct {
	Success bool   `json:"success"`
	Payload string `json:"payload"`
}

type FailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
