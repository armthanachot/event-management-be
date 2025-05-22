package model

type Response[T any, E any] struct {
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Errors  E      `json:"errors"`
	Success bool   `json:"success"`
}

type NoOmiteDataResponse[T any, E any] struct {
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
	Errors  E      `json:"errors,omitempty"`
	Success bool   `json:"success"`
}

type ResponseWithOffset[T any, E any] struct {
	Data    T      `json:"data"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Total   int64  `json:"total"`
	Message string `json:"message"`
	Errors  E      `json:"errors,omitempty"`
	Success bool   `json:"success"`
}

type ResponseMany[T any] struct {
	Data    []T    `json:"data"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Total   int    `json:"total"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ResponseWithoutOffset[T any] struct {
	Data    []T    `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}