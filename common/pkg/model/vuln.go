package model

type Vuln struct {
	ID          uint32   `json:"id"`
	Title       string   `json:"title"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
	Product     *Product `json:"product"`
}
