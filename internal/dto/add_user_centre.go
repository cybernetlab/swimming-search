package dto

import "github.com/cybernetlab/course-progress/internal/domain"

type AddUserCentreInput struct {
	User     *domain.User
	CentreID uint
}

type AddUserCentreOutput []domain.Centre
