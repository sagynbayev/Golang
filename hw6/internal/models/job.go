package models

type Job struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`

	Price       int    `json:"price"`
	Description string `json:"description"`
	Date        string `json:"date"`
}
