package dto

import "github.com/cybernetlab/course-progress/internal/domain"

type DeleteUserCentreInput struct {
	User     *domain.User
	CentreID uint
}

type DeleteUserCentreOutput []domain.Centre
