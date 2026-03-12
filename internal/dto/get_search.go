package dto

import "github.com/cybernetlab/swimming-search/internal/domain"

type GetSearchInput struct {
	UserName string
}

type GetSearchOutput struct {
	domain.Search
}
