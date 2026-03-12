package commands

import (
	"errors"
	"regexp"
	"strings"

	"github.com/cybernetlab/swimming-search/internal/controller/bot"
	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/cybernetlab/swimming-search/internal/dto"
	"github.com/cybernetlab/swimming-search/internal/usecase"
)

func users(uc *usecase.UseCase, cmd *bot.Command) {
	// get current user
	user, err := cmd.User()
	if err != nil || !user.Authorize("user:create") {
		cmd.Error("You are not authorized for this command", err, "cmd.User")
	}
	if cmd.Args == "" {
		// "/users" command
		users, err := uc.GetUsers(cmd.Context())
		if err != nil {
			cmd.Error("Can't retrieve users list", err, "uc.GetUsers")
			return
		}
		cmd.Reply(formatUsers(users))
		return
	}
	args := regexp.MustCompile(`\s+`).Split(cmd.Args, 3)
	if len(args) < 2 {
		help(uc, cmd)
		return
	}
	name := args[1]
	name = strings.TrimLeft(name, "@")
	switch args[0] {
	case "add":
		// "/users add NAME [admin]" command
		isAdmin := len(args) >= 3 && args[3] == "admin"
		_, err := uc.CreateUser(cmd.Context(), dto.CreateUserInput{UserName: name, IsAdmin: isAdmin})
		if err != nil {
			var existsErr domain.ErrAlreadyExists
			if errors.As(err, &existsErr) {
				cmd.Replyf("The user %s is already exists", name)
				return
			}
			cmd.Error("Can't create user", err, "uc.CreateUser")
			return
		}
		cmd.Replyf("User %s succesfully created", name)
		return
	case "delete":
		// "/users delete NAME" command
		if name == user.Name {
			cmd.Reply("You can't delete yourself")
			return
		}
		_, err := uc.DeleteUser(cmd.Context(), dto.DeleteUserInput{UserName: name})
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				cmd.Replyf("There are no user %s", name)
				return
			}
			cmd.Error("Can't create user", err, "uc.CreateUser")
			return
		}
		cmd.Replyf("User %s successfully deleted", name)
		return
	}
	help(uc, cmd)
}

func formatUsers(users []domain.User) string {
	msg := "List of users:\n"
	for _, user := range users {
		msg += formatUser(user) + "\n"
	}
	return msg
}

func formatUser(user domain.User) string {
	text := user.Name
	if user.IsAdmin {
		text += " (admin)"
	}
	return text
}
