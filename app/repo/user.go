package repo

import (
	"database/sql"
	"github.com/stealcash/AgentFlow/app/entity"
)

func CreateUser(db *sql.DB, u *entity.User) error {

	err := db.QueryRow(`
		INSERT INTO users (email, password_hash, company_name, user_type, parent_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, u.Email, u.PasswordHash, u.CompanyName, u.UserType, u.ParentID).Scan(&u.ID)

	return err
}

func GetUserByEmail(db *sql.DB, email string) (*entity.User, error) {
	row := db.QueryRow(`
		SELECT id, email, password_hash, company_name, user_type, parent_id
		FROM users WHERE email = $1
	`, email)

	var u entity.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CompanyName, &u.UserType, &u.ParentID)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
