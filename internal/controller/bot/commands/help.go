package commands

import (
	"github.com/cybernetlab/swimming-search/internal/controller/bot"
	"github.com/cybernetlab/swimming-search/internal/usecase"
)

const usageHeader = `
I can help you with the following commands:
`

const startUsage = `
<b>/start [CENTRES] DAY_OF_WEEK SEARCH_QUERY</b> - Start search for courses.

Where

<i>CENTRES:</i> list of centre IDs separated with ',' (please don't add extra
spaces between IDs) or 'all' word (In case of 'all' only user selected
centres will be matched - see '/centres' command). This is an optional
parameter. It have 'all' value by default

<i>DAY_OF_WEEK:</i> day of week in short (both 'Mon' and 'Mo' supported) or full 
('Monday') form or a list of week days separated with ',' (please don't add
extra spaces between days).

<i>SEARCH_QUERY:</i> a string that should be present in course name to match.

Examples:
<code>
/start Friday Stage 5
/start Mon,Tu Stage 3
/start 3,5 Fri,Sat Stage 1
</code>
`

const currentUsage = `
<b>/current</b> - Show current search
`

const stopUsage = `
<b>/stop</b> - Stop current search
`

const centresUsage = `
<b>/centres</b> - List of user selected centres

<b>/centres all</b> - List of all available centres

<b>/centres add ID</b> - Add centre to user selection

<b>/centres delete ID</b> - Remove centre from user selection
`

const adminUsage = `

<b>/users</b> - Show users list

<b>/users add NAME</b> - Add user with telegram name NAME

<b>/users add NAME admin</b> - Add admin with telegram name NAME

<b>/users delete NAME</b> - Delete user with telegram name NAME
`

const helpUsage = `

<b>/help</b> - Display this help message

<b>/help COMMAND</b> - Display help about specified command
`

func help(_ *usecase.UseCase, cmd *bot.Command) {
	text := usageHeader
	isAdmin := false
	user, err := cmd.User()
	if err != nil {
		cmd.Reply(text + helpUsage)
		cmd.Error("You are not authorized for this command", err, "cmd.User")
		return
	}
	if user.Authorize("user:create") {
		isAdmin = true
	}
	switch cmd.Args {
	case "start":
		text += startUsage
	case "current":
		text += currentUsage
	case "stop":
		text += stopUsage
	case "centres":
		text += centresUsage
	case "admin":
		if isAdmin {
			text += adminUsage
		}
	default:
		text += startUsage + currentUsage + stopUsage + centresUsage
		if isAdmin {
			text += adminUsage
		}
		text += helpUsage
	}
	cmd.Reply(text)
}
