package dto

import "github.com/cybernetlab/course-progress/internal/domain"

type StopSearchInput struct {
	UserName string
}

type StopSearchOutput struct {
	domain.Search
}
