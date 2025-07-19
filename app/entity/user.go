package entity

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	CompanyName  string `json:"company_name"`
	UserType     string `json:"user_type"` // admin/editor
	ParentID     *int   `json:"parent_id"`
}
