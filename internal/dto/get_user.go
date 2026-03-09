package dto

import "github.com/cybernetlab/course-progress/internal/domain"

type GetUserInput struct {
	UserName string
}

type GetUserOutput struct {
	domain.User
}
