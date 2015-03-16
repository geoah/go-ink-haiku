package main

type Node struct {
	Model
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
}
