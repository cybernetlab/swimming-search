package dto

import "github.com/cybernetlab/course-progress/internal/domain"

type GetSearchInput struct {
	UserName string
}

type GetSearchOutput struct {
	domain.Search
}
