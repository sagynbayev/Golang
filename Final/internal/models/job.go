package models

type(
	Job struct {
		ID       int    `json:"id" db:"id"`
		Name     string `json:"name" db:"name"`
		CategoryID string `json:"category_id" db:"category_id"`

		Price       int    `json:"price" db:"price"`
		Description string `json:"description" db:"description"`
		Date        string `json:"date" db:"date"`

	}

	JobsFilter struct{
		Query *string `json:"query"`
	}
)
