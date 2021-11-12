package models

type Category struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
}
type Job struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	CategoryID string `json:"category_id" db:"category_id"`

	Price       int    `json:"price" db:"price"`
	Description string `json:"description" db:"description"`
	Date        string `json:"date" db:"date"`

}
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`

	Email    string `json:"email"`
	Password string `json:"password"`
}
