package dto

import "github.com/cybernetlab/swimming-search/internal/domain"

type StopSearchInput struct {
	UserName string
}

type StopSearchOutput struct {
	domain.Search
}
