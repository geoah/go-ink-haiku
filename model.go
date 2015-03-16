package main

type ModelInterface interface {
	String() string
}

type Model struct {
	ID string `json:"id"`
}
