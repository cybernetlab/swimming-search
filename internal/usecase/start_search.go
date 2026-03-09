package usecase

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/internal/dto"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) StartSearch(ctx context.Context, input dto.StartSearchInput) error {
	for _, search := range u.searches {
		if search.UserName == input.UserName {
			return domain.NewErrAlreadyExists("search")
		}
	}
	if len(input.CentreIDs) == 0 {
		return domain.ErrEmptyCentreIDs
	}
	err := u.store.PutSearch(ctx, input.Search)
	if err != nil {
		return fmt.Errorf("u.store.PutSearch: %w", err)
	}

	searchCtx, cancel := context.WithCancel(ctx)
	input.Search.Cancel = cancel
	u.searches = append(u.searches, input.Search)
	context.AfterFunc(searchCtx, func() {
		for i, search := range u.searches {
			if search.UserName == input.UserName {
				log.Debug().Msgf("Deleting search for user: %s", search.UserName)
				u.searches = slices.Delete(u.searches, i, i+1)
				err = u.store.DeleteSearch(ctx, search)
				if err != nil {
					log.Error().Err(err).Msg("Can't delete search")
				}
			}
		}
	})

	courses := make(chan domain.Course, 10)
	go func() {
		for {
			course, ok := <-courses
			if !ok {
				cancel()
				return
			}
			msg := fmt.Sprintf("Course found: %s", courseToString(course))
			err := u.bot.Send(ctx, msg)
			if err != nil {
				log.Error().Err(err).Msg("u.bot.Send")
			}
		}
	}()

	u.booking.StartSearchCourses(ctx, input.Search, courses)
	return nil
}

func courseToString(c domain.Course) string {
	sched := schedToString(c.Schedule)
	return fmt.Sprintf("<b>'%s'</b> on %s, spaces left: <b>%d</b>", c.Name, sched, c.Availability.Spaces.Free)
}

func schedToString(s domain.CourseSchedule) string {
	if s.Type == "dayOfWeek" {
		t := ""
		if s.Time.Start != "" {
			x := strings.Split(s.Time.Start, ":")
			if len(x) > 1 {
				t = fmt.Sprintf("%s:%s", x[0], x[1])
			}
		}
		if s.Time.End != "" {
			x := strings.Split(s.Time.End, ":")
			if len(x) > 1 {
				t = fmt.Sprintf("%s-%s:%s", t, x[0], x[1])
			}
		}
		return fmt.Sprintf("%s %s", s.DayOfWeek, t)
	}
	return s.DateDescr
}
