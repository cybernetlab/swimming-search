package usecase

import (
	"context"
	"slices"

	"github.com/cybernetlab/swimming-search/internal/domain"
)

//go:generate mockery

type Store interface {
	GetUser(ctx context.Context, name string) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	PutUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, name string) error
	GetCentres(ctx context.Context) ([]domain.Centre, error)
	PutCentres(ctx context.Context, centres []domain.Centre) error
	GetSearch(ctx context.Context, userName string) (domain.Search, domain.NodeID, error)
	GetSearches(ctx context.Context, nodeID domain.NodeID) ([]domain.Search, error)
	PutSearch(ctx context.Context, search domain.Search) error
	DeleteSearch(ctx context.Context, search domain.Search) error
}

type Booking interface {
	GetCentres(ctx context.Context) ([]domain.Centre, error)
	StartSearchCourses(ctx context.Context, search domain.Search, out chan<- domain.Course)
}

type Bot interface {
	Send(ctx context.Context, text string) error
}

type UseCase struct {
	store    Store
	booking  Booking
	bot      Bot
	searches []domain.Search
}

func New(store Store, booking Booking, bot Bot) *UseCase {
	return &UseCase{
		store:    store,
		booking:  booking,
		bot:      bot,
		searches: []domain.Search{},
	}
}

// private helper functions

func findCentre(id uint, centres []domain.Centre) (domain.Centre, bool) {
	idx := slices.IndexFunc(centres, func(c domain.Centre) bool { return c.ID == id })
	if idx >= 0 {
		return centres[idx], true
	}
	return domain.Centre{}, false
}

func collectCentres(ids []uint, centres []domain.Centre) []domain.Centre {
	userCentres := []domain.Centre{}
	for _, id := range ids {
		centre, ok := findCentre(id, centres)
		if ok {
			userCentres = append(userCentres, centre)
		}
	}
	return userCentres
}
