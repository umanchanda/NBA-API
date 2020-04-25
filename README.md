# NBA-API

To get started:

`go get -u github.com/gorilla/mux`

`go get -u github.com/PuerkitoBio/goquery`

run: `go run main.go` in your terminal

Head to `localhost:8000` in your browser

To view all the boxscores for a given date, head over to `/boxscore/year/month/date`

To view team totals for a given game, head over to `/boxscore/year/month/date/awayTeam/homeTeam`

Note that `awayTeam` and `homeTeam` must be the three letter code for the team

So, for example, the `New York Knicks` are `NYK`, the Cleveland Cavaliers, are `CLE`, and so on so forth

For example, if I wanted to view the boxscores for February 11th, 2020, I would head to `localhost:8000/boxscore/2020/02/11`
