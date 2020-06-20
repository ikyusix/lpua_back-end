package models

type ResponseWrapper struct {
	Success bool
	Message string      `json:"successMessage"`
	Data    interface{} `json:"data,omitempty"`
}

type Error struct {
	Message1 interface{} `json:"message"`
}
