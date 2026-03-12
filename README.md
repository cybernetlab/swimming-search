![Tests](https://github.com)
![Coverage](https://img.shields.io)
![Go Version](https://img.shields.io)

Example telegram bot that can crawl some external HTTP API for swimming
courses schedule and report if there is any courses that satisfies search
conditions like selected day of week, matching course name pattern etc.

This bot was creatred for testing purposes only.

## TODO

- [x] add timeout to HTTP requests
- [ ] fix bot.Cmd.Error to print correct error place
- [x] add unit tests
- [ ] add integration tests
- [x] add centre option in /start command into help
- [x] help for separate commands
- [x] add centre name in search output
- [ ] add tracing
- [ ] add metrics