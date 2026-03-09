package commands

import (
	"errors"
	"fmt"

	"github.com/cybernetlab/course-progress/internal/controller/bot"
	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/internal/dto"
	"github.com/cybernetlab/course-progress/internal/usecase"
)

func stop(uc *usecase.UseCase, cmd *bot.Command) {
	user, err := cmd.User()
	if err != nil {
		cmd.Error("Can't get active search", err, "cmd.User")
		return
	}
	search, err := uc.StopSearch(cmd.Context(), dto.StopSearchInput{UserName: user.Name})
	if errors.Is(err, domain.ErrNotFound) {
		cmd.Reply("There are no active searches")
		return
	}
	if err != nil {
		cmd.Error("Can't get active search", err, "uc.StopSearch")
	}
	msg := fmt.Sprintf("%s was stopped", filterToString(search.Search))
	cmd.Reply(msg)
}
