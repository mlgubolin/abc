package application

import (
	"fmt"
)

// User representa um usuário do sistema
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"user_name"`
	Email    string `json:"user_email"`
	Password string `json:"user_password"`
	Document string `json:"user_document"`
}

// DTO de saída sem o campo Password
type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"user_name"`
	Email    string `json:"user_email"`
	Document string `json:"user_document"`
}

// Função para converter um único User para UserResponse
func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Document: user.Document,
	}
}

// Função para converter um slice de User para um slice de UserResponse
func ToUserResponses(users []User) []UserResponse {
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ToUserResponse(user))
	}
	return userResponses
}

// Validate verifica se os campos obrigatórios do usuário estão preenchidos
func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if u.Password == "" {
		return fmt.Errorf("password is required")
	}
	if u.Document == "" {
		return fmt.Errorf("document is required")
	}
	return nil
}

// ValidateUpdate verifica se os campos obrigatórios do usuário estão preenchidos
func (u *User) ValidateUpdate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// UserService define as operações disponíveis para um usuário
type UserService interface {
	CreateUser(user User) (User, error)
	GetUser(id int) (User, error)
	GetUsers() ([]User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id int) error
	GetUserByEmail(email string) (User, error)
}
