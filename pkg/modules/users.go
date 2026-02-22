package modules

type User struct {
	ID        int     `db:"id" json:"id"`
	Name      string  `db:"name" json:"name"`
	Email     string  `db:"email" json:"email"`
	Age       int     `db:"age" json:"age"`
	Password  string  `db:"password" json:"password,omitempty"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	DeletedAt *string `db:"deleted_at" json:"deleted_at,omitempty"`
}
