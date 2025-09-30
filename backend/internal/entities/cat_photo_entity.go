package entities

type CatPhoto struct {
	ID        int    `json:"id" db:"id"`
	Url       string `json:"url" db:"url"`
	FileName  string `json:"file_name" db:"file_name"`
	FileSize  int    `json:"file_size" db:"file_size"`
	MimeType  string `json:"mime_type" db:"mime_type"`
	IsPrimary bool   `json:"is_primary" db:"is_primary"`
	CreatedAt string `db:"created_at"`
}
