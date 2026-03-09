package commands

import (
	"fmt"

	"github.com/cybernetlab/course-progress/internal/controller/bot"
	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/internal/usecase"
)

func Handler(uc *usecase.UseCase, cmd *bot.Command) {
	switch cmd.Name {
	case "start":
		start(uc, cmd)
	case "current":
		current(uc, cmd)
	case "stop":
		stop(uc, cmd)
	case "centres":
		centres(uc, cmd)
	case "users":
		users(uc, cmd)
	case "help":
		help(uc, cmd)
	}
}

// private helper functions

func formatCentres(centres []domain.Centre) string {
	msg := ""
	for _, centre := range centres {
		msg += formatCentre(centre) + "\n"
	}
	return msg
}

func formatCentre(centre domain.Centre) string {
	return fmt.Sprintf("<b>%d</b> %s", centre.ID, centre.Name)
}
