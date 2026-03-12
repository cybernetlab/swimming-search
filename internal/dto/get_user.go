package dto

import "github.com/cybernetlab/swimming-search/internal/domain"

type GetUserInput struct {
	UserName string
}

type GetUserOutput struct {
	domain.User
}
