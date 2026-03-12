package dto

import "github.com/cybernetlab/swimming-search/internal/domain"

type DeleteUserCentreInput struct {
	User     *domain.User
	CentreID uint
}

type DeleteUserCentreOutput []domain.Centre
