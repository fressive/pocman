package model

type Product struct {
	ID     uint32 `json:"id"`
	Name   string `json:"name"`
	Vendor string `json:"vendor"`
}
