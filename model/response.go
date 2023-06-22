package model

type Response struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

type FailureResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}
