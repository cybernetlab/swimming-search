package dto

import "github.com/cybernetlab/swimming-search/internal/domain"

type CreateUserInput struct {
	UserName string
	IsAdmin  bool
}

type CreateUserOutput struct {
	domain.User
}
