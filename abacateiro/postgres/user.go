package postgres

import (
	"application"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(user application.User) (application.User, error) {

	if err := user.Validate(); err != nil {
		return application.User{}, fmt.Errorf("invalid user: %w", err)
	}

	query := `INSERT INTO users (name, email, password, document) VALUES ($1, $2, $3, $4) RETURNING id`

	err := s.db.QueryRow(context.Background(), query, user.Name, user.Email, user.Password, user.Document).Scan(&user.ID)

	if err != nil {
		return application.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(id int) (application.User, error) {

	var user application.User

	query := `SELECT id, name, email FROM users WHERE id = $1`

	err := s.db.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return application.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUsers() ([]application.User, error) {

	query := `SELECT id, name, email, document FROM users`

	rows, err := s.db.Query(context.Background(), query)

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	defer rows.Close()

	var users []application.User

	for rows.Next() {
		var user application.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Document); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) UpdateUser(user application.User) (application.User, error) {

	if err := user.ValidateUpdate(); err != nil {
		return application.User{}, fmt.Errorf("invalid user: %w", err)
	}

	query := `UPDATE users SET name = $1 WHERE id = $2`

	_, err := s.db.Exec(context.Background(), query, user.Name, user.ID)

	if err != nil {
		return application.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *UserService) DeleteUser(id int) error {

	query := `DELETE FROM users WHERE id = $1`

	_, err := s.db.Exec(context.Background(), query, id)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) GetUserByEmail(email string) (application.User, error) {
	var user application.User
	query := "SELECT id, name, email, password, document FROM users WHERE email = ?"

	err := s.db.QueryRow(context.Background(), query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Document)
	if err != nil {
		if err == sql.ErrNoRows {
			return application.User{}, errors.New("user not found")
		}
		return application.User{}, err
	}
	return user, nil
}
