package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cybernetlab/course-progress/internal/controller/bot"
	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/internal/dto"
	"github.com/cybernetlab/course-progress/internal/usecase"
)

func current(uc *usecase.UseCase, cmd *bot.Command) {
	ctx := cmd.Context()
	user, err := cmd.User()
	if err != nil {
		cmd.Error("Can't get current search", err, "cmd.User")
		return
	}
	search, err := uc.GetSearch(ctx, dto.GetSearchInput{UserName: user.Name})
	if errors.Is(err, domain.ErrNotFound) {
		cmd.Reply("There are no active searches")
		return
	}
	if err != nil {
		cmd.Error("Can't get current search", err, "uc.GetSearch")
	}
	msg := fmt.Sprintf("Now we are %s", filterToString(search.Search))
	cmd.Reply(msg)
}

func filterToString(s domain.Search) string {
	return fmt.Sprintf("searching for '%s' on %s", s.NameQuery, daysToString(s.Days))
}

func daysToString(days []uint8) string {
	daysStr := []string{}
	for _, day := range days {
		switch day {
		case 1:
			daysStr = append(daysStr, "Monday")
		case 2:
			daysStr = append(daysStr, "Tuesday")
		case 3:
			daysStr = append(daysStr, "Wednesday")
		case 4:
			daysStr = append(daysStr, "Thursday")
		case 5:
			daysStr = append(daysStr, "Friday")
		case 6:
			daysStr = append(daysStr, "Saturday")
		case 7:
			daysStr = append(daysStr, "Sunday")
		}
	}
	return strings.Join(daysStr, " or ")
}
