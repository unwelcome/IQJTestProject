package entities

type Cat struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Age         int    `json:"age" db:"age"`
	Description string `json:"description" db:"description"`
	Photos      []*CatPhoto
}

type CatPhoto struct {
	ID  int
	Url string
}
