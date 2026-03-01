package modules

type Movie struct {
	ID          int     `db:"id" json:"id"`
	Title       string  `db:"title" json:"title"`
	Description string  `db:"description" json:"description"`
	Year        int     `db:"year" json:"year"`
	Rating      float64 `db:"rating" json:"rating"`
	CreatedAt   string  `db:"created_at" json:"created_at"`
	DeletedAt   *string `db:"deleted_at" json:"deleted_at,omitempty"`
}
