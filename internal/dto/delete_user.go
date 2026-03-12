package dto

import "github.com/cybernetlab/swimming-search/internal/domain"

type DeleteUserInput struct {
	UserName string
}

type DeleteUserOutput struct {
	domain.User
}
