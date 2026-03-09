package commands

import (
	"github.com/cybernetlab/course-progress/internal/controller/bot"
	"github.com/cybernetlab/course-progress/internal/usecase"
)

const usage = `
I can help you with the following commands:

<b>/start centre dayOfWeek stringToSearch</b> - Start search for courses.

Where

<i>centre:</i> list of centre IDs separated with ',' (please don't add extra
spaces between IDs) or 'all' word (In case of 'all' only user selected
centres will be matched - see '/centres' command).

<i>dayOfWeek:</i> day of week in short (both 'Mon' and 'Mo' supported) or full 
('Monday') form or a list of week days separated with ',' (please don't add
extra spaces between days).

<i>stringToSearch:</i> a string that should be present in coutrse name to match.

Examples:
<code>
/start Friday Stage 5
/start Mon,Tu Stage 3
</code>

<b>/current</b> - Show current search

<b>/stop</b> - Stop current search

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
`

func help(_ *usecase.UseCase, cmd *bot.Command) {
	text := usage
	user, err := cmd.User()
	if err != nil {
		cmd.Reply(text + helpUsage)
		cmd.Error("You are not authorized for this command", err, "cmd.User")
	}
	if user.Authorize("users:create") {
		text += adminUsage
	}
	cmd.Reply(text + helpUsage)
}
