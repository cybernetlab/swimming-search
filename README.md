![Tests](https://github.com/cybernetlab/swimming-search/actions/workflows/build.yml/badge.svg)
![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/cybernetlab/b0d91574542f919600a3af6e1c74223a/raw/coverage.json)
![Go Version](https://img.shields.io/github/go-mod/go-version/cybernetlab/swimming-search)

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