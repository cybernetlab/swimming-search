package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/cybernetlab/course-progress/internal/controller/bot"
	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/internal/dto"
	"github.com/cybernetlab/course-progress/internal/usecase"
)

func centres(uc *usecase.UseCase, cmd *bot.Command) {
	// get current user
	user, err := cmd.User()
	if err != nil {
		cmd.Error("You are not authorized for this command", err, "cmd.User")
	}
	if cmd.Args == "all" {
		// "/centres all" command
		centres, err := uc.GetCentres(cmd.Context())
		if err != nil {
			cmd.Error("Can't retrieve centers list", err, "uc.GetCentres")
			return
		}
		msg := "List of all available centres:\n"
		for _, centre := range centres {
			msg += formatCentre(centre) + "\n"
		}
		cmd.Reply(msg)
		return
	}
	if cmd.Args == "" {
		// "/centres" command
		centres, err := uc.GetUserCentres(cmd.Context(), dto.GetUserCentresInput{User: user})
		if err != nil {
			cmd.Error("Can't retrieve selected centers list", err, "uc.GetUserCentres")
			return
		}
		if len(centres) == 0 {
			cmd.Reply("You have no selected centres yet. Select one with `/centres select ID` command")
			return
		}
		cmd.Reply("List of selected centres:\n" + formatCentres(centres))
		return
	}
	args := regexp.MustCompile(`\s+`).Split(cmd.Args, 2)
	if len(args) < 2 {
		help(uc, cmd)
		return
	}
	id, err := strconv.Atoi(args[1])
	if err != nil || id < 0 {
		cmd.Reply("Invalid centre ID")
		return
	}
	switch args[0] {
	case "add":
		// "/centres add ID" command
		centres, err := uc.AddUserCentre(cmd.Context(), dto.AddUserCentreInput{User: user, CentreID: uint(id)})
		if err != nil {
			var existsErr domain.ErrAlreadyExists
			if errors.As(err, &existsErr) {
				cmd.Replyf("The centre with ID %d is already selected", id)
				return
			}
			if errors.Is(err, domain.ErrNotFound) {
				cmd.Replyf("There are no centre with ID %d", id)
				return
			}
			cmd.Error("Can't add to selected centres", err, "uc.AddUserCentre")
		}
		cmd.Reply("List of selected centres:\n" + formatCentres(centres))
		return
	case "delete":
		// "/centres delete ID" command
		centres, err := uc.DeleteUserCentre(cmd.Context(), dto.DeleteUserCentreInput{User: user, CentreID: uint(id)})
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				cmd.Reply(fmt.Sprintf("There are no selected centre with ID %d", id))
				return
			}
			cmd.Error("Can't delete from selected centres", err, "uc.DeleteUserCentre")
		}
		cmd.Reply("List of selected centres:\n" + formatCentres(centres))
		return
	}
	help(uc, cmd)
}
