package commands

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/cybernetlab/course-progress/internal/controller/bot"
	"github.com/cybernetlab/course-progress/internal/domain"
	"github.com/cybernetlab/course-progress/internal/dto"
	"github.com/cybernetlab/course-progress/internal/usecase"
	"github.com/rs/zerolog/log"
)

func start(uc *usecase.UseCase, cmd *bot.Command) {
	ctx := cmd.Context()
	args := regexp.MustCompile(`\s+`).Split(cmd.Args, 2)
	centreIDs, days, query, ok := parseStartArgs(args)
	if ok {
		user, err := cmd.User()
		if err != nil {
			cmd.Error("Can't start search", err, "cmd.User")
			return
		}
		if len(centreIDs) == 0 {
			centreIDs = user.CentreIDs
		}
		chatID, err := cmd.ChatID()
		if err != nil {
			cmd.Error("Can't start search", err, "cmd.ChatID")
			return
		}
		search := domain.Search{
			NameQuery: query,
			Days:      days,
			CentreIDs: centreIDs,
			UserName:  user.Name,
			ChatID:    chatID,
		}
		err = uc.StartSearch(ctx, dto.StartSearchInput{Search: search})
		if err != nil {
			if errors.Is(err, domain.ErrEmptyCentreIDs) {
				cmd.Reply(
					"Center IDs wasn't specified and where are no selected center IDs\n" +
						"Either specify IDs in `/start` command or select with `centres` command",
				)
				return
			}
			cmd.Error("Can't start search", err, "uc.StartSearch")
			return
		}
		msg := fmt.Sprintf("Searching for '%s' on %s", query, daysToString(days))
		centres, err := uc.GetCentres(ctx)
		if err == nil {
			centres = slices.DeleteFunc(centres, func(c domain.Centre) bool { return slices.Index(centreIDs, c.ID) < 0 })
			msg += " in following centres:\n" + formatCentres(centres)
		} else {
			log.Error().Err(err).Msg("uc.GetCentres")
		}
		cmd.Reply(msg)
	} else {
		help(uc, cmd)
	}
}

func parseStartArgs(args []string) ([]uint, []uint8, string, bool) {
	centres := []uint{}
	days := []uint8{}
	if len(args) == 0 || len(args[0]) == 0 {
		return centres, days, "", false
	}
	if unicode.IsDigit(rune(args[0][0])) {
		for _, centre := range strings.Split(args[0], ",") {
			id, err := strconv.Atoi(centre)
			if err != nil || id < 0 {
				log.Debug().Msg("2")
				return centres, days, "", false
			}
			centres = append(centres, uint(id))
		}
		slices.Sort(centres)
		centres = slices.Compact(centres)
		args = regexp.MustCompile(`\s+`).Split(args[1], 2)
	}
	for _, day := range strings.Split(args[0], ",") {
		switch strings.ToLower(day) {
		case "mo", "mon", "monday":
			days = append(days, 1)
		case "tu", "tue", "tuesday":
			days = append(days, 2)
		case "we", "wed", "wednesday":
			days = append(days, 3)
		case "th", "thu", "thursday":
			days = append(days, 4)
		case "fr", "fri", "friday":
			days = append(days, 5)
		case "sa", "sat", "saturday":
			days = append(days, 6)
		case "su", "sun", "sunday":
			days = append(days, 7)
		default:
		}
	}
	slices.Sort(days)
	days = slices.Compact(days)
	if len(days) == 0 {
		return centres, days, "", false
	}
	query := strings.Join(args[1:], " ")
	return centres, days, query, true
}
