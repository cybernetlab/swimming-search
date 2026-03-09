package dto

import "github.com/cybernetlab/course-progress/internal/domain"

type CreateUserInput struct {
	UserName string
	IsAdmin  bool
}

type CreateUserOutput struct {
	domain.User
}
