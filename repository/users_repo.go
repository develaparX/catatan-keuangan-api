package repository

import (
	"database/sql"
	"fmt"
	"livecode-catatan-keuangan/models"
)

type UserRepository interface {
	CreateNew(payload models.User) (models.User, error)
	GetByEmail(email string) (models.User, error)
}

// struct : menaruh depedency/fungsi/library yang akan digunakan
type userRepository struct {
	db *sql.DB
}

// CreateNew implements UserRepository.
func (p *userRepository) CreateNew(payload models.User) (models.User, error) {
	var user models.User
	err := p.db.QueryRow(`
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, password, created_at
	`, payload.Email, payload.Password).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// GetByUsername implements UserRepository.
func (u *userRepository) GetByEmail(email string) (models.User, error) {
	var user models.User

	query := `
        SELECT id, email, password, created_at
        FROM users
        WHERE email =$1
    `

	err := u.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password,
		&user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("email with %s not found", email)
		}
		return models.User{}, err
	}

	return user, nil
}

func NewUserRepository(database *sql.DB) UserRepository { // mereturn user repo
	return &userRepository{db: database}

}
