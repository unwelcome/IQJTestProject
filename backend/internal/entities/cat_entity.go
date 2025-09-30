package entities

type Cat struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Age         int    `json:"age" db:"age"`
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	CreatedBy   int    `json:"created_by" db:"created_by"`
}
